// app/handlers.go
package app

import (
	"context"
	"log"

	"github.com/AbolfazlZarei-dev/codemeet-go/codemeet"
	"github.com/AbolfazlZarei-dev/codemeet-go/methods"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

func RegisterHandlers(bot *codemeet.Bot) {

	// هندلر /start (فقط یک پیام ارسال می‌شود + کیبورد پایین صفحه)
	bot.OnCommand("start", func(ctx context.Context, msg *models.Message) {
		if msg.From == nil || msg.Chat == nil {
			return
		}
		chat, err := bot.API().Chat().GetChat(ctx, msg.Chat.ID)
		if err != nil {
			chat = msg.Chat
		}

		// ترکیب متن سلام و اطلاعات در یک پیام
		text := RenderUserInfoHTML(msg.From, chat, "👋 سلام! به ربات اطلاعات خوش آمدید.\nاز منوی پایین صفحه برای دسترسی سریع استفاده کنید.")

		// ارسال فقط یک پیام همراه با کیبورد پایین صفحه (ReplyKeyboard)
		_, _ = bot.API().Messages().Send(ctx, &methods.SendMessageRequest{
			ChatID:      msg.Chat.ID,
			Text:        text,
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: BuildMainMenu(),
		})
	})

	// هندلر /info (با دکمه شیشه‌ای)
	bot.OnCommand("info", func(ctx context.Context, msg *models.Message) {
		if msg.From == nil || msg.Chat == nil {
			return
		}
		chat, err := bot.API().Chat().GetChat(ctx, msg.Chat.ID)
		if err != nil {
			chat = msg.Chat
		}
		sendUserInfo(ctx, bot, msg.From, chat)
	})

	// هندلر /help
	bot.OnCommand("help", func(ctx context.Context, msg *models.Message) {
		if msg.Chat == nil {
			return
		}
		_, _ = bot.API().Messages().Send(ctx, &methods.SendMessageRequest{
			ChatID:      msg.Chat.ID,
			Text:        RenderHelpHTML(),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: BuildHelpKeyboard(),
		})
	})

	// هندلر /stats
	bot.OnCommand("stats", func(ctx context.Context, msg *models.Message) {
		if msg.Chat == nil {
			return
		}
		sendBotStats(ctx, bot, msg.Chat.ID)
	})

	// هندلر دکمه‌های معمولی (Reply Keyboard)

	// وقتی کاربر دکمه "اطلاعات من" رو میزنه، اطلاعات با دکمه شیشه‌ای ارسال میشه
	bot.OnText("👤 اطلاعات من", func(ctx context.Context, msg *models.Message) {
		if msg.From == nil || msg.Chat == nil {
			return
		}
		chat, err := bot.API().Chat().GetChat(ctx, msg.Chat.ID)
		if err != nil {
			chat = msg.Chat
		}
		sendUserInfo(ctx, bot, msg.From, chat)
	})

	bot.OnText("❓ راهنما", func(ctx context.Context, msg *models.Message) {
		if msg.Chat == nil {
			return
		}
		_, _ = bot.API().Messages().Send(ctx, &methods.SendMessageRequest{
			ChatID:      msg.Chat.ID,
			Text:        RenderHelpHTML(),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: BuildHelpKeyboard(),
		})
	})

	bot.OnText("📊 آمار ربات", func(ctx context.Context, msg *models.Message) {
		if msg.Chat == nil {
			return
		}
		sendBotStats(ctx, bot, msg.Chat.ID)
	})

	// هندلر دکمه‌های شیشه‌ای (Inline Keyboard)
	bot.OnCallback(func(ctx context.Context, cq *models.CallbackQuery) {
		if cq.Message == nil || cq.From == nil {
			return
		}

		switch cq.Data {
		case "refresh_info":
			_ = bot.AnswerCallback(ctx, cq.ID, "اطلاعات به‌روزرسانی شد", false)
			chat, err := bot.API().Chat().GetChat(ctx, cq.Message.Chat.ID)
			if err != nil {
				chat = cq.Message.Chat
			}
			text := RenderUserInfoHTML(cq.From, chat, "")
			_ = bot.API().Messages().EditText(ctx, cq.Message.Chat.ID, cq.Message.MessageID, text, models.ParseModeHTML, BuildInfoKeyboard())

		case "get_my_info":
			_ = bot.AnswerCallback(ctx, cq.ID, "", false)
			chat, err := bot.API().Chat().GetChat(ctx, cq.Message.Chat.ID)
			if err != nil {
				chat = cq.Message.Chat
			}
			sendUserInfo(ctx, bot, cq.From, chat)

		case "bot_stats":
			_ = bot.AnswerCallback(ctx, cq.ID, "", false)
			text := RenderStatsHTML(bot.Uptime(), bot.Stats().Requests)
			_ = bot.API().Messages().EditText(ctx, cq.Message.Chat.ID, cq.Message.MessageID, text, models.ParseModeHTML, BuildInfoKeyboard())
		}
	})

	// هندلر Fallback
	bot.Fallback(func(ctx context.Context, u *models.Update) {
		if u.Message == nil || u.Message.Chat == nil {
			return
		}
		text := "⚠️ متوجه دستور شما نشدم!\nلطفاً از منوی پایین صفحه یا دستور /help استفاده کنید." + RenderFooter()
		_, _ = bot.API().Messages().Send(ctx, &methods.SendMessageRequest{
			ChatID:    u.Message.Chat.ID,
			Text:      text,
			ParseMode: models.ParseModeHTML,
		})
	})

	log.Println("Handlers registered successfully.")
}

// تابع کمکی برای ارسال پیام اطلاعات همراه با دکمه شیشه‌ای
func sendUserInfo(ctx context.Context, bot *codemeet.Bot, user *models.User, chat *models.Chat) {
	text := RenderUserInfoHTML(user, chat, "")
	_, _ = bot.API().Messages().Send(ctx, &methods.SendMessageRequest{
		ChatID:      chat.ID,
		Text:        text,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: BuildInfoKeyboard(),
	})
}

// تابع کمکی برای ارسال آمار ربات
func sendBotStats(ctx context.Context, bot *codemeet.Bot, chatID string) {
	text := RenderStatsHTML(bot.Uptime(), bot.Stats().Requests)
	_, _ = bot.API().Messages().Send(ctx, &methods.SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: BuildInfoKeyboard(),
	})
}
