# Profanity Filter

```go
pf := profanityfilter.New(
    profanityfilter.DefaultConfig(),
)
bot.Use(pf.Middleware())
```

## Config

```text
BannedWords
DeleteMessage
WarnUser
WarnText
Action
```

پیش‌فرض نمونه شامل:

```text
idiot
stupid
fuck
dumb
```

Normalization برای برخی شکل‌های leetspeak نیز وجود دارد، مانند:

```text
1d10t → idiot
f@ck  → fuck
```

## Action

```go
func(ctx context.Context, userID, chatID string, messageID int, reason string)
```

را می‌توان برای اجرای رفتار سفارشی استفاده کرد.
