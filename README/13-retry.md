# Retry

Package: `retry`

## Policy

```go
type Policy struct {
    MaxAttempts int
    InitialDelay time.Duration
    MaxDelay time.Duration
    Multiplier float64
    Jitter bool
    MaxTotalTime time.Duration
}
```

## Policies

### Default

```text
3 attempts
500ms initial
10s max delay
2.0 multiplier
jitter=true
60s total
```

### Aggressive

```text
5 attempts
200ms initial
5s max
1.5 multiplier
jitter=true
30s total
```

### Conservative

```text
2 attempts
1s initial
20s max
2.0 multiplier
jitter=false
120s total
```

## اجرا

```go
err := retry.DefaultPolicy().Do(ctx, func(ctx context.Context) error {
    return operation(ctx)
})
```

فقط خطاهای `IsRetryable` دوباره اجرا می‌شوند.

## 429

اگر `APIError.RetryAfter` وجود داشته باشد، Policy همان زمان را رعایت کرده و 500ms به انتظار اضافه می‌کند.

## Backoff

```text
InitialDelay × Multiplier^(attempt-1)
```

با سقف `MaxDelay`.
