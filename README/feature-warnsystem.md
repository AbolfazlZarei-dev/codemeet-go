# Warn System

سیستم شمارش اخطار per-chat/per-user.

```go
ws := warnsystem.New(warnsystem.DefaultConfig())
defer ws.Close()
```

پیش‌فرض:

```text
MaxWarnings=3
```

## API

```go
current := ws.AddWarning(ctx, chatID, userID)
count := ws.GetWarnings(ctx, chatID, userID)
reset := ws.ResetWarnings(chatID, userID)
stats := ws.Stats()
```

## Config

- `MaxWarnings`
- `WarnAction`
- `MaxWarnAction`

وقتی تعداد اخطار به سقف برسد، `MaxWarnAction` اجرا می‌شود.
