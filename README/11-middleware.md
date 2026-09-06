# Middleware

Package: `middleware`

## Recovery

```go
bot.Use(middleware.Recovery(bot.Logger()))
```

Panic را recover و در صورت وجود Logger، stack را ثبت می‌کند.

## Logging

```go
bot.Use(middleware.Logging(bot.Logger()))
```

شروع، پایان و مدت پردازش Update را log می‌کند.

## RateLimit

```go
bot.Use(middleware.RateLimit(10, time.Second))
```

محدودیت per-user دارد و از ساختار sharded برای state کاربران استفاده می‌کند.

## MetricsCounter

```go
c := middleware.NewMetricsCounter()
c.Inc("messages")
snapshot := c.Snapshot()
```

## Timeout

برای محدود کردن زمان اجرای Handler.

## BotOnly / UserOnly / AdminOnly

برای محدود کردن اجرای Handler بر اساس نوع حساب یا دسترسی.

## Blacklist / Whitelist

برای جلوگیری یا اجازه دادن به کاربران خاص.

## Middleware سفارشی

```go
mw := func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
    return func(ctx context.Context, u *models.Update) {
        // before
        next(ctx, u)
        // after
    }
}
bot.Use(mw)
```
