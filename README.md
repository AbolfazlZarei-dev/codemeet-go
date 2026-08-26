بسیار عالی! با اضافه شدن پکیج‌های قدرتمندی مثل کپچا (`gatekeeper`)، تشخیص VPN (`vpndetector`)، سیستم اخطار (`warnsystem`)، صفحه‌بندی (`pagination`)، فیلتر کلمات (`profanityfilter`)، حالت تعمیرات (`maintenancemode`) و عضویت اجباری (`forcejoin`)، فایل README شما باید کاملاً به‌روزرسانی شود تا این قابلیت‌های جدید و سطح سازمانی (Enterprise) به نمایش گذاشته شوند.

این فایل `README.md` نهایی و فوق‌العاده حرفه‌ای است. آن را کپی کرده و جایگزین محتوای قبلی کنید:

```markdown
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
- [🛡 پکیج‌های کاربردی و امنیتی (Contrib)](#-پکیج‌های-کاربردی-و-امنیتی-contrib)
- [🛠 ساختار پکیج‌ها](#-ساختار-پکیج‌ها)
- [📊 داشبورد مانیتورینگ](#-داشبورد-مانیتورینگ)
- [🤝 مشارکت و سازنده](#-مشارکت-و-سازنده)

---

## 📌 معرفی کتابخانه
این کتابخانه با تمرکز بر **عملکرد بالا (High Performance)**، ایمنی در کانکارنسی (Concurrency) و معماری ماژولار طراحی شده است. اگر به دنبال ساخت ربات‌های پرمصرف، سریع و مقیاس‌پذیر در پلتفرم کدمیت هستید، این SDK تمام نیازهای زیرساختی شما را پوشش می‌دهد. 

این کتابخانه با استفاده از تکنیک‌هایی مانند `sync.Pool`، `Sharded Maps` و `Worker Pools` کمترین فشار را روی رم و CPU وارد می‌کند و به عنوان یک ابزار Enterprise-Grade شناخته می‌شود.

---

## ✨ ویژگی‌های کلیدی

| دسته‌بندی | قابلیت‌ها | توضیحات |
|-----------|-----------|---------|
| 🌐 **شبکه و ارتباطات** | Connection Pooling، HTTP/2، Compression | ارتباطات سریع و پایدار با سرور کدمیت با کمترین Latency و مکانیزم `Circuit Breaker`. |
| ⚡ **کنترل ترافیک** | Rate Limiter (Token Bucket) | جلوگیری دقیق از خطای `429` با صف‌بندی درخواست‌ها و کنترل همزمانی (Concurrency). |
| ⚙️ **پردازش همزمان** | Dispatcher & Worker Pool | مسیریاب مرکزی، پردازش همزمان آپدیت‌ها با صف باند شده (Bounded Queue) برای جلوگیری از Memory Leak. |
| 💾 **مدیریت حافظه** | Sharded In-Memory Cache | کش بسیار سریع با TTL، بخش‌بندی شده (Sharded) و Singleflight برای جلوگیری از Cache Stampede. |
| 🛡 **موتور امنیتی** | Anti-Spam, Anti-Link, Profanity Filter | تشخیص هوشمند فلود، لینک‌های تبلیغاتی، هایپرلینک‌های مخفی، کلمات رکیک (با پشتیبانی Leetspeak) و بن خودکار. |
| 🚪 **مدیریت گروه** | Gatekeeper, ForceJoin, WarnSystem | کپچای ریاضی، عضویت اجباری در کانال، سیستم اخطار و اخراج خودکار با معماری Lock-Free. |
| 🔍 **تشخیص تهدید** | VPN & Proxy Detector | اسکن هوشمند فایل‌ها، APKها و کانفیگ‌ها (V2ray, Clash, etc.) برای جلوگیری از دور زدن تحریم‌ها. |
| 📄 **رابط کاربری** | Pagination, Maintenance Mode | صفحه‌بندی خودکار دکمه‌های شیشه‌ای و فعال کردن حالت تعمیرات ربات در زمان آپدیت سرور. |
| 📊 **مانیتورینگ و لاگ** | Web Dashboard، Live Logs | رابط گرافیکی برای مشاهده آمار، Stop/Start کردن لاگ‌ها و مانیتورینگ زنده سیستم. |

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
	
	// مدیریت سیگنال‌های سیستم برای توقف امن ربات (Graceful Shutdown)
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
		bot.Reply(ctx, msg, "👋 سلام! پیام شما دریافت شد.")
	})

	// هندلر دستور /start
	bot.OnCommand("start", func(ctx context.Context, msg *models.Message) {
		bot.SendHTML(ctx, msg.Chat.ID, "<b>🤖 ربات با موفقیت استارت شد!</b>")
	})

	fmt.Println("🚀 ربات در حال اجراست... (Ctrl+C برای خروج)")
	
	pollingCfg := polling.DefaultConfig()
	pollingCfg.Timeout = 20
	
	if err := bot.StartPolling(ctx, pollingCfg); err != nil {
		log.Fatal(err)
	}
}
```

---

## 🛡 پکیج‌های کاربردی و امنیتی (Contrib)

این کتابخانه دارای پکیج‌های مستقل و قدرتمندی برای مدیریت گروه‌هاست که به صورت اختیاری ایمپورت می‌شوند:

### ۱. سیستم دربان گروه و کپچا (`contrib/gatekeeper`)
سیستم کپچای بومی و فوق‌سریع با استفاده از دکمه‌های شیشه‌ای (Inline Keyboards) که از ورود ربات‌های اسپمر به گروه‌ها جلوگیری می‌کند. نیازی به API خارجی ندارد.

### ۲. سیستم تشخیص VPN و کانفیگ (`contrib/vpndetector`)
موتور تشخیص دور زدن تحریم‌ها. این سیستم فایل‌های APK و فایل‌های متنی (JSON, YAML) را دانلود کرده و محتوای داخل آن‌ها را برای پیدا کردن کانفیگ‌های V2ray, Clash, Trojan و... با دقت بسیار بالا اسکن می‌کند.

### ۳. سیستم اخطار و اخراج (`contrib/warnsystem`)
به ادمین‌ها اجازه می‌دهد با دستور `/warn` به کاربران اخطار دهند. پس از رسیدن به سقف اخطارها، کاربر به صورت خودکار اخراج می‌شود. با استفاده از `atomic` و معماری Sharded، کمترین فشار را به CPU در ترافیک بالا وارد می‌کند.

### ۴. سیستم عضویت اجباری (`contrib/forcejoin`)
کاربر را مجبور می‌کند قبل از استفاده از ربات، در کانال شما عضو شود. وضعیت عضویت کاربران در حافظه کش می‌شود تا درخواست‌های مداوم به API کدمیت ارسال نشود.

### ۵. سیستم صفحه‌بندی (`contrib/pagination`)
ابزاری حرفه‌ای برای ساخت لیست‌های طولانی (مثل لیست کاربران، محصولات و...) با دکمه‌های «صفحه بعد» و «صفحه قبل». این پکیج به صورت خودکار پیام را `Edit` می‌کند.

### ۶. سایر پکیج‌های مفید
- **`contrib/antilink`**: مسدودسازی لینک‌های تبلیغاتی و هایپرلینک‌های مخفی.
- **`contrib/antispam`**: تشخیص فلود و کلمات ممنوعه.
- **`contrib/profanityfilter`**: فیلتر کلمات رکیک با تشخیص Leetspeak (مثل `f@ck` یا `1d10t`).
- **`contrib/maintenancemode`**: فعال کردن حالت "ربات در حال به‌روزرسانی" برای کاربران عادی.

---

## 🛠 ساختار پکیج‌ها
کتابخانه به صورت ماژولار طراحی شده تا توسعه‌دهنده تنها از بخش‌های مورد نیاز استفاده کند:

```text
codemeet-go/
├── api/           # کلاینت HTTP، مدیریت شبکه، Circuit Breaker و Stats
├── methods/       # پیاده‌سازی تمام متدهای ربات (Messages, Media, Chat, Webhook, Bot)
├── models/        # مدل‌های داده‌ای (User, Chat, Message, Keyboards)
├── dispatcher/    # مسیریاب مرکزی با Bounded Queue برای ارجاع آپدیت‌ها
├── middleware/    # میان‌افزارهای آماده (Recovery, Logging, RateLimit, Metrics)
├── contrib/       # پکیج‌های مستقل امنیتی و کاربردی (antispam, vpndetector, gatekeeper, etc.)
├── ratelimit/     # محدودیت نرخ توکن‌محور (Token Bucket) و Semaphore
├── retry/         # سیاست‌های تلاش مجدد (Exponential Backoff)
├── cache/         # کش درون‌حافظه‌ای، Sharded Cache و Singleflight
├── polling/       # دریافت آپدیت با Long Polling
├── webhook/       # سرور وب‌هوک امن
├── logger/        # لاگر رنگی، سریع و Async با پشتیبانی از JSON
└── errors/        # مدیریت خطاها و کدهای وضعیت HTTP
```

---

## 📊 داشبورد مانیتورینگ
با استفاده از تابع `StartDashboard`، یک سرور وب سبک روی پورت دلخواه شما اجرا می‌شود که آمار دقیق درخواست‌ها، وضعیت سیستم و لاگ‌های زنده و رنگی را در یک رابط کاربری گرافیکی مدرن نمایش می‌دهد. همچنین قابلیت Stop/Start کردن لاگ‌ها از طریق داشبورد فراهم است.

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
```
