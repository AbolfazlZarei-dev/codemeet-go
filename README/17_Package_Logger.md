
# 17 — Package `logger`

Logger داخلی برای logging ساخت‌یافته طراحی شده است. این پکیج به صورت کامل Cross-Platform است و بدون هیچ خطایی روی سیستم‌عامل‌های ویندوز، لینوکس و مک کامپایل و اجرا می‌شود.

## Constructors

```go
New(level)
NewJSON(level)
NewAsync(level)
```

## Configuration

```go
SetIncludeCaller()
SetEnabled()
IsEnabled()
SetOutput()
SetLevel()
SetFormat()
WithFields()
```

## Levels

- Debug
- Info
- Warn
- Error
- Fatal

## Logging

```go
log.Debug("debug")
log.Info("started")
log.Warn("warning")
log.Error("failed", "error", err)
log.Fatal("fatal")
```

## Output

```go
log.Output()
```

Logger از text و JSON format پشتیبانی می‌کند.

## Cross-Platform Support (Windows, Linux, Mac)

برای جلوگیری از خطاهای کامپایل کراس (Cross-Compilation) و پشتیبانی صحیح از رنگ‌ها در ترمینال‌های مختلف، تابع `enableWindowsANSI` با استفاده از **Build Tags** در فایل‌های مجزا تفکیک شده است:

- **ویندوز (`logger_windows.go`):** از `syscall` مربوط به `kernel32.dll` برای فعال‌سازی ANSI Colors در کنسول ویندوز استفاده می‌کند.
- **لینوکس و مک (`logger_unix.go`):** به صورت خودکار از رنگ‌ها پشتیبانی می‌کند و تابع مربوطه به صورت خالی (No-op) قرار داده شده است تا در زمان کامپایل هیچ خطایی رخ ندهد.

## Async

`NewAsync` برای کاهش هزینه‌ی logging در مسیرهای پرترافیک مناسب است.

## Sync و Close

```go
log.Sync()
log.Close()
```

## خاموش کردن

```go
codemeet.WithoutLogger()
```

در این حالت output به `io.Discard` هدایت می‌شود.
