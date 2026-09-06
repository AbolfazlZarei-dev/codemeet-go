# BotFather

## ساخت Bot

فرآیند کلی:

1. ورود به `@BotFather`
2. ایجاد Bot
3. دریافت Token
4. نگهداری Token در Secret/Environment
5. ساخت `codemeet.Bot`

## Token URL

لایه API URL را با الگوی زیر می‌سازد:

```text
https://botapi.codemeet.chat/bot<TOKEN>/<method>
```

## مدیریت پروفایل

Methodهای مربوط به پروفایل:

```text
GetMe
SetName / GetName / GetNameWithLang
SetDescription / GetDescription / ...WithLang
SetShortDescription / GetShortDescription / ...WithLang
SetCommands / GetCommands / DeleteCommands
```

> Method بودن این wrapperها تضمین‌کننده فعال بودن Endpoint سمت سرور نیست.
