package antispam

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

// Config تنظیمات سیستم ضد اسپم
type Config struct {
	MaxMessages        int           // حداکثر پیام در پنجره زمانی
	Window             time.Duration // پنجره زمانی (مثلا 5 ثانیه)
	Cooldown           time.Duration // زمان خنک‌شدن بعد از اخطار
	MaxWarnings        int           // حداکثر اخطار قبل از بن
	BanDuration        time.Duration // مدت بن
	DetectFlood        bool          // تشخیص فلود (پیام‌های تکراری پشت سر هم)
	FloodThreshold     int           // آستانه تشخیص فلود
	DetectSpamKeywords bool          // تشخیص کلمات اسپم
	SpamKeywords       []string      // لیست کلمات ممنوعه
	BlockBots          bool          // مسدودسازی بات‌ها
	MaxCommandLength   int           // حداکثر طول پیام
	WarnAction         func(ctx context.Context, userID string, reason string)
	BanAction          func(ctx context.Context, userID string, reason string)
}

// DefaultConfig تنظیمات پیش‌فرض
func DefaultConfig() Config {
	return Config{
		MaxMessages:        8,
		Window:             5 * time.Second,
		Cooldown:           10 * time.Second,
		MaxWarnings:        3,
		BanDuration:        30 * time.Minute,
		DetectFlood:        true,
		FloodThreshold:     5,
		DetectSpamKeywords: true,
		SpamKeywords: []string{
			"http://", "https://", "www.", "t.me/", "telegram.me",
			"joinchat", "promo", "ad.", "تبلیغ", "آگهی",
		},
		BlockBots:        true,
		MaxCommandLength: 512,
	}
}

type userState struct {
	mu             sync.Mutex
	messages       []time.Time
	warnings       int
	lastMessage    string
	lastMessageCnt int
	cooldownUntil  time.Time
	bannedUntil    time.Time
	lastActivity   time.Time
	totalMessages  int64
}

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

// AntiSpam ساختار اصلی پکیج ضد اسپم
type AntiSpam struct {
	cfg           Config
	shardedStates [64]map[string]*userState
	shardedMu     [64]sync.Mutex
	blocked       sync.Map // بن‌های دائمی
	stats         stats
}

// New ساخت یک نمونه جدید از سیستم ضد اسپم
func New(cfg Config) *AntiSpam {
	if cfg.MaxMessages <= 0 {
		cfg.MaxMessages = 8
	}
	if cfg.Window <= 0 {
		cfg.Window = 5 * time.Second
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

	as := &AntiSpam{cfg: cfg}
	for i := range as.shardedStates {
		as.shardedStates[i] = make(map[string]*userState)
	}

	// اجرای goroutine برای پاکسازی دوره‌ای حافظه
	go as.cleanupLoop()

	return as
}

func (as *AntiSpam) cleanupLoop() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		for i := range as.shardedStates {
			as.shardedMu[i].Lock()
			for uid, s := range as.shardedStates[i] {
				s.mu.Lock()
				expired := now.Sub(s.lastActivity) > 30*time.Minute
				s.mu.Unlock()
				if expired {
					delete(as.shardedStates[i], uid)
				}
			}
			as.shardedMu[i].Unlock()
		}
	}
}

func (as *AntiSpam) hash(s string) int {
	h := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		h *= 16777619
		h ^= uint32(s[i])
	}
	return int(h) % 64
}

func (as *AntiSpam) getState(userID string) *userState {
	idx := as.hash(userID)
	as.shardedMu[idx].Lock()
	defer as.shardedMu[idx].Unlock()
	s, ok := as.shardedStates[idx][userID]
	if !ok {
		s = &userState{}
		as.shardedStates[idx][userID] = s
	}
	return s
}

// BanUser بن دستی کاربر (دائمی)
func (as *AntiSpam) BanUser(userID string) {
	as.blocked.Store(userID, true)
	as.stats.usersBanned.Add(1)
}

// UnbanUser رفع بن کاربر
func (as *AntiSpam) UnbanUser(userID string) {
	as.blocked.Delete(userID)
}

// Stats گرفتن آمار سیستم ضد اسپم
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

// Middleware این متد میدل‌ور را برمی‌گرداند تا به ربات متصل شود
func (as *AntiSpam) Middleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			var userID string
			var isBot bool
			var text string

			if u.Message != nil {
				if u.Message.From != nil {
					userID = u.Message.From.ID
					isBot = u.Message.From.IsBot
				}
				text = u.Message.Text
			} else if u.CallbackQuery != nil {
				userID = u.CallbackQuery.From.ID
				isBot = u.CallbackQuery.From.IsBot
				text = u.CallbackQuery.Data
			}

			if userID == "" {
				next(ctx, u)
				return
			}

			// بررسی بن دائمی
			if _, blocked := as.blocked.Load(userID); blocked {
				as.stats.blockedHits.Add(1)
				return
			}

			// مسدودسازی بات‌ها
			if as.cfg.BlockBots && isBot {
				as.stats.botBlocked.Add(1)
				return
			}

			state := as.getState(userID)
			state.mu.Lock()
			defer state.mu.Unlock()

			now := time.Now()
			state.lastActivity = now
			state.totalMessages++

			// بررسی بن موقت
			if now.Before(state.bannedUntil) {
				as.stats.bannedHits.Add(1)
				if as.cfg.BanAction != nil {
					remaining := state.bannedUntil.Sub(now).Round(time.Second)
					as.cfg.BanAction(ctx, userID, fmt.Sprintf("banned (%v remaining)", remaining))
				}
				return
			}

			// بررسی Cooldown
			if now.Before(state.cooldownUntil) {
				as.stats.cooldownHits.Add(1)
				return
			}

			// تشخیص فلود
			if as.cfg.DetectFlood && text != "" {
				if text == state.lastMessage {
					state.lastMessageCnt++
					if state.lastMessageCnt >= as.cfg.FloodThreshold {
						state.warnings++
						state.cooldownUntil = now.Add(as.cfg.Cooldown)
						state.lastMessageCnt = 0
						as.stats.floodDetected.Add(1)
						if as.cfg.WarnAction != nil {
							as.cfg.WarnAction(ctx, userID, "flood (repeated message)")
						}
						if state.warnings >= as.cfg.MaxWarnings {
							state.bannedUntil = now.Add(as.cfg.BanDuration)
							as.stats.usersBanned.Add(1)
							if as.cfg.BanAction != nil {
								as.cfg.BanAction(ctx, userID, "auto-ban: flood")
							}
						}
						return
					}
				} else {
					state.lastMessage = text
					state.lastMessageCnt = 1
				}
			}

			// تشخیص کلمات اسپم
			if as.cfg.DetectSpamKeywords && text != "" {
				lowerText := strings.ToLower(text)
				for _, kw := range as.cfg.SpamKeywords {
					if strings.Contains(lowerText, strings.ToLower(kw)) {
						state.warnings++
						as.stats.spamKeywordHits.Add(1)
						if as.cfg.WarnAction != nil {
							as.cfg.WarnAction(ctx, userID, "spam keyword: "+kw)
						}
						if state.warnings >= as.cfg.MaxWarnings {
							state.bannedUntil = now.Add(as.cfg.BanDuration)
							as.stats.usersBanned.Add(1)
							if as.cfg.BanAction != nil {
								as.cfg.BanAction(ctx, userID, "auto-ban: spam keyword")
							}
						}
						return
					}
				}
			}

			// بررسی طول پیام
			if as.cfg.MaxCommandLength > 0 && len(text) > as.cfg.MaxCommandLength {
				state.warnings++
				if as.cfg.WarnAction != nil {
					as.cfg.WarnAction(ctx, userID, "message too long")
				}
				if state.warnings >= as.cfg.MaxWarnings {
					state.bannedUntil = now.Add(as.cfg.BanDuration)
					as.stats.usersBanned.Add(1)
				}
				return
			}

			// Sliding Window Rate Limiting
			cutoff := now.Add(-as.cfg.Window)
			valid := state.messages[:0]
			for _, t := range state.messages {
				if t.After(cutoff) {
					valid = append(valid, t)
				}
			}
			state.messages = valid
			state.messages = append(state.messages, now)

			if len(state.messages) > as.cfg.MaxMessages {
				state.warnings++
				state.cooldownUntil = now.Add(as.cfg.Cooldown)
				as.stats.rateLimitHits.Add(1)
				if as.cfg.WarnAction != nil {
					as.cfg.WarnAction(ctx, userID, "rate limit exceeded")
				}
				if state.warnings >= as.cfg.MaxWarnings {
					state.bannedUntil = now.Add(as.cfg.BanDuration)
					as.stats.usersBanned.Add(1)
					if as.cfg.BanAction != nil {
						as.cfg.BanAction(ctx, userID, "auto-ban: rate limit")
					}
				}
				return
			}

			// کاربر مجاز است
			as.stats.allowed.Add(1)
			next(ctx, u)
		}
	}
}
