

# ۴. مدیریت گروه‌ها و کانال‌ها

متدهای مربوط به مدیریت چت‌ها، اعضا و تنظیمات گروه از طریق `bot.API().Chat()` قابل دسترسی هستند.

## اطلاعات چت و اعضا
شما می‌توانید اطلاعات کامل یک چت، تعداد اعضا و لیست مدیران را استعلام کنید:

```go
// دریافت اطلاعات چت
chat, _ := bot.API().Chat().GetChat(ctx, chatID)

// دریافت تعداد اعضا
count, _ := bot.API().Chat().GetChatMemberCount(ctx, chatID)

// دریافت لیست ادمین‌ها
admins, _ := bot.API().Chat().GetChatAdministrators(ctx, chatID)

// بررسی وضعیت عضویت یک کاربر خاص
member, _ := bot.API().Chat().GetChatMember(ctx, chatID, userID)
```

## مدیریت اعضا (بن، اخراج، ارتقا)
ربات‌ها در صورت داشتن دسترسی کافی می‌توانند اعضا را مدیریت کنند:

```go
// اخراج (بن) کاربر
bot.API().Chat().BanChatMember(ctx, chatID, userID, 0, false)

// لغو اخراج (آنبن)
bot.API().Chat().UnbanChatMember(ctx, chatID, userID, true)

// ارتقا به ادمین
rights := &models.ChatAdministratorRights{
    CanDeleteMessages: true,
    CanInviteUsers:    true,
}
bot.API().Chat().PromoteChatMember(ctx, chatID, userID, rights)

// محدود کردن کاربر (مثلاً ارسال پیام ممنوع شود)
perms := &models.ChatPermissions{
    CanSendMessages: false,
}
bot.API().Chat().RestrictChatMember(ctx, chatID, userID, perms, time.Now().Add(1*time.Hour).Unix())
```

## تنظیمات گروه و کانال
تغییر عنوان، توضیحات، عکس و مدیریت پیام‌های سنجاق شده:

```go
// تغییر عنوان گروه
bot.API().Chat().SetChatTitle(ctx, chatID, "عنوان جدید گروه")

// تغییر توضیحات گروه
bot.API().Chat().SetChatDescription(ctx, chatID, "توضیحات جدید")

// تغییر عکس گروه (آپلود فایل)
bot.API().Chat().SetChatPhoto(ctx, chatID, "/path/to/photo.jpg")

// سنجاق کردن یک پیام در بالای چت
bot.API().Chat().PinMessage(ctx, chatID, messageID, true) // true = بدون صدای نوتیفیکیشن

// حذف سنجاق پیام خاص
bot.API().Chat().UnpinMessage(ctx, chatID, messageID)

// حذف سنجاق تمام پیام‌ها
bot.API().Chat().UnpinAllMessages(ctx, chatID)
```

## لینک‌های دعوت
ساخت و مدیریت لینک‌های دعوت عضویت در گروه:

```go
// ساخت لینک دعوت ساده
link, _ := bot.API().Chat().ExportChatInviteLink(ctx, chatID)

// ساخت لینک دعوت با محدودیت اعضا و زمان انقضا
req := &models.CreateChatInviteLinkRequest{
    ChatID:      chatID,
    MemberLimit: 10, // حداکثر ۱۰ نفر
    ExpireDate:  time.Now().Add(24 * time.Hour).Unix(), // ۲۴ ساعت دیگر منقضی شود
}
newLink, _ := bot.API().Chat().CreateChatInviteLink(ctx, req)

