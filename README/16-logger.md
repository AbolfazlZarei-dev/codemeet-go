# Logger

Package: `logger`

## Levels

```go
logger.LevelDebug
logger.LevelInfo
logger.LevelWarn
logger.LevelError
logger.LevelFatal
```

## ساخت

```go
log := logger.New(logger.LevelInfo)
defer log.Close()
```

## JSON

```go
log := logger.NewJSON(logger.LevelInfo)
```

## Async

```go
log := logger.NewAsync(...)
```

## تنظیمات

```go
log.SetLevel(logger.LevelDebug)
log.SetEnabled(true)
log.SetIncludeCaller(true)
log.SetOutput(os.Stdout)
log.SetFormat(logger.FormatText)
```

## Fields

```go
l := log.WithFields(map[string]interface{}{
    "chat_id": chatID,
    "component": "bot",
})
```

## Log

```go
log.Debug("debug")
log.Info("started")
log.Warn("warning")
log.Error("failed", "error", err)
log.Fatal("fatal")
```

## Sync

```go
_ = log.Sync()
```

برای Loggerهای buffered مناسب است.
