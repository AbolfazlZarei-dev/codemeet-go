# 26 — Reliability، Performance و Observability

## Defaults مهم

Bot با این اجزای runtime ساخته می‌شود:

- Dispatcher با 200 worker
- Retry Policy پیش‌فرض
- Rate Limiter با 30 RPS
- Logger سطح Info
- API Client
- Cache در صورت فعال‌سازی

## Retry

Retry از exponential backoff و در صورت فعال بودن jitter استفاده می‌کند.

## Rate Limit

Rate limiter از token bucket و burst پشتیبانی می‌کند.

## Dispatcher

Bounded queue و Worker Pool از رشد نامحدود queue جلوگیری می‌کنند.

## Cache

Sharded cache با یک cleanup scheduler مرکزی تعداد goroutineهای پس‌زمینه را کاهش می‌دهد.

## API performance

- sync.Pool برای buffer
- streaming JSON
- multipart streaming
- connection pooling
- HTTP/2
- atomic metrics
- Circuit Breaker

## Singleflight

GetMe و cache GetOrSet برای جلوگیری از duplicate work از singleflight استفاده می‌کنند.

## Graceful shutdown

```go
defer bot.Close()
```

این کار منابع runtime را آزاد می‌کند.

## Observability

Bot:

```go
bot.Stats()
bot.Uptime()
bot.HealthCheck(ctx)
```

API:

```go
bot.API().Client().StatsSnapshot()
```

Webhook:

```go
server.Stats()
```

Contribها نیز بسته به ماژول statistics اختصاصی دارند.
