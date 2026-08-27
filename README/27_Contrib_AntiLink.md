# 27 — `contrib/antilink`

سیستم AntiLink یک Middleware برای جلوگیری از لینک‌های ناخواسته در پیام‌هاست.

## Config

```go
type Config struct {
    AllowedDomains []string
    BlockUsernames bool
    BlockInvites bool
    Action func(ctx context.Context, userID, chatID string, messageID int, reason string)
}
```

## رفتار

Middleware:

1. Message را بررسی می‌کند.
2. اگر Text خالی باشد Caption را بررسی می‌کند.
3. hidden text links را از MessageEntity تشخیص می‌دهد.
4. Invite link را بررسی می‌کند.
5. URL/domain را بررسی می‌کند.
6. در صورت ممنوع بودن `Action` را اجرا می‌کند.

## Hidden links

Text linkهای داخل entity نیز بررسی می‌شوند؛ فقط regex روی متن کافی نیست.

## Allowed domains

`AllowedDomains` دامنه‌های مجاز را تعریف می‌کند.

## Username blocking

با `BlockUsernames` می‌توان usernameهای تشخیص‌داده‌شده را نیز مسدود کرد.

## Invite blocking

```go
BlockInvites: true
```

برای لینک‌های دعوت.

## مثال

```go
al := antilink.New(antilink.Config{
    AllowedDomains: []string{"example.com", "codemeet.chat"},
    BlockInvites: true,
    Action: func(ctx context.Context, userID, chatID string, messageID int, reason string) {
        log.Printf("blocked: %s", reason)
    },
})

bot.Use(al.Middleware())
```

## API

```text
DefaultConfig
New
Middleware
```

تشخیص دامنه با normalize کردن host و بررسی domain suffix انجام می‌شود.
