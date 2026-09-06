# Update

مدل `models.Update` رویداد ورودی Bot است.

## فیلدها

```text
UpdateID
Message
EditedMessage
CallbackQuery
ChannelPost
EditedChannelPost
MyChatMember
ChatMember
ChatJoinRequest
```

## Helperها

```go
u.Type()
u.EffectiveMessage()
u.EffectiveUser()
u.EffectiveChat()
```

## Message

برای پیام جدید:

```go
bot.OnMessage(func(ctx context.Context, msg *models.Message) {
    // msg.Chat, msg.From, msg.Text, ...
})
```

## CallbackQuery

```go
bot.OnCallback(func(ctx context.Context, cq *models.CallbackQuery) {
    fmt.Println(cq.ID, cq.Data)
})
```

## انواع دریافت

```text
Long Polling → getUpdates → Update
Webhook      → HTTP POST  → Update
```

## Offset

در Polling، بعد از پردازش معمولاً:

```text
nextOffset = update_id + 1
```

استفاده می‌شود.
