# 19 — Package `polling`

## Config

```go
type Config struct {
    Timeout int
    PollInterval time.Duration
    Limit int
    AllowedUpdates []string
    BufferSize int
    DeleteWebhookFirst bool
    MaxRetries int
}
```

## Defaults

```text
Timeout = 10
PollInterval = 2s
Limit = 100
BufferSize = 1000
DeleteWebhookFirst = true
MaxRetries = 5
```

## Poller

```go
New(client, dispatcher, logger, cfg)
Start(ctx)
Offset()
ResetOffset()
```

## رفتار

1. در صورت نیاز Webhook قبلی حذف می‌شود.
2. `getUpdates` با offset اجرا می‌شود.
3. Updateها به Dispatcher تحویل داده می‌شوند.
4. offset جلو می‌رود.
5. خطاها با retry/backoff مدیریت می‌شوند.
6. 429 با `retry_after` مدیریت می‌شود.
7. Context cancellation باعث توقف Polling می‌شود.
