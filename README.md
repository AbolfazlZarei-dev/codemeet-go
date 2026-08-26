
<div align="center">

<h1>🚀 CodeMeet Go SDK</h1>

**یک SDK قدرتمند، مدرن و کاملاً غیرهمزمان (Async) برای توسعه ربات‌های پلتفرم کدمیت (CodeMeet) به زبان Go.**

<br>

[![Go Reference](https://img.shields.io/badge/go-reference-00ADD8?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/AbolfazlZarei-dev/codemeet-go)
[![Release](https://img.shields.io/github/v/release/AbolfazlZarei-dev/codemeet-go?style=for-the-badge&logo=github&color=blue)]()
[![License](https://img.shields.io/github/license/AbolfazlZarei-dev/codemeet-go?style=for-the-badge&color=green)]()
[![CodeMeet](https://img.shields.io/badge/CodeMeet-Official-7B2FBE?style=for-the-badge&logo=telegram)](https://codemeet.chat)

</div>

---

## 📖 فهرست مطالب
- [معرفی کتابخانه](#-معرفی-کتابخانه)
- [✨ ویژگی‌های کلیدی](#-ویژگی‌های-کلیدی)
- [📦 نصب و راه‌اندازی](#-نصب-و-راه‌اندازی)
- [🚀 شروع سریع (Quick Start)](#-شروع-سریع-quick-start)
- [🛡 سیستم امنیتی (ضد اسپم و ضد لینک)](#-سیستم-امنیتی-ضد-اسپم-و-ضد-لینک)
- [🛠 ساختار پکیج‌ها](#-ساختار-پکیج‌ها)
- [📊 داشبورد مانیتورینگ](#-داشبورد-مانیتورینگ)
- [🤝 مشارکت و سازنده](#-مشارکت-و-سازنده)

---

## 📌 معرفی کتابخانه
این کتابخانه با تمرکز بر **عملکرد بالا (High Performance)**، ایمنی در کانکارنسی (Concurrency) و معماری ماژولار طراحی شده است. اگر به دنبال ساخت ربات‌های پرمصرف، سریع و مقیاس‌پذیر در پلتفرم کدمیت هستید، این SDK تمام نیازهای زیرساختی شما را پوشش می‌دهد. این کتابخانه با استفاده از تکنیک‌هایی مانند `sync.Pool` و `Sharded Maps` کمترین فشار را روی رم و CPU وارد می‌کند.

---

## ✨ ویژگی‌های کلیدی

| دسته‌بندی | قابلیت‌ها | توضیحات |
|-----------|-----------|---------|
| 🌐 **شبکه و ارتباطات** | Connection Pooling، HTTP/2، Compression | ارتباطات سریع و پایدار با سرور کدمیت با کمترین Latency. |
| 🛡 **پایداری و پوشش خطا** | Circuit Breaker، Automatic Retry | پیشگیری از قطعی سرور، تلاش مجدد هوشمند با Exponential Backoff و Jitter. |
| ⚡ **کنترل ترافیک** | Rate Limiter (Token Bucket) | جلوگیری دقیق از خطای `429 Too Many Requests` با صف‌بندی درخواست‌ها. |
| ⚙️ **پردازش همزمان** | Dispatcher & Worker Pool | مسیریاب مرکزی، پردازش همزمان آپدیت‌ها و پشتیبانی کامل از Middlewares. |
| 💾 **مدیریت حافظه** | Sharded In-Memory Cache | کش بسیار سریع با TTL و بخش‌بندی شده (Sharded) برای کاهش Lock Contention. |
| 🛡 **امنیت ربات** | Anti-Spam & Anti-Link Engine | تشخیص هوشمند فلود، کلمات ممنوعه، لینک‌های تبلیغاتی و بن خودکار کاربران. |
| 📊 **مانیتورینگ و لاگ** | Web Dashboard، Live Logs | رابط کاربری گرافیکی برای مشاهده آمار، سرعت و لاگ‌های رنگی زنده ربات. |

---

## 📦 نصب و راه‌اندازی
برای افزودن کتابخانه به پروژه Go خود، دستور زیر را در ترمینال اجرا کنید:

```bash
go get github.com/AbolfazlZarei-dev/codemeet-go
```

---

## 🚀 شروع سریع (Quick Start)

در کمتر از ۱ دقیقه اولین ربات خود را با Long Polling اجرا کنید:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/AbolfazlZarei-dev/codemeet-go/codemeet"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
	"github.com/AbolfazlZarei-dev/codemeet-go/polling"
)

func main() {
	token := "YOUR_BOT_TOKEN"
	
	// مدیریت سیگنال‌های سیستم برای توقف امن ربات
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// ساخت ربات با تنظیمات پیش‌فرض بهینه
	bot, err := codemeet.New(token)
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Close()

	// هندلر پیام‌های متنی
	bot.OnMessage(func(ctx context.Context, msg *models.Message) {
		fmt.Printf("📥 پیام از %s: %s\n", msg.From.FullName(), msg.Text)
		
		// ارسال پاسخ به کاربر
		bot.Reply(ctx, msg, "👋 سلام! پیام شما دریافت شد.")
	})

	// هندلر دستور /start
	bot.OnCommand("start", func(ctx context.Context, msg *models.Message) {
		bot.SendHTML(ctx, msg.Chat.ID, "<b>🤖 ربات با موفقیت استارت شد!</b>")
	})

	fmt.Println("🚀 ربات در حال اجراست... (Ctrl+C برای خروج)")
	
	// شروع دریافت رویدادها با Long Polling
	pollingCfg := polling.DefaultConfig()
	pollingCfg.Timeout = 20
	
	if err := bot.StartPolling(ctx, pollingCfg); err != nil {
		log.Fatal(err)
	}
}
```

---

## 🛡 سیستم امنیتی (ضد اسپم و ضد لینک)

این کتابخانه دارای پکیج‌های مستقل و قدرتمند برای امنیت ربات شماست که به صورت اختیاری ایمپورت می‌شوند:

### ۱. سیستم ضد اسپم (`contrib/antispam`)
تشخیص فلود، کلمات ممنوعه، محدودیت نرخ پیام کاربران، اخطار خودکار و در نهایت بن کردن موقت یا دائمی کاربر متخلف.

### ۲. سیستم ضد لینک (`contrib/antilink`)
جلوگیری از ارسال لینک‌های تبلیغاتی (لینک‌های ربات‌های دیگر، لینک‌های دعوت و...) با قابلیت تعریف لیست سفید (Whitelist) برای دامنه‌های مجاز. این سیستم می‌تواند پیام کاربر را به صورت خودکار حذف کند.

```go
import (
	"github.com/AbolfazlZarei-dev/codemeet-go/contrib/antispam"
	"github.com/AbolfazlZarei-dev/codemeet-go/contrib/antilink"
)

// راه‌اندازی ضد لینک
linkCfg := antilink.DefaultConfig()
linkCfg.AllowedDomains = []string{"codemeet.chat"} // لینک‌های مجاز
linkCfg.Action = func(ctx context.Context, userID, chatID string, messageID int, reason string) {
    bot.API().Messages().Delete(ctx, chatID, messageID) // حذف پیام
    bot.Send(ctx, chatID, "⚠️ ارسال لینک مجاز نیست!")
}
antiLinkEngine := antilink.New(linkCfg)

// اتصال به ربات
bot.Use(antiLinkEngine.Middleware())
```

---

## 🛠 ساختار پکیج‌ها
کتابخانه به صورت ماژولار طراحی شده تا توسعه‌دهنده تنها از بخش‌های مورد نیاز استفاده کند:

```text
codemeet-go/
├── api/           # کلاینت HTTP، مدیریت شبکه، Circuit Breaker و Stats
├── methods/       # پیاده‌سازی تمام متدهای ربات (Messages, Media, Chat, Webhook, Bot)
├── models/        # مدل‌های داده‌ای (User, Chat, Message, Keyboards)
├── dispatcher/    # مسیریاب مرکزی برای ارجاع آپدیت‌ها به هندلرها
├── middleware/    # میان‌افزارهای آماده (Recovery, Logging, RateLimit, Metrics)
├── contrib/       # پکیج‌های کمکی و مستقل (antispam, antilink)
├── ratelimit/     # محدودیت نرخ توکن‌محور (Token Bucket)
├── retry/         # سیاست‌های تلاش مجدد (Exponential Backoff)
├── cache/         # کش درون‌حافظه‌ای و Sharded Cache
├── polling/       # دریافت آپدیت با Long Polling
├── webhook/       # سرور وب‌هوک امن
├── logger/        # لاگر رنگی و سریع با پشتیبانی از JSON
└── errors/        # مدیریت خطاها و کدهای وضعیت HTTP
```

---

## 📊 داشبورد مانیتورینگ
با استفاده از تابع `StartDashboard`، یک سرور وب سبک روی پورت دلخواه شما اجرا می‌شود که آمار دقیق درخواست‌ها، وضعیت سیستم و لاگ‌های زنده و رنگی را در یک رابط کاربری گرافیکی مدرن نمایش می‌دهد.

```go
go func() {
    dashCtx, cancel := context.WithCancel(context.Background())
    defer cancel()
    bot.StartDashboard(dashCtx, ":9090")
}()
```
سپس به آدرس `http://localhost:9090` بروید.

---

## 🤝 مشارکت و سازنده
پروژه‌های Open Source با مشارکت شما بزرگ می‌شوند. اگر باگ پیدا کردید یا قابلیت جدیدی مد نظرتان است، خوشحال می‌شویم Pull Request بزنید.

<div align="center">

**Abolfazl Zarei**
<br>
<a href="https://github.com/AbolfazlZarei-dev"><img src="https://img.shields.io/badge/GitHub-AbolfazlZarei--dev-181717?style=flat-square&logo=github" alt="GitHub"></a>
<br>
© 2026 CodeMeet Go SDK. All rights reserved.

</div>
