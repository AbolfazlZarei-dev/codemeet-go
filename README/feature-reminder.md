# Reminder

Package: `reminder`

این ماژول یک Middleware برای Reminder ارائه می‌کند.

## Config

```go
type Config struct {
    SendAction func(ctx context.Context, chatID, text string) error
}
```

## استفاده

```go
r := reminder.New(reminder.Config{
    SendAction: func(ctx context.Context, chatID, text string) error {
        _, err := bot.Send(ctx, chatID, text)
        return err
    },
})

bot.Use(r.Middleware())
```

`SendAction` مسئول ارسال پیام Reminder است.
