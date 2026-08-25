


# ۱. شروع به کار و پیکربندی

برای ساخت ربات، ابتدا باید توکن خود را از [@BotFather](https://codemeet.chat) در پیام‌رسان کدمیت دریافت کنید.

## ساخت ربات
ساخت ربات با استفاده از تابع `New` انجام می‌شود. این تابع به‌صورت پیش‌فرض مجموعه‌ای از تنظیمات بهینه (شامل Rate Limiter، Retry Policy و Logger) را اعمال می‌کند.

```go
bot, err := codemeet.New("YOUR_BOT_TOKEN")
if err != nil {
    log.Fatal(err)
}
```

## پیکربندی با Options
شما می‌توانید رفتار ربات را با استفاده از Option Pattern سفارشی‌سازی کنید:

```go
bot, _ := codemeet.New(token,
    codemeet.WithTimeout(30*time.Second),           // تایم‌اوت درخواست‌های HTTP
    codemeet.WithRetry(retry.AggressivePolicy()),  // سیاست تلاش مجدد تهاجمی‌تر
    codemeet.WithRateLimitBurst(50, 100),          // محدودیت ۵۰ درخواست در ثانیه با Burst ۱۰۰
    codemeet.WithShardedCache(64, 10*time.Minute), // کش ۶۴ بخشی با TTL ۱۰ دقیقه
    codemeet.WithLogger(logger.New(logger.LevelDebug)), // لاگر در حالت دیباگ
)
```

## بستن اتصالات
همیشه یادتان باشد پس از اتمام کار، اتصالات را ببندید تا منابع سیستم آزاد شوند:

```go
defer bot.Close()
```
