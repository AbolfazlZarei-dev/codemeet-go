# 03 — Messages and Media

## Send Text

```go
msg, err := bot.API().Messages().SendText(ctx, chatID, "سلام")
```

## HTML

```go
msg, err := bot.API().Messages().SendHTML(
    ctx,
    chatID,
    "<b>سلام</b>",
)
```

## Markdown

```go
msg, err := bot.API().Messages().SendMarkdown(
    ctx,
    chatID,
    "*سلام*",
)
```

## Request کامل

```go
req := &methods.SendMessageRequest{
    ChatID: chatID,
    Text: "پیام",
    ParseMode: "HTML",
    ReplyToMessageID: messageID,
}

msg, err := bot.API().Messages().Send(ctx, req)
```

## Keyboard

```go
markup := models.NewInlineKeyboard(
    models.InlineRow(
        models.Btn("تأیید", "confirm"),
        models.Btn("لغو", "cancel"),
    ),
)

_, err := bot.API().Messages().SendWithKeyboard(
    ctx,
    chatID,
    "انتخاب کنید:",
    markup,
)
```

## Forward

```go
_, err := bot.API().Messages().Forward(
    ctx,
    destinationChatID,
    sourceChatID,
    messageID,
)
```

## Copy

```go
_, err := bot.API().Messages().Copy(
    ctx,
    destinationChatID,
    sourceChatID,
    messageID,
    "کپشن جدید",
)
```

## Edit

- `EditText`
- `EditTextInline`
- `EditCaption`
- `EditReplyMarkup`

## Delete

- `Delete`
- `DeleteMessages`

## Chat Action

```go
err := bot.API().Messages().SendChatAction(
    ctx,
    chatID,
    "typing",
)
```

## Callback

```go
err := bot.AnswerCallback(ctx, callbackID, "انجام شد", true)
```

## Media

`Media` متدهای زیر را پوشش می‌دهد:

- `SendPhoto`
- `SendPhotoWithParams`
- `SendVideo`
- `SendDocument`
- `SendVoice`
- `SendAudio`
- `SendAnimation`
- `SendSticker`
- `SendVideoNote`
- `SendMediaGroup`
- `SendLocation`
- `SendVenue`
- `SendContact`
- `SendPoll`
- `SendDice`
- `GetStickerSet`
- `UploadStickerFile`
- `GetFile`
- `DownloadFile`
- `SetMessageReaction`

## فایل

کتابخانه برای upload از multipart streaming استفاده می‌کند.

نمونه:

```go
_, err := bot.API().Media().SendDocument(
    ctx,
    chatID,
    "./manual.pdf",
)
```

> امضای دقیق متدهای Media را از package reference نسخه‌ی `1.0.0` استفاده کنید؛ این مستندات نام و رفتار سطح API را پوشش می‌دهند و از ساختن پارامترهای غیرموجود خودداری می‌کنند.
