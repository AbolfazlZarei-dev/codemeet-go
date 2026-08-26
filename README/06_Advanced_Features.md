
# ۶. قابلیت‌های پیشرفته

این کتابخانه ابزارهای پیشرفته‌ای برای محیط‌های پروداکشن، ترافیک بالا و مانیتورینگ ارائه می‌دهد.

## داشبورد وب (Web Dashboard)
یک رابط کاربری گرافیکی (GUI) در مرورگر شما فراهم می‌کند تا آمار ربات، لاگ‌های زنده و وضعیت سرویس‌ها را به‌صورت لحظه‌ای مانیتور کنید.

```go
go func() {
    dashCtx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // داشبورد روی پورت 9090 اجرا می‌شود
    if err := bot.StartDashboard(dashCtx, ":9090"); err != nil {
        log.Printf("Dashboard error: %v", err)
    }
}()
```
سپس مرورگر را روی `http://localhost:9090` باز کنید.
**امکانات داشبورد:**
- مشاهده آمار درخواست‌ها (Total, Success, Error, Avg Latency, Bytes In/Out)
- مشاهده لاگ‌های زنده (Live Logs) با رنگ‌بندی استاندارد در یک ترمینال شبیه‌سازی شده
- دکمه توقف اسکرول لاگ‌ها (Pause/Resume) و دکمه کپی لاگ‌ها
- مشاهده وضعیت سرویس‌ها (System Status) و ویژگی‌های فعال ربات

## میان‌افزارها (Middlewares)
برای کنترل جریان آپدیت‌ها قبل از رسیدن به هندلر اصلی می‌توانید از Middlewares آماده استفاده کنید. متد `bot.Use` به شما اجازه می‌دهد چندین میدل‌ور را به سادگی متصل کنید:

```go
import (
    "time"
    "github.com/AbolfazlZarei-dev/codemeet-go/middleware"
)

// اتصال میدل‌ورها به صورت زنجیره‌ای
bot.Use(
    // ۱. جلوگیری از کرش در پانیک (Panic Recovery)
    middleware.Recovery(bot.Logger()),

    // ۲. لاگ‌گیری زمان پردازش هر آپدیت
    middleware.Logging(bot.Logger()),

    // ۳. محدودیت نرخ برای هر کاربر (مثلا نهایتاً ۵ پیام در ۱ دقیقه)
    middleware.RateLimit(5, time.Minute),

    // ۴. فیلتر کاربران (لیست سیاه)
    middleware.Blacklist(func(userID string) bool {
        return userID == "banned-user-uuid"
    }),
)
```

همچنین می‌توانید از پکیج‌های مستقل امنیتی مانند **ضد اسپم (Anti-Spam)** و **ضد لینک (Anti-Link)** که در مسیر `contrib/` قرار دارند استفاده کنید.

## مدیریت خطاها و Circuit Breaker
کتابخانه از الگوی Circuit Breaker در سطح شبکه استفاده می‌کند. اگر سرور کدمیت دچار مشکل شود (خطاهای 5xx یا 429)، پس از چندین خطای متوالی، ارسال درخواست‌ها موقتاً متوقف می‌شود تا سرور فشار کمتری داشته باشد و سپس در حالت Half-Open تست می‌شود.

شما می‌توانید خطاهای برگشتی از API را به دقت مدیریت کنید:

```go
msg, err := bot.Send(ctx, chatID, "test")
if err != nil {
    // بررسی نوع خطای رخ داده
    if apiErr, ok := errors.AsAPIError(err); ok {
        switch apiErr.Code {
        case errors.CodeTooManyRequests:
            log.Printf("خطای 429! باید %d ثانیه صبر کنید.", apiErr.RetryAfter)
        case errors.CodeForbidden:
            log.Println("کاربر ربات را بلاک کرده است یا ربات ادمین نیست.")
        case errors.CodeBadRequest:
            log.Println("درخواست نامعتبر:", apiErr.Description)
        }
    } else if netErr, ok := errors.AsNetworkError(err); ok {
        // خطاهای شبکه (قطعی اینترنت، timeout و غیره)
        log.Println("خطای شبکه:", netErr)
    } else {
        log.Println("خطای ناشناخته:", err)
    }
}
```

برای بررسی دستی وضعیت Circuit Breaker نیز می‌توانید به این شکل عمل کنید:

```go
// 0 = Closed, 1 = Open, 2 = Half-Open
if bot.API().Client().Breaker().State() == 1 {
    log.Println("هشدار: سرور در دسترس نیست و Circuit Breaker باز شده است!")
}
```

## لاگر پیشرفته و رنگی
کتابخانه دارای یک لاگر داخلی بسیار سبک و سریع است که در تمام سیستم‌عامل‌ها (ویندوز، لینوکس، مک) از کدهای ANSI برای رنگی کردن خروجی ترمینال استفاده می‌کند. شما می‌توانید سطح لاگ‌ها را کنترل کنید:

```go
// تغییر سطح لاگ به Debug (برای دیباگ)
bot.Logger().SetLevel(logger.LevelDebug)

// لاگ‌گیری با فیلدهای ساختاریافته
bot.Logger().Info("Server started", "port", 9090, "status", "running")
bot.Logger().Warn("Rate limit approaching", "user_id", "12345")
bot.Logger().Error("Failed to send message", "error", err)
```
اگر قصد دارید لاگ‌ها را در فایل ذخیره کنید یا فرمت آن‌ها را به JSON تغییر دهید:
```go
bot.Logger().SetFormat(logger.FormatJSON)
bot.Logger().SetOutput(myLogFile)
```
