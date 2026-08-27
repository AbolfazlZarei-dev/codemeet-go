# پیام و رسانه

لایه‌ی `methods.Messages` برای عملیات پیام و `methods.Media` برای ارسال رسانه طراحی شده است.

## ارسال پیام

```go
_, err := bot.API().Messages().SendText(
    ctx,
    chatID,
    "سلام!",
)
```

برای HTML:

```go
_, err := bot.API().Messages().SendHTML(
    ctx,
    chatID,
    "<b>سلام</b>",
)
```

برای MarkdownV2:

```go
_, err := bot.API().Messages().SendMarkdown(
    ctx,
    chatID,
    "*سلام*",
)
```

برای کنترل کامل:

```go
req := &methods.SendMessageRequest{
    ChatID:           chatID,
    Text:             "پیام",
    ParseMode:        models.ParseModeHTML,
    ReplyToMessageID: 25,
    DisableNotification: true,
}

msg, err := bot.API().Messages().Send(ctx, req)
```

فیلدهای اصلی SendMessageRequest شامل `chat_id`، `text`، `parse_mode`، `entities`، `reply_to_message_id`، `disable_notification`، `protect_content` و `reply_markup` هستند.

## رسانه

کتابخانه برای ارسال عکس، ویدیو، سند و Voice از multipart و streaming استفاده می‌کند. فایل‌ها می‌توانند از مسیر فایل محلی ارسال شوند.

نمونه‌ی ساده:

```go
_, err := bot.API().Media().SendPhoto(
    ctx,
    chatID,
    "./photo.jpg",
    "تصویر جدید",
)
```

در لایه‌ی API، multipart با `io.Pipe` ساخته می‌شود و فایل با `io.Copy` به request stream می‌شود.

## Forward و Copy

```go
_, err := bot.API().Messages().Forward(
    ctx,
    destinationChatID,
    sourceChatID,
    messageID,
)
```

کپی پیام بدون برچسب Forward:

```go
_, err := bot.API().Messages().Copy(
    ctx,
    destinationChatID,
    sourceChatID,
    messageID,
    "کپشن جدید",
)
```

## ویرایش

متدهای اصلی:

- `EditText`
- `EditTextInline`
- `EditCaption`
- `EditReplyMarkup`

مثال:

```go
err := bot.API().Messages().EditText(
    ctx,
    chatID,
    messageID,
    "وضعیت: تکمیل شد",
    models.ParseModeHTML,
    nil,
)
```

## حذف

```go
err := bot.API().Messages().Delete(ctx, chatID, messageID)
```

برای حذف چند پیام نیز `DeleteMessages` وجود دارد.

## پاسخ به Callback

```go
err := bot.AnswerCallback(
    ctx,
    callbackID,
    "انجام شد",
    true,
)
```

## Chat Action

Bot API برای نمایش وضعیت‌هایی مثل typing و upload action متد `sendChatAction` را ارائه می‌کند.
