# Examples

## /start + Inline Keyboard

```go
bot.OnCommand("start", func(ctx context.Context, msg *models.Message) {
    markup := models.NewInlineKeyboard(
        models.InlineRow(
            models.Btn("راهنما", "help"),
            models.Btn("سایت", "https://example.com"),
        ),
    )

    _, err := bot.SendWithKeyboard(
        ctx, msg.Chat.ID, "سلام 👋", markup,
    )
    if err != nil { log.Println(err) }
})
```

## Callback

```go
bot.OnCallback(func(ctx context.Context, cq *models.CallbackQuery) {
    switch cq.Data {
    case "help":
        _ = bot.AnswerCallback(ctx, cq.ID, "راهنما باز شد", false)
    }
})
```

## Regex

```go
bot.OnRegex(
    `^/echo\s+(.+)$`,
    func(ctx context.Context, msg *models.Message, matches []string) {
        _, _ = bot.Reply(ctx, msg, matches[1])
    },
)
```

## Middleware stack

```go
bot.Use(
    middleware.Recovery(bot.Logger()),
    middleware.Logging(bot.Logger()),
    middleware.RateLimit(10, time.Second),
)
```

## Media

```go
_, err := bot.API().Media().SendDocument(
    ctx, chatID, "./file.pdf", "سند",
)
```
