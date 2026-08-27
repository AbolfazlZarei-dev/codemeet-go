# 06 — Bot Profile و Commands

## GetMe

```go
me, err := bot.API().Bot().GetMe(ctx)
```

در هسته، نتیجه‌ی GetMe cache می‌شود و برای جلوگیری از درخواست‌های همزمان تکراری از singleflight استفاده می‌شود.

## Name

- `SetName`
- `GetName`
- `GetNameWithLang`

## Description

- `SetDescription`
- `SetDescriptionWithLang`
- `GetDescription`
- `GetDescriptionWithLang`

## Short Description

- `SetShortDescription`
- `SetShortDescriptionWithLang`
- `GetShortDescription`
- `GetShortDescriptionWithLang`

## Commands

```go
commands := []models.BotCommand{
    {Command: "start", Description: "شروع"},
    {Command: "help", Description: "راهنما"},
}

err := bot.API().Bot().SetCommands(ctx, commands, "fa")
```

متدها:

- `SetCommands`
- `GetCommands`
- `DeleteCommands`

## Service

- `LogOut`
- `Close`

## Helperهای Bot

- `OnCommand`
- `OnMessage`
- `OnCallback`
- `OnText`
- `OnRegex`
- `Fallback`
- `Use`
- `Send`
- `SendHTML`
- `SendWithKeyboard`
- `Reply`
- `AnswerCallback`

## Dashboard

Bot قابلیت Dashboard داخلی دارد و Option زیر برای احراز هویت آن وجود دارد:

```go
codemeet.WithDashboardAuth("admin", "strong-password")
```

در صورت فعال بودن authentication، session مربوط به Dashboard مدیریت می‌شود.
