# CodeMeet Go — Complete Documentation

> **Version:** `1.0.0`  
> **Author:** Abolfazl Zarei  
> **Repository:** `github.com/AbolfazlZarei-dev/codemeet-go`  
> **Bot API:** `https://botapi.codemeet.chat`

CodeMeet Go یک کتابخانه‌ی Go برای ساخت Botهای CodeMeet است که لایه‌ی Bot API را با typeهای Go، متدهای سطح بالا، Dispatcher، Middleware، Polling، Webhook، Retry، Rate Limit، Cache، Logging، Metrics و ابزارهای کمکی قابل استفاده در پروژه‌های واقعی ارائه می‌کند.

این مستندات بر اساس سورس کتابخانه و مستندات Bot API ارائه‌شده برای نسخه `1.0.0` تهیه شده‌اند. ساختار پروژه شامل هسته‌ی `codemeet`، لایه‌ی `api`، `methods`، `models`، `dispatcher`، `middleware`، `polling`، `webhook`، `cache`، `ratelimit`، `retry`، `logger`، `errors` و مجموعه‌ی `contrib` است.

## فهرست مستندات

### شروع و مفاهیم
- [01 — Getting Started](01_Getting_Started.md)
- [02 — Bot API و معماری کتابخانه](02_Bot_API_and_Architecture.md)
- [03 — Messages and Media](03_Messages_and_Media.md)
- [04 — Updates, Polling and Webhook](04_Updates_Polling_Webhook.md)
- [05 — Chats, Groups and Channels](05_Chat_Management.md)
- [06 — Bot Profile and Commands](06_Bot_Profile_and_Commands.md)
- [07 — Keyboards and Buttons](07_Keyboards_and_Buttons.md)
- [08 — Models and Update](08_Models.md)

### Package Reference
- [09 — codemeet](09_Package_Codemeet.md)
- [10 — api](10_Package_API.md)
- [11 — methods](11_Package_Methods.md)
- [12 — dispatcher](12_Package_Dispatcher.md)
- [13 — middleware](13_Package_Middleware.md)
- [14 — cache](14_Package_Cache.md)
- [15 — ratelimit](15_Package_RateLimit.md)
- [16 — retry](16_Package_Retry.md)
- [17 — logger](17_Package_Logger.md)
- [18 — errors](18_Package_Errors.md)
- [19 — polling](19_Package_Polling.md)
- [20 — webhook](20_Package_Webhook.md)

### Models و API Methods
- [21 — Bot Methods](21_Methods_Bot.md)
- [22 — Chat Methods](22_Methods_Chat.md)
- [23 — Message Methods](23_Methods_Messages.md)
- [24 — Media Methods](24_Methods_Media.md)
- [25 — Updates and Webhook Methods](25_Methods_Updates_Webhook.md)

### قابلیت‌های پیشرفته و Contrib
- [26 — Reliability, Performance and Observability](26_Reliability_Performance.md)
- [27 — contrib/antilink](27_Contrib_AntiLink.md)
- [28 — contrib/antispam](28_Contrib_AntiSpam.md)
- [29 — contrib/forcejoin](29_Contrib_ForceJoin.md)
- [30 — contrib/gatekeeper](30_Contrib_Gatekeeper.md)
- [31 — contrib/maintenancemode](31_Contrib_MaintenanceMode.md)
- [32 — contrib/pagination](32_Contrib_Pagination.md)
- [33 — contrib/profanityfilter](33_Contrib_ProfanityFilter.md)
- [34 — contrib/vpndetector](34_Contrib_VPNDetector.md)
- [35 — contrib/warnsystem](35_Contrib_WarnSystem.md)

### توسعه و استفاده واقعی
- [36 — Examples and Production Patterns](36_Examples.md)
- [37 — Errors and Troubleshooting](37_Errors_and_Troubleshooting.md)
- [38 — Package Map and API Index](38_API_Index.md)

## Quick Start

```go
package main

import (
    "context"
    "log"

    codemeet "github.com/AbolfazlZarei-dev/codemeet-go"
    "github.com/AbolfazlZarei-dev/codemeet-go/models"
    "github.com/AbolfazlZarei-dev/codemeet-go/polling"
)

func main() {
    bot, err := codemeet.New("YOUR_BOT_TOKEN")
    if err != nil {
        log.Fatal(err)
    }
    defer bot.Close()

    bot.OnCommand("start", func(ctx context.Context, msg *models.Message) {
        _, err := bot.Reply(ctx, msg, "سلام! ربات CodeMeet شما فعال است.")
        if err != nil {
            log.Println(err)
        }
    })

    cfg := polling.DefaultConfig()
    if err := bot.StartPolling(context.Background(), cfg); err != nil {
        log.Println(err)
    }
}
```

## Bot API خام

الگوی Endpoint:

```text
https://botapi.codemeet.chat/bot{token}/{method}
```

کتابخانه این لایه‌ی HTTP را پشت `api.Client` و `methods` قرار می‌دهد.

## فلسفه کتابخانه

- **Typed API:** مدل‌های JSON به structهای Go تبدیل می‌شوند.
- **Separation of concerns:** Transport، Methods، Dispatch و Runtime از هم جدا هستند.
- **Production features:** Retry، Rate Limit، Cache، Circuit Breaker و Graceful Shutdown در هسته وجود دارند.
- **Streaming:** فایل‌ها و multipart بدون نیاز به نگهداری کامل داده در RAM پردازش می‌شوند.
- **Extensibility:** Middleware و `contrib` امکان افزودن رفتارهای سطح بالاتر را می‌دهند.
- **Observability:** Logger، API statistics، Bot statistics و Webhook metrics در دسترس هستند.

## نکته درباره `contrib`

پکیج‌های `contrib` قابلیت‌های کمکی‌اند و بخشی از API پایه‌ی Bot نیستند. هر کدام مستقل هستند و از Middleware یا Callbackهای مخصوص خود استفاده می‌کنند.

> فایل `contrib/pagination/pagination.go` در فهرست ساختار پروژه‌ی ارائه‌شده وجود دارد، اما متن سورس آن در dump کد 292KB که در اختیار قرار گرفت، قابل بازیابی نبود؛ بنابراین در مستندات Pagination از ساختن API خیالی خودداری شده و این محدودیت صریحاً ثبت شده است.
