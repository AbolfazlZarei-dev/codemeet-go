# 04 — Updates، Polling و Webhook

CodeMeet Bot API برای دریافت رویداد دو مدل دارد:

1. Long Polling با `getUpdates`
2. Webhook با `setWebhook`

## Long Polling

```go
cfg := polling.DefaultConfig()
if err := bot.StartPolling(ctx, cfg); err != nil {
    log.Println(err)
}
```

Default:

| Field | Default |
|---|---:|
| Timeout | 10 |
| PollInterval | 2s |
| Limit | 100 |
| BufferSize | 1000 |
| DeleteWebhookFirst | true |
| MaxRetries | 5 |

Polling offset را نگه می‌دارد و Updateهای جدید را از `getUpdates` می‌گیرد.

در صورت فعال بودن `DeleteWebhookFirst`، قبل از شروع polling، Webhook قبلی حذف می‌شود.

برای 429، Poller مقدار `retry_after` را رعایت می‌کند و یک buffer کوتاه نیز اضافه می‌کند.

## Webhook

```go
cfg := webhook.DefaultConfig()
cfg.ListenAddr = ":8443"
cfg.Path = "/webhook"
cfg.SecretToken = "secret"

if err := bot.StartWebhook(ctx, cfg); err != nil {
    log.Println(err)
}
```

Config شامل:

- ListenAddr
- Path
- SecretToken
- ReadTimeout
- WriteTimeout
- IdleTimeout
- ReadHeaderTimeout
- MaxHeaderBytes
- MaxBodySize
- HTTPS
- CertFile
- KeyFile

است.

## امنیت

Secret Token با مقایسه‌ی constant-time بررسی می‌شود.

Body با `MaxBodySize` محدود می‌شود.

## Endpointها

Webhook server:

```text
/webhook
/health
/metrics
```

را فراهم می‌کند.

## Update

Update می‌تواند شامل:

- `message`
- `callback_query`
- `my_chat_member`
- `chat_join_request`

باشد.

## Effective helpers

```go
msg := update.EffectiveMessage()
user := update.EffectiveUser()
chat := update.EffectiveChat()
```

این helperها برای handlerهای عمومی بسیار مفیدند.

## Webhook API methods

- `Set`
- `GetInfo`
- `Delete`
- `DeleteWithDrop`

## انتخاب روش

**Polling:** توسعه، localhost و محیط بدون endpoint عمومی.

**Webhook:** deploymentهای production با endpoint عمومی و HTTPS.
