

# 🚀 CodeMeet Go SDK (codemeet-go)

یک SDK قدرتمند، مدرن و کاملاً غیرهمزمان (Async) برای توسعه ربات‌های پلتفرم **کدمیت (CodeMeet)** به زبان Go. این کتابخانه با تمرکز بر عملکرد بالا (High Performance)، ایمنی در کانکارنسی و معماری ماژولار طراحی شده است.

## ✨ ویژگی‌های کلیدی
- **معماری شبکه بهینه:** Connection Pooling، HTTP/2، Circuit Breaker و Compression.
- **پیشگیری از قطعی:** سیستم Retry خودکار با Exponential Backoff و Jitter.
- **کنترل ترافیک:** Rate Limiter مبتنی بر Token Bucket برای جلوگیری از خطای 429.
- **پردازش همزمان:** Dispatcher با Worker Pool و پشتیبانی کامل از Middlewares.
- **داشبورد وب لحظه‌ای:** مانیتورینگ آمار درخواست‌ها، لاگ‌های زنده و وضعیت سرویس‌ها.
- **کش حافظه (In-Memory Cache):** ذخیره‌سازی سریع با TTL و پشتیبانی از Sharded Cache.
- **پشتیبانی از WebSocket:** مدیریت اتصالات بلادرنگ (Real-time).
- **پوشش کامل API:** پیاده‌سازی تمام متدهای رسمی کدمیت به علاوه متدهای پیشرفته مدیریت گروه‌ها.

## 📦 نصب
```bash
go get github.com/AbolfazlZarei-dev/codemeet-go
```

## 🚀 شروع سریع
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bot, err := codemeet.New(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.OnMessage(func(ctx context.Context, msg *models.Message) {
		fmt.Printf("پیام از %s: %s\n", msg.From.FullName(), msg.Text)
		_, _ = bot.Send(ctx, msg.Chat.ID, "شما گفتید: "+msg.Text)
	})

	bot.StartPolling(ctx, codemeet.DefaultConfig())
}
```

## 📚 مستندات
برای مشاهده راهنمای جامع هر بخش، به فایل‌های زیر مراجعه کنید:

- [۱. شروع به کار و پیکربندی](README/01_Getting_Started.md)
- [۲. مدیریت پیام‌ها و رسانه‌ها](README/02_Messages_and_Media.md)
- [۳. دریافت رویدادها (Polling & Webhook)](README/03_Updates_Polling_Webhook.md)
- [۴. مدیریت گروه‌ها و کانال‌ها](README/04_Chat_Management.md)
- [۵. تنظیمات بات و دستورات](README/05_Bot_Profile_and_Commands.md)
- [۶. قابلیت‌های پیشرفته (کش، داشبورد، Middleware)](README/06_Advanced_Features.md)

## 👤 سازنده
- **نام:** Abolfazl Zarei
- **گیت‌هاب:** [github.com/AbolfazlZarei-dev](https://github.com/AbolfazlZarei-dev)

