# AntiSpam

محافظت در برابر flood، keyword spam، Bot و commandهای بیش از حد طولانی.

## Defaults

```text
MaxMessages=8
Window=5s
Cooldown=10s
MaxWarnings=3
BanDuration=30m
DetectFlood=true
FloodThreshold=5
DetectSpamKeywords=true
BlockBots=true
MaxCommandLength=512
```

## استفاده

```go
as := antispam.New(antispam.DefaultConfig())
bot.Use(as.Middleware())
```

## مدیریت state

```go
as.BanUser(userID)
as.UnbanUser(userID)
stats := as.Stats()
```

## Config

`WarnAction` و `BanAction` callbackهای عملیاتی هستند و می‌توانید آن‌ها را به سیستم مدیریت Bot متصل کنید.
