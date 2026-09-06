# AntiLink

برای تشخیص URL، domain، invite و username.

```go
al := antilink.New(antilink.DefaultConfig())
bot.Use(al.Middleware())
```

## Config

```go
type Config struct {
    AllowedDomains []string
    BlockUsernames bool
    BlockInvites bool
    Action func(ctx context.Context, userID, chatID string, messageID int, reason string)
}
```

پیش‌فرض:

```text
AllowedDomains = ["codemeet.chat"]
BlockUsernames = false
BlockInvites = true
```

## Direct check

```go
blocked, reason := al.IsBlocked(text, entities)
```

دامنه‌ها normalize می‌شوند و Detector از regex برای URL/domain/Telegram invite استفاده می‌کند.
