package groupmanager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

func (gm *GroupManager) handlePublicCommand(ctx context.Context, msg *models.Message, userID, cmd string) bool {
	chatID := msg.Chat.ID
	switch cmd {
	case gm.cfg.Commands.Help:
		gm.sendHelpMessage(ctx, chatID)
		return true
	case gm.cfg.Commands.ID:
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, fmt.Sprintf("🆔 شما: `%s`\n🆔 گروه: `%s`", userID, chatID), nil)
		return true
	case gm.cfg.Commands.Rules:
		if val, ok := gm.rulesStore.Get(chatID); ok {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, val.(string), nil)
		} else {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, "❌ قوانین تنظیم نشده.", nil)
		}
		return true
	case gm.cfg.Commands.Report:
		if msg.ReplyToMessage == nil {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, "❌ روی پیام شخص ریپلای بزنید.", nil)
			return true
		}
		target := msg.ReplyToMessage.From
		targetName := "نامشخص"
		targetID := ""
		if target != nil {
			targetName = target.FullName()
			targetID = target.ID
		}
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, fmt.Sprintf("🚨 گزارش تخلف!\nکاربر: %s\nآیدی: `%s`", targetName, targetID), nil)
		return true
	}
	return false
}

func (gm *GroupManager) handleAdminCommand(ctx context.Context, msg *models.Message, userID, cmd, args string) bool {
	chatID := msg.Chat.ID
	var targetID, targetName string
	if msg.ReplyToMessage != nil && msg.ReplyToMessage.From != nil {
		targetID = msg.ReplyToMessage.From.ID
		targetName = msg.ReplyToMessage.From.FullName()
	}

	switch cmd {
	case gm.cfg.Commands.Ban, gm.cfg.Commands.Kick:
		if targetID == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.NotReplied, nil)
			return true
		}
		_ = gm.cfg.BanUserAction(ctx, chatID, targetID)
		if cmd == gm.cfg.Commands.Kick {
			_ = gm.cfg.UnbanUserAction(ctx, chatID, targetID)
		}
		finalMsg := strings.ReplaceAll(gm.cfg.Messages.BanMsg, "{user}", targetName)
		finalMsg = strings.ReplaceAll(finalMsg, "{reason}", args)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, finalMsg, nil)
		return true

	case gm.cfg.Commands.Mute:
		if targetID == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.NotReplied, nil)
			return true
		}
		duration, _ := ParseDuration(args)
		muteUntil := time.Now().Add(duration).Unix()
		if duration == 0 {
			muteUntil = 0
		}
		_ = gm.cfg.RestrictUserAction(ctx, chatID, targetID, muteUntil, &models.ChatPermissions{})
		gm.mutedUsers.Set(chatID+"_"+targetID, muteUntil)
		finalMsg := strings.ReplaceAll(gm.cfg.Messages.MuteMsg, "{user}", targetName)
		finalMsg = strings.ReplaceAll(finalMsg, "{duration}", duration.String())
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, finalMsg, nil)
		return true

	case gm.cfg.Commands.Unmute:
		if targetID == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.NotReplied, nil)
			return true
		}
		unmutePerms := &models.ChatPermissions{CanSendMessages: true, CanSendAudios: true, CanSendDocuments: true, CanSendPhotos: true, CanSendVideos: true, CanSendVideoNotes: true, CanSendVoiceNotes: true, CanSendOtherMessages: true, CanAddWebPagePreviews: true}
		_ = gm.cfg.RestrictUserAction(ctx, chatID, targetID, 0, unmutePerms)
		gm.mutedUsers.Delete(chatID + "_" + targetID)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, "🔊 کاربر آزاد شد.", nil)
		return true

	case gm.cfg.Commands.Warn:
		if targetID == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.NotReplied, nil)
			return true
		}
		warns := gm.warnManager.AddWarning(chatID, targetID)
		warnMsg := strings.ReplaceAll(gm.cfg.Messages.WarnMsg, "{user}", targetName)
		warnMsg = strings.ReplaceAll(warnMsg, "{current}", fmt.Sprintf("%d", warns))
		warnMsg = strings.ReplaceAll(warnMsg, "{max}", fmt.Sprintf("%d", gm.cfg.MaxWarns))
		warnMsg = strings.ReplaceAll(warnMsg, "{reason}", args)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, warnMsg, nil)
		if warns >= gm.cfg.MaxWarns {
			_ = gm.cfg.BanUserAction(ctx, chatID, targetID)
		}
		return true

	case gm.cfg.Commands.Unwarn:
		if targetID == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.NotReplied, nil)
			return true
		}
		gm.warnManager.ResetWarnings(chatID, targetID)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, fmt.Sprintf("✅ اخطارهای %s پاک شد.", targetName), nil)
		return true

	case gm.cfg.Commands.Lock:
		if args == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, "❌ نوع قفل؟ مثال: قفل لینک", nil)
			return true
		}
		gm.lockManager.SetLock(chatID, args)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, strings.ReplaceAll(gm.cfg.Messages.LockSuccess, "{lockType}", args), nil)
		return true

	case gm.cfg.Commands.Unlock:
		if args == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, "❌ نوع قفل؟", nil)
			return true
		}
		gm.lockManager.RemoveLock(chatID, args)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, strings.ReplaceAll(gm.cfg.Messages.UnlockSuccess, "{lockType}", args), nil)
		return true

	case gm.cfg.Commands.LockAll:
		gm.lockManager.SetLock(chatID, "media")
		gm.lockManager.SetLock(chatID, "links")
		gm.lockManager.SetLock(chatID, "forward")
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, "🔒 تمام رسانه‌ها قفل شدند!", nil)
		return true

	case gm.cfg.Commands.UnlockAll:
		gm.lockManager.RemoveLock(chatID, "media")
		gm.lockManager.RemoveLock(chatID, "links")
		gm.lockManager.RemoveLock(chatID, "forward")
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, "🔓 تمام رسانه‌ها باز شدند!", nil)
		return true

	case gm.cfg.Commands.SlowMode:
		if args == "" {
			gm.slowModeStore.Delete(chatID + "_" + userID)
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, "❌ زمان را مشخص کنید. مثال: اسلوومود 30s", nil)
			return true
		}
		dur, _ := ParseDuration(args)
		gm.slowModeStore.Set(chatID+"_"+userID, time.Now().Add(-dur))
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, fmt.Sprintf("🐢 حالت آهسته فعال شد: هر %s یک پیام.", dur.String()), nil)
		return true

	case gm.cfg.Commands.Settings:
		locks := gm.lockManager.GetLocksList(chatID)
		text := "<b>⚙ تنظیمات گروه:</b>\n\n"
		if len(locks) > 0 {
			text += "<b>قفل‌های فعال:</b> " + strings.Join(locks, ", ")
		} else {
			text += "<b>قفل‌های فعال:</b> هیچ قفلی فعال نیست."
		}
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, text, nil)
		return true

	case gm.cfg.Commands.Pin:
		if msg.ReplyToMessage == nil {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.NotReplied, nil)
			return true
		}
		_ = gm.cfg.PinMessageAction(ctx, chatID, msg.ReplyToMessage.MessageID, false)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, "📌 پیام سنجاق شد.", nil)
		return true

	case gm.cfg.Commands.Unpin:
		if msg.ReplyToMessage == nil {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.NotReplied, nil)
			return true
		}
		_ = gm.cfg.UnpinMessageAction(ctx, chatID, msg.ReplyToMessage.MessageID)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, "🗑 سنجاق برداشته شد.", nil)
		return true

	case gm.cfg.Commands.Admins:
		if gm.cfg.GetChatAdminsAction == nil {
			return true
		}
		admins, err := gm.cfg.GetChatAdminsAction(ctx, chatID)
		if err != nil {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, "❌ خطا در دریافت ادمین‌ها.", nil)
			return true
		}
		text := "👑 لیست ادمین‌ها:\n\n"
		for _, admin := range admins {
			name := "نامشخص"
			if admin.User != nil {
				name = admin.User.FullName()
			}
			text += fmt.Sprintf("▫️ %s (%s)\n", name, admin.Status)
		}
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, text, nil)
		return true

	case gm.cfg.Commands.Members:
		if gm.cfg.GetChatMemberCount == nil {
			return true
		}
		count, err := gm.cfg.GetChatMemberCount(ctx, chatID)
		if err != nil {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, "❌ خطا در دریافت تعداد اعضا.", nil)
			return true
		}
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, fmt.Sprintf("👥 اعضا: %d نفر", count), nil)
		return true

	case gm.cfg.Commands.Info:
		if msg.ReplyToMessage == nil || msg.ReplyToMessage.From == nil {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.NotReplied, nil)
			return true
		}
		target := msg.ReplyToMessage.From
		warns := gm.warnManager.GetWarnings(chatID, target.ID)
		text := fmt.Sprintf("👤 اطلاعات کاربر\n\n▫️ نام: %s\n▫️ آیدی: `%s`\n⚠️ اخطارها: %d از %d", target.FullName(), target.ID, warns, gm.cfg.MaxWarns)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, text, nil)
		return true

	case gm.cfg.Commands.SetRules:
		if args == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, "❌ متن قوانین؟", nil)
			return true
		}
		gm.rulesStore.Set(chatID, args)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, "✅ قوانین ذخیره شد.", nil)
		return true

	case gm.cfg.Commands.Purge:
		if msg.ReplyToMessage == nil {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.NotReplied, nil)
			return true
		}
		startID := msg.ReplyToMessage.MessageID
		endID := msg.MessageID
		for i := startID; i <= endID; i++ {
			_ = gm.cfg.DeleteMessageAction(ctx, chatID, i)
		}
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, "🧹 پاکسازی شد.", nil)
		return true

	case gm.cfg.Commands.Gban:
		if targetID == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.NotReplied, nil)
			return true
		}
		gm.gbans.Set(targetID, true)
		_ = gm.cfg.BanUserAction(ctx, chatID, targetID)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, fmt.Sprintf("🌍 %s گبن شد!", targetName), nil)
		return true

	case gm.cfg.Commands.Ungban:
		if targetID == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.NotReplied, nil)
			return true
		}
		gm.gbans.Delete(targetID)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, fmt.Sprintf("✅ گبن %s لغو شد.", targetName), nil)
		return true

	case gm.cfg.Commands.LockWord:
		if args == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, "❌ کلمه را مشخص کنید.", nil)
			return true
		}
		gm.blacklistManager.AddWord(chatID, args)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, fmt.Sprintf("🚫 کلمه '%s' مسدود شد.", args), nil)
		return true

	case gm.cfg.Commands.UnlockWord:
		if args == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, "❌ کلمه را مشخص کنید.", nil)
			return true
		}
		gm.blacklistManager.RemoveWord(chatID, args)
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, fmt.Sprintf("✅ کلمه '%s' آزاد شد.", args), nil)
		return true

	case gm.cfg.Commands.Note:
		if args == "" {
			_, _ = gm.cfg.SendMessageAction(ctx, chatID, "❌ متن یادداشت؟", nil)
			return true
		}
		parts := strings.SplitN(args, " ", 2)
		if len(parts) < 2 {
			return true
		}
		gm.noteManager.AddNote(chatID, parts[0], parts[1])
		_, _ = gm.cfg.SendMessageAction(ctx, chatID, fmt.Sprintf("📝 یادداشت '%s' ذخیره شد.", parts[0]), nil)
		return true
	}
	return false
}

func (gm *GroupManager) sendHelpMessage(ctx context.Context, chatID string) {
	helpText := "📚 راهنمای ربات:\n\n🚫 بن / اخراج\n🤐 سکوت 1h / آزاد\n⚠️ اخطار / حذف اخطار / اطلاعات\n🌍 گبن / حذف گبن\n🧹 پاکسازی\n🔒 قفل لینک / باز کردن قفل لینک\n🔒 قفل همه / باز کردن همه\n🚫 قفل کلمه word / باز کردن کلمه word\n📝 یادداشت keyword text\n📌 سنجاق / حذف سنجاق\n👑 ادمین ها / اعضا\n📝 تنظیم قوانین text / قوانین\n⚙ تنظیمات (نمایش قفل‌ها)\n🐢 اسلوومود 30s\n🆔 آیدی / گزارش\n🔌 /start / /stop"
	_, _ = gm.cfg.SendMessageAction(ctx, chatID, helpText, nil)
}

func (gm *GroupManager) isBotAdmin(userID string) bool { _, ok := gm.admins[userID]; return ok }

func (gm *GroupManager) isChatAdmin(ctx context.Context, userID, chatID string) bool {
	if gm.isBotAdmin(userID) {
		return true
	}
	if val, ok := gm.adminCache.Load(chatID + "_" + userID); ok {
		return val.(bool)
	}
	if gm.cfg.GetChatMemberAction != nil {
		member, err := gm.cfg.GetChatMemberAction(ctx, chatID, userID)
		isAdmin := err == nil && (member.Status == "creator" || member.Status == "administrator")
		gm.adminCache.Store(chatID+"_"+userID, isAdmin)
		time.AfterFunc(5*time.Minute, func() { gm.adminCache.Delete(chatID + "_" + userID) })
		return isAdmin
	}
	return false
}

// CheckLockViolation بررسی تخلفات قفل‌ها
func (gm *GroupManager) CheckLockViolation(msg *models.Message) string {
	chatID := msg.Chat.ID
	text := msg.Text
	entities := msg.Entities

	// اگر متن خالی بود، کپشن و انتیتی‌های کپشن را بررسی می‌کنیم
	if text == "" {
		text = msg.Caption
		entities = msg.CaptionEntities
	}

	if gm.lockManager.IsLocked(chatID, "links") {
		// اگر ضد لینک تنظیم شده باشد، از آن استفاده می‌کنیم
		if gm.antiLink != nil {
			if blocked, _ := gm.antiLink.IsBlocked(text, entities); blocked {
				return "ارسال لینک ممنوع است"
			}
		} else {
			// در غیر این صورت به روش ساده‌ی قبلی بررسی می‌کنیم
			if strings.Contains(text, "http://") || strings.Contains(text, "https://") || strings.Contains(text, "t.me") || strings.Contains(text, ".com") || strings.Contains(text, ".ir") {
				return "ارسال لینک ممنوع است"
			}
		}
	}

	if gm.lockManager.IsLocked(chatID, "forward") {
		if msg.ForwardDate > 0 || msg.ForwardFrom != nil || len(msg.ForwardFromMessageID) > 0 {
			return "فوروارد ممنوع است"
		}
	}
	if len(msg.Photo) > 0 && (gm.lockManager.IsLocked(chatID, "photo") || gm.lockManager.IsLocked(chatID, "media")) {
		return "ارسال عکس ممنوع است"
	}
	if len(msg.Video) > 0 && (gm.lockManager.IsLocked(chatID, "video") || gm.lockManager.IsLocked(chatID, "media")) {
		return "ارسال ویدیو ممنوع است"
	}
	if len(msg.Voice) > 0 && (gm.lockManager.IsLocked(chatID, "voice") || gm.lockManager.IsLocked(chatID, "media")) {
		return "ارسال ویس ممنوع است"
	}
	if len(msg.Document) > 0 && (gm.lockManager.IsLocked(chatID, "document") || gm.lockManager.IsLocked(chatID, "media")) {
		return "ارسال فایل ممنوع است"
	}
	if len(msg.Sticker) > 0 && gm.lockManager.IsLocked(chatID, "sticker") {
		return "ارسال استیکر ممنوع است"
	}
	if len(msg.Animation) > 0 && (gm.lockManager.IsLocked(chatID, "gif") || gm.lockManager.IsLocked(chatID, "media")) {
		return "ارسال گیف ممنوع است"
	}
	if msg.Text != "" && gm.lockManager.IsLocked(chatID, "text") {
		return "ارسال متن ممنوع است"
	}
	if gm.lockManager.IsLocked(chatID, "fosh") && containsProfanity(msg.Text) {
		return "فحش ممنوع است"
	}
	if gm.lockManager.IsLocked(chatID, "arabic") && containsArabic(msg.Text) {
		return "عربی ممنوع است"
	}
	return ""
}
