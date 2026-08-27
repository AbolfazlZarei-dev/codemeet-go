# کیبوردها و دکمه‌ها

CodeMeet Bot API از عناصر تعاملی برای ساخت UI داخل گفتگو پشتیبانی می‌کند.

## Inline Keyboard

نمونه:

```go
markup := &models.InlineKeyboardMarkup{
    InlineKeyboard: [][]models.InlineKeyboardButton{
        {
            {Text: "تأیید", CallbackData: "confirm"},
            {Text: "لغو", CallbackData: "cancel"},
        },
        {
            {Text: "وب‌سایت", URL: "https://example.com"},
        },
    },
}

_, err := bot.API().Messages().SendWithKeyboard(
    ctx,
    chatID,
    "یک گزینه را انتخاب کنید:",
    markup,
)
```

هر `InlineKeyboardButton` باید `text` و یک action داشته باشد.

Actionهای مستندشده:

- `callback_data`
- `url`
- `switch_inline_query`

محدودیت‌های اعلام‌شده:

- `callback_data`: بین 1 تا 64 بایت UTF-8
- حداکثر 100 row
- در هر row بین 1 تا 8 دکمه

## Callback Query

بعد از کلیک کاربر، Update شامل `callback_query` می‌شود.

```go
bot.OnCallback("confirm", func(ctx context.Context, q *models.CallbackQuery) {
    bot.AnswerCallback(ctx, q.ID, "تأیید شد", false)
})
```

اگر در نسخه‌ی پروژه‌ی شما helper مستقیم برای callback ثبت نشده باشد، می‌توان Update را از Dispatcher پردازش کرد.

## Answer Callback

برای پایان دادن به حالت loading دکمه:

```go
err := bot.AnswerCallback(
    ctx,
    callbackID,
    "عملیات انجام شد",
    true,
)
```

`show_alert=true` پیام را به صورت Alert نمایش می‌دهد.

## Reply Keyboard

Reply Keyboard در پایین صفحه‌ی چت نمایش داده می‌شود:

```go
markup := models.ReplyKeyboardMarkup{
    Keyboard: [][]models.KeyboardButton{
        {
            {Text: "پروفایل"},
            {Text: "تنظیمات"},
        },
    },
    ResizeKeyboard: true,
}
```

## حذف Reply Keyboard

از `ReplyKeyboardRemove` استفاده می‌شود:

```go
markup := models.ReplyKeyboardRemove{
    RemoveKeyboard: true,
}
```

## Force Reply

API همچنین مفهوم `ForceReply` را برای وادار کردن UI به حالت reply پشتیبانی می‌کند.

> ساختار دقیق مدل‌ها را از پکیج `models` نسخه‌ی نصب‌شده‌ی پروژه‌ی خود بگیرید تا با نسخه‌ی کتابخانه هماهنگ بماند.
