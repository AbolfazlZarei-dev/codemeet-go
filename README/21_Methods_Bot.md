# 21 — BotMethods

## Methods

```text
GetMe
SetName
GetName
GetNameWithLang
SetDescription
SetDescriptionWithLang
GetDescription
GetDescriptionWithLang
SetShortDescription
SetShortDescriptionWithLang
GetShortDescription
GetShortDescriptionWithLang
SetCommands
GetCommands
DeleteCommands
LogOut
Close
```

## مثال

```go
me, err := bot.API().Bot().GetMe(ctx)
```

## Commands

```go
commands := []models.BotCommand{
    {Command: "start", Description: "شروع"},
    {Command: "settings", Description: "تنظیمات"},
}

err := bot.API().Bot().SetCommands(ctx, commands, "fa")
```

## Language-aware methods

برای Name، Description و Short Description نسخه‌های `WithLang` وجود دارند تا مقدار زبان مشخص تنظیم/دریافت شود.
