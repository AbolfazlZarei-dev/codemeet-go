

# ۳. دریافت رویدادها (Polling & Webhook)

برای اینکه ربات شما از تعامل کاربران (ارسال پیام، کلیک روی دکمه‌ها و ...) مطلع شود، دو روش استاندارد وجود دارد: Long Polling و Webhook.

## روش Long Polling
این روش مناسب برای توسعه محلی (Local) و سرورهای بدون دامنه عمومی است. ربات به‌صورت دوره‌ای از سرور کدمیت درخواست رویدادهای جدید می‌کند.

```go
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
defer cancel()

// تنظیمات پیش‌فرض دارای محدودیت‌های محتاطانه برای جلوگیری از خطای 429 است
cfg := polling.DefaultConfig() 
cfg.Timeout = 10 // ثانیه

if err := bot.StartPolling(ctx, cfg); err != nil {
    log.Fatal(err)
}
```

## روش Webhook
برای محیط‌های پروداکشن با ترافیک بالا، استفاده از Webhook بهترین انتخاب است. سرور کدمیت به‌محض وقوع رویداد، آن را به سرور شما ارسال می‌کند. (نیازمند سرور با HTTPS است)

```go
whCfg := webhook.DefaultConfig()
whCfg.ListenAddr = ":8443"
whCfg.Path = "/webhook"
whCfg.SecretToken = "my_super_secret_token" // برای امنیت سرور شما الزامی است

// ثبت وب‌هوک در سرور کدمیت
bot.SetWebhook(ctx, "https://yourdomain.com:8443/webhook", whCfg.SecretToken)

// اجرای سرور وب‌هوک
if err := bot.StartWebhook(ctx, whCfg); err != nil {
    log.Fatal(err)
}
```

## روتر رویدادها (Dispatcher)
شما می‌توانید با استفاده از `Dispatcher` برای انواع مختلف رویدادها هندلر (مدیریت‌کننده) ثبت کنید:

```go
// هندلر پیام‌های متنی و رسانه‌ای
bot.OnMessage(func(ctx context.Context, msg *models.Message) {
    // کدهای شما
})

// هندلر دستورات (مثلاً /start)
bot.OnCommand("start", func(ctx context.Context, msg *models.Message) {
    // کدهای شما
})

// هندلر کلیک روی دکمه‌های شیشه‌ای (Callback Query)
bot.OnCallback(func(ctx context.Context, cq *models.CallbackQuery) {
    // کدهای شما
})

// هندلر تطبیق متن دقیق
bot.OnText("سلام", func(ctx context.Context, msg *models.Message) {
    // کدهای شما
})

// هندلر Regex
bot.OnRegex(`^/help (\d+)$`, func(ctx context.Context, msg *models.Message, matches []string) {
    // کدهای شما
})

// هندلر پیش‌فرض (Fallback) - وقتی هیچ هندلری با رویداد تطابق ندارد اجرا می‌شود
bot.Fallback(func(ctx context.Context, u *models.Update) {
    bot.Send(ctx, u.EffectiveChat().ID, "دستور نامشخص است!")
})
```
