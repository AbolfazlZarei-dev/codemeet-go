# شروع کار

## CodeMeet Go چیست؟

CodeMeet Go یک کتابخانه‌ی Go برای ساخت ربات‌های CodeMeet است. هدف کتابخانه این است که بخش بزرگی از کارهای تکراری Bot API را به APIهای typed و قابل استفاده در Go تبدیل کند.

نسخه‌ی فعلی کتابخانه `1.0.0` است و نویسنده‌ی پروژه **Abolfazl Zarei** است.

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

`New` توکن خالی را قبول نمی‌کند. Base URL پیش‌فرض:

```text
https://botapi.codemeet.chat
```

## Options

هنگام ساخت Bot می‌توان Option اضافه کرد:

```go
bot, err := codemeet.New(
    "YOUR_BOT_TOKEN",
    codemeet.WithTimeout(60*time.Second),
    codemeet.WithRateLimit(30),
    codemeet.WithCache(5*time.Minute),
)
```

Options موجود در هسته شامل تغییر Base URL، HTTP Client، Timeout، Retry Policy، Rate Limit، Burst Rate Limit، Cache، Sharded Cache، Logger، احراز هویت Dashboard و Middleware است.

## ساختار داخلی

Bot این اجزای اصلی را در کنار هم قرار می‌دهد:

- `api.Client`: ارتباط HTTP با Bot API
- `methods.Methods`: facade متدهای API
- `dispatcher.Dispatcher`: مسیریابی Update به Handler
- `polling.Poller`: دریافت Update با Long Polling
- `webhook.Server`: دریافت Update با HTTP Webhook
- `ratelimit.Limiter`: کنترل نرخ درخواست
- `retry.Policy`: تکرار درخواست‌های قابل retry
- `cache.Cache` / `ShardedCache`: نگهداری داده‌های موقت
- `logger.Logger`: logging
- `models`: مدل‌های داده‌ی API

## ارسال یک پیام

```go
_, err := bot.API().Messages().SendText(
    ctx,
    chatID,
    "سلام از CodeMeet Go!",
)
```

یا از helper سطح Bot:

```go
_, err := bot.Reply(ctx, msg, "پاسخ شما آماده است.")
```

## بستن منابع

در پایان اجرای برنامه:

```go
defer bot.Close()
```

Shutdown کتابخانه منابع Cache، Dispatcher، Rate Limiter، API Client و Logger را می‌بندد.
