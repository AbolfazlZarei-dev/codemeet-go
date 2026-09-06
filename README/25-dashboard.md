# Dashboard

## Start

```go
err := bot.StartDashboard(ctx, ":8080")
```

## Auth

```go
codemeet.WithDashboardAuth("admin", "strong-password")
```

## Polling Dashboard

```go
bot.StartPollingDashboard(polling.DefaultConfig())
bot.StopPollingDashboard()
```

## State

```go
bot.RunMode()
bot.Stats()
bot.Uptime()
```

DashboardWriter لاگ‌های اخیر را نگهداری می‌کند و در implementation فعلی سقف 200 رکورد دارد.

> Dashboard عمومی را بدون authentication مناسب در Internet قرار ندهید.
