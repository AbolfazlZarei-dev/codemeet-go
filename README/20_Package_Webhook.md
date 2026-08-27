# 20 — Package `webhook`

## Config

```go
type Config struct {
    ListenAddr string
    Path string
    SecretToken string
    ReadTimeout time.Duration
    WriteTimeout time.Duration
    IdleTimeout time.Duration
    ReadHeaderTimeout time.Duration
    MaxHeaderBytes int
    MaxBodySize int64
    HTTPS bool
    CertFile string
    KeyFile string
}
```

## DefaultConfig

برای production می‌توان Config را با مقادیر مناسب deployment تغییر داد.

## Server

```go
server := webhook.New(apiClient, dispatcher, logger, cfg)
err := server.Start(ctx)
```

## Handler

Webhook body به `models.Update` decode می‌شود و به Dispatcher ارسال می‌شود.

## Security

- Secret token با constant-time comparison
- محدودیت body size
- محدودیت Header
- timeoutهای HTTP
- HTTPS اختیاری از طریق CertFile/KeyFile

## Health

```text
/health
```

## Metrics

```text
/metrics
```

Server آمار request و error را نگه می‌دارد.

## Stats

```go
server.Stats()
```
