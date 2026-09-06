# Message Methods

دسترسی:

```go
m := bot.API().Messages()
```

## SendMessageRequest

```go
type SendMessageRequest struct {
    ChatID string
    Text string
    ParseMode models.ParseMode
    Entities []models.MessageEntity
    ReplyToMessageID int
    DisableNotification bool
    ProtectContent bool
    ReplyMarkup interface{}
}
```

## ارسال

```go
msg, err := m.Send(ctx, &methods.SendMessageRequest{
    ChatID: chatID,
    Text: "سلام",
})
```

Helperها:

```go
m.SendText(ctx, chatID, text)
m.SendHTML(ctx, chatID, "<b>text</b>")
m.SendMarkdown(ctx, chatID, "**text**")
m.SendWithKeyboard(ctx, chatID, text, markup)
```

## Forward / Copy

```go
m.Forward(ctx, destination, source, messageID)
m.Copy(ctx, destination, source, messageID, caption)
```

## Edit

```go
m.EditText(ctx, chatID, messageID, text, parseMode, markup)
m.EditTextInline(ctx, inlineMessageID, text, parseMode, markup)
m.EditCaption(ctx, chatID, messageID, caption, parseMode, markup)
m.EditReplyMarkup(ctx, chatID, messageID, markup)
```

## Delete

```go
m.Delete(ctx, chatID, messageID)
m.DeleteMessages(ctx, chatID, []int{1,2,3})
```

## Action

```go
m.SendChatAction(ctx, chatID, models.ActionTyping)
```

## Callback

```go
m.AnswerCallback(ctx, &models.AnswerCallbackRequest{
    CallbackQueryID: callbackID,
    Text: "ثبت شد",
    ShowAlert: false,
})
```

Helper:

```go
m.AnswerCallbackSimple(ctx, callbackID, "ثبت شد", false)
```

## سایر

`SetMyDescription` و `ParseError` نیز در `Messages` وجود دارند.

> **یادآوری:** وجود Method الزاماً به معنی فعال بودن Endpoint متناظر در Bot API نسخه فعلی نیست.
