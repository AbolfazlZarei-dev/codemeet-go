# Advanced API

## Accessors

```go
bot.API()
bot.Dispatcher()
bot.Cache()
bot.Logger()
bot.RateLimiter()
bot.RetryPolicy()
bot.Token()
bot.BaseURL()
```

## Options

```go
bot, err := codemeet.New(
    token,
    codemeet.WithBaseURL("https://botapi.codemeet.chat"),
    codemeet.WithTimeout(30*time.Second),
    codemeet.WithRetry(retry.DefaultPolicy()),
    codemeet.WithRateLimitBurst(30, 60),
    codemeet.WithCache(10*time.Minute),
)
```

## Health

```go
if err := bot.HealthCheck(ctx); err != nil {
    log.Println(err)
}
```

## Request tracing

```go
ctx = api.WithRequestID(ctx, "order-123")
```

در transport این ID به `X-Request-ID` تبدیل می‌شود.
