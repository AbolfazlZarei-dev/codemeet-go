# 29 — `contrib/forcejoin`

ForceJoin اجازه‌ی عبور کاربر را به عضویت او در Channelهای مشخص وابسته می‌کند.

## Config

```go
type Config struct {
    RequiredChannels []string
    CacheTTL time.Duration
    AdminIDs []string
    CheckMembership func(ctx context.Context, userID, chatID string) (bool, error)
    NotJoinedAction func(ctx context.Context, userID, chatID string)
}
```

## رفتار

- عضویت کاربر در RequiredChannels بررسی می‌شود.
- وضعیت عضویت برای `CacheTTL` cache می‌شود.
- AdminIDs می‌توانند از بررسی عبور کنند.
- اگر کاربر عضو نباشد `NotJoinedAction` اجرا می‌شود.
- Statistics تعداد checks، blocked و passed را نگه می‌دارد.

## مثال

```go
fj := forcejoin.New(forcejoin.Config{
    RequiredChannels: []string{"channel-a", "channel-b"},
    CacheTTL: 5 * time.Minute,
    CheckMembership: func(ctx context.Context, userID, chatID string) (bool, error) {
        return true, nil
    },
    NotJoinedAction: func(ctx context.Context, userID, chatID string) {
        // ارسال پیام join
    },
})

bot.Use(fj.Middleware())
```

## مدیریت cache

```go
fj.ClearUserCache(userID)
```

## Stats

```go
fj.Stats()
```
