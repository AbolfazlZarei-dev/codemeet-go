# Webhook

Package: `webhook`

## تنظیم

```go
err := bot.SetWebhook(
    ctx,
    "https://example.com/webhook",
    "secret-token",
)
```

## Server داخلی

```go
cfg := webhook.DefaultConfig()

server := webhook.New(
    bot.API().Client(),
    bot.Dispatcher(),
    bot.Logger(),
    cfg,
)

err := server.Start(ctx)
```

## Config

```go
type Config struct {
    ListenAddr        string
    Path              string
    SecretToken       string
    ReadTimeout       time.Duration
    WriteTimeout      time.Duration
    IdleTimeout       time.Duration
    ReadHeaderTimeout time.Duration
    MaxHeaderBytes    int
    MaxBodySize       int64
    HTTPS             bool
    CertFile          string
    KeyFile           string
}
```

پیش‌فرض‌ها:

```text
ListenAddr = :8443
Path = /webhook
ReadTimeout = 10s
WriteTimeout = 10s
IdleTimeout = 120s
ReadHeaderTimeout = 5s
MaxHeaderBytes = 1MB
MaxBodySize = 10MB
```

## مسیرهای Server

```text
<cfg.Path>  POST  webhook
/health     GET   health check
/metrics    GET   metrics
```

## Secret

در صورت تنظیم Secret، Header زیر بررسی می‌شود:

```text
X-CodeMeet-Bot-Api-Secret-Token
```

## HTTPS

می‌توانید TLS را مستقیماً در Server یا در Reverse Proxy terminate کنید.

## وضعیت Webhook

```go
info, err := bot.API().Webhook().GetInfo(ctx)
```

## حذف

```go
bot.DeleteWebhook(ctx)
bot.API().Webhook().DeleteWithDrop(ctx, true)
```
