
# ۵. تنظیمات بات و دستورات

مدیریت پروفایل بات (نام، توضیحات، بیوگرافی) و منوی دستورات (Commands Menu) از طریق `bot.API().Bot()` انجام می‌شود.

## دریافت اطلاعات بات (GetMe)
توجه: این متد به‌صورت خودکار در سیستم کش (Cache) ذخیره می‌شود تا درخواست‌های اضافی به سرور ارسال نشود و سرعت بالایی داشته باشد.
```go
me, _ := bot.API().Bot().GetMe(ctx)
fmt.Println("Bot Name:", me.FullName())
fmt.Println("Bot Username:", me.Username)
```

## مدیریت نام و توضیحات
شما می‌توانید نام نمایشی ربات و متنی که کاربران قبل از زدن دکمه Start می‌بینند را تغییر دهید:

```go
// تنظیم نام نمایشی (حداکثر ۱۰۰ کاراکتر)
bot.API().Bot().SetName(ctx, "ربات تستی من")

// تنظیم توضیحات (متن قبل از زدن Start - حداکثر ۵۱۲ کاراکتر)
bot.API().Bot().SetDescription(ctx, "این ربات برای تست ساخته شده است.")

// تنظیم بیوگرافی (About - حداکثر ۱۲۰ کاراکتر)
bot.API().Bot().SetShortDescription(ctx, "تست‌کننده ربات‌های کدمیت")
```

## مدیریت دستورات (Commands)
شما می‌توانید دستوراتی که کاربران با زدن `/` در منوی کدمیت می‌بینند را تنظیم کنید. این بخش از زبان‌های مختلف (Language Code) پشتیبانی کامل می‌کند.

```go
cmds := []models.BotCommand{
    {Command: "start", Description: "شروع کار با ربات"},
    {Command: "help", Description: "دریافت راهنما"},
    {Command: "support", Description: "ارتباط با پشتیبانی"},
}

// تنظیم دستورات برای زبان فارسی
bot.API().Bot().SetCommands(ctx, cmds, "fa")

// دریافت دستورات ثبت شده برای زبان فارسی
savedCmds, _ := bot.API().Bot().GetCommands(ctx, "fa")

// حذف دستورات زبان فارسی
bot.API().Bot().DeleteCommands(ctx, "fa")

