
<div align="center">

# 🚀 CodeMeet Go SDK

### یک SDK قدرتمند، مدرن و Production-Ready برای توسعه ربات‌های CodeMeet با زبان Go

<p>
  <strong>High Performance</strong>
  •
  <strong>Type-Safe</strong>
  •
  <strong>Concurrent</strong>
  •
  <strong>Modular</strong>
  •
  <strong>Production Ready</strong>
</p>

<br>

[![Go Reference](https://img.shields.io/badge/Go%20Reference-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/AbolfazlZarei-dev/codemeet-go)
[![GitHub](https://img.shields.io/badge/GitHub-181717?style=for-the-badge&logo=github&logoColor=white)](https://github.com/AbolfazlZarei-dev/codemeet-go)
[![CodeMeet](https://img.shields.io/badge/CodeMeet-Official-7B2FBE?style=for-the-badge&logo=telegram&logoColor=white)](https://codemeet.chat)
[![License](https://img.shields.io/github/license/AbolfazlZarei-dev/codemeet-go?style=for-the-badge&color=22c55e)](https://github.com/AbolfazlZarei-dev/codemeet-go/blob/main/LICENSE)

<br>

[![Documentation](https://img.shields.io/badge/📚%20Documentation-Read%20the%20Docs-6366F1?style=for-the-badge)](README/)
[![Release](https://img.shields.io/github/v/release/AbolfazlZarei-dev/codemeet-go?style=for-the-badge&logo=github&color=3b82f6)](https://github.com/AbolfazlZarei-dev/codemeet-go/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/AbolfazlZarei-dev/codemeet-go?style=for-the-badge&logo=go&color=00ADD8)](https://github.com/AbolfazlZarei-dev/codemeet-go)

</div>

---

# 📖 فهرست مطالب

- [📌 معرفی کتابخانه](#-معرفی-کتابخانه)
- [✨ ویژگی‌های کلیدی](#-ویژگیهای-کلیدی)
- [📦 نصب و راه‌اندازی](#-نصب-و-راهاندازی)
- [🚀 شروع سریع](#-شروع-سریع)
- [🌐 Bot API](#-bot-api)
- [🧩 معماری کتابخانه](#-معماری-کتابخانه)
- [🛡 پکیج‌های Contrib](#-پکیجهای-contrib)
- [📊 مانیتورینگ و Dashboard](#-مانیتورینگ-و-dashboard)
- [⚡ Performance و Concurrency](#-performance-و-concurrency)
- [📚 مستندات کامل](#-مستندات-کامل)
- [🤝 مشارکت](#-مشارکت)
- [👨‍💻 توسعه‌دهنده](#-توسعهدهنده)

---

# 📌 معرفی کتابخانه

**CodeMeet Go SDK** یک کتابخانه مدرن و ماژولار برای توسعه ربات‌های پلتفرم **CodeMeet** با زبان Go است.

هدف اصلی این پروژه فراهم کردن یک لایه‌ی قدرتمند، type-safe و قابل استفاده در پروژه‌های واقعی روی Bot API کدمیت است؛ به‌طوری که توسعه‌دهنده بتواند بدون درگیر شدن مستقیم با جزئیات HTTP و JSON، ربات‌های سریع، پایدار و قابل توسعه ایجاد کند.

کتابخانه علاوه بر API اصلی، مجموعه‌ای از ابزارهای زیرساختی و پکیج‌های کمکی را نیز ارائه می‌کند:

- HTTP Client
- Bot API Methods
- Typed Models
- Dispatcher
- Middleware
- Long Polling
- Webhook
- Retry
- Rate Limiting
- Concurrency Limiting
- In-Memory Cache
- Sharded Cache
- Circuit Breaker
- Structured Logging
- Error Handling
- Statistics
- Health Check
- Dashboard
- ابزارهای امنیتی و مدیریتی در `contrib`

---

# ✨ ویژگی‌های کلیدی

| دسته‌بندی | قابلیت‌ها | توضیح |
|---|---|---|
| 🌐 ارتباط با API | HTTP Client، Request، Multipart، Form | لایه‌ی ارتباطی مرکزی با Bot API |
| 🛡 پایداری شبکه | Retry، Circuit Breaker | مدیریت خطاها و درخواست‌های ناموفق |
| ⚡ کنترل ترافیک | Rate Limiter، Concurrency Limiter | کنترل نرخ و تعداد درخواست‌های همزمان |
| 🚀 پردازش آپدیت | Dispatcher، Worker | مدیریت و پردازش Updateها |
| 🔌 دریافت Update | Long Polling، Webhook | دو روش اصلی دریافت رویدادها |
| 💾 Cache | Cache، Sharded Cache، TTL | ذخیره سریع داده‌های موقت در حافظه |
| 🧩 Middleware | Recovery، Logging، Metrics، Timeout | کنترل رفتار پردازش Update |
| 📨 پیام | Text، HTML، Markdown، Edit، Delete | مدیریت پیام‌ها |
| 🎞 Media | Photo، Video، Document، Audio، Voice و... | ارسال و دریافت انواع Media |
| 🎛 Keyboard | Inline Keyboard، Reply Keyboard | ساخت رابط تعاملی برای ربات |
| 👥 مدیریت Chat | Member، Admin، Ban، Restrict، Invite | مدیریت گروه‌ها و کانال‌ها |
| 🤖 Bot Management | Profile، Commands | مدیریت اطلاعات و دستورات Bot |
| 📊 Observability | Statistics، Logs، Metrics، Health Check | مشاهده وضعیت و عملکرد سیستم |
| 🛡 امنیت | AntiSpam، AntiLink، Gatekeeper و... | ابزارهای امنیتی آماده |
| 🧰 ابزارهای مدیریتی | ForceJoin، WarnSystem، Maintenance Mode | امکانات مدیریتی برای ربات‌ها |
| 🔍 تحلیل فایل | VPN Detector | بررسی فایل‌ها و کانفیگ‌های مشکوک |

---

# 📦 نصب و راه‌اندازی

برای نصب کتابخانه:

```bash
go get github.com/AbolfazlZarei-dev/codemeet-go
````

سپس می‌توانید آن را در پروژه‌ی Go خود import کنید.

---

# 🚀 شروع سریع

یک ربات ساده با Long Polling:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	codemeet "github.com/AbolfazlZarei-dev/codemeet-go/codemeet"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
	"github.com/AbolfazlZarei-dev/codemeet-go/polling"
)

func main() {
	token := "YOUR_BOT_TOKEN"

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	bot, err := codemeet.New(token)
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Close()

	bot.OnMessage(func(
		ctx context.Context,
		msg *models.Message,
	) {
		fmt.Printf(
			"📥 پیام از %s: %s\n",
			msg.From.FullName(),
			msg.Text,
		)

		_, err := bot.Reply(
			ctx,
			msg,
			"👋 سلام! پیام شما دریافت شد.",
		)

		if err != nil {
			log.Println(err)
		}
	})

	bot.OnCommand(
		"start",
		func(
			ctx context.Context,
			msg *models.Message,
		) {
			_, err := bot.SendHTML(
				ctx,
				msg.Chat.ID,
				"<b>🤖 ربات با موفقیت فعال شد!</b>",
			)

			if err != nil {
				log.Println(err)
			}
		},
	)

	fmt.Println("🚀 ربات در حال اجراست...")

	cfg := polling.DefaultConfig()
	cfg.Timeout = 20

	if err := bot.StartPolling(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
```

---

# 🌐 Bot API

CodeMeet Go روی Bot API پلتفرم CodeMeet ساخته شده است.

Endpoint اصلی API:

```text
https://botapi.codemeet.chat/bot{TOKEN}/{METHOD}
```

کتابخانه جزئیات HTTP Request، Response، Encoding، Multipart، Error Handling و Retry را در لایه‌ی `api` و `methods` مدیریت می‌کند.

بنابراین توسعه‌دهنده معمولاً به جای ساخت دستی HTTP Request می‌تواند از API سطح بالای کتابخانه استفاده کند.

مثال:

```go
msg, err := bot.API().
	Messages().
	SendText(
		ctx,
		chatID,
		"سلام 👋",
	)
```

---

# 🧩 معماری کتابخانه

ساختار پروژه به شکل ماژولار طراحی شده است:

```text
codemeet-go/
│
├── api/
│   └── HTTP Client و لایه ارتباط با Bot API
│
├── methods/
│   ├── Messages
│   ├── Media
│   ├── Bot
│   ├── Chat
│   ├── Updates
│   └── Webhook
│
├── models/
│   ├── User
│   ├── Chat
│   ├── Message
│   ├── Update
│   ├── Keyboard
│   └── Media
│
├── dispatcher/
│   └── مدیریت و Dispatch کردن Updateها
│
├── middleware/
│   └── Middlewareهای آماده
│
├── polling/
│   └── Long Polling
│
├── webhook/
│   └── Webhook Server
│
├── cache/
│   └── In-Memory و Sharded Cache
│
├── ratelimit/
│   └── Rate و Concurrency Limiting
│
├── retry/
│   └── Retry Policy
│
├── logger/
│   └── Logging و JSON Logging
│
├── errors/
│   └── Error Handling
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

# 🛡 پکیج‌های Contrib

پکیج‌های `contrib` قابلیت‌های اختیاری و سطح بالاتری هستند که می‌توانند در کنار هسته‌ی کتابخانه استفاده شوند.

## 🚪 Gatekeeper

مسئول مدیریت ورود کاربران و اجرای Challenge/Captcha.

قابلیت‌ها:

* Math Captcha
* Number Captcha
* Inline Button Challenge
* مدیریت کاربران در انتظار
* مدیریت Challengeها
* Statistics
* Cleanup خودکار

مسیر:

```text
contrib/gatekeeper
```

---

## 🔍 VPN Detector

سیستم تحلیل فایل و متن برای شناسایی نشانه‌های مربوط به VPN، Proxy و برخی فرمت‌های کانفیگ.

قابلیت‌های موجود در package شامل:

* بررسی متن
* بررسی نام APK
* بررسی نام فایل کانفیگ
* بررسی فایل‌های متنی
* بررسی APK
* بررسی ZIP
* تشخیص Base64 Candidate
* بررسی JSON
* بررسی YAML
* تشخیص الگوهای Host/Port
* تشخیص Proxy Indicator
* تحلیل Binary Signature
* Statistics

مسیر:

```text
contrib/vpndetector
```

---

## ⚠️ Warn System

سیستم مدیریت اخطار کاربران.

قابلیت‌ها:

* افزودن Warning
* دریافت تعداد Warning
* Reset کردن Warning
* نگهداری وضعیت کاربران
* Cleanup
* Statistics
* پشتیبانی از Sharding

مسیر:

```text
contrib/warnsystem
```

---

## 🔐 AntiLink

Middleware برای شناسایی و مدیریت لینک‌ها.

قابلیت‌ها:

* بررسی URL
* بررسی Host
* Normalization دامنه
* تشخیص دامنه‌های مجاز
* مدیریت لینک‌های مجاز
* Block کردن پیام‌های دارای لینک غیرمجاز

مسیر:

```text
contrib/antilink
```

---

## 🛡 AntiSpam

سیستم کنترل Spam و Flood.

قابلیت‌ها:

* Rate Limit کاربران
* تشخیص Spam
* Keyword Detection
* Warning
* Ban / Unban
* User State
* Cleanup
* Statistics

مسیر:

```text
contrib/antispam
```

---

## 👥 ForceJoin

سیستم بررسی عضویت اجباری.

قابلیت‌ها:

* Middleware
* بررسی وضعیت کاربر
* Cache کاربران
* پاک‌سازی Cache
* Statistics

مسیر:

```text
contrib/forcejoin
```

---

## 🔧 Maintenance Mode

امکان قرار دادن ربات در حالت تعمیرات.

قابلیت‌ها:

* فعال / غیرفعال کردن Maintenance Mode
* بررسی وضعیت
* تشخیص Admin
* Middleware
* مدیریت Notification
* Cleanup
* Stop

مسیر:

```text
contrib/maintenancemode
```

---

## 🔤 Profanity Filter

Middleware برای بررسی و فیلتر کلمات ممنوعه.

قابلیت‌ها:

* Normalization متن
* Normalization کلمات ممنوعه
* Scan متن
* Middleware

مسیر:

```text
contrib/profanityfilter
```

---

# 📊 مانیتورینگ و Dashboard

کتابخانه دارای سیستم Statistics و Dashboard است.

در سطح Bot می‌توان از قابلیت‌هایی مانند:

```go
bot.HealthCheck()
```

و:

```go
bot.Stats()
```

و:

```go
bot.Uptime()
```

استفاده کرد.

همچنین Dashboard از طریق:

```go
bot.StartDashboard(ctx, ":9090")
```

قابل اجرا است.

نمونه:

```go
go func() {
	err := bot.StartDashboard(
		context.Background(),
		":9090",
	)

	if err != nil {
		log.Println(err)
	}
}()
```

پس از اجرا:

```text
http://localhost:9090
```

قابل دسترسی خواهد بود.

---

# ⚡ Performance و Concurrency

یکی از اهداف اصلی CodeMeet Go SDK، ارائه‌ی زیرساخت مناسب برای Botهای پرترافیک است.

بخش‌های مختلف کتابخانه برای کار با محیط‌های Concurrent طراحی شده‌اند.

### موارد مهم:

* Worker-based Dispatcher
* Bounded Dispatch Queue
* Sharded Cache
* TTL Cache
* Rate Limiting
* Concurrency Limiting
* Retry Policy
* Circuit Breaker
* Async Logger
* Multipart Streaming
* Statistics
* Graceful Shutdown

این معماری اجازه می‌دهد بخش‌های مختلف سیستم بدون وابستگی شدید به یکدیگر توسعه داده شوند.

---

# 📨 Messages و Media

کتابخانه API سطح بالایی برای ارسال و مدیریت پیام ارائه می‌کند.

### پیام متنی

```go
msg, err := bot.API().
	Messages().
	SendText(
		ctx,
		chatID,
		"سلام",
	)
```

### HTML

```go
msg, err := bot.API().
	Messages().
	SendHTML(
		ctx,
		chatID,
		"<b>سلام</b>",
	)
```

### Markdown

```go
msg, err := bot.API().
	Messages().
	SendMarkdown(
		ctx,
		chatID,
		"*سلام*",
	)
```

### Keyboard

```go
markup := models.NewInlineKeyboard(
	models.InlineRow(
		models.Btn("تأیید", "confirm"),
		models.Btn("لغو", "cancel"),
	),
)

_, err := bot.API().
	Messages().
	SendWithKeyboard(
		ctx,
		chatID,
		"انتخاب کنید:",
		markup,
	)
```

کتابخانه همچنین API مربوط به:

* Forward
* Copy
* Edit
* Delete
* Chat Action
* Callback
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
* File
* Reaction

را در لایه‌ی Methods/Media پوشش می‌دهد.

---

# 🎛 Keyboard و Buttons

مدل‌های Keyboard در package `models` قرار دارند.

### Inline Keyboard

```go
markup := models.NewInlineKeyboard(
	models.InlineRow(
		models.Btn("تأیید", "confirm"),
		models.Btn("لغو", "cancel"),
	),
)
```

همچنین Builderهای مربوط به:

* URL Button
* WebApp Button
* Switch Inline Button
* Reply Keyboard
* Contact Button
* Location Button

وجود دارند.

---

# 🧱 Middleware

کتابخانه Middlewareهای مختلفی را ارائه می‌کند:

```text
Recovery
Logging
RateLimit
Metrics
Timeout
BotOnly
UserOnly
AdminOnly
Blacklist
Whitelist
```

نمونه:

```go
bot.Use(
	middleware.Logging(),
	middleware.Recovery(),
)
```

Middlewareها را می‌توان برای کنترل جریان پردازش Updateها و اضافه کردن رفتارهای مشترک استفاده کرد.

---

# 🔄 Polling و Webhook

CodeMeet Go از دو روش اصلی دریافت Update پشتیبانی می‌کند.

### Long Polling

```go
cfg := polling.DefaultConfig()

if err := bot.StartPolling(
	ctx,
	cfg,
); err != nil {
	log.Fatal(err)
}
```

### Webhook

```go
if err := bot.StartWebhook(
	ctx,
	config,
); err != nil {
	log.Fatal(err)
}
```

همچنین API مربوط به Set/Get/Delete Webhook در package `methods` وجود دارد.

---

# 🔁 Retry و Rate Limit

برای محیط‌های Production، مدیریت خطاهای موقتی اهمیت زیادی دارد.

Package `retry` سیاست‌های مختلفی ارائه می‌کند:

```text
DefaultPolicy
AggressivePolicy
ConservativePolicy
```

و package `ratelimit` شامل:

```text
Limiter
ConcurrencyLimiter
Token Bucket
Burst
Wait
TryWait
WaitTimeout
```

است.

---

# 💾 Cache

Package `cache` شامل Cache درون‌حافظه‌ای و Sharded Cache است.

قابلیت‌ها:

```text
Get
Set
SetWithTTL
SetForever
Delete
Len
Keys
Clear
GetOrSet
GetOrSetWithTTL
Close
```

برای پروژه‌هایی که به دسترسی سریع به داده‌های موقت نیاز دارند، Cache می‌تواند فشار روی API و سایر منابع را کاهش دهد.

---

# 📝 Logging

Package `logger` چند حالت مختلف Logging ارائه می‌کند.

از جمله:

```text
New
NewJSON
NewAsync
Debug
Info
Warn
Error
Fatal
SetLevel
SetFormat
SetOutput
WithFields
```

همچنین قابلیت Async Logging و JSON Logging در API این package وجود دارد.

---

# ❌ Error Handling

Package `errors` برای مدیریت خطاهای کتابخانه طراحی شده است.

انواع اصلی:

```text
APIError
ValidationError
NetworkError
MultiError
```

همچنین قابلیت‌هایی مانند:

```text
IsRetryable
ParseError
AsAPIError
AsNetworkError
AsValidationError
```

در اختیار توسعه‌دهنده قرار دارد.

---

# 📚 مستندات کامل

مستندات پروژه به صورت بخش‌بندی‌شده داخل پوشه‌ی `README/` قرار گرفته‌اند.

برای مشاهده مستندات کامل:

### 📖 Documentation

[![Documentation](https://img.shields.io/badge/📚%20Open-Complete%20Documentation-6366F1?style=for-the-badge)](README/README.md)

### بخش‌های اصلی

* [01 — Getting Started](README/01_Getting_Started.md)
* [02 — Bot API and Architecture](README/02_Bot_API_and_Architecture.md)
* [03 — Messages and Media](README/03_Messages_and_Media.md)
* [04 — Updates, Polling and Webhook](README/04_Updates_Polling_Webhook.md)
* [05 — Chat Management](README/05_Chat_Management.md)
* [06 — Bot Profile and Commands](README/06_Bot_Profile_and_Commands.md)
* [07 — Keyboards and Buttons](README/07_Keyboards_and_Buttons.md)
* [08 — Models](README/08_Models.md)

### Package Reference

* [09 — Package: codemeet](README/09_Package_Codemeet.md)
* [10 — Package: api](README/10_Package_API.md)
* [11 — Package: methods](README/11_Package_Methods.md)
* [12 — Package: dispatcher](README/12_Package_Dispatcher.md)
* [13 — Package: middleware](README/13_Package_Middleware.md)
* [14 — Package: cache](README/14_Package_Cache.md)
* [15 — Package: ratelimit](README/15_Package_RateLimit.md)
* [16 — Package: retry](README/16_Package_Retry.md)
* [17 — Package: logger](README/17_Package_Logger.md)
* [18 — Package: errors](README/18_Package_Errors.md)
* [19 — Package: polling](README/19_Package_Polling.md)
* [20 — Package: webhook](README/20_Package_Webhook.md)

### Methods

* [21 — Bot Methods](README/21_Methods_Bot.md)
* [22 — Chat Methods](README/22_Methods_Chat.md)
* [23 — Message Methods](README/23_Methods_Messages.md)
* [24 — Media Methods](README/24_Methods_Media.md)
* [25 — Updates and Webhook Methods](README/25_Methods_Updates_Webhook.md)

### قابلیت‌های پیشرفته

* [26 — Reliability, Performance and Observability](README/26_Reliability_Performance.md)
* [27 — AntiLink](README/27_Contrib_AntiLink.md)
* [28 — AntiSpam](README/28_Contrib_AntiSpam.md)
* [29 — ForceJoin](README/29_Contrib_ForceJoin.md)
* [30 — Gatekeeper](README/30_Contrib_Gatekeeper.md)
* [31 — Maintenance Mode](README/31_Contrib_MaintenanceMode.md)
* [32 — Profanity Filter](README/32_Contrib_ProfanityFilter.md)
* [33 — VPN Detector](README/33_Contrib_VPNDetector.md)
* [34 — Warn System](README/34_Contrib_WarnSystem.md)

### توسعه و استفاده واقعی

* [35 — Examples and Production Patterns](README/35_Examples.md)
* [36 — Errors and Troubleshooting](README/36_Errors_and_Troubleshooting.md)
* [37 — Package Map and API Index](README/37_API_Index.md)

---

# ⚠️ وضعیت API و سازگاری نسخه‌ها

این مستندات بر اساس API و ساختار نسخه‌ی `1.0.0` کتابخانه تهیه شده‌اند.

برخی قابلیت‌ها و Methodهایی که در لایه‌ی مستندات Bot API تعریف یا برای توسعه‌های آینده در نظر گرفته شده‌اند، ممکن است در نسخه‌ی فعلی کتابخانه هنوز به صورت کامل در Bot API یا Implementation نهایی CodeMeet در دسترس نباشند.

بنابراین:

> **وجود نام یک Method در Documentation به معنی تضمین در دسترس بودن همان Method در تمامی نسخه‌های Bot API نیست.**

در صورت اضافه شدن یا تغییر API سمت CodeMeet، کتابخانه نیز می‌تواند در نسخه‌های بعدی به‌روزرسانی شود تا سازگاری کامل برقرار شود.

این موضوع به صورت عمدی در طراحی Documentation در نظر گرفته شده تا ساختار API آینده نیز قابل توسعه و نگهداری باشد.

---

# 📦 نسخه

```text
CodeMeet Go SDK
Version: 1.0.0
```

Repository:

```text
github.com/AbolfazlZarei-dev/codemeet-go
```

---

# 🤝 مشارکت

اگر در کتابخانه با Bug، مشکل Performance یا ناسازگاری با Bot API مواجه شدید، می‌توانید Issue ایجاد کنید یا Pull Request ارسال کنید.

پیشنهادهای مربوط به:

* Performance
* API Design
* Security
* Concurrency
* Documentation
* Developer Experience

نیز مورد استقبال هستند.

---

# 👨‍💻 توسعه‌دهنده

<div align="center">

## Abolfazl Zarei

Developer & Creator of CodeMeet Go SDK

<br>

[![GitHub](https://img.shields.io/badge/GitHub-AbolfazlZarei--dev-181717?style=for-the-badge\&logo=github\&logoColor=white)](https://github.com/AbolfazlZarei-dev)

<br><br>

**CodeMeet Go SDK**

Built with ❤️ and Go

<br>

© 2026 CodeMeet Go SDK

</div>

---

<div align="center">

### 🚀 Build Fast. Scale Further. Code with Go.

<br>

[![Get Started](https://img.shields.io/badge/🚀%20Get%20Started-Documentation-6366F1?style=for-the-badge)](README/README.md)
[![Go Reference](https://img.shields.io/badge/📦%20Go-Reference-00ADD8?style=for-the-badge\&logo=go\&logoColor=white)](https://pkg.go.dev/github.com/AbolfazlZarei-dev/codemeet-go)
[![GitHub](https://img.shields.io/badge/⭐%20Star-on%20GitHub-181717?style=for-the-badge\&logo=github\&logoColor=white)](https://github.com/AbolfazlZarei-dev/codemeet-go)

</div>