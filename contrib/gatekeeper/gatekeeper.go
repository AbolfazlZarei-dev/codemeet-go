package gatekeeper

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

const (
	cbPrefix = "cm_cap_"

	defaultChallengeTimeout = 120 * time.Second
	defaultWrongAnswers     = 2
	defaultVerifiedTTL      = 24 * time.Hour

	defaultWorkerCount = 16
	defaultQueueSize   = 256
)

const (
	CaptchaMath    = "math"
	CaptchaNumbers = "numbers"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type CaptchaConfig struct {
	Type        string
	Options     int
	MinNumber   int
	MaxNumber   int
	Title       string
	CorrectText string
	WrongText   string
}

type Config struct {
	BotID             string
	ChallengeTimeout  time.Duration
	WrongAnswersLimit int
	VerifiedTTL       time.Duration

	Captcha CaptchaConfig

	WorkerCount int
	QueueSize   int

	SendCaptchaAction    func(ctx context.Context, chatID, userID, text string, keyboard *models.InlineKeyboardMarkup, replyToMessageID int) (int, error)
	EditMessageAction    func(ctx context.Context, chatID string, messageID int, text string, markup *models.InlineKeyboardMarkup) error
	AnswerCallbackAction func(ctx context.Context, callbackID, text string, showAlert bool) error
	VerifyAction         func(ctx context.Context, chatID, userID string)
	KickAction           func(ctx context.Context, chatID, userID string)
	DeleteMessageAction  func(ctx context.Context, chatID string, messageID int)
}

func DefaultConfig() Config {
	return Config{
		ChallengeTimeout:  defaultChallengeTimeout,
		WrongAnswersLimit: defaultWrongAnswers,
		VerifiedTTL:       defaultVerifiedTTL,
		WorkerCount:       defaultWorkerCount,
		QueueSize:         defaultQueueSize,
		Captcha: CaptchaConfig{
			Type:        CaptchaMath,
			Options:     4,
			MinNumber:   2,
			MaxNumber:   9,
			Title:       "🤖 برای تأیید انسان بودن، پاسخ درست را انتخاب کنید:",
			CorrectText: "پاسخ درست",
			WrongText:   "پاسخ اشتباه",
		},
	}
}

type pendingUser struct {
	chatID        string
	userID        string
	correctAnswer int
	expiresAt     time.Time
	wrongAnswers  atomic.Int32
	captchaMsgID  int
}

type stats struct {
	challengesSent   atomic.Int64
	challengesPassed atomic.Int64
	challengesFailed atomic.Int64
	usersKicked      atomic.Int64
}

type challengeJob struct {
	ctx              context.Context
	chatID           string
	userID           string
	user             *models.User
	config           CaptchaConfig
	replyToMessageID int
}

type Gatekeeper struct {
	cfg       Config
	pending   sync.Map
	verified  sync.Map
	inflight  sync.Map
	stats     stats
	jobs      chan challengeJob
	workers   sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	closeOnce sync.Once
}

func New(cfg Config) *Gatekeeper {
	def := DefaultConfig()

	if cfg.ChallengeTimeout <= 0 {
		cfg.ChallengeTimeout = def.ChallengeTimeout
	}
	if cfg.WrongAnswersLimit <= 0 {
		cfg.WrongAnswersLimit = def.WrongAnswersLimit
	}
	if cfg.VerifiedTTL <= 0 {
		cfg.VerifiedTTL = def.VerifiedTTL
	}
	if cfg.WorkerCount <= 0 {
		cfg.WorkerCount = def.WorkerCount
	}
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = def.QueueSize
	}

	cfg.Captcha = normalizeCaptchaConfig(def.Captcha, cfg.Captcha)

	ctx, cancel := context.WithCancel(context.Background())

	gk := &Gatekeeper{
		cfg:    cfg,
		jobs:   make(chan challengeJob, cfg.QueueSize),
		ctx:    ctx,
		cancel: cancel,
	}

	for i := 0; i < cfg.WorkerCount; i++ {
		gk.workers.Add(1)
		go gk.worker()
	}

	go gk.cleanupLoop()

	return gk
}

func (gk *Gatekeeper) worker() {
	defer gk.workers.Done()
	for {
		select {
		case job := <-gk.jobs:
			_ = gk.sendChallenge(job.ctx, job.chatID, job.userID, job.user, job.config, job.replyToMessageID)
		case <-gk.ctx.Done():
			return
		}
	}
}

func (gk *Gatekeeper) Middleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			if u == nil {
				next(ctx, u)
				return
			}

			// مدیریت ورود کاربر جدید
			if u.Message != nil && len(u.Message.NewChatMembers) > 0 {
				chatID := u.Message.Chat.ID
				msgID := u.Message.MessageID
				for i := range u.Message.NewChatMembers {
					user := &u.Message.NewChatMembers[i]
					if !user.IsBot {
						// حذف پیام "کاربر جدید به گروه پیوست"
						if gk.cfg.DeleteMessageAction != nil {
							gk.cfg.DeleteMessageAction(ctx, chatID, msgID)
						}
						gk.enqueueChallenge(chatID, user.ID, user, gk.cfg.Captcha, 0)
					}
				}
				return
			}

			// مدیریت کلیک روی دکمه کپچا
			if u.CallbackQuery != nil && strings.HasPrefix(u.CallbackQuery.Data, cbPrefix) {
				gk.handleButton(ctx, u.CallbackQuery)
				return
			}

			// بررسی پیام‌های معمولی
			if u.Message != nil && u.Message.From != nil {
				chatID := u.Message.Chat.ID
				userID := u.Message.From.ID
				msgID := u.Message.MessageID

				// جلوگیری از لوپ ربات
				if userID == gk.cfg.BotID {
					next(ctx, u)
					return
				}

				key := makeKey(chatID, userID)

				// اگر کاربر قبلاً تأیید شده باشد
				if value, ok := gk.verified.Load(key); ok {
					if expiresAt, valid := value.(time.Time); valid && time.Now().Before(expiresAt) {
						next(ctx, u)
						return
					}
					gk.verified.Delete(key)
				}

				// اگر کاربر در حال انجام کپچا باشد
				if value, ok := gk.pending.Load(key); ok {
					pu := value.(*pendingUser)

					// 1. پیام کاربر را حذف کن
					if gk.cfg.DeleteMessageAction != nil {
						gk.cfg.DeleteMessageAction(ctx, chatID, msgID)
					}

					// 2. پیام کپچای قبلی را حذف کن
					if pu.captchaMsgID != 0 && gk.cfg.DeleteMessageAction != nil {
						gk.cfg.DeleteMessageAction(ctx, chatID, pu.captchaMsgID)
					}

					// 3. یک کپچای جدید برایش بفرست
					gk.enqueueChallenge(chatID, userID, u.Message.From, gk.cfg.Captcha, 0)
					return
				}

				if gk.cfg.DeleteMessageAction != nil {
					gk.cfg.DeleteMessageAction(ctx, chatID, msgID)
				}
				gk.enqueueChallenge(chatID, userID, u.Message.From, gk.cfg.Captcha, 0)
				return
			}

			next(ctx, u)
		}
	}
}

func (gk *Gatekeeper) enqueueChallenge(chatID, userID string, user *models.User, config CaptchaConfig, replyToMessageID int) {
	if chatID == "" || userID == "" {
		return
	}

	key := makeKey(chatID, userID)

	if _, ok := gk.verified.Load(key); ok {
		return
	}

	if _, loaded := gk.inflight.LoadOrStore(key, struct{}{}); loaded {
		return
	}

	job := challengeJob{
		ctx:              context.Background(),
		chatID:           chatID,
		userID:           userID,
		user:             user,
		config:           config,
		replyToMessageID: replyToMessageID,
	}

	select {
	case gk.jobs <- job:
		return
	default:
		gk.inflight.Delete(key)
	}
}

func (gk *Gatekeeper) sendChallenge(ctx context.Context, chatID, userID string, user *models.User, config CaptchaConfig, replyToMessageID int) error {
	key := makeKey(chatID, userID)
	defer gk.inflight.Delete(key)

	if err := ctx.Err(); err != nil {
		return err
	}

	switch config.Type {
	case CaptchaMath:
		return gk.sendMathCaptcha(ctx, chatID, userID, key, user, config, replyToMessageID)
	case CaptchaNumbers:
		return gk.sendNumbersCaptcha(ctx, chatID, userID, key, user, config, replyToMessageID)
	default:
		return errors.New("unsupported captcha type")
	}
}

func (gk *Gatekeeper) sendMathCaptcha(ctx context.Context, chatID, userID, key string, user *models.User, config CaptchaConfig, replyToMessageID int) error {
	a := randomInt(config.MinNumber, config.MaxNumber)
	b := randomInt(config.MinNumber, config.MaxNumber)

	correctAnswer := a + b
	title := config.Title
	if title == "" {
		title = "🤖 برای تأیید، پاسخ درست را انتخاب کنید:"
	}

	userName := "کاربر"
	if user != nil {
		userName = user.FullName()
	}

	// متن شامل نام و آیدی کاربر
	text := fmt.Sprintf("👤 کاربر گرامی: <b>%s</b>\n🆔 آیدی: <code>%s</code>\n\n%s\n\n<b>%d + %d = ?</b>", userName, userID, title, a, b)

	options := buildMathOptions(correctAnswer, config.Options)

	buttons := make([]models.InlineKeyboardButton, 0, len(options))
	for _, answer := range options {
		buttons = append(buttons, models.InlineKeyboardButton{
			Text:         strconv.Itoa(answer),
			CallbackData: cbPrefix + strconv.Itoa(answer),
		})
	}

	keyboard := models.NewInlineKeyboard(models.InlineRow(buttons...))

	return gk.storeChallenge(ctx, chatID, userID, key, user, correctAnswer, text, keyboard, replyToMessageID)
}

func (gk *Gatekeeper) sendNumbersCaptcha(ctx context.Context, chatID, userID, key string, user *models.User, config CaptchaConfig, replyToMessageID int) error {
	correctAnswer := randomInt(config.MinNumber, config.MaxNumber)
	title := config.Title
	if title == "" {
		title = "🔢 عدد درست را انتخاب کنید:"
	}

	userName := "کاربر"
	if user != nil {
		userName = user.FullName()
	}

	text := fmt.Sprintf("👤 کاربر گرامی: <b>%s</b>\n🆔 آیدی: <code>%s</code>\n\n%s\n\n<b>عدد: %d</b>", userName, userID, title, correctAnswer)

	options := buildNumberOptions(correctAnswer, config.Options, config.MinNumber, config.MaxNumber)

	buttons := make([]models.InlineKeyboardButton, 0, len(options))
	for _, answer := range options {
		buttons = append(buttons, models.InlineKeyboardButton{
			Text:         strconv.Itoa(answer),
			CallbackData: cbPrefix + strconv.Itoa(answer),
		})
	}

	keyboard := models.NewInlineKeyboard(models.InlineRow(buttons...))

	return gk.storeChallenge(ctx, chatID, userID, key, user, correctAnswer, text, keyboard, replyToMessageID)
}

func (gk *Gatekeeper) storeChallenge(ctx context.Context, chatID, userID, key string, user *models.User, correctAnswer int, text string, keyboard *models.InlineKeyboardMarkup, replyToMessageID int) error {
	var (
		messageID int
		err       error
	)

	if gk.cfg.SendCaptchaAction != nil {
		messageID, err = gk.cfg.SendCaptchaAction(ctx, chatID, userID, text, keyboard, replyToMessageID)
		if err != nil {
			return err
		}
	}

	// حفظ تعداد پاسخ‌های اشتباه قبلی کاربر
	var wrongAnswers int32 = 0
	if val, exists := gk.pending.Load(key); exists {
		wrongAnswers = val.(*pendingUser).wrongAnswers.Load()
		gk.pending.Delete(key)
	}

	pu := &pendingUser{
		chatID:        chatID,
		userID:        userID,
		correctAnswer: correctAnswer,
		expiresAt:     time.Now().Add(gk.cfg.ChallengeTimeout),
		captchaMsgID:  messageID,
	}
	pu.wrongAnswers.Store(wrongAnswers)

	if _, loaded := gk.pending.LoadOrStore(key, pu); loaded {
		val, _ := gk.pending.Load(key)
		existing := val.(*pendingUser)
		existing.captchaMsgID = messageID
		existing.correctAnswer = correctAnswer
		existing.expiresAt = time.Now().Add(gk.cfg.ChallengeTimeout)
	}

	gk.stats.challengesSent.Add(1)
	return nil
}

func (gk *Gatekeeper) handleButton(ctx context.Context, cq *models.CallbackQuery) {
	if cq == nil || cq.From == nil || cq.Message == nil {
		return
	}

	data := cq.Data
	if !strings.HasPrefix(data, cbPrefix) {
		return
	}

	answerText := strings.TrimPrefix(data, cbPrefix)
	answer, err := strconv.Atoi(answerText)
	if err != nil {
		gk.answerCallback(ctx, cq.ID, "دکمه معتبر نیست.", true)
		return
	}

	chatID := cq.Message.Chat.ID
	userID := cq.From.ID
	key := makeKey(chatID, userID)

	value, ok := gk.pending.Load(key)
	if !ok {
		gk.answerCallback(ctx, cq.ID, "این کپچا منقضی شده است.", true)
		return
	}

	pu, ok := value.(*pendingUser)
	if !ok || pu == nil {
		gk.pending.Delete(key)
		gk.answerCallback(ctx, cq.ID, "اطلاعات کپچا پیدا نشد.", true)
		return
	}

	if !time.Now().Before(pu.expiresAt) {
		gk.failChallenge(ctx, key, pu)
		gk.answerCallback(ctx, cq.ID, "⏰ زمان کپچا تمام شده است.", true)
		return
	}

	// اگر پاسخ درست بود
	if answer == pu.correctAnswer {
		gk.pending.Delete(key)
		gk.verified.Store(key, time.Now().Add(gk.cfg.VerifiedTTL))
		gk.stats.challengesPassed.Add(1)
		gk.answerCallback(ctx, cq.ID, "✅ تأیید شدید.", false)

		// پیام کپچا ادیت شود و دکمه‌ها حذف شوند
		if gk.cfg.EditMessageAction != nil {
			userName := "کاربر"
			if cq.From != nil {
				userName = cq.From.FullName()
			}
			newText := fmt.Sprintf("✅ %s عزیز، احراز هویت شما با موفقیت انجام شد.\nاکنون می‌توانید در گروه پیام ارسال کنید.", userName)
			emptyMarkup := &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{}}
			_ = gk.cfg.EditMessageAction(ctx, chatID, pu.captchaMsgID, newText, emptyMarkup)
		}

		if gk.cfg.VerifyAction != nil {
			gk.cfg.VerifyAction(ctx, chatID, userID)
		}
		return
	}

	// اگر پاسخ اشتباه بود
	count := pu.wrongAnswers.Add(1)
	if int(count) >= gk.cfg.WrongAnswersLimit {
		gk.failChallenge(ctx, key, pu)
		if pu.captchaMsgID != 0 && gk.cfg.DeleteMessageAction != nil {
			gk.cfg.DeleteMessageAction(ctx, pu.chatID, pu.captchaMsgID)
		}
	} else {
		gk.answerCallback(ctx, cq.ID, "❌ پاسخ اشتباه است. کپچای جدید برای شما ارسال شد.", false)

		// پیام کپچای قبلی حذف شود
		if pu.captchaMsgID != 0 && gk.cfg.DeleteMessageAction != nil {
			gk.cfg.DeleteMessageAction(ctx, pu.chatID, pu.captchaMsgID)
		}

		// کپچای جدید ارسال شود
		gk.enqueueChallenge(pu.chatID, pu.userID, cq.From, gk.cfg.Captcha, 0)
	}
}

func (gk *Gatekeeper) answerCallback(ctx context.Context, callbackID, text string, showAlert bool) {
	if gk.cfg.AnswerCallbackAction == nil {
		return
	}
	_ = gk.cfg.AnswerCallbackAction(ctx, callbackID, text, showAlert)
}

func (gk *Gatekeeper) failChallenge(ctx context.Context, key string, pu *pendingUser) {
	if pu == nil {
		return
	}
	gk.pending.Delete(key)

	gk.stats.challengesFailed.Add(1)
	gk.stats.usersKicked.Add(1)

	if pu.captchaMsgID != 0 && gk.cfg.DeleteMessageAction != nil {
		gk.cfg.DeleteMessageAction(ctx, pu.chatID, pu.captchaMsgID)
	}

	if gk.cfg.KickAction != nil {
		gk.cfg.KickAction(ctx, pu.chatID, pu.userID)
	}
}

func (gk *Gatekeeper) cleanupLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			gk.cleanup(now)
		case <-gk.ctx.Done():
			return
		}
	}
}

func (gk *Gatekeeper) cleanup(now time.Time) {
	gk.verified.Range(func(key, value any) bool {
		expiresAt, ok := value.(time.Time)
		if !ok || !now.Before(expiresAt) {
			gk.verified.Delete(key)
		}
		return true
	})

	gk.pending.Range(func(key, value any) bool {
		pu, ok := value.(*pendingUser)
		if !ok || pu == nil {
			gk.pending.Delete(key)
			return true
		}
		if now.Before(pu.expiresAt) {
			return true
		}

		k, ok := key.(string)
		if !ok {
			return true
		}
		gk.failChallenge(gk.ctx, k, pu)
		return true
	})
}

func (gk *Gatekeeper) Stats() map[string]int64 {
	return map[string]int64{
		"challenges_sent":   gk.stats.challengesSent.Load(),
		"challenges_passed": gk.stats.challengesPassed.Load(),
		"challenges_failed": gk.stats.challengesFailed.Load(),
		"users_kicked":      gk.stats.usersKicked.Load(),
	}
}

func (gk *Gatekeeper) Close() {
	if gk == nil {
		return
	}
	gk.closeOnce.Do(func() {
		gk.cancel()
		gk.workers.Wait()
	})
}

func makeKey(chatID, userID string) string {
	return chatID + "\x00" + userID
}

func randomInt(min, max int) int {
	if max <= min {
		return min
	}
	return rand.Intn(max-min+1) + min
}

func buildMathOptions(correct, count int) []int {
	if count < 2 {
		count = 4
	}
	if count > 8 {
		count = 8
	}

	options := make([]int, 0, count)
	seen := make(map[int]struct{}, count)

	options = append(options, correct)
	seen[correct] = struct{}{}

	for len(options) < count {
		offset := rand.Intn(11) - 5
		if offset == 0 {
			continue
		}
		value := correct + offset
		if value < 0 {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		options = append(options, value)
	}

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return options
}

func buildNumberOptions(correct, count, min, max int) []int {
	if count < 2 {
		count = 4
	}
	if count > 8 {
		count = 8
	}
	if max-min+1 < count {
		count = max - min + 1
	}

	options := make([]int, 0, count)
	seen := make(map[int]struct{}, count)

	options = append(options, correct)
	seen[correct] = struct{}{}

	for len(options) < count {
		value := randomInt(min, max)
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		options = append(options, value)
	}

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return options
}

func normalizeCaptchaConfig(base, override CaptchaConfig) CaptchaConfig {
	if override.Type != "" {
		base.Type = override.Type
	}
	if override.Options > 0 {
		base.Options = override.Options
	}
	if override.MinNumber != 0 {
		base.MinNumber = override.MinNumber
	}
	if override.MaxNumber != 0 {
		base.MaxNumber = override.MaxNumber
	}
	if override.Title != "" {
		base.Title = override.Title
	}
	if override.CorrectText != "" {
		base.CorrectText = override.CorrectText
	}
	if override.WrongText != "" {
		base.WrongText = override.WrongText
	}

	if base.Options < 2 {
		base.Options = 4
	}
	if base.MinNumber < 0 {
		base.MinNumber = 0
	}
	if base.MaxNumber <= base.MinNumber {
		base.MaxNumber = base.MinNumber + 9
	}

	return base
}
