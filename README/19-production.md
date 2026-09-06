# Production

## چک‌لیست

- Token در Secret/Environment
- Timeout محدود
- Retry برای خطاهای transient
- Rate Limit متناسب با quota
- Cache برای داده‌های پرتکرار
- Recovery در Dispatcher
- Webhook با HTTPS و Secret
- Health Check
- Metrics
- Graceful Shutdown
- عدم log کردن token/secret
- تست integration برای Endpointهای استفاده‌شده

## Retry + Rate Limit

```go
bot, err := codemeet.New(
    token,
    codemeet.WithRetry(retry.DefaultPolicy()),
    codemeet.WithRateLimitBurst(30, 60),
    codemeet.WithTimeout(30*time.Second),
)
```

## Shutdown

```go
defer bot.Close()
```

برای ماژول‌هایی که lifecycle جدا دارند نیز `Close`/`Stop` مربوطه را اجرا کنید.

## Reverse Proxy

```text
Client → HTTPS → IIS/Nginx/Apache → CodeMeet-Go
```

TLS می‌تواند در Proxy terminate شود.
