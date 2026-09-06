# معرفی

`CodeMeet-Go` یک SDK برای ساخت Bot با CodeMeet Bot API در Go است. نقطه ورود اصلی `codemeet.Bot` است و کنار آن لایه‌های تخصصی برای API، دریافت Update، routing، امنیت، cache، retry، rate limiting، logging و persistence ارائه شده‌اند.

## مشخصات

```text
Version: 1.1.0
Module: github.com/AbolfazlZarei-dev/codemeet-go
Author: Abolfazl Zarei
```

## مهم‌ترین قابلیت‌ها

- API wrapperهای گروه‌بندی‌شده
- `context.Context` در عملیات شبکه
- JSON و Multipart
- Request ID
- Statistics
- Circuit Breaker
- Retry با Exponential Backoff و Jitter
- Rate Limiter و Concurrency Limiter
- Cache ساده و Sharded
- Dispatcher و Handlerهای Command/Text/Regex/Callback
- Middlewareهای آماده
- Long Polling و Webhook
- Dashboard
- Database
- ماژول‌های مدیریت گروه و moderation
- AI Client
- Shop Engine

## مدل ذهنی

```text
Bot
 ├─ API/Methods ──> Bot API
 ├─ Dispatcher ───> Handler
 ├─ Polling ──────> Update
 ├─ Webhook ──────> Update
 ├─ Middleware ───> Policy/Security
 ├─ Cache
 ├─ Retry
 ├─ RateLimit
 └─ Logger/Stats
```

## اصل سازگاری

Method موجود در SDK و Endpoint موجود روی Server دو چیز متفاوت‌اند. **وجود یک Method در کتابخانه الزاماً به معنی فعال بودن Endpoint متناظر آن در Bot API نسخه‌ی فعلی نیست.**

بنابراین تست integration را با API واقعی انجام دهید.
