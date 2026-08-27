# پروفایل و دستورات Bot

`BotMethods` مدیریت هویت، نام، توضیحات و Commandهای ربات را ساده می‌کند.

## اطلاعات Bot

```go
me, err := bot.API().Bot().GetMe(ctx)
```

همچنین سطح Bot دارای `GetMe` با cache و singleflight است تا درخواست‌های همزمان غیرضروری کاهش پیدا کنند.

## نام Bot

```go
err := bot.API().Bot().SetName(ctx, "دستیار کدمیت")
```

دریافت:

```go
name, err := bot.API().Bot().GetName(ctx)
```

نسخه‌ی language-aware:

```go
name, err := bot.API().Bot().GetNameWithLang(ctx, "fa")
```

## Description

```go
err := bot.API().Bot().SetDescription(
    ctx,
    "این ربات دستیار شماست.",
)
```

برای زبان مشخص:

```go
err := bot.API().Bot().SetDescriptionWithLang(
    ctx,
    "توضیحات فارسی",
    "fa",
)
```

## Short Description

```go
err := bot.API().Bot().SetShortDescription(
    ctx,
    "پشتیبانی و خدمات هوشمند",
)
```

## Commands

```go
commands := []models.BotCommand{
    {Command: "start", Description: "شروع"},
    {Command: "help", Description: "راهنما"},
}

err := bot.API().Bot().SetCommands(ctx, commands, "fa")
```

دریافت:

```go
commands, err := bot.API().Bot().GetCommands(ctx, "fa")
```

حذف:

```go
err := bot.API().Bot().DeleteCommands(ctx, "fa")
```

## عملیات سرویس

`BotMethods` متدهای `LogOut` و `Close` را نیز در اختیار می‌گذارد که نتیجه‌ی Boolean برمی‌گردانند.

## Helperهای سطح Bot

```go
bot.OnCommand("help", func(ctx context.Context, msg *models.Message) {
    bot.Reply(ctx, msg, "راهنمای ربات")
})
```

این helperها ثبت handler و آمار اجرای command را با لایه‌ی داخلی Dispatcher هماهنگ می‌کنند.
