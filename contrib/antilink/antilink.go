package antilink

import (
	"context"
	"net/url"
	"regexp"
	"strings"

	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

type Config struct {
	AllowedDomains []string
	BlockUsernames bool
	BlockInvites   bool
	Action         func(
		ctx context.Context,
		userID string,
		chatID string,
		messageID int,
		reason string,
	)
}

func DefaultConfig() Config {
	return Config{
		AllowedDomains: []string{"codemeet.chat"},
		BlockUsernames: false,
		BlockInvites:   true,
		Action:         nil,
	}
}

type AntiLink struct {
	cfg            Config
	allowedDomains map[string]struct{}
	urlRegex       *regexp.Regexp
	domainRegex    *regexp.Regexp
	telegramRegex  *regexp.Regexp
	inviteRegex    *regexp.Regexp
	usernameRegex  *regexp.Regexp
}

func New(cfg Config) *AntiLink {
	if cfg.AllowedDomains == nil {
		cfg.AllowedDomains = []string{}
	}

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
		urlRegex:       regexp.MustCompile(`(?i)(?:https?://|www\.)[^\s<>"']+`),
		domainRegex:    regexp.MustCompile(`(?i)(?:^|[\s(])((?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,})(?:/[^\s<>"']*)?`),
		telegramRegex:  regexp.MustCompile(`(?i)(?:https?://)?(?:www\.)?(?:t\.me|telegram\.me)/[^\s<>"']+`),
		inviteRegex:    regexp.MustCompile(`(?i)(?:https?://)?(?:www\.)?(?:t\.me|telegram\.me)/(?:joinchat/|\+)`),
		usernameRegex:  regexp.MustCompile(`(?i)(?:^|[^a-z0-9._%+-])@([a-z0-9_]{5,32})(?:$|[^a-z0-9_])`),
	}
}

// IsBlocked بررسی می‌کند که آیا متن یا انتیتی‌ها حاوی لینک مسدود شده هستند یا خیر
// این متد برای استفاده در سایر پکیج‌ها (مثل GroupManager) اضافه شده است.
func (al *AntiLink) IsBlocked(text string, entities []models.MessageEntity) (bool, string) {
	if text == "" {
		return false, ""
	}

	// 1. بررسی هایپرلینک‌های مخفی (Text Links)
	for _, entity := range entities {
		if entity.Type == "text_link" && entity.URL != "" {
			if !al.isAllowedURL(entity.URL) {
				return true, "hyperlink detected (hidden text link)"
			}
		}
	}

	// 2. بررسی لینک‌های دعوت
	if al.cfg.BlockInvites && al.inviteRegex.MatchString(text) {
		return true, "invite link detected"
	}

	// 3. بررسی لینک‌های تلگرام و کدمیت
	if al.telegramRegex.MatchString(text) {
		if !al.hasAllowedTelegramDomain(text) {
			return true, "telegram/codemeet link detected"
		}
	}

	// ۴. بررسی لینک‌های دارای پروتکل یا www
	if match := al.urlRegex.FindString(text); match != "" {
		if !al.isAllowedURL(match) {
			return true, "url link detected"
		}
	}

	// 5. بررسی دامنه‌های بدون پروتکل (مثل google.com)
	if match := al.domainRegex.FindStringSubmatch(text); len(match) > 1 {
		host := match[1]
		if !al.isAllowedHost(host) {
			return true, "domain link detected"
		}
	}

	// 6. بررسی منشن‌های کاربری
	if al.cfg.BlockUsernames && al.usernameRegex.MatchString(text) {
		return true, "username mention detected"
	}

	return false, ""
}

func (al *AntiLink) Middleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			if u == nil || u.Message == nil {
				next(ctx, u)
				return
			}

			text := u.Message.Text
			entities := u.Message.Entities

			if text == "" {
				text = u.Message.Caption
				entities = u.Message.CaptionEntities
			}

			if text == "" {
				next(ctx, u)
				return
			}

			// استفاده از متد IsBlocked برای بررسی یکپارچه
			blocked, reason := al.IsBlocked(text, entities)
			if blocked {
				al.block(ctx, u, reason)
				return
			}

			next(ctx, u)
		}
	}
}

func (al *AntiLink) block(ctx context.Context, u *models.Update, reason string) {
	if al.cfg.Action == nil {
		return
	}
	userID := ""
	if u.Message.From != nil {
		userID = u.Message.From.ID
	}
	al.cfg.Action(ctx, userID, u.Message.Chat.ID, u.Message.MessageID, reason)
}

func (al *AntiLink) isAllowedURL(rawURL string) bool {
	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	host := parsed.Hostname()
	if host == "" {
		return false
	}
	return al.isAllowedHost(host)
}

func (al *AntiLink) isAllowedHost(host string) bool {
	host = normalizeDomain(host)
	if host == "" {
		return false
	}
	if _, ok := al.allowedDomains[host]; ok {
		return true
	}
	for domain := range al.allowedDomains {
		if strings.HasSuffix(host, "."+domain) {
			return true
		}
	}
	return false
}

func (al *AntiLink) hasAllowedTelegramDomain(text string) bool {
	matches := al.telegramRegex.FindAllString(text, -1)
	for _, match := range matches {
		if strings.Contains(strings.ToLower(match), "codemeet.chat") {
			return false
		}
	}
	return true
}

func normalizeDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	domain = strings.ToLower(domain)
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "www.")
	if index := strings.IndexByte(domain, '/'); index >= 0 {
		domain = domain[:index]
	}
	if index := strings.IndexByte(domain, '?'); index >= 0 {
		domain = domain[:index]
	}
	if index := strings.IndexByte(domain, '#'); index >= 0 {
		domain = domain[:index]
	}
	domain = strings.TrimSuffix(domain, ".")
	return domain
}
