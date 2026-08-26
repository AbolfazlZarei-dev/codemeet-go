package profanityfilter

import (
	"context"
	"strings"
	"unicode"

	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

// Config تنظیمات فیلتر کلمات نامناسب را مشخص می‌کند.
type Config struct {
	// فهرست کلمات یا عبارات ممنوعه.
	BannedWords []string

	// در صورت فعال بودن، پیام حاوی کلمه ممنوعه باید حذف شود.
	DeleteMessage bool

	// در صورت فعال بودن، برای کاربر پیام هشدار ارسال می‌شود.
	WarnUser bool

	// متن هشدار برای کاربر.
	WarnText string

	// تابع اختیاری برای اجرای عملیات سفارشی هنگام شناسایی تخلف.
	Action func(ctx context.Context, userID, chatID string, messageID int, reason string)
}

// DefaultConfig تنظیمات پیش‌فرض فیلتر را برمی‌گرداند.
func DefaultConfig() Config {
	return Config{
		BannedWords: []string{
			"idiot",
			"stupid",
			"fuck",
			"dumb",
		},
		DeleteMessage: true,
		WarnUser:      true,
		WarnText:      "🚫 استفاده از کلمات نامناسب در این گروه ممنوع است. لطفاً ادب را رعایت کنید.",
		Action:        nil,
	}
}

// ProfanityFilter فیلتر اصلی کلمات نامناسب است.
type ProfanityFilter struct {
	cfg       Config
	bannedMap map[string]struct{}
	replacer  *strings.Replacer // برای جایگزینی Leetspeak
}

// New یک فیلتر جدید با تنظیمات مشخص‌شده ایجاد می‌کند.
func New(cfg Config) *ProfanityFilter {
	if cfg.BannedWords == nil {
		cfg.BannedWords = []string{}
	}

	if cfg.WarnText == "" {
		cfg.WarnText = DefaultConfig().WarnText
	}

	bannedMap := make(map[string]struct{}, len(cfg.BannedWords))

	for _, word := range cfg.BannedWords {
		word = normalizeBannedWord(word)
		if word != "" {
			bannedMap[word] = struct{}{}
		}
	}

	// تعریف جایگزین‌های Leetspeak برای تشخیص کلمات مخفی شده
	// مثلا: 1d10t -> idiot, f@ck -> fuck
	replacer := strings.NewReplacer(
		"@", "a",
		"4", "a",
		"8", "b",
		"(", "c",
		"{", "c",
		"[", "c",
		"3", "e",
		"€", "e",
		"6", "g",
		"9", "g",
		"#", "h",
		"1", "i",
		"!", "i",
		"|", "i",
		"0", "o",
		"$", "s",
		"5", "s",
		"7", "t",
		"2", "z",
	)

	return &ProfanityFilter{
		cfg:       cfg,
		bannedMap: bannedMap,
		replacer:  replacer,
	}
}

// Middleware میدل‌ور فیلتر کلمات نامناسب را ایجاد می‌کند.
func (pf *ProfanityFilter) Middleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			if u == nil || u.Message == nil {
				next(ctx, u)
				return
			}

			// بررسی هم متن و هم کپشن (برای عکس‌ها و فایل‌ها)
			text := u.Message.Text
			if text == "" {
				text = u.Message.Caption
			}

			if text == "" {
				next(ctx, u)
				return
			}

			if len(pf.bannedMap) == 0 {
				next(ctx, u)
				return
			}

			// بررسی متن با سیستم تشخیص پیشرفته
			word := pf.scanText(text)

			if word == "" {
				next(ctx, u)
				return
			}

			userID := ""
			if u.Message.From != nil {
				userID = u.Message.From.ID
			}

			if u.Message.Chat == nil {
				return
			}

			chatID := u.Message.Chat.ID
			messageID := u.Message.MessageID

			if pf.cfg.Action != nil {
				pf.cfg.Action(
					ctx,
					userID,
					chatID,
					messageID,
					"banned word: "+word,
				)
			}
		}
	}
}

// scanText متن را نرمال‌سازی کرده و کلمات ممنوعه را به صورت هوشمند تشخیص می‌دهد.
// این روش حتی اگر کاربر از کاراکترهای خاص، فاصله یا Leetspeak استفاده کند، کلمه را شناسایی می‌کند.
func (pf *ProfanityFilter) scanText(text string) string {
	// ۱. نرمال‌سازی متن (Leetspeak و حروف کوچک)
	cleaned := pf.normalizeText(text)

	// ۲. بررسی کلمات به صورت مجزا (Word Boundary)
	// این کار از False Positive جلوگیری می‌کند (مثلا کلمه class برای ass شناسایی نمی‌شود)
	words := strings.Fields(cleaned)
	for _, w := range words {
		if _, banned := pf.bannedMap[w]; banned {
			return w
		}
	}

	// ۳. بررسی متن فشرده شده (حذف فاصله‌ها)
	// برای تشخیص کلماتی که با فاصله یا کاراکترهای خاص جدا شده‌اند (مثل f u c k یا f*ck یا f.u.c.k)
	// فقط برای کلمات طولانی‌تر از ۳ کاراکتر استفاده می‌شود تا False Positive کم شود.
	squashed := strings.ReplaceAll(cleaned, " ", "")
	for bannedWord := range pf.bannedMap {
		if len(bannedWord) > 3 {
			if strings.Contains(squashed, bannedWord) {
				return bannedWord
			}
		}
	}

	return ""
}

// normalizeText متن را برای بررسی نرمال می‌کند.
// - حروف به کوچک تبدیل می‌شوند
// - Leetspeak جایگزین می‌شود (1->i, 3->e, @->a, ...)
// - کاراکترهای غیر حرفی به فاصله (Space) تبدیل می‌شوند
func (pf *ProfanityFilter) normalizeText(text string) string {
	text = strings.ToLower(text)
	text = pf.replacer.Replace(text)

	var b strings.Builder
	b.Grow(len(text))
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}
	return b.String()
}

// normalizeBannedWord کلمات ممنوعه را در زمان راه‌اندازی نرمال می‌کند.
func normalizeBannedWord(word string) string {
	word = strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	if word == "" {
		return ""
	}

	return strings.ToLower(word)
}
