# 15 — Package `ratelimit`

## Token Bucket

```go
limiter := ratelimit.New(30)
```

یا:

```go
limiter := ratelimit.NewWithBurst(30, 60)
```

## Operations

```go
Wait(ctx)
TryWait()
WaitTimeout(ctx, timeout)
Available()
Rate()
Total()
Dropped()
Close()
```

## Concurrency Limiter

```go
cl := ratelimit.NewConcurrencyLimiter(50)

if err := cl.Acquire(ctx); err != nil {
    return
}
defer cl.Release()
```

این limiter تعداد عملیات همزمان را با semaphore محدود می‌کند.

## استفاده در Bot

```go
codemeet.WithRateLimit(30)
codemeet.WithRateLimitBurst(30, 60)
```

Rate limiting در لایه‌ی methods نیز استفاده می‌شود.
