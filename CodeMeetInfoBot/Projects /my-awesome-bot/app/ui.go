// app/ui.go
package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

// RenderFooter
func RenderFooter() string {
	var sb strings.Builder
	sb.WriteString("\n<blockquote>")
	sb.WriteString("━━━━━━━━━━━━━━━━━\n")
	sb.WriteString("⚡️ <b>قدرت گرفته از کتابخانه CodeMeet Go</b>\n")
	sb.WriteString("📦 مخزن کتابخانه: <a href=\"https://github.com/AbolfazlZarei-dev/codemeet-go\">CodeMeet Go</a>\n")
	sb.WriteString("👤 سازنده ربات و کتابخانه: <a href=\"https://codemeet.chat/fullstack\">Abolfazl Zarei</a>\n")
	sb.WriteString("📢 کانال CodeMeet Go: <a href=\"https://codemeet.chat/gocodemeet\">CodeMeet Go</a>\n")
	sb.WriteString("💻 دریافت سورس ربات: <a href=\"https://github.com/AbolfazlZarei-dev/codemeet-go/examples/my-awesome-bot\">my-awesome-bot</a>\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━")
	sb.WriteString("</blockquote>")
	return sb.String()
}

// BuildInfoKeyboard
func BuildInfoKeyboard() *models.InlineKeyboardMarkup {
	return models.NewInlineKeyboard(
		models.InlineRow(
			models.Btn("🔄 بروزرسانی", "refresh_info"),
			models.Btn("📊 آمار ربات", "bot_stats"),
		),
	)
}

// BuildHelpKeyboard
func BuildHelpKeyboard() *models.InlineKeyboardMarkup {
	return models.NewInlineKeyboard(
		models.InlineRow(
			models.Btn("👤 اطلاعات من", "get_my_info"),
			models.Btn("📊 آمار ربات", "bot_stats"),
		),
	)
}

// BuildMainMenu
func BuildMainMenu() *models.ReplyKeyboardMarkup {
	return models.NewReplyKeyboard(
		models.ReplyRow(models.KBtn("👤 اطلاعات من"), models.KBtn("❓ راهنما")),
		models.ReplyRow(models.KBtn("📊 آمار ربات")),
	)
}

// RenderUserInfoHTML
// welcomeMsg
func RenderUserInfoHTML(user *models.User, chat *models.Chat, welcomeMsg string) string {
	var sb strings.Builder

	if welcomeMsg != "" {
		sb.WriteString("<blockquote><b>👤 پنل اطلاعات کاربر</b>\n")
		sb.WriteString(welcomeMsg)
		sb.WriteString("</blockquote>\n\n")
	} else {
		sb.WriteString("<blockquote><b>👤 پنل اطلاعات کاربر</b></blockquote>\n\n")
	}

	// بخش اطلاعات کاربر
	sb.WriteString("<b>🧑 مشخصات فردی</b>\n")
	sb.WriteString("<pre>")
	sb.WriteString(fmt.Sprintf("🆔 آیدی        : <code>%s</code>\n", user.ID))

	fullName := user.FirstName
	if user.LastName != "" {
		fullName += " " + user.LastName
	}
	sb.WriteString(fmt.Sprintf("👤 نام         : %s\n", escapeHTML(fullName)))

	if user.Username != "" {
		sb.WriteString(fmt.Sprintf("👥 یوزرنیم     : @%s\n", escapeHTML(user.Username)))
	} else {
		sb.WriteString("👥 یوزرنیم     : <i>تنظیم نشده</i>\n")
	}

	if user.LanguageCode != "" {
		sb.WriteString(fmt.Sprintf("🌐 زبان        : %s\n", user.LanguageCode))
	} else {
		sb.WriteString("🌐 زبان        : <i>نامشخص</i>\n")
	}

	sb.WriteString(fmt.Sprintf("🤖 ربات است؟   : %s\n", boolToText(user.IsBot)))
	sb.WriteString(fmt.Sprintf("💎 پریمیوم     : %s\n", boolToText(user.IsPremium)))
	sb.WriteString("</pre>\n")

	// بخش اطلاعات چت
	sb.WriteString("<b>💬 اطلاعات گفتگو</b>\n")
	sb.WriteString("<pre>")
	sb.WriteString(fmt.Sprintf("🆔 چت آیدی     : <code>%s</code>\n", chat.ID))

	typeMap := map[string]string{
		"private": "👤 خصوصی",
		"group":   "👥 گروه",
		"channel": "📢 کانال",
	}
	chatType := typeMap[chat.Type]
	if chatType == "" {
		chatType = chat.Type
	}
	sb.WriteString(fmt.Sprintf("📦 نوع چت     : %s\n", chatType))

	if chat.Title != "" {
		sb.WriteString(fmt.Sprintf("🏷️ عنوان       : %s\n", escapeHTML(chat.Title)))
	}
	if chat.MembersCount > 0 {
		sb.WriteString(fmt.Sprintf("👨‍👩‍👧‍👦 اعضا       : %d\n", chat.MembersCount))
	}
	sb.WriteString("</pre>\n")

	sb.WriteString(RenderFooter())
	return sb.String()
}

// RenderHelpHTML متن راهنما
func RenderHelpHTML() string {
	var sb strings.Builder
	sb.WriteString("<blockquote><b>❓ مرکز راهنما و پشتیبانی</b></blockquote>\n\n")
	sb.WriteString("سلام! من یک ربات نمایش اطلاعات هستم. با استفاده از من می‌توانید مشخصات کاربری خود را ببینید.\n\n")
	sb.WriteString("<b>📝 دستورات موجود:</b>\n")
	sb.WriteString("<pre>")
	sb.WriteString("/start - شروع کار با ربات و نمایش اطلاعات\n")
	sb.WriteString("/info  - نمایش مجدد اطلاعات کاربری شما\n")
	sb.WriteString("/help  - نمایش همین پیام راهنما\n")
	sb.WriteString("/stats - نمایش آمار و وضعیت ربات\n")
	sb.WriteString("</pre>\n")
	sb.WriteString("همچنین می‌توانید از دکمه‌های منوی پایین صفحه یا دکمه‌های شیشه‌ای زیر پیام‌ها استفاده کنید.")

	sb.WriteString(RenderFooter())
	return sb.String()
}

// RenderStatsHTML متن آمار ربات
func RenderStatsHTML(uptime time.Duration, apiRequests int64) string {
	var sb strings.Builder
	sb.WriteString("<blockquote><b>📊 آمار و وضعیت سیستم</b></blockquote>\n\n")
	sb.WriteString("<b>⚙️ اطلاعات ربات</b>\n")
	sb.WriteString("<pre>")
	sb.WriteString(fmt.Sprintf("⏱️ آپتایم      : %s\n", uptime.String()))
	sb.WriteString(fmt.Sprintf("📨 درخواست‌ها : %d\n", apiRequests))
	sb.WriteString("🛡️ وضعیت      : ✅ آنلاین و پایدار\n")
	sb.WriteString("</pre>\n")

	sb.WriteString(RenderFooter())
	return sb.String()
}

// تابع کمکی برای تبدیل بولین به متن
func boolToText(b bool) string {
	if b {
		return "✅ بله"
	}
	return "❌ خیر"
}

// تابع فرار از کاراکترهای خاص HTML
func escapeHTML(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
