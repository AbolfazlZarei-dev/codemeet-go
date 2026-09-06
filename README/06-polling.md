# Long Polling

Package: `polling`

## استفاده

```go
err := bot.StartPolling(
    ctx,
    polling.DefaultConfig(),
)
```

## Config

```go
type Config struct {
    Timeout            int
    PollInterval       time.Duration
    Limit              int
    AllowedUpdates     []string
    BufferSize         int
    DeleteWebhookFirst bool
    MaxRetries         int
}
```

## مقادیر پیش‌فرض v1.1.0

```text
Timeout = 10
PollInterval = 2s
Limit = 100
BufferSize = 1000
DeleteWebhookFirst = true
MaxRetries = 5
```

## Config سفارشی

```go
cfg := polling.Config{
    Timeout: 20,
    PollInterval: time.Second,
    Limit: 100,
    BufferSize: 2000,
    DeleteWebhookFirst: true,
    MaxRetries: 5,
}
err := bot.StartPolling(ctx, cfg)
```

## Poller مستقل

```go
p := polling.New(
    client,
    dispatcher,
    logger,
    cfg,
)
```

برای اشتراک Limiter و Retry:

```go
p := polling.NewWithLimiter(
    client,
    dispatcher,
    logger,
    cfg,
    limiter,
    retryPolicy,
)
```

## کنترل Offset

```go
p.Offset()
p.ResetOffset()
```

## Conflict

Polling و Webhook را همزمان فعال نکنید مگر اینکه معماری سرویس شما صراحتاً آن را پشتیبانی کند.
