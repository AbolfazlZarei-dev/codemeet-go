# 36 — Examples و الگوهای Production

## Bot ساده

```go
bot, err := codemeet.New(token)
if err != nil {
    log.Fatal(err)
}
defer bot.Close()

bot.OnCommand("start", func(ctx context.Context, msg *models.Message) {
    _, _ = bot.Reply(ctx, msg, "سلام!")
})

_ = bot.StartPolling(context.Background(), polling.DefaultConfig())
```

## Bot با reliability

```go
bot, err := codemeet.New(
    token,
    codemeet.WithTimeout(30*time.Second),
    codemeet.WithRateLimitBurst(30, 60),
    codemeet.WithRetry(retry.DefaultPolicy()),
    codemeet.WithShardedCache(32, 5*time.Minute),
)
```

## Middleware chain

```go
bot.Use(
    middleware.Recovery(bot.Logger()),
    middleware.Logging(bot.Logger()),
    middleware.RateLimit(10, time.Minute),
)
```

## ترکیب contrib

```go
bot.Use(antilink.New(antilink.Config{
    BlockInvites: true,
}).Middleware())

bot.Use(profanityfilter.New(profanityfilter.Config{
    DeleteMessage: true,
    WarnUser: true,
}).Middleware())
```

## Production checklist

- Token را در environment نگه دارید.
- Webhook را فقط با HTTPS استفاده کنید.
- SecretToken تنظیم کنید.
- Rate Limit را متناسب با workload تنظیم کنید.
- Retry را بیش از حد تهاجمی نکنید.
- برای فایل‌های بزرگ streaming را حفظ کنید.
- `defer bot.Close()` داشته باشید.
- Health check و metrics را expose کنید.
- Contribها را با policy واقعی گروه تنظیم کنید.
- برای VPNDetector نتیجه را heuristic در نظر بگیرید.
