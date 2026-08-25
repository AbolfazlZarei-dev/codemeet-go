
```markdown
# ۲. مدیریت پیام‌ها و رسانه‌ها

تمام متدهای ارسال محتوا از طریق `bot.API()` در دسترس هستند. این بخش شامل ارسال متن، عکس، ویدیو، فایل، ویرایش پیام و مدیریت کیبوردها است.

## ارسال پیام متنی
شما می‌توانید پیام‌های ساده، فرمت‌دار (HTML یا Markdown) و همراه با کیبورد ارسال کنید:

```go
// ارسال متن ساده
bot.API().Messages().SendText(ctx, chatID, "سلام!")

// ارسال متن با فرمت HTML
bot.API().Messages().SendHTML(ctx, chatID, "<b>متن ضخیم</b> و <code>کد</code>")

// ارسال با کیبورد شیشه‌ای (InlineKeyboard)
keyboard := models.NewInlineKeyboard(
    models.InlineRow(
        models.Btn("تایید", "confirm_yes"),
        models.URLBtn("وب‌سایت", "https://codemeet.chat"),
    ),
)
bot.API().Messages().SendWithKeyboard(ctx, chatID, "یک گزینه انتخاب کنید:", keyboard)
```

## ارسال رسانه (عکس، ویدیو، فایل)
کتابخانه به‌صورت هوشمند تشخیص می‌دهد که اگر مسیر فایل محلی (مثل `/path/to/img.jpg`) دادید، آن را به‌صورت Multipart آپلود می‌کند و اگر آیدی فایل (media_id) دادید، به‌صورت JSON ارسال می‌کند.

```go
// ارسال عکس از روی هارد دیسک با کپشن
bot.API().Media().SendPhoto(ctx, chatID, "/path/to/image.jpg", "توضیحات عکس")

// ارسال فایل با پارامترهای کامل (ریپلای و کیبورد)
req := &methods.SendPhotoRequest{
    ChatID:           chatID,
    Photo:            "/path/to/doc.pdf",
    Caption:          "فایل شما آماده دانلود است",
    ReplyToMessageID: receivedMsgID,
    ReplyMarkup:      keyboard,
}
bot.API().Media().SendPhotoWithParams(ctx, req)
```

## ویرایش و حذف پیام
ربات‌ها می‌توانند پیام‌هایی که خودشان ارسال کرده‌اند را ویرایش یا حذف کنند:

```go
// ویرایش متن پیام (با حفظ فرمت HTML)
bot.API().Messages().EditText(ctx, chatID, msgID, "متن جدید ویرایش شده", models.ParseModeHTML, nil)

// ویرایش کیبورد شیشه‌ای یک پیام بدون تغییر متن آن
bot.API().Messages().EditReplyMarkup(ctx, chatID, msgID, newKeyboard)

// حذف پیام
bot.API().Messages().Delete(ctx, chatID, msgID)

// پاسخ به کلیک دکمه شیشه‌ای (Callback Query)
// اگر showAlert برابر true باشد پیام به‌صورت پاپ‌آپ (Alert) نمایش داده می‌شود
bot.API().Messages().AnswerCallbackSimple(ctx, callbackQueryID, "دکمه کلیک شد!", false)
```

