# Maintenance Mode

برای فعال کردن حالت تعمیرات:

```go
mm := maintenancemode.New(
    maintenancemode.DefaultConfig(),
)

mm.SetEnabled(true)
bot.Use(mm.Middleware())
```

## Config

```text
AdminIDs
IsEnabled
MaintenanceMsg
NotifyCooldown
Action
CleanupInterval
```

پیش‌فرض:

```text
IsEnabled=false
NotifyCooldown=300
CleanupInterval=10m
```

## API

```go
mm.SetEnabled(true)
mm.IsEnabled()
mm.Stop()
mm.Middleware()
```

Adminهای تعیین‌شده در `AdminIDs` برای استثنای مدیریتی استفاده می‌شوند.
