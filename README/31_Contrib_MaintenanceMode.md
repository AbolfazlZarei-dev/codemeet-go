# 31 — `contrib/maintenancemode`

MaintenanceMode برای محدود کردن ربات در زمان تعمیرات است.

## Config

```go
type Config struct {
    AdminIDs []string
    IsEnabled bool
    MaintenanceMsg string
    NotifyCooldown int64
    Action ActionFunc
    CleanupInterval time.Duration
}
```

`ActionFunc`:

```go
type ActionFunc func(ctx context.Context, userID string, update *models.Update)
```

## رفتار

- Adminها از محدودیت عبور می‌کنند.
- سایر کاربران در صورت فعال بودن Maintenance متوقف می‌شوند.
- Action می‌تواند پیام یا رفتار دلخواه اجرا کند.
- NotifyCooldown از ارسال notification تکراری جلوگیری می‌کند.
- cleanup دوره‌ای notification state را پاک می‌کند.

## Runtime

```go
mm.SetEnabled(true)
if mm.IsEnabled() {
    // maintenance
}
```

## Stop

```go
mm.Stop()
```

## مثال

```go
mm := maintenancemode.New(maintenancemode.Config{
    AdminIDs: []string{"admin-id"},
    MaintenanceMsg: "ربات موقتاً در حال بروزرسانی است.",
    Action: func(ctx context.Context, userID string, update *models.Update) {
        // اطلاع‌رسانی
    },
})

bot.Use(mm.Middleware())
```
