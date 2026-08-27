# 28 — `contrib/antispam`

AntiSpam برای کنترل flood، spam keyword، bot abuse، cooldown و ban طراحی شده است.

## Config

```go
type Config struct {
    MaxMessages int
    Window time.Duration
    Cooldown time.Duration
    MaxWarnings int
    BanDuration time.Duration
    DetectFlood bool
    FloodThreshold int
    DetectSpamKeywords bool
    SpamKeywords []string
    BlockBots bool
    MaxCommandLength int
    WarnAction func(ctx context.Context, userID string, reason string)
    BanAction func(ctx context.Context, userID string, reason string)
}
```

## Defaults

مقادیر پیش‌فرض مهم:

```text
MaxMessages = 8
Window = 5s
Cooldown = 10s
MaxWarnings = 3
BanDuration = 30m
FloodThreshold = 5
```

## معماری

State کاربران به صورت sharded نگه‌داری می‌شود و هر shard mutex مستقل دارد.

Spam keywords در startup به map normalize می‌شوند تا lookup سریع باشد.

## تشخیص

سیستم می‌تواند:

- flood
- repeated messages
- spam keywords
- excessive commands
- bot users
- rate-limit abuse

را بررسی کند.

## مدیریت Ban

```go
as.BanUser(userID)
as.UnbanUser(userID)
```

## Stats

```go
as.Stats()
```

آمار شامل allowed، blocked، flood، spam keyword، rate-limit و users banned است.

## مثال

```go
as := antispam.New(antispam.Config{
    DetectFlood: true,
    DetectSpamKeywords: true,
    SpamKeywords: []string{"spam", "buy now"},
})

bot.Use(as.Middleware())
```

Stateهای قدیمی به صورت دوره‌ای cleanup می‌شوند.
