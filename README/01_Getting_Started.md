# 01 — Getting Started

## پیش‌نیاز

- Go
- Token ربات CodeMeet
- دسترسی به Bot API

## نصب

```bash
go get github.com/AbolfazlZarei-dev/codemeet-go
```

## ساخت Bot

```go
bot, err := codemeet.New("YOUR_BOT_TOKEN")
if err != nil {
    log.Fatal(err)
}
defer bot.Close()
```

توکن خالی باعث `ValidationError` می‌شود.

## Optionها

```go
bot, err := codemeet.New(
    token,
    codemeet.WithBaseURL("https://botapi.codemeet.chat"),
    codemeet.WithTimeout(30*time.Second),
    codemeet.WithRateLimit(30),
    codemeet.WithCache(5*time.Minute),
)
```

Optionهای هسته:

- `WithBaseURL`
- `WithHTTPClient`
- `WithTimeout`
- `WithRetry`
- `WithRateLimit`
- `WithRateLimitBurst`
- `WithCache`
- `WithShardedCache`
- `WithLogger`
- `WithDashboardAuth`
- `WithMiddleware`
- `WithoutLogger`

## اولین Handler

```go
bot.OnMessage(func(ctx context.Context, msg *models.Message) {
    if msg == nil {
        return
    }
    _, _ = bot.Reply(ctx, msg, "پیامت دریافت شد.")
})
```

## Command

```go
bot.OnCommand("start", func(ctx context.Context, msg *models.Message) {
    _, _ = bot.Reply(ctx, msg, "شروع شد.")
})
```

## متن مشخص

```go
bot.OnText("سلام", func(ctx context.Context, msg *models.Message) {
    _, _ = bot.Reply(ctx, msg, "سلام 👋")
})
```

## Regex

```go
bot.OnRegex(`^/id(?:\s+(.+))?$`, func(ctx context.Context, msg *models.Message) {
    _, _ = bot.Reply(ctx, msg, "Regex matched.")
})
```

## Fallback

```go
bot.Fallback(func(ctx context.Context, update *models.Update) {
    // رویدادهایی که handler اختصاصی ندارند.
})
```

## استفاده از API

```go
msg, err := bot.API().Messages().SendText(
    ctx,
    chatID,
    "Hello CodeMeet",
)
```

## Health Check

```go
if err := bot.HealthCheck(ctx); err != nil {
    log.Println(err)
}
```

## Shutdown

```go
defer bot.Close()
```

`Close` به ترتیب Cache، Dispatcher، Rate Limiter، API Client و Logger را می‌بندد.

## Build

ساختار پروژه Makefile نیز دارد:

```bash
make deps
make tidy
make run
make build
make fmt
make vet
make clean
```
