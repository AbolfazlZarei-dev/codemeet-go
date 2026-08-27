# معماری و قابلیت‌های کتابخانه

## معماری کلی

CodeMeet Go چند لایه را از هم جدا می‌کند:

```text
Bot
 ├── API Client
 │    ├── HTTP Transport
 │    ├── Stats
 │    └── Circuit Breaker
 │
 ├── Methods
 │    ├── Messages
 │    ├── Media
 │    ├── Bot
 │    ├── Chat
 │    ├── Updates
 │    └── Webhook
 │
 ├── Dispatcher
 ├── Polling
 ├── Webhook Server
 ├── Retry
 ├── Rate Limiter
 ├── Cache
 ├── Logger
 └── Middleware
```

## API Client

Client درخواست‌ها را با `context.Context` اجرا می‌کند و Headerهای استاندارد مانند `Accept` و `User-Agent` را تنظیم می‌کند.

برای JSON از streaming decode استفاده شده و response به 10MB محدود می‌شود.

برای multipart، کتابخانه از `io.Pipe` و streaming استفاده می‌کند.

## Connection Pool

Transport پیش‌فرض با مقادیر بالای idle connection برای workloadهای رباتی تنظیم شده و HTTP/2 را فعال می‌کند.

## Circuit Breaker

Circuit Breaker سه وضعیت دارد:

- Closed
- Open
- Half-Open

در Half-Open فقط یک Probe اجازه‌ی عبور دارد.

آستانه‌ی پیش‌فرض 10 failure و reset timeout پیش‌فرض 30 ثانیه است.

## Dispatcher

Dispatcher برای جدا کردن دریافت Update از اجرای Handler استفاده می‌شود و Worker Pool دارد.

در ساخت Bot، Dispatcher با 200 worker ساخته می‌شود.

## Middleware

Middlewareهای موجود شامل:

- `Recovery`
- `Logging`
- `RateLimit`

هستند.

Recovery panic handler را می‌گیرد و stack trace را log می‌کند.

Logging زمان و نوع Update را ثبت می‌کند.

RateLimit با 64 shard برای شمارنده‌ی کاربر طراحی شده است.

## Dashboard

کد هسته شامل Dashboard داخلی برای نمایش وضعیت/لاگ‌ها است و امکان authentication برای Dashboard نیز در Optionها وجود دارد.

## Request ID

کتابخانه می‌تواند Request ID را داخل Context قرار دهد:

```go
ctx = api.WithRequestID(ctx, "request-123")
```

در صورت وجود Request ID، آن را در Header `X-Request-ID` ارسال می‌کند.

## Graceful Shutdown

`Bot.Close()` اجزای داخلی را به شکل متمرکز می‌بندد تا goroutine و connectionهای باز باقی نمانند.
