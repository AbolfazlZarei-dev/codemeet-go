# Bot Methods

دسترسی:

```go
bot.API().Bot()
```

> **هشدار:** وجود Method الزاماً به معنی فعال بودن Endpoint متناظر در Bot API نسخه فعلی نیست.

| Method | Endpoint | توضیح |
|---|---|---|
| GetMe | getMe | اطلاعات Bot |
| SetName | setMyName | تنظیم نام |
| GetName | getMyName | دریافت نام |
| GetNameWithLang | getMyName | نام با زبان |
| SetDescription | setMyDescription | تنظیم توضیحات |
| SetDescriptionWithLang | setMyDescription | توضیحات با زبان |
| GetDescription | getMyDescription | دریافت توضیحات |
| GetDescriptionWithLang | getMyDescription | توضیحات با زبان |
| SetShortDescription | setMyShortDescription | متن کوتاه |
| SetShortDescriptionWithLang | setMyShortDescription | متن کوتاه با زبان |
| GetShortDescription | getMyShortDescription | دریافت متن کوتاه |
| GetShortDescriptionWithLang | getMyShortDescription | دریافت متن کوتاه با زبان |
| SetCommands | setMyCommands | ثبت Commandها |
| GetCommands | getMyCommands | دریافت Commandها |
| DeleteCommands | deleteMyCommands | حذف Commandها |
| LogOut | logout | خروج |
| Close | — | بستن لایه BotMethods |

## نمونه

```go
me, err := bot.API().Bot().GetMe(ctx)

err = bot.API().Bot().SetName(ctx, "CodeMeet Bot")
name, err := bot.API().Bot().GetName(ctx)

commands := []models.BotCommand{
    {Command: "start", Description: "شروع"},
    {Command: "help", Description: "راهنما"},
}
err = bot.API().Bot().SetCommands(ctx, commands, "fa")
```

## زبان

Methodهای `WithLang` مقدار زبان را جداگانه دریافت می‌کنند.
