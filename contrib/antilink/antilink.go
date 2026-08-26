package antilink

import (
	"context"
	"net/url"
	"regexp"
	"strings"

	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

// Config تنظیمات سیستم ضد لینک را مشخص می‌کند.
type Config struct {
	// دامنه‌هایی که ارسال لینک آن‌ها مجاز است.
	// ساب‌دامین‌های این دامنه‌ها نیز مجاز خواهند بود.
	AllowedDomains []string

	// در صورت فعال بودن، منشن‌های کاربری مانند @username مسدود می‌شوند.
	BlockUsernames bool

	// در صورت فعال بودن، لینک‌های دعوت مانند t.me/+xxxx مسدود می‌شوند.
	BlockInvites bool

	// تابعی که هنگام شناسایی لینک یا منشن مسدودشده اجرا می‌شود.
	Action func(
		ctx context.Context,
		userID string,
		chatID string,
		messageID int,
		reason string,
	)
}

// DefaultConfig تنظیمات پیش‌فرض سیستم ضد لینک را برمی‌گرداند.
func DefaultConfig() Config {
	return Config{
		AllowedDomains: []string{
			"codemeet.chat",
		},
		BlockUsernames: false,
		BlockInvites:   true,
		Action:         nil,
	}
}

// AntiLink ساختار اصلی سیستم ضد لینک است.
type AntiLink struct {
	cfg Config

	// دامنه‌های مجاز برای جست‌وجوی سریع O(1).
	allowedDomains map[string]struct{}

	// تشخیص URLهایی که دارای پروتکل یا www هستند.
	urlRegex *regexp.Regexp

	// تشخیص دامنه‌هایی که بدون پروتکل ارسال شده‌اند.
	domainRegex *regexp.Regexp

	// تشخیص لینک‌های تلگرام و کدمیت.
	telegramRegex *regexp.Regexp

	// تشخیص لینک‌های دعوت.
	inviteRegex *regexp.Regexp

	// تشخیص منشن کاربران.
	usernameRegex *regexp.Regexp
}

// New یک نمونه جدید از سیستم ضد لینک ایجاد می‌کند.
func New(cfg Config) *AntiLink {
	// در صورت مشخص نشدن دامنه‌ها، از لیست خالی استفاده می‌کنیم.
	if cfg.AllowedDomains == nil {
		cfg.AllowedDomains = []string{}
	}

	// Map را با ظرفیت مناسب ایجاد می‌کنیم تا Resize کمتری داشته باشد.
	allowedDomains := make(map[string]struct{}, len(cfg.AllowedDomains))

	for _, domain := range cfg.AllowedDomains {
		domain = normalizeDomain(domain)

		if domain == "" {
			continue
		}

		allowedDomains[domain] = struct{}{}
	}

	return &AntiLink{
		cfg:            cfg,
		allowedDomains: allowedDomains,

		// URLهای دارای http، https یا www را تشخیص می‌دهد.
		urlRegex: regexp.MustCompile(
			`(?i)(?:https?://|www\.)[^\s<>"']+`,
		),

		// دامنه‌هایی که بدون پروتکل نوشته شده‌اند را تشخیص می‌دهد.
		//
		// مثال:
		// google.com
		// example.ir/path
		// codemeet.chat
		domainRegex: regexp.MustCompile(
			`(?i)(?:^|[\s(])((?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,})(?:/[^\s<>"']*)?`,
		),

		// لینک‌های تلگرام و کدمیت را تشخیص می‌دهد.
		telegramRegex: regexp.MustCompile(
			`(?i)(?:https?://)?(?:www\.)?(?:t\.me|telegram\.me)/[^\s<>"']+`,
		),

		// لینک‌های دعوت تلگرام را تشخیص می‌دهد.
		inviteRegex: regexp.MustCompile(
			`(?i)(?:https?://)?(?:www\.)?(?:t\.me|telegram\.me)/(?:joinchat/|\+)`,
		),

		// منشن کاربر را تشخیص می‌دهد.
		//
		// با مرزگذاری مناسب، ایمیل‌هایی مانند:
		// test@gmail.com
		//
		// به‌عنوان منشن شناسایی نمی‌شوند.
		usernameRegex: regexp.MustCompile(
			`(?i)(?:^|[^a-z0-9._%+-])@([a-z0-9_]{5,32})(?:$|[^a-z0-9_])`,
		),
	}
}

// Middleware میدل‌ور سیستم ضد لینک را ایجاد می‌کند.
func (al *AntiLink) Middleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			// اگر Update یا پیام وجود نداشته باشد،
			// بدون پردازش اضافه به هندلر بعدی منتقل می‌شویم.
			if u == nil || u.Message == nil {
				next(ctx, u)
				return
			}

			text := u.Message.Text
			entities := u.Message.Entities

			// اگر متن خالی بود، کپشن را بررسی می‌کنیم (برای عکس‌ها و فایل‌ها)
			if text == "" {
				text = u.Message.Caption
				entities = u.Message.CaptionEntities
			}

			// اگر هم متن و هم کپشن خالی بودند، عبور کن
			if text == "" {
				next(ctx, u)
				return
			}

			// ----------------------------------------------------
			// ۱. بررسی هایپرلینک‌های مخفی (Text Links) از طریق Entities
			// ----------------------------------------------------
			// هایپرلینک‌ها مثل [متن](لینک) در آرایه Entities قرار می‌گیرند
			// و با regex به راحتی قابل استخراج نیستند، بنابراین مستقیماً بررسی می‌شوند.
			for _, entity := range entities {
				if entity.Type == "text_link" && entity.URL != "" {
					if !al.isAllowedURL(entity.URL) {
						al.block(ctx, u, "hyperlink detected (hidden text link)")
						return
					}
				}
			}

			// ----------------------------------------------------
			// ۲. بررسی لینک‌های دعوت
			// ----------------------------------------------------
			// این بررسی قبل از لینک معمولی انجام می‌شود تا
			// لینک‌های دعوت دلیل دقیق‌تری داشته باشند.
			if al.cfg.BlockInvites && al.inviteRegex.MatchString(text) {
				al.block(ctx, u, "invite link detected")
				return
			}

			// ----------------------------------------------------
			// ۳. بررسی لینک‌های تلگرام و کدمیت
			// ----------------------------------------------------
			if al.telegramRegex.MatchString(text) {
				if !al.hasAllowedTelegramDomain(text) {
					al.block(ctx, u, "telegram/codemeet link detected")
					return
				}
			}

			// ----------------------------------------------------
			// ۴. بررسی لینک‌های دارای پروتکل یا www
			// ----------------------------------------------------
			if match := al.urlRegex.FindString(text); match != "" {
				if !al.isAllowedURL(match) {
					al.block(ctx, u, "url link detected")
					return
				}
			}

			// ----------------------------------------------------
			// ۵. بررسی دامنه‌های بدون پروتکل
			// ----------------------------------------------------
			// مثال:
			// google.com
			// example.ir/test
			if match := al.domainRegex.FindStringSubmatch(text); len(match) > 1 {
				host := match[1]

				if !al.isAllowedHost(host) {
					al.block(ctx, u, "domain link detected")
					return
				}
			}

			// ----------------------------------------------------
			// ۶. بررسی منشن‌های کاربری
			// ----------------------------------------------------
			// در صورت فعال بودن، منشن‌های کاربری را بررسی می‌کنیم.
			if al.cfg.BlockUsernames && al.usernameRegex.MatchString(text) {
				al.block(ctx, u, "username mention detected")
				return
			}

			// هیچ مورد غیرمجازی پیدا نشد؛
			// پیام به سایر میدل‌ورها و هندلرها منتقل می‌شود.
			next(ctx, u)
		}
	}
}

// block عملیات مربوط به مسدود کردن پیام را انجام می‌دهد.
func (al *AntiLink) block(
	ctx context.Context,
	u *models.Update,
	reason string,
) {
	// اگر اکشن تعریف نشده باشد، فقط جلوی ادامه پیام را می‌گیریم.
	if al.cfg.Action == nil {
		return
	}

	// شناسه کاربر را در صورت وجود دریافت می‌کنیم.
	userID := ""

	if u.Message.From != nil {
		userID = u.Message.From.ID
	}

	// اجرای اکشن سفارشی.
	al.cfg.Action(
		ctx,
		userID,
		u.Message.Chat.ID,
		u.Message.MessageID,
		reason,
	)
}

// isAllowedURL مشخص می‌کند که URL ارسال‌شده مجاز است یا خیر.
func (al *AntiLink) isAllowedURL(rawURL string) bool {
	// URLهایی که پروتکل ندارند، برای Parse کردن نیاز به پروتکل دارند.
	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}

	parsed, err := url.Parse(rawURL)

	if err != nil {
		// URL خراب را غیرمجاز در نظر می‌گیریم.
		return false
	}

	host := parsed.Hostname()

	if host == "" {
		return false
	}

	return al.isAllowedHost(host)
}

// isAllowedHost مشخص می‌کند که یک Host در لیست دامنه‌های مجاز قرار دارد یا خیر.
func (al *AntiLink) isAllowedHost(host string) bool {
	host = normalizeDomain(host)

	if host == "" {
		return false
	}

	// ابتدا تطبیق مستقیم انجام می‌دهیم.
	if _, ok := al.allowedDomains[host]; ok {
		return true
	}

	// سپس ساب‌دامین‌ها را بررسی می‌کنیم.
	//
	// برای مثال:
	// cdn.codemeet.chat
	//
	// در صورت مجاز بودن:
	// codemeet.chat
	//
	// مجاز خواهد بود.
	for domain := range al.allowedDomains {
		if strings.HasSuffix(host, "."+domain) {
			return true
		}
	}

	return false
}

// hasAllowedTelegramDomain بررسی می‌کند که لینک تلگرام شامل
// دامنه‌ای باشد که در لیست دامنه‌های مجاز قرار دارد یا خیر.
func (al *AntiLink) hasAllowedTelegramDomain(text string) bool {
	matches := al.telegramRegex.FindAllString(text, -1)

	for _, match := range matches {
		if strings.Contains(
			strings.ToLower(match),
			"t.me/",
		) {
			// t.me یک دامنه اختصاصی است و در صورت وجود
			// لینک تلگرام، آن را لینک خارجی در نظر می‌گیریم.
			return false
		}

		if strings.Contains(
			strings.ToLower(match),
			"telegram.me/",
		) {
			return false
		}
	}

	return true
}

// normalizeDomain دامنه را برای مقایسه نرمال می‌کند.
func normalizeDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	domain = strings.ToLower(domain)

	// پروتکل در تنظیم دامنه مجاز لازم نیست.
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")

	// www نیز بخشی از نام اصلی دامنه محسوب نمی‌شود.
	domain = strings.TrimPrefix(domain, "www.")

	// مسیر، Query و Fragment حذف می‌شوند.
	if index := strings.IndexByte(domain, '/'); index >= 0 {
		domain = domain[:index]
	}

	if index := strings.IndexByte(domain, '?'); index >= 0 {
		domain = domain[:index]
	}

	if index := strings.IndexByte(domain, '#'); index >= 0 {
		domain = domain[:index]
	}

	// نقطه انتهایی دامنه حذف می‌شود.
	domain = strings.TrimSuffix(domain, ".")

	return domain
}
