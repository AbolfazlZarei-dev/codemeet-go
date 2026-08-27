# Reliability، Retry، Rate Limit، Cache و Performance

## Retry

Policy پیش‌فرض:

- MaxAttempts: 3
- InitialDelay: 500ms
- MaxDelay: 10s
- Multiplier: 2
- Jitter: فعال
- MaxTotalTime: 60s

نمونه‌ی Policy سفارشی:

```go
policy := retry.DefaultPolicy()
policy.MaxAttempts = 5

bot, err := codemeet.New(
    token,
    codemeet.WithRetry(policy),
)
```

Policy تهاجمی (`AggressivePolicy`) نیز در پکیج Retry وجود دارد.

Retry با context کار می‌کند و برای خطاهای قابل retry اعمال می‌شود.

## Rate Limiter

Limiter پیش‌فرض Bot برابر 30 request در ثانیه ساخته می‌شود.

Option:

```go
codemeet.WithRateLimit(30)
```

Burst:

```go
codemeet.WithRateLimitBurst(30, 60)
```

متدهای Limiter شامل `Wait`، `TryWait`، `WaitTimeout`، `Available`، `Rate`، `Total`، `Dropped` و `Close` است.

## Concurrency Limiter

برای کنترل تعداد عملیات همزمان:

```go
limiter := ratelimit.NewConcurrencyLimiter(50)

if err := limiter.Acquire(ctx); err != nil {
    return
}
defer limiter.Release()
```

## Cache

Cache ساده از TTL پشتیبانی می‌کند:

```go
c := cache.New(5 * time.Minute)
defer c.Close()

c.Set("key", value)

v, ok := c.Get("key")
```

متدهای مهم:

- `Get`
- `GetTyped`
- `Set`
- `SetWithTTL`
- `SetForever`
- `Delete`
- `Len`
- `Keys`
- `Clear`
- `GetOrSet`
- `GetOrSetWithTTL`

`GetOrSet` از `singleflight` استفاده می‌کند تا محاسبه‌ی همزمان یک key تکراری نشود.

## Sharded Cache

`NewSharded` کش را به shardهای متعدد تقسیم می‌کند و یک scheduler مرکزی برای cleanup دارد؛ در نتیجه به جای goroutine اختصاصی برای هر shard، یک cleanup loop مرکزی استفاده می‌شود.

## آمار API

Client آمار زیر را نگهداری می‌کند:

- Requests
- SuccessCount
- ErrorCount
- BytesIn
- BytesOut
- AvgLatency

```go
stats := bot.Stats()
fmt.Println(stats.Requests, stats.AvgLatency)
```

## Memory و Streaming

Response JSON با `json.Decoder` decode می‌شود و Body با `io.LimitReader` تا 10MB محدود شده است.

Multipart upload نیز streaming است؛ این طراحی برای جلوگیری از نگه داشتن کل فایل در RAM مناسب است.

## نکته‌ی مهم

اعداد Rate Limit، Timeout و Worker Pool بخشی از defaultهای نسخه‌ی فعلی هستند و در صورت تغییر نسخه باید دوباره بررسی شوند.
