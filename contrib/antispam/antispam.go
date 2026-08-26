package antispam

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

const (
	// تعداد شاردهای داخلی برای کاهش رقابت روی قفل‌ها.
	shardCount = 64

	// مدت زمانی که State بدون فعالیت در حافظه نگه داشته می‌شود.
	stateTTL = 30 * time.Minute

	// فاصله اجرای پاک‌سازی Stateهای قدیمی.
	cleanupInterval = 2 * time.Minute
)

// Config تنظیمات سیستم ضد اسپم را مشخص می‌کند.
type Config struct {
	// حداکثر تعداد پیام مجاز در پنجره زمانی.
	MaxMessages int

	// مدت پنجره زمانی Rate Limit.
	Window time.Duration

	// مدت زمان جلوگیری از ارسال پیام پس از شناسایی اسپم.
	Cooldown time.Duration

	// حداکثر تعداد اخطار قبل از اعمال بن موقت.
	MaxWarnings int

	// مدت زمان بن موقت.
	BanDuration time.Duration

	// فعال کردن تشخیص ارسال چندباره یک پیام.
	DetectFlood bool

	// تعداد دفعات تکرار پیام برای شناسایی Flood.
	FloodThreshold int

	// فعال کردن بررسی کلمات و عبارات اسپم.
	DetectSpamKeywords bool

	// فهرست کلمات و عبارات اسپم.
	SpamKeywords []string

	// در صورت فعال بودن، حساب‌های Bot مسدود می‌شوند.
	BlockBots bool

	// حداکثر طول مجاز متن پیام بر اساس تعداد کاراکتر Unicode.
	// اگر مقدار صفر یا منفی باشد، محدودیت طول غیرفعال است.
	MaxCommandLength int

	// تابعی که هنگام دادن اخطار اجرا می‌شود.
	WarnAction func(ctx context.Context, userID string, reason string)

	// تابعی که هنگام بن شدن کاربر اجرا می‌شود.
	BanAction func(ctx context.Context, userID string, reason string)
}

// DefaultConfig تنظیمات پیش‌فرض سیستم ضد اسپم را برمی‌گرداند.
func DefaultConfig() Config {
	return Config{
		MaxMessages: 8,
		Window:      5 * time.Second,
		Cooldown:    10 * time.Second,
		MaxWarnings: 3,
		BanDuration: 30 * time.Minute,

		DetectFlood:    true,
		FloodThreshold: 5,

		DetectSpamKeywords: true,
		SpamKeywords: []string{
			"http://",
			"https://",
			"www.",
			"t.me/",
			"telegram.me",
			"joinchat",
			"promo",
			"ad.",
			"تبلیغ",
			"آگهی",
		},

		BlockBots:        true,
		MaxCommandLength: 512,
	}
}

// userState وضعیت ضد اسپم یک کاربر را نگه می‌دارد.
type userState struct {
	mu sync.Mutex

	// Ring Buffer زمان پیام‌ها را نگه می‌دارد.
	// این روش نسبت به ایجاد Slice جدید در هر پیام، فشار بسیار کمتری روی GC ایجاد می‌کند.
	messages []time.Time

	// موقعیت بعدی برای ثبت زمان پیام در Ring Buffer.
	messageIndex int

	// تعداد پیام‌های ثبت‌شده در Ring Buffer.
	messageCount int

	// تعداد اخطارهای کاربر.
	warnings int

	// آخرین پیام کاربر برای تشخیص Flood.
	lastMessage string

	// تعداد تکرار متوالی آخرین پیام.
	lastMessageCount int

	// زمان پایان Cooldown.
	cooldownUntil time.Time

	// زمان پایان بن موقت.
	bannedUntil time.Time

	// آخرین زمان فعالیت کاربر.
	lastActivity time.Time

	// تعداد کل پیام‌های پردازش‌شده کاربر.
	totalMessages int64
}

// stats آمار سیستم ضد اسپم را به‌صورت اتمیک نگه می‌دارد.
type stats struct {
	allowed         atomic.Int64
	blockedHits     atomic.Int64
	botBlocked      atomic.Int64
	bannedHits      atomic.Int64
	cooldownHits    atomic.Int64
	floodDetected   atomic.Int64
	spamKeywordHits atomic.Int64
	rateLimitHits   atomic.Int64
	usersBanned     atomic.Int64
}

// AntiSpam ساختار اصلی سیستم ضد اسپم است.
type AntiSpam struct {
	cfg Config

	// State کاربران به‌صورت Sharded نگه‌داری می‌شود
	// تا کاربران مختلف مجبور نباشند یک Mutex مشترک داشته باشند.
	shardedStates [shardCount]map[string]*userState

	// قفل هر Shard به‌صورت مستقل است.
	shardedMu [shardCount]sync.Mutex

	// کاربران دارای بن دائمی.
	blocked sync.Map

	// کلمات اسپم از قبل Normalize شده‌اند.
	// استفاده از Map باعث می‌شود در بررسی‌های سریع، به‌جای پیمایش Slice، جست‌وجوی مستقیم داشته باشیم.
	spamKeywords map[string]struct{}

	stats stats
}

// New یک نمونه جدید از سیستم ضد اسپم ایجاد می‌کند.
func New(cfg Config) *AntiSpam {
	// تنظیم مقادیر نامعتبر با مقادیر پیش‌فرض.
	if cfg.MaxMessages <= 0 {
		cfg.MaxMessages = 8
	}
	if cfg.Window <= 0 {
		cfg.Window = 5 * time.Second
	}
	if cfg.Cooldown <= 0 {
		cfg.Cooldown = 10 * time.Second
	}
	if cfg.MaxWarnings <= 0 {
		cfg.MaxWarnings = 3
	}
	if cfg.BanDuration <= 0 {
		cfg.BanDuration = 30 * time.Minute
	}
	if cfg.FloodThreshold <= 0 {
		cfg.FloodThreshold = 5
	}

	// Map کلمات اسپم را با ظرفیت مناسب ایجاد می‌کنیم.
	spamKeywords := make(map[string]struct{}, len(cfg.SpamKeywords))
	for _, keyword := range cfg.SpamKeywords {
		keyword = strings.TrimSpace(keyword)
		if keyword == "" {
			continue
		}
		// فقط یک بار هنگام ساخت AntiSpam انجام می‌شود.
		keyword = strings.ToLower(keyword)
		spamKeywords[keyword] = struct{}{}
	}

	as := &AntiSpam{
		cfg:          cfg,
		spamKeywords: spamKeywords,
	}

	// ساخت Mapهای Shardها.
	for i := range as.shardedStates {
		as.shardedStates[i] = make(map[string]*userState)
	}

	// پاک‌سازی دوره‌ای Stateهای بدون استفاده.
	go as.cleanupLoop()

	return as
}

// cleanupLoop Stateهای قدیمی را از حافظه حذف می‌کند.
func (as *AntiSpam) cleanupLoop() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for now := range ticker.C {
		as.cleanup(now)
	}
}

// cleanup کاربران بدون فعالیت را حذف می‌کند.
func (as *AntiSpam) cleanup(now time.Time) {
	for i := range as.shardedStates {
		as.shardedMu[i].Lock()

		for userID, state := range as.shardedStates[i] {
			state.mu.Lock()
			expired := now.Sub(state.lastActivity) > stateTTL
			state.mu.Unlock()

			if expired {
				delete(as.shardedStates[i], userID)
			}
		}

		as.shardedMu[i].Unlock()
	}
}

// hash شناسه کاربر را به یکی از Shardها تبدیل می‌کند.
// الگوریتم FNV-1a برای این کاربرد سریع و کم‌هزینه است.
func (as *AntiSpam) hash(value string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(value); i++ {
		hash ^= uint32(value[i])
		hash *= 16777619
	}
	return hash
}

// getState State مربوط به یک کاربر را دریافت یا ایجاد می‌کند.
func (as *AntiSpam) getState(userID string) *userState {
	index := as.hash(userID) & (shardCount - 1)

	as.shardedMu[index].Lock()

	state := as.shardedStates[index][userID]

	if state == nil {
		state = &userState{
			// Ring Buffer را فقط به اندازه Rate Limit ایجاد می‌کنیم.
			messages: make([]time.Time, as.cfg.MaxMessages),
		}
		as.shardedStates[index][userID] = state
	}

	as.shardedMu[index].Unlock()

	return state
}

// BanUser کاربر را به‌صورت دائمی مسدود می‌کند.
func (as *AntiSpam) BanUser(userID string) {
	if userID == "" {
		return
	}

	// اگر کاربر از قبل بن باشد، شمارنده افزایش پیدا نمی‌کند.
	_, loaded := as.blocked.LoadOrStore(userID, struct{}{})
	if !loaded {
		as.stats.usersBanned.Add(1)
	}
}

// UnbanUser بن دائمی کاربر را حذف می‌کند.
func (as *AntiSpam) UnbanUser(userID string) {
	if userID == "" {
		return
	}
	as.blocked.Delete(userID)
}

// Stats آمار فعلی سیستم ضد اسپم را برمی‌گرداند.
func (as *AntiSpam) Stats() map[string]int64 {
	return map[string]int64{
		"allowed":           as.stats.allowed.Load(),
		"blocked_hits":      as.stats.blockedHits.Load(),
		"bot_blocked":       as.stats.botBlocked.Load(),
		"banned_hits":       as.stats.bannedHits.Load(),
		"cooldown_hits":     as.stats.cooldownHits.Load(),
		"flood_detected":    as.stats.floodDetected.Load(),
		"spam_keyword_hits": as.stats.spamKeywordHits.Load(),
		"rate_limit_hits":   as.stats.rateLimitHits.Load(),
		"users_banned":      as.stats.usersBanned.Load(),
	}
}

// Middleware میدل‌ور ضد اسپم را ایجاد می‌کند.
func (as *AntiSpam) Middleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			if u == nil {
				next(ctx, u)
				return
			}

			userID, isBot, text := extractUserData(u)

			// پیام بدون کاربر مشخص را محدود نمی‌کنیم.
			if userID == "" {
				next(ctx, u)
				return
			}

			// بررسی بن دائمی.
			if _, blocked := as.blocked.Load(userID); blocked {
				as.stats.blockedHits.Add(1)
				return
			}

			// مسدودسازی Botها.
			if as.cfg.BlockBots && isBot {
				as.stats.botBlocked.Add(1)
				return
			}

			state := as.getState(userID)
			now := time.Now()

			state.mu.Lock()
			state.lastActivity = now
			state.totalMessages++

			// بررسی بن موقت.
			if now.Before(state.bannedUntil) {
				remaining := state.bannedUntil.Sub(now)
				state.mu.Unlock()

				as.stats.bannedHits.Add(1)

				if as.cfg.BanAction != nil {
					as.cfg.BanAction(ctx, userID, fmt.Sprintf("بن موقت فعال است؛ %v باقی مانده", remaining.Round(time.Second)))
				}
				return
			}

			// بررسی Cooldown.
			if now.Before(state.cooldownUntil) {
				state.mu.Unlock()
				as.stats.cooldownHits.Add(1)
				return
			}

			// بررسی طول پیام.
			if as.cfg.MaxCommandLength > 0 && utf8.RuneCountInString(text) > as.cfg.MaxCommandLength {
				banned := as.addWarningLocked(state, now)
				state.mu.Unlock()

				if as.cfg.WarnAction != nil {
					as.cfg.WarnAction(ctx, userID, "طول پیام بیش از حد مجاز است")
				}

				if banned {
					as.executeBan(ctx, userID, "بن خودکار به دلیل طول بیش از حد پیام")
				}
				return
			}

			// تشخیص Flood پیام‌های تکراری.
			if as.cfg.DetectFlood && text != "" {
				if text == state.lastMessage {
					state.lastMessageCount++
				} else {
					state.lastMessage = text
					state.lastMessageCount = 1
				}

				if state.lastMessageCount >= as.cfg.FloodThreshold {
					state.lastMessageCount = 0
					state.cooldownUntil = now.Add(as.cfg.Cooldown)

					banned := as.addWarningLocked(state, now)
					state.mu.Unlock()

					as.stats.floodDetected.Add(1)

					if as.cfg.WarnAction != nil {
						as.cfg.WarnAction(ctx, userID, "ارسال مکرر یک پیام شناسایی شد")
					}

					if banned {
						as.executeBan(ctx, userID, "بن خودکار به دلیل Flood")
					}
					return
				}
			}

			// بررسی کلمات و عبارات اسپم.
			if as.cfg.DetectSpamKeywords && text != "" {
				if keyword := as.findSpamKeyword(text); keyword != "" {
					banned := as.addWarningLocked(state, now)
					state.mu.Unlock()

					as.stats.spamKeywordHits.Add(1)

					if as.cfg.WarnAction != nil {
						as.cfg.WarnAction(ctx, userID, "کلمه یا عبارت اسپم شناسایی شد: "+keyword)
					}

					if banned {
						as.executeBan(ctx, userID, "بن خودکار به دلیل استفاده از عبارت اسپم")
					}
					return
				}
			}

			// بررسی Rate Limit با Ring Buffer.
			if as.rateLimitExceeded(state, now) {
				state.cooldownUntil = now.Add(as.cfg.Cooldown)

				banned := as.addWarningLocked(state, now)
				state.mu.Unlock()

				as.stats.rateLimitHits.Add(1)

				if as.cfg.WarnAction != nil {
					as.cfg.WarnAction(ctx, userID, "تعداد پیام‌ها در مدت کوتاه بیش از حد مجاز است")
				}

				if banned {
					as.executeBan(ctx, userID, "بن خودکار به دلیل عبور از محدودیت ارسال پیام")
				}
				return
			}

			state.mu.Unlock()

			// پیام مجاز است.
			as.stats.allowed.Add(1)
			next(ctx, u)
		}
	}
}

// extractUserData اطلاعات کاربر را از Update استخراج می‌کند.
// ارور های مربوط به پرانتز در Return در اینجا رفع شدند.
func extractUserData(u *models.Update) (string, bool, string) {
	if u.Message != nil {
		if u.Message.From == nil {
			return "", false, u.Message.Text
		}
		return u.Message.From.ID, u.Message.From.IsBot, u.Message.Text
	}

	if u.CallbackQuery != nil {
		if u.CallbackQuery.From == nil {
			return "", false, u.CallbackQuery.Data
		}
		return u.CallbackQuery.From.ID, u.CallbackQuery.From.IsBot, u.CallbackQuery.Data
	}

	return "", false, ""
}

// findSpamKeyword اولین عبارت اسپم پیدا شده را برمی‌گرداند.
// متن فقط یک بار lowercase می‌شود. کلمات اسپم قبلاً در New نرمال شده‌اند.
func (as *AntiSpam) findSpamKeyword(text string) string {
	lowerText := strings.ToLower(text)

	for keyword := range as.spamKeywords {
		if strings.Contains(lowerText, keyword) {
			return keyword
		}
	}

	return ""
}

// rateLimitExceeded محدودیت Rate Limit را بررسی می‌کند.
// از Ring Buffer استفاده می‌شود تا به‌جای ساختن Slice جدید، همان حافظه قبلی مجدداً استفاده شود.
func (as *AntiSpam) rateLimitExceeded(state *userState, now time.Time) bool {
	cutoff := now.Add(-as.cfg.Window)

	// تعداد پیام‌های معتبر داخل پنجره زمانی.
	valid := 0

	for i := 0; i < state.messageCount; i++ {
		index := (state.messageIndex - state.messageCount + i + as.cfg.MaxMessages) % as.cfg.MaxMessages
		if state.messages[index].After(cutoff) {
			valid++
		}
	}

	// اگر تعداد پیام‌های قبلی به سقف رسیده باشد، پیام فعلی باعث عبور از Rate Limit می‌شود.
	if valid >= as.cfg.MaxMessages {
		return true
	}

	// ثبت پیام جدید در Ring Buffer.
	state.messages[state.messageIndex] = now
	state.messageIndex++

	if state.messageIndex >= as.cfg.MaxMessages {
		state.messageIndex = 0
	}

	if state.messageCount < as.cfg.MaxMessages {
		state.messageCount++
	}

	return false
}

// addWarningLocked یک اخطار به کاربر اضافه می‌کند.
// این تابع باید در حالی که state.mu قفل است فراخوانی شود.
func (as *AntiSpam) addWarningLocked(state *userState, now time.Time) bool {
	state.warnings++

	if state.warnings >= as.cfg.MaxWarnings {
		// فقط وضعیت داخلی را تغییر می‌دهیم. اجرای Callback خارج از Lock انجام می‌شود.
		if !now.Before(state.bannedUntil) {
			state.bannedUntil = now.Add(as.cfg.BanDuration)
			return true
		}
	}

	return false
}

// executeBan عملیات بن را خارج از Lock کاربر اجرا می‌کند.
// این موضوع مهم است چون BanAction ممکن است عملیات شبکه‌ای یا نسبتاً سنگینی انجام دهد.
func (as *AntiSpam) executeBan(ctx context.Context, userID string, reason string) {
	as.stats.usersBanned.Add(1)

	if as.cfg.BanAction != nil {
		as.cfg.BanAction(ctx, userID, reason)
	}
}
