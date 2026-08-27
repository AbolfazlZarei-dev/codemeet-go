# 35 — `contrib/warnsystem`

WarnSystem سیستم هشدار per-user برای Chat است.

## Config

```go
type Config struct {
    MaxWarnings int
    WarnAction func(ctx context.Context, chatID, userID string, current, max int)
    MaxWarnAction func(ctx context.Context, chatID, userID string)
}
```

## API

```go
ws := warnsystem.New(warnsystem.Config{
    MaxWarnings: 3,
    WarnAction: func(ctx context.Context, chatID, userID string, current, max int) {
        // warning
    },
    MaxWarnAction: func(ctx context.Context, chatID, userID string) {
        // رسیدن به سقف
    },
})
```

## عملیات

```text
AddWarning
ResetWarnings
GetWarnings
Stats
Close
```

## رفتار

Warningها به صورت sharded نگه‌داری می‌شوند تا contention روی یک lock واحد کاهش پیدا کند.

## Stats

- totalWarns
- totalMaxWarns

## استفاده

```go
count := ws.AddWarning(ctx, chatID, userID)
```

> signature دقیق callbackها و return value را با نسخه‌ی نصب‌شده‌ی package هماهنگ کنید؛ رفتار اصلی package ثبت warning، reset، query و callback هنگام رسیدن به سقف است.
