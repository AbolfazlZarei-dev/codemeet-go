# 37 — Errors و Troubleshooting

## 400

معمولاً request یا parameter مشکل دارد:

- chat ID
- text
- markup
- media parameter
- command configuration

## 403

ممکن است Bot permission کافی نداشته باشد یا operation نیازمند role خاص باشد.

## 429

Rate limit.

راهکار:

1. `retry_after` را رعایت کنید.
2. Rate limiter را تنظیم کنید.
3. concurrency را کاهش دهید.
4. retry را با policy مناسب انجام دهید.

## 5xx

خطای سرویس. Retry و Circuit Breaker در چنین شرایطی اهمیت دارند.

## Network errors

`NetworkError` را می‌توان با:

```go
errors.AsNetworkError(err)
```

تشخیص داد.

## API errors

```go
if apiErr, ok := errors.AsAPIError(err); ok {
    log.Println(apiErr.Code)
    log.Println(apiErr.Description)
    log.Println(apiErr.RetryAfter)
}
```

## Panic

```go
bot.Use(middleware.Recovery(bot.Logger()))
```

## Polling متوقف می‌شود

Context cancellation را بررسی کنید:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

## Webhook دریافت نمی‌شود

بررسی کنید:

- URL عمومی باشد.
- HTTPS درست باشد.
- Path با تنظیم Bot یکسان باشد.
- Secret Token درست باشد.
- Body size محدودیت را رد نکرده باشد.
- health endpoint پاسخ دهد.

## Shutdown

```go
if err := bot.Close(); err != nil {
    log.Println(err)
}
```
