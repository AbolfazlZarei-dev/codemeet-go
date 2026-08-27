
# CodeMeet Go — Complete Documentation

> **Version:** `1.0.0`  
> **Author:** Abolfazl Zarei  
> **Repository:** `github.com/AbolfazlZarei-dev/codemeet-go`  
> **Bot API:** `https://botapi.codemeet.chat`

CodeMeet Go یک کتابخانه‌ی Go برای ساخت Botهای CodeMeet است که لایه‌ی Bot API را با typeهای Go، متدهای سطح بالا، Dispatcher، Middleware، Polling، Webhook، Retry، Rate Limit، Cache، Logging، Metrics و ابزارهای کمکی قابل استفاده در پروژه‌های واقعی ارائه می‌کند.

این مستندات بر اساس سورس کتابخانه و مستندات Bot API ارائه‌شده برای نسخه `1.0.0` تهیه شده‌اند. ساختار پروژه شامل هسته‌ی `codemeet`، لایه‌ی `api`، `methods`، `models`، `dispatcher`، `middleware`، `polling`، `webhook`، `cache`، `ratelimit`، `retry`، `logger`، `errors` و مجموعه‌ی `contrib` است.

---

# 📚 Documentation Coverage

مستندات CodeMeet Go به‌صورت بخش‌بندی‌شده طراحی شده‌اند تا هم برای توسعه‌دهنده‌ی تازه‌کار قابل استفاده باشند و هم برای توسعه‌ی Production مرجع مناسبی برای API کتابخانه باشند.

| بخش | وضعیت | توضیح |
|---|---|---|
| Bot API Client | ✅ Supported | ارتباط HTTP و اجرای درخواست‌های Bot API |
| Typed Models | ✅ Supported | تبدیل ساختارهای API به Typeهای Go |
| Messages | ✅ Supported | ارسال، ویرایش، حذف، Forward و Copy پیام |
| Media | ⚠️ Mixed | بخشی از قابلیت‌ها پیاده‌سازی شده و بخشی برای API آینده آماده شده‌اند |
| Keyboards | ⚠️ API Dependent | ساخت و ارسال Keyboard بر اساس مدل‌های کتابخانه |
| Callback Query | ✅ Supported | پاسخ به Callbackهای دریافتی |
| Chat Management | ✅ Supported | مدیریت اطلاعات و عملیات Chat |
| Bot Profile | ✅ Supported | اطلاعات و تنظیمات پروفایل Bot |
| Commands | ✅ Supported | مدیریت Commandهای Bot |
| Updates | ✅ Supported | دریافت و پردازش Updateها |
| Long Polling | ✅ Supported | دریافت Update با `getUpdates` |
| Webhook | ✅ Supported | دریافت Update از طریق HTTP Webhook |
| Dispatcher | ✅ Supported | توزیع Update بین Handlerها |
| Middleware | ✅ Supported | افزودن رفتارهای مشترک به Pipeline |
| Retry | ✅ Supported | Retry با Backoff و Jitter |
| Rate Limit | ✅ Supported | کنترل نرخ درخواست‌ها |
| Concurrency Limit | ✅ Supported | کنترل تعداد عملیات همزمان |
| Cache | ✅ Supported | Cache با TTL |
| Sharded Cache | ✅ Supported | Cache تقسیم‌شده برای کاهش contention |
| Circuit Breaker | ✅ Supported | محافظت از API در برابر خطاهای متوالی |
| Logger | ✅ Supported | Logging سطح‌بندی‌شده |
| API Statistics | ✅ Supported | آمار Request، Success، Error، Latency و حجم داده |
| Bot Statistics | ✅ Supported | آمار Update، Command، Message و Error |
| Health Check | ✅ Supported | بررسی ارتباط Bot با API |
| Graceful Shutdown | ✅ Supported | بستن کنترل‌شده‌ی منابع |
| Dashboard | ✅ Supported | نمایش اطلاعات و لاگ‌های Bot |
| Dashboard Authentication | ✅ Supported | محافظت از Dashboard با Authentication |
| Request ID | ✅ Supported | Trace کردن Requestها |
| `contrib/antilink` | 🧩 Contrib | تشخیص و کنترل لینک‌ها |
| `contrib/antispam` | 🧩 Contrib | تشخیص و کنترل Spam |
| `contrib/forcejoin` | 🧩 Contrib | بررسی عضویت اجباری |
| `contrib/gatekeeper` | 🧩 Contrib | سیستم Verification و CAPTCHA |
| `contrib/maintenancemode` | 🧩 Contrib | حالت تعمیرات Bot |
| `contrib/profanityfilter` | 🧩 Contrib | تشخیص کلمات نامناسب |
| `contrib/vpndetector` | 🧩 Contrib | تشخیص VPN/Proxy و داده‌های مشکوک |
| `contrib/warnsystem` | 🧩 Contrib | سیستم Warning و مدیریت تخلف |

> وضعیت‌های جدول بالا بر اساس وضعیت API و سورس نسخه `1.0.0` تعریف شده‌اند.
>
> `🧩 Contrib` به معنی قابلیت مستقل و کمکی است و جزو هسته‌ی Bot API محسوب نمی‌شود.

---

# ⚠️ API Compatibility & Future Features

CodeMeet Go با هدف فراهم کردن یک API پایدار و قابل توسعه طراحی شده است.

در نسخه `1.0.0` ممکن است بعضی از Methodها در سطح کتابخانه تعریف شده باشند، اما Endpoint یا قابلیت متناظر آن‌ها هنوز در نسخه‌ی فعلی **CodeMeet Bot API** ارائه نشده باشد.

این موارد به‌صورت **Future-Ready / Forward-Compatible API** در نظر گرفته شده‌اند.

به عبارت دیگر:

> **وجود یک Method در کتابخانه الزاماً به معنی فعال بودن Endpoint متناظر آن در Bot API نسخه‌ی فعلی نیست.**

برخی APIها از هم‌اکنون در ساختار کتابخانه تعریف شده‌اند تا با اضافه شدن قابلیت مربوطه به Bot API، نیاز به تغییر اساسی در معماری کتابخانه وجود نداشته باشد.

در صورتی که یک Method در نسخه‌ی فعلی Bot API پشتیبانی نشود، اجرای آن ممکن است با خطای API یا خطای مربوط به Endpoint مواجه شود.

این وضعیت یک **Compatibility Limitation** محسوب می‌شود و به معنی خرابی کتابخانه نیست.

پس از انتشار قابلیت مربوطه در Bot API، در صورت نیاز implementation، request schema، response handling یا signature مربوطه در یک Release جدید با specification رسمی API هماهنگ خواهد شد.

### Source of Truth

برای تشخیص پشتیبانی واقعی یک قابلیت، اولویت با این موارد است:

1. نسخه‌ی فعال CodeMeet Bot API
2. Specification رسمی Bot API
3. Implementation نسخه‌ی فعلی کتابخانه
4. مستندات این Repository

بنابراین مستندات کتابخانه نباید به‌تنهایی به‌عنوان تضمین فعال بودن یک Endpoint سمت سرور در نظر گرفته شوند.

---

# 🚀 Quick Start

ساده‌ترین Bot ممکن:

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

    bot.OnCommand("start", func(
        ctx context.Context,
        msg *models.Message,
    ) {
        _, err := bot.Reply(
            ctx,
            msg,
            "سلام! ربات CodeMeet شما فعال است.",
        )

        if err != nil {
            log.Println(err)
        }
    })

    cfg := polling.DefaultConfig()

    if err := bot.StartPolling(
        context.Background(),
        cfg,
    ); err != nil {
        log.Println(err)
    }
}
````

---

# 🤖 Bot API خام

الگوی Endpoint:

```text
https://botapi.codemeet.chat/bot{token}/{method}
```

کتابخانه این لایه‌ی HTTP را پشت `api.Client` و `methods` قرار می‌دهد.

---

# 🧠 فلسفه کتابخانه

* **Typed API:** مدل‌های JSON به structهای Go تبدیل می‌شوند.
* **Separation of concerns:** Transport، Methods، Dispatch و Runtime از هم جدا هستند.
* **Production features:** Retry، Rate Limit، Cache، Circuit Breaker و Graceful Shutdown در هسته وجود دارند.
* **Streaming:** فایل‌ها و multipart بدون نیاز به نگهداری کامل داده در RAM پردازش می‌شوند.
* **Extensibility:** Middleware و `contrib` امکان افزودن رفتارهای سطح بالاتر را می‌دهند.
* **Observability:** Logger، API statistics، Bot statistics و Webhook metrics در دسترس هستند.

---

# 🏗️ Architecture

ساختار کلی CodeMeet Go:

```text
                         CodeMeet Bot
                              │
                              ▼
                       ┌─────────────┐
                       │ API Client  │
                       └──────┬──────┘
                              │
              ┌───────────────┼────────────────┐
              ▼               ▼                ▼
          Methods           Models          Runtime
              │                                │
      ┌───────┼────────┐             ┌─────────┼─────────┐
      ▼       ▼        ▼             ▼         ▼         ▼
    Bot     Chat    Messages     Dispatcher  Polling   Webhook
                     │
                     ▼
                   Media

Runtime Services
──────────────────────────────────────────────
Retry
Rate Limiter
Concurrency Limiter
Cache
Sharded Cache
Circuit Breaker
Logger
Metrics
Dashboard
Middleware
```

---

# 📦 Package Structure

کتابخانه از چند Package مستقل تشکیل شده است:

```text
codemeet-go/
│
├── codemeet
│
├── api
├── methods
├── models
│
├── dispatcher
├── middleware
│
├── polling
├── webhook
│
├── cache
├── ratelimit
├── retry
├── logger
├── errors
│
└── contrib/
    ├── antilink
    ├── antispam
    ├── forcejoin
    ├── gatekeeper
    ├── maintenancemode
    ├── profanityfilter
    ├── vpndetector
    └── warnsystem
```

---

# ⚡ Core Features

## Typed API

مدل‌های JSON مربوط به Bot API در Package `models` به Typeهای مشخص Go تبدیل شده‌اند.

این کار باعث می‌شود توسعه‌دهنده به‌جای کار مستقیم با JSON، با Structها و Typeهای Go کار کند.

---

## API Client

Package `api` مسئول ارتباط مستقیم با Bot API است.

وظایف اصلی:

* HTTP Request
* Authentication
* JSON Encoding
* JSON Decoding
* Multipart Upload
* File Download
* Timeout
* Request ID
* Statistics
* Circuit Breaker

---

## Methods

Package `methods` یک لایه‌ی سطح بالاتر روی API Client ایجاد می‌کند.

بخش‌های اصلی:

```text
Bot
Chat
Messages
Media
Updates
Webhook
```

این معماری باعث می‌شود استفاده از API برای توسعه‌دهنده ساده‌تر و خواناتر باشد.

---

## Dispatcher

Dispatcher مسئول دریافت Update و انتقال آن به Handler مناسب است.

کتابخانه برای Dispatcher از Worker Pool استفاده می‌کند و Bot در زمان ساخت به‌صورت پیش‌فرض Dispatcher را با Workerهای متعدد ایجاد می‌کند.

---

## Middleware

Middleware امکان اضافه کردن رفتارهای مشترک به جریان پردازش Update را فراهم می‌کند.

نمونه قابلیت‌های Middleware:

* Recovery
* Logging
* Rate Limiting
* Custom Middleware

---

## Retry

کتابخانه دارای Retry Policy داخلی است.

قابلیت‌ها:

* Maximum Attempts
* Initial Delay
* Maximum Delay
* Exponential Backoff
* Multiplier
* Jitter
* Context Cancellation

---

## Rate Limiting

برای جلوگیری از ارسال بیش از حد Request می‌توان Rate Limiter را فعال کرد.

```go
bot, err := codemeet.New(
    token,
    codemeet.WithRateLimit(30),
)
```

Burst Rate Limiting نیز قابل تنظیم است.

---

## Cache

Cache برای نگهداری داده‌های موقت با TTL استفاده می‌شود.

کتابخانه علاوه بر Cache معمولی، `ShardedCache` نیز دارد.

Sharded Cache برای کاهش contention و افزایش concurrency طراحی شده است.

---

## Observability

CodeMeet Go قابلیت‌های داخلی برای مشاهده وضعیت Bot دارد:

* API Request Statistics
* Bot Statistics
* Logging
* Dashboard
* Health Check
* Webhook Metrics
* Request ID

Bot آمار مربوط به Updateها، Commandها، Messageها و Errorها را نیز نگهداری می‌کند.

---

# 📖 Documentation Index

## شروع و مفاهیم

* [01 — Getting Started](01_Getting_Started.md)
* [02 — Bot API و معماری کتابخانه](02_Bot_API_and_Architecture.md)
* [03 — Messages and Media](03_Messages_and_Media.md)
* [04 — Updates, Polling and Webhook](04_Updates_Polling_Webhook.md)
* [05 — Chats, Groups and Channels](05_Chat_Management.md)
* [06 — Bot Profile and Commands](06_Bot_Profile_and_Commands.md)
* [07 — Keyboards and Buttons](07_Keyboards_and_Buttons.md)
* [08 — Models and Update](08_Models.md)

## Package Reference

* [09 — codemeet](09_Package_Codemeet.md)
* [10 — api](10_Package_API.md)
* [11 — methods](11_Package_Methods.md)
* [12 — dispatcher](12_Package_Dispatcher.md)
* [13 — middleware](13_Package_Middleware.md)
* [14 — cache](14_Package_Cache.md)
* [15 — ratelimit](15_Package_RateLimit.md)
* [16 — retry](16_Package_Retry.md)
* [17 — logger](17_Package_Logger.md)
* [18 — errors](18_Package_Errors.md)
* [19 — polling](19_Package_Polling.md)
* [20 — webhook](20_Package_Webhook.md)

## Models و API Methods

* [21 — Bot Methods](21_Methods_Bot.md)
* [22 — Chat Methods](22_Methods_Chat.md)
* [23 — Message Methods](23_Methods_Messages.md)
* [24 — Media Methods](24_Methods_Media.md)
* [25 — Updates and Webhook Methods](25_Methods_Updates_Webhook.md)

## قابلیت‌های پیشرفته و Contrib

* [26 — Reliability, Performance and Observability](26_Reliability_Performance.md)
* [27 — contrib/antilink](27_Contrib_AntiLink.md)
* [28 — contrib/antispam](28_Contrib_AntiSpam.md)
* [29 — contrib/forcejoin](29_Contrib_ForceJoin.md)
* [30 — contrib/gatekeeper](30_Contrib_Gatekeeper.md)
* [31 — contrib/maintenancemode](31_Contrib_MaintenanceMode.md)
* [32 — contrib/profanityfilter](32_Contrib_ProfanityFilter.md)
* [33 — contrib/vpndetector](33_Contrib_VPNDetector.md)
* [34 — contrib/warnsystem](34_Contrib_WarnSystem.md)

## توسعه و استفاده واقعی

* [35 — Examples and Production Patterns](35_Examples.md)
* [36 — Errors and Troubleshooting](36_Errors_and_Troubleshooting.md)
* [37 — Package Map and API Index](37_API_Index.md)

---

# 🧩 Contrib Packages

پکیج‌های `contrib` قابلیت‌های اختیاری و سطح بالاتر هستند که برای ساخت Botهای واقعی و Production-oriented ارائه شده‌اند.

آن‌ها بخشی از هسته‌ی اجباری Bot API نیستند و توسعه‌دهنده می‌تواند بر اساس نیاز پروژه از آن‌ها استفاده کند.

| Package           | کاربرد                     |
| ----------------- | -------------------------- |
| `antilink`        | تشخیص و کنترل لینک‌ها      |
| `antispam`        | تشخیص Flood و Spam         |
| `forcejoin`       | بررسی عضویت اجباری         |
| `gatekeeper`      | Verification / CAPTCHA     |
| `maintenancemode` | فعال‌سازی حالت Maintenance |
| `profanityfilter` | تشخیص کلمات نامناسب        |
| `vpndetector`     | تشخیص VPN / Proxy          |
| `warnsystem`      | مدیریت Warning کاربران     |

برای جزئیات implementation و API هر Package به مستندات همان Package مراجعه کنید.

---

# 🛡️ Production Ready Features

CodeMeet Go فقط یک Wrapper ساده برای HTTP API نیست.

این کتابخانه برای استفاده در Botهای واقعی، سرویس‌های طولانی‌مدت و Workloadهای concurrent طراحی شده است.

قابلیت‌های Production:

```text
✓ Context-aware API
✓ Retry Policy
✓ Exponential Backoff
✓ Jitter
✓ Rate Limiting
✓ Concurrency Limiting
✓ Cache
✓ Sharded Cache
✓ Circuit Breaker
✓ Worker Pool
✓ Streaming Upload
✓ Streaming JSON Decode
✓ HTTP Timeout
✓ Graceful Shutdown
✓ Structured Logging
✓ Metrics
✓ Health Check
✓ Request ID
✓ Dashboard
```

---

# 🔐 Security

لایه‌های امنیتی و حفاظتی کتابخانه شامل مواردی مانند:

* محدودیت Body در Webhook
* محدودیت Header
* Secret Token برای Webhook
* Constant-time Secret Comparison
* Timeoutهای HTTP
* Context Cancellation
* محدودیت حجم File Download
* Recovery Middleware
* Rate Limiting
* Circuit Breaker
* Dashboard Authentication

است.

---

# 📤 Media & File Handling

Media API شامل قابلیت‌هایی مانند:

* Photo
* Video
* Document
* Voice
* Audio
* Animation
* Sticker
* Video Note
* Media Group
* Location
* Venue
* Contact
* Poll
* Dice
* Sticker Set
* Sticker Upload
* File Information
* File Download
* Message Reaction

است.

برخی از این قابلیت‌ها ممکن است بسته به وضعیت فعلی CodeMeet Bot API در نسخه‌ی `1.0.0` فعال یا غیرفعال باشند.

---

# 💬 Messages

Message API قابلیت‌های اصلی زیر را پوشش می‌دهد:

* Send Text
* HTML
* Markdown
* Reply
* Forward
* Copy
* Edit
* Delete
* Chat Action
* Reply Markup
* Callback Response

---

# 🔄 Polling

Polling برای دریافت Update بدون نیاز به Public Webhook Endpoint استفاده می‌شود.

Default Configuration شامل:

```text
Timeout:             10s
Poll Interval:       2s
Limit:               100
Buffer Size:         1000
Delete Webhook First: true
Max Retries:         5
```

این مقادیر در `polling.DefaultConfig()` تعریف شده‌اند.

---

# 🌐 Webhook

Webhook Server داخلی کتابخانه:

```text
POST /webhook
GET  /health
GET  /metrics
```

را مدیریت می‌کند.

Default:

```text
Listen Address: :8443
Path:           /webhook
Read Timeout:   10s
Write Timeout:  10s
Idle Timeout:   120s
Max Header:     1MB
Max Body:       10MB
```

در صورت فعال بودن Secret Token، Header امنیتی نیز بررسی می‌شود.

---

# 🧹 Graceful Shutdown

برای آزاد کردن منابع:

```go
defer bot.Close()
```

`Close()` اجزای داخلی کتابخانه را به‌صورت کنترل‌شده می‌بندد.

---

# ❤️ Health Check

برای بررسی اتصال Bot به API:

```go
if err := bot.HealthCheck(ctx); err != nil {
    log.Println("bot is unhealthy:", err)
}
```

Health Check برای بررسی وضعیت ارتباط Bot با API استفاده می‌شود.

---

# 📊 Statistics

آمار API:

```go
stats := bot.Stats()
```

در دسترس است و برای Monitoring و Performance Analysis قابل استفاده است.

همچنین Bot اطلاعات runtime مانند:

* Updates Processed
* Commands Executed
* Messages Sent
* Errors
* Uptime

را نگهداری می‌کند.

---

# 👨‍💻 Developer

**Abolfazl Zarei**

CodeMeet Go توسط Abolfazl Zarei توسعه داده شده است.

```text
Author:
Abolfazl Zarei

GitHub:
github.com/AbolfazlZarei-dev

Repository:
github.com/AbolfazlZarei-dev/codemeet-go
```

نسخه‌ی فعلی کتابخانه:

```text
v1.0.0
```

است.

---

# 📌 Documentation Policy

این Documentation با هدف مستندسازی API، معماری و قابلیت‌های موجود در نسخه‌ی `1.0.0` تهیه شده است.

در مواردی که API کتابخانه قابلیت‌هایی را برای نسخه‌های آینده آماده کرده باشد، وضعیت آن‌ها در مستندات با برچسب‌هایی مانند:

```text
⚠️ Compatibility Notice
🧪 Future / Prepared API
```

مشخص خواهد شد.

این کار باعث می‌شود توسعه‌دهنده بتواند تفاوت بین:

* قابلیت واقعاً پشتیبانی‌شده
* قابلیت موجود در Library ولی وابسته به Bot API
* قابلیت آماده‌شده برای آینده
* قابلیت‌های اختیاری `contrib`

را به‌وضوح تشخیص دهد.

---

# ⭐ Summary

CodeMeet Go یک Wrapper ساده‌ی HTTP نیست.

این پروژه یک SDK لایه‌بندی‌شده برای ساخت Botهای CodeMeet است که از لایه‌ی Transport و API Client شروع می‌شود و تا Dispatcher، Runtime Services، Middleware و ابزارهای Production ادامه پیدا می‌کند.

```text
CodeMeet Bot API
       │
       ▼
   API Client
       │
       ▼
    Methods
       │
       ▼
     Models
       │
       ▼
   Dispatcher
       │
 ┌─────┼─────┐
 ▼     ▼     ▼
Bot  Polling Webhook
       │
       ▼
Middleware
       │
       ▼
Production Services
────────────────────────
Retry
Rate Limit
Cache
Circuit Breaker
Logging
Metrics
Dashboard
Health Check
```

> **CodeMeet Go v1.0.0 — Built for CodeMeet Bot Development in Go.**

