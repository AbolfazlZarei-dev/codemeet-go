
```markdown
# ۶. قابلیت‌های پیشرفته

این کتابخانه ابزارهای پیشرفته‌ای برای محیط‌های پروداکشن، ترافیک بالا و مانیتورینگ ارائه می‌دهد.

## داشبورد وب (Web Dashboard)
یک رابط کاربری گرافیکی (GUI) در مرورگر شما فراهم می‌کند تا آمار ربات، لاگ‌های زنده و وضعیت سرویس‌ها را به‌صورت لحظه‌ای مانیتور کنید.

```go
go func() {
    if err := bot.StartDashboard(ctx, ":8080"); err != nil {
        log.Fatal(err)
    }
}()
```
سپس مرورگر را روی `http://localhost:8080` باز کنید.
**امکانات داشبورد:**
- مشاهده آمار درخواست‌ها (Total, Success, Error, Avg Latency, Bytes In/Out)
- مشاهده لاگ‌های زنده (Live Logs) با رنگ‌بندی استاندارد در یک ترمینال شبیه‌سازی شده
- مشاهده وضعیت سرویس‌ها (System Status) و ویژگی‌های فعال ربات

## میان‌افزارها (Middlewares)
برای کنترل جریان آپدیت‌ها قبل از رسیدن به هندلر اصلی می‌توانید از Middlewares استفاده کنید:

```go
// ۱. جلوگیری از کرش در پانیک (Panic Recovery)
bot.Use(func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
    return middleware.Recovery(bot.Logger())(func(ctx context.Context, u *models.Update) {
        next(ctx, u)
    })
})

// ۲. لاگ‌گیری زمان پردازش هر آپدیت
bot.Use(func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
    return middleware.Logging(bot.Logger())(func(ctx context.Context, u *models.Update) {
        next(ctx, u)
    })
})

// ۳. محدودیت نرخ برای هر کاربر (مثلا نهایتاً ۵ پیام در ۱ دقیقه)
bot.Use(func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
    return middleware.RateLimit(5, time.Minute)(func(ctx context.Context, u *models.Update) {
        next(ctx, u)
    })
})

// ۴. فیلتر کاربران (لیست سیاه)
bot.Use(func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
    return middleware.Blacklist(func(userID string) bool {
        return userID == "banned-user-uuid"
    })(func(ctx context.Context, u *models.Update) {
        next(ctx, u)
    })
})
```

## مدیریت خطاها و Circuit Breaker
کتابخانه از الگوی Circuit Breaker استفاده می‌کند. اگر سرور کدمیت دچار مشکل شود، پس از چندین خطای متوالی، ارسال درخواست‌ها موقتاً متوقف می‌شود تا سرور فشار کمتری داشته باشد و سپس در حالت Half-Open تست می‌شود.

```go
// بررسی وضعیت Circuit Breaker
if bot.API().Client().Breaker().State() == api.StateOpen {
    log.Println("هشدار: سرور در دسترس نیست و Circuit Breaker باز شده است!")
}

// مدیریت خطاهای API به‌صورت دستی
msg, err := bot.Send(ctx, chatID, "test")
if err != nil {
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
        log.Println("خطای شبکه:", netErr)
    }
}
```

## WebSocket Hub
برای ارتباطات Real-time و دریافت رویدادهای آنی (بدون نیاز به Polling یا Webhook عمومی):

```go
hub := bot.WS()
// اتصال به سرور وب‌سوکت کدمیت
err := hub.Connect(ctx, "wss://botapi.codemeet.chat/ws/your_endpoint")
if err != nil {
    log.Fatal(err)
}

// اشتراک در رویدادهای خاص (مثلاً پیام جدید)
eventChan := hub.Subscribe("message")

go func() {
    for event := range eventChan {
        log.Println("رویداد WS دریافت شد:", event.Type)
        // پردازش payload رویداد
    }
}()

// ارسال پیام از طریق وب‌سوکت
hub.Send(ctx, "wss://botapi.codemeet.chat/ws/your_endpoint", map[string]string{
    "type": "ping",
})

