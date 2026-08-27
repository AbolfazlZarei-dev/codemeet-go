# دریافت Update: Polling و Webhook

Update رویداد ورودی ربات است. کتابخانه دو مسیر اصلی برای دریافت آن دارد.

## Long Polling

برای شروع ساده یا محیطی که Endpoint عمومی ندارد:

```go
cfg := polling.DefaultConfig()
err := bot.StartPolling(context.Background(), cfg)
```

تنظیمات اصلی Polling:

| گزینه | مقدار پیش‌فرض |
|---|---:|
| Timeout | 10 ثانیه |
| PollInterval | 2 ثانیه |
| Limit | 100 |
| BufferSize | 1000 |
| DeleteWebhookFirst | true |
| MaxRetries | 5 |

Polling از `getUpdates` استفاده می‌کند و offset را برای جلوگیری از پردازش تکراری Updateها نگه می‌دارد.

## Webhook

برای Production می‌توان HTTP Server داخلی کتابخانه را اجرا کرد:

```go
cfg := webhook.DefaultConfig()
cfg.ListenAddr = ":8443"
cfg.Path = "/webhook"
cfg.SecretToken = "YOUR_SECRET"

err := bot.StartWebhook(context.Background(), cfg)
```

تنظیمات Webhook شامل Listen Address، Path، Secret Token، timeoutهای HTTP، محدودیت Header/Body و در صورت نیاز TLS است.

## امنیت Webhook

اگر `SecretToken` تنظیم شده باشد، سرور Header زیر را بررسی می‌کند:

```text
X-CodeMeet-Bot-Api-Secret-Token
```

مقایسه‌ی Secret با `subtle.ConstantTimeCompare` انجام می‌شود.

Body نیز با `http.MaxBytesReader` محدود می‌شود؛ مقدار پیش‌فرض `MaxBodySize` برابر 10MB است.

## Health و Metrics

Webhook Server دو endpoint داخلی دارد:

```text
/health
/metrics
```

`/health` وضعیت ساده‌ی `ok` می‌دهد و `/metrics` تعداد request و error را برمی‌گرداند.

## مدل Update

Update می‌تواند شامل مواردی مانند:

- `message`
- `callback_query`
- `my_chat_member`
- `chat_join_request`

باشد.

Helperهایی مانند `EffectiveMessage` و `EffectiveChat` برای استخراج پیام یا چت موثر وجود دارند.

## Handler

برای Command:

```go
bot.OnCommand("start", func(ctx context.Context, msg *models.Message) {
    bot.Reply(ctx, msg, "خوش آمدید")
})
```

برای پردازش عمومی Update نیز Dispatcher و Middleware قابل استفاده هستند.
