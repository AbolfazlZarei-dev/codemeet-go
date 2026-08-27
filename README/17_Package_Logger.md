# 17 — Package `logger`

Logger داخلی برای logging ساخت‌یافته طراحی شده است.

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
