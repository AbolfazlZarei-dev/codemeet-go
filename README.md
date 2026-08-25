
<div align="center">

<h1>🚀 CodeMeet Go SDK</h1>

**یک SDK قدرتمند، مدرن و کاملاً غیرهمزمان (Async) برای توسعه ربات‌های پلتفرم کدمیت (CodeMeet) به زبان Go.**

<br>

[![Go Reference](https://img.shields.io/badge/go-reference-00ADD8?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/AbolfazlZarei-dev/codemeet-go)
[![Release](https://img.shields.io/github/v/release/AbolfazlZarei-dev/codemeet-go?style=for-the-badge&logo=github&color=blue)]()
[![License](https://img.shields.io/github/license/AbolfazlZarei-dev/codemeet-go?style=for-the-badge&color=green)]()
[![SLSA 3](https://img.shields.io/badge/SLSA-Level%203-brightgreen?style=for-the-badge&logo=security)]()
[![CodeMeet](https://img.shields.io/badge/CodeMeet-Official-7B2FBE?style=for-the-badge&logo=telegram)](https://codemeet.chat)

</div>

---

## 📖 فهرست مطالب
- [معرفی کتابخانه](#-معرفی-کتابخانه)
- [✨ ویژگی‌های کلیدی](#-ویژگی‌های-کلیدی)
- [📦 نصب و راه‌اندازی](#-نصب-و-راه‌اندازی)
- [🚀 شروع سریع (Quick Start)](#-شروع-سریع-quick-start)
- [🛠 ساختار پکیج‌ها](#-ساختار-پکیج‌ها)
- [📚 مستندات جامع](#-مستندات-جامع)
- [🔐 امنیت و SLSA](#-امنیت-و-slsa)
- [🤝 مشارکت و سازنده](#-مشارکت-و-سازنده)

---

## 📌 معرفی کتابخانه
این کتابخانه با تمرکز بر **عملکرد بالا (High Performance)**، ایمنی در کانکارنسی (Concurrency) و معماری ماژولار طراحی شده است. اگر به دنبال ساخت ربات‌های پرمصرف، سریع و مقیاس‌پذیر در پلتفرم کدمیت هستید، این SDK تمام نیازهای زیرساختی شما را پوشش می‌دهد.

---

## ✨ ویژگی‌های کلیدی

| دسته‌بندی | قابلیت‌ها | توضیحات |
|-----------|-----------|---------|
| 🌐 **شبکه و ارتباطات** | Connection Pooling، HTTP/2، Compression | ارتباطات سریع و پایدار با سرور کدمیت با کمترین Latency. |
| 🛡 **پایداری و پوشش خطا** | Circuit Breaker، Automatic Retry | پیشگیری از قطعی سرور، تلاش مجدد هوشمند با Exponential Backoff و Jitter. |
| ⚡ **کنترل ترافیک** | Rate Limiter (Token Bucket) | جلوگیری دقیق از خطای `429 Too Many Requests` با صف‌بندی درخواست‌ها. |
| ⚙️ **پردازش همزمان** | Dispatcher & Worker Pool | مسیریاب مرکزی، پردازش همزمان آپدیت‌ها و پشتیبانی کامل از Middlewares. |
| 📊 **مانیتورینگ و لاگ** | Web Dashboard، Live Logs | رابط کاربری گرافیکی برای مشاهده آمار، سرعت و لاگ‌های زنده ربات. |
| 💾 **مدیریت حافظه** | Sharded In-Memory Cache | کش بسیار سریع با TTL و بخش‌بندی شده (Sharded) برای کاهش Lock Contention. |
| 🔌 **ارتباط بلادرنگ** | WebSocket Hub | مدیریت اتصالات Real-time بدون نیاز به Polling مداوم. |

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
	"os"
	"os/signal"

	"github.com/AbolfazlZarei-dev/codemeet-go"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

func main() {
	token := "YOUR_BOT_TOKEN"
	
	// مدیریت سیگنال‌های سیستم برای توقف امن ربات
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// ساخت ربات با تنظیمات پیش‌فرض بهینه
	bot, err := codemeet.New(token)
	if err != nil {
		log.Fatal(err)
	}

	// هندلر پیام‌های متنی
	bot.OnMessage(func(ctx context.Context, msg *models.Message) {
		fmt.Printf("📥 پیام از %s: %s\n", msg.From.FullName(), msg.Text)
		
		// ارسال پاسخ به کاربر
		_, _ = bot.Send(ctx, msg.Chat.ID, "👋 سلام! پیام شما دریافت شد.")
	})

	// هندلر دستور /start
	bot.OnCommand("start", func(ctx context.Context, msg *models.Message) {
		_, _ = bot.SendHTML(ctx, msg.Chat.ID, "<b>🤖 ربات با موفقیت استارت شد!</b>")
	})

	fmt.Println("🚀 ربات در حال اجراست... (Ctrl+C برای خروج)")
	
	// شروع دریافت رویدادها
	if err := bot.StartPolling(ctx, codemeet.DefaultConfig()); err != nil {
		log.Fatal(err)
	}
}
```

---

## 🛠 ساختار پکیج‌ها
کتابخانه به صورت ماژولار طراحی شده تا توسعه‌دهنده تنها از بخش‌های مورد نیاز استفاده کند:

```text
codemeet-go/
├── api/          # کلاینت HTTP، مدیریت شبکه، Circuit Breaker و Stats
├── methods/      # پیاده‌سازی تمام متدهای ربات (Messages, Media, Chat, Webhook, Bot)
├── models/       # مدل‌های داده‌ای (User, Chat, Message, Keyboards)
├── dispatcher/   # مسیریاب مرکزی برای ارجاع آپدیت‌ها به هندلرها
├── middleware/   # میان‌افزارهای آماده (Recovery, Logging, RateLimit, Metrics)
├── ratelimit/    # محدودیت نرخ توکن‌محور
├── retry/        # سیاست‌های تلاش مجدد
├── cache/        # کش درون‌حافظه‌ای و Sharded Cache
├── polling/      # دریافت آپدیت با Long Polling
├── webhook/      # سرور وب‌هوک امن
└── ws/           # مدیریت اتصالات WebSocket
```

---

## 📚 مستندات جامع
برای راهنمای دقیق و کاربردی هر بخش، روی دکمه‌های زیر کلیک کنید:

<table>
  <tr>
    <td align="center"><a href="README/01_Getting_Started.md">1️⃣ شروع به کار و پیکربندی</a></td>
    <td align="center"><a href="README/02_Messages_and_Media.md">2️⃣ مدیریت پیام‌ها و رسانه‌ها</a></td>
    <td align="center"><a href="README/03_Updates_Polling_Webhook.md">3️⃣ دریافت رویدادها</a></td>
  </tr>
  <tr>
    <td align="center"><a href="README/04_Chat_Management.md">4️⃣ مدیریت گروه‌ها و کانال‌ها</a></td>
    <td align="center"><a href="README/05_Bot_Profile_and_Commands.md">5️⃣ تنظیمات بات و دستورات</a></td>
    <td align="center"><a href="README/06_Advanced_Features.md">6️⃣ قابلیت‌های پیشرفته</a></td>
  </tr>
</table>

---

## 🔐 امنیت و SLSA
این پروژه با استفاده از استانداردهای **SLSA (Supply-chain Levels for Software Artifacts)** توسعه و انتشار می‌یابد. فایل‌های باینری و خروجی‌های هر Release به‌صورت خودکار توسط GitHub Actions تولید شده‌اند و دارای گواهی Provenance برای تایید اصالت و امنیت زنجیره تامین (Supply Chain) هستند.

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
