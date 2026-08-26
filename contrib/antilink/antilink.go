package antilink

import (
	"context"
	"regexp"
	"strings"

	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

// Config تنظیمات سیستم ضد لینک
type Config struct {
	// AllowedDomains لیست دامنه‌های مجاز (مثلاً codemeet.chat)
	AllowedDomains []string
	// BlockUsernames مسدودسازی آیدی کاربران (@username)
	BlockUsernames bool
	// BlockInvites مسدودسازی لینک‌های دعوت (joinchat)
	BlockInvites bool
	// Action اکشن هنگام تشخیص لینک (برای حذف پیام، اخطار یا بن)
	Action func(ctx context.Context, userID, chatID string, messageID int, reason string)
}

// DefaultConfig تنظیمات پیش‌فرض
func DefaultConfig() Config {
	return Config{
		AllowedDomains: []string{"codemeet.chat"},
		BlockUsernames: false, // معمولاً در گروه‌ها فعال می‌شود
		BlockInvites:   true,
		Action:         nil,
	}
}

// AntiLink ساختار اصلی پکیج
type AntiLink struct {
	cfg           Config
	urlRegex      *regexp.Regexp
	tgLinkRegex   *regexp.Regexp
	inviteRegex   *regexp.Regexp
	usernameRegex *regexp.Regexp
}

// New ساخت یک نمونه جدید از سیستم ضد لینک
func New(cfg Config) *AntiLink {
	if cfg.AllowedDomains == nil {
		cfg.AllowedDomains = []string{}
	}

	return &AntiLink{
		cfg: cfg,
		// تشخیص لینک‌های HTTP و WWW
		urlRegex: regexp.MustCompile(`(?i)\b((https?://|www\.)[^\s]+)`),
		// تشخیص لینک‌های کدمیت/تلگرام (t.me/...)
		tgLinkRegex: regexp.MustCompile(`(?i)\b(t\.me|telegram\.me)/[^\s]+`),
		// تشخیص لینک‌های دعوت
		inviteRegex: regexp.MustCompile(`(?i)(t\.me/|telegram\.me/)(joinchat|\+)`),
		// تشخیص یوزرنیم‌ها (@username)
		usernameRegex: regexp.MustCompile(`(?i)@([a-zA-Z0-9_]{5,})`),
	}
}

// Middleware این متد میدل‌ور را برمی‌گرداند تا به ربات متصل شود
func (al *AntiLink) Middleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			// فقط پیام‌های متنی را بررسی کن
			if u.Message == nil || u.Message.Text == "" {
				next(ctx, u)
				return
			}

			text := u.Message.Text
			userID := ""
			if u.Message.From != nil {
				userID = u.Message.From.ID
			}
			chatID := u.Message.Chat.ID
			msgID := u.Message.MessageID

			var reason string
			var isBlocked bool

			// ۱. بررسی لینک‌های دعوت (joinchat)
			if al.cfg.BlockInvites && al.inviteRegex.MatchString(text) {
				isBlocked = true
				reason = "invite link detected"
			}

			// ۲. بررسی لینک‌های کدمیت/تلگرام
			if !isBlocked && al.tgLinkRegex.MatchString(text) {
				isBlocked = true
				reason = "telegram/codemeet link detected"
			}

			// ۳. بررسی لینک‌های معمولی (HTTP/WWW)
			if !isBlocked {
				matches := al.urlRegex.FindAllStringSubmatch(text, -1)
				for _, match := range matches {
					if len(match) > 1 {
						urlStr := strings.ToLower(match[1])
						allowed := false
						// بررسی لیست سفید دامنه‌ها
						for _, domain := range al.cfg.AllowedDomains {
							if strings.Contains(urlStr, strings.ToLower(domain)) {
								allowed = true
								break
							}
						}
						if !allowed {
							isBlocked = true
							reason = "url link detected"
							break
						}
					}
				}
			}

			// ۴. بررسی منشن‌ها (@username)
			if !isBlocked && al.cfg.BlockUsernames {
				if al.usernameRegex.MatchString(text) {
					isBlocked = true
					reason = "username mention detected"
				}
			}

			// اگر لینکی شناسایی شد
			if isBlocked {
				if al.cfg.Action != nil {
					al.cfg.Action(ctx, userID, chatID, msgID, reason)
				}
				// جلوگیری از اجرای بقیه هندلرها
				return
			}

			// اگر لینکی نبود، ادامه هندلرها اجرا شود
			next(ctx, u)
		}
	}
}
