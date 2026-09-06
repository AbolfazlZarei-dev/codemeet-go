
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

[![Documentation](https://img.shields.io/badge/📚%20Documentation-Read%20the%20Docs-6366F1?style=for-the-badge)](README/00-index.md)
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

برای مشاهده فهرست مستندات:

### 📖 Documentation

[![Documentation](https://img.shields.io/badge/📚%20Open-Complete%20Documentation-6366F1?style=for-the-badge)](README/00-index.md)

### شروع کار و معماری

* [00 — Index](README/00-index.md)
* [01 — Introduction](README/01-introduction.md)
* [02 — Installation](README/02-installation.md)
* [03 — Quick Start](README/03-quick-start.md)
* [04 — BotFather](README/04-botfather.md)
* [05 — Updates](README/05-updates.md)
* [06 — Polling](README/06-polling.md)
* [07 — Webhook](README/07-webhook.md)
* [08 — Architecture](README/08-architecture.md)

### پکیج‌های پایه

* [09 — API Client](README/09-api-client.md)
* [10 — Dispatcher](README/10-dispatcher.md)
* [11 — Middleware](README/11-middleware.md)
* [12 — Cache](README/12-cache.md)
* [13 — Retry](README/13-retry.md)
* [14 — Rate Limit](README/14-ratelimit.md)
* [15 — Errors](README/15-errors.md)
* [16 — Logger](README/16-logger.md)
* [17 — Database](README/17-database.md)

### توسعه و رفرنس

* [18 — Examples](README/18-examples.md)
* [19 — Production](README/19-production.md)
* [20 — Troubleshooting](README/20-troubleshooting.md)
* [21 — Reference](README/21-reference.md)
* [22 — Changelog 1.1.0](README/22-changelog-1.1.0.md)
* [23 — Keyboards and Callbacks](README/23-keyboards-and-callbacks.md)
* [24 — Sending and Media](README/24-sending-and-media.md)
* [25 — Dashboard](README/25-dashboard.md)
* [26 — Advanced API](README/26-advanced-api.md)
* [27 — Exported API Index](README/27-exported-api-index.md)
* [28 — Endpoint Map](README/28-endpoint-map.md)
* [29 — Method Signatures](README/29-method-signatures.md)

### متدها (Methods)

* [Methods — Bot](README/methods-bot.md)
* [Methods — Chat](README/methods-chat.md)
* [Methods — Media](README/methods-media.md)
* [Methods — Messages](README/methods-messages.md)
* [Methods — Updates](README/methods-updates.md)
* [Methods — Webhook](README/methods-webhook.md)

### مدل‌ها (Models)

* [Models — Bot](README/models-bot.md)
* [Models — Chat](README/models-chat.md)
* [Models — Keyboard](README/models-keyboard.md)
* [Models — Message](README/models-message.md)
* [Models — Update](README/models-update.md)
* [Models — User](README/models-user.md)
* [Models — Webhook](README/models-webhook.md)

### قابلیت‌های پیشرفته (Contrib)

* [Feature — AI](README/feature-ai.md)
* [Feature — AntiLink](README/feature-antilink.md)
* [Feature — AntiSpam](README/feature-antispam.md)
* [Feature — Console & Broadcast](README/feature-console-broadcast.md)
* [Feature — ForceJoin](README/feature-forcejoin.md)
* [Feature — Gatekeeper](README/feature-gatekeeper.md)
* [Feature — Group Manager](README/feature-groupmanager.md)
* [Feature — Maintenance](README/feature-maintenance.md)
* [Feature — Profanity Filter](README/feature-profanity.md)
* [Feature — Reminder](README/feature-reminder.md)
* [Feature — Shop Engine](README/feature-shopengine.md)
* [Feature — VPN Detector](README/feature-vpndetector.md)
* [Feature — Warn System](README/feature-warnsystem.md)

### سایر

* [README — Contributing](README/README-CONTRIBUTING.md)

---

# ⚠️ وضعیت API و سازگاری نسخه‌ها

این مستندات بر اساس API و ساختار نسخه‌ی `1.1.0` کتابخانه تهیه شده‌اند.

برخی قابلیت‌ها و Methodهایی که در لایه‌ی مستندات Bot API تعریف یا برای توسعه‌های آینده در نظر گرفته شده‌اند، ممکن است در نسخه‌ی فعلی کتابخانه هنوز به صورت کامل در Bot API یا Implementation نهایی CodeMeet در دسترس نباشند.

بنابراین:

> **وجود نام یک Method در Documentation به معنی تضمین در دسترس بودن همان Method در تمامی نسخه‌های Bot API نیست.**

در صورت اضافه شدن یا تغییر API سمت CodeMeet، کتابخانه نیز می‌تواند در نسخه‌های بعدی به‌روزرسانی شود تا سازگاری کامل برقرار شود.

این موضوع به صورت عمدی در طراحی Documentation در نظر گرفته شده تا ساختار API آینده نیز قابل توسعه و نگهداری باشد.

---

# 📦 نسخه

```text
CodeMeet Go SDK
Version: 1.1.0
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

[![GitHub](https://img.shields.io/badge/GitHub-AbolfazlZarei--dev-181717?style=for-the-badge&logo=github&logoColor=white)](https://github.com/AbolfazlZarei-dev)

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

[![Get Started](https://img.shields.io/badge/🚀%20Get%20Started-Documentation-6366F1?style=for-the-badge)](README/00-index.md)
[![Go Reference](https://img.shields.io/badge/📦%20Go-Reference-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/AbolfazlZarei-dev/codemeet-go)
[![GitHub](https://img.shields.io/badge/⭐%20Star-on%20GitHub-181717?style=for-the-badge&logo=github&logoColor=white)](https://github.com/AbolfazlZarei-dev/codemeet-go)

</div>
