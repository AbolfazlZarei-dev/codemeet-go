package groupmanager

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseDuration تبدیل رشته‌هایی مثل "1h", "2d", "10m" یا عدد خام به time.Duration
func ParseDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	// اگر فقط عدد بود، فرض کن دقیقه است (پیش‌فرض بهتر برای سکوت کاربران)
	if val, err := strconv.Atoi(s); err == nil {
		return time.Duration(val) * time.Minute, nil
	}

	// جدا کردن واحد آخر (مثل s, m, h, d)
	unit := s[len(s)-1:]
	valueStr := s[:len(s)-1]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		// اگر استاندارد Go بود (مثل 1h30m) آن را پردازش کن
		return time.ParseDuration(s)
	}

	switch strings.ToLower(unit) {
	case "s":
		return time.Duration(value) * time.Second, nil
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	case "w":
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	default:
		return time.ParseDuration(s)
	}
}

// کلمات نامناسب برای فیلتر
var badWords = []string{
	"fuck", "shit", "idiot", "bitch", "asshole",
	"کیرو", "کص", "کون", "جنده", "خارکسه",
}

// Regex برای تشخیص حروف عربی و کلمات نامناسب
var arabicRegex = regexp.MustCompile(`[\x{0600}-\x{06FF}]`)
var profanityRegex *regexp.Regexp

func init() {
	profanityRegex = regexp.MustCompile("(?i)(" + strings.Join(badWords, "|") + ")")
}

// توابع بررسی متن
func containsProfanity(text string) bool {
	return profanityRegex.MatchString(text)
}

func containsArabic(text string) bool {
	return arabicRegex.MatchString(text)
}
