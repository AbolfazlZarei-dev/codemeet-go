# 16 — Package `retry`

## Policy

```go
type Policy struct {
    MaxAttempts  int
    InitialDelay time.Duration
    MaxDelay     time.Duration
    Multiplier   float64
    Jitter       bool
    MaxTotalTime time.Duration
}
```

## Presets

### DefaultPolicy

Policy استاندارد کتابخانه.

### AggressivePolicy

برای عملیات حساس به latency با retry بیشتر/فاصله‌ی کمتر.

### ConservativePolicy

```text
MaxAttempts: 2
InitialDelay: 1s
MaxDelay: 20s
Multiplier: 2
Jitter: false
MaxTotalTime: 120s
```

## Execute

```go
err := policy.Do(ctx, func(ctx context.Context) error {
    return operation(ctx)
})
```

فقط خطاهای retryable دوباره اجرا می‌شوند.

اگر API error دارای `retry_after` باشد، Policy همان زمان را رعایت می‌کند و buffer کوتاهی اضافه می‌کند.

## Backoff

Delay بر اساس:

```text
InitialDelay × Multiplier^(attempt-1)
```

تا `MaxDelay` محدود می‌شود.

در صورت فعال بودن Jitter، delay کمی تصادفی می‌شود.
