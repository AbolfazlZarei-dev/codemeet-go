# 07 — Keyboards و Buttons

## Inline Keyboard

سازنده‌های helper:

- `NewInlineKeyboard`
- `InlineRow`
- `Btn`
- `URLBtn`
- `WebAppBtn`
- `SwitchInlineBtn`

نمونه:

```go
keyboard := models.NewInlineKeyboard(
    models.InlineRow(
        models.Btn("تأیید", "confirm"),
        models.Btn("لغو", "cancel"),
    ),
    models.InlineRow(
        models.URLBtn("وب‌سایت", "https://example.com"),
    ),
)
```

## Reply Keyboard

- `NewReplyKeyboard`
- `ReplyRow`
- `KBtn`
- `ContactBtn`
- `LocationBtn`

## مدل‌های Keyboard

- `InlineKeyboardMarkup`
- `InlineKeyboardButton`
- `WebAppInfo`
- `LoginURL`
- `ReplyKeyboardMarkup`
- `KeyboardButton`
- `KeyboardButtonPollType`
- `ReplyKeyboardRemove`
- `ForceReply`

## Callback

وقتی کاربر روی Inline Button با `callback_data` کلیک می‌کند، Update دارای `CallbackQuery` می‌شود.

```go
bot.OnCallback("confirm", func(ctx context.Context, q *models.CallbackQuery) {
    _ = bot.AnswerCallback(ctx, q.ID, "تأیید شد", false)
})
```

## Reply markup

Markup را می‌توان هنگام ارسال پیام با `SendWithKeyboard` به پیام متصل کرد.

> ساختار دقیق JSON مدل‌های keyboard در package `models` تعریف شده است.
