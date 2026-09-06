# Dispatcher

Package: `dispatcher`

## Types

```go
type HandlerFunc func(context.Context, *models.Update)
type MiddlewareFunc func(HandlerFunc) HandlerFunc
```

## ساخت

```go
d := dispatcher.New(32)
```

`maxWorkers` ظرفیت workerها را تعیین می‌کند.

## Handlerها

```go
d.OnMessage(...)
d.OnCallback(...)
d.OnCommand("start", ...)
d.OnText("سلام", ...)
d.OnRegex(`^/id\s+(.+)$`, ...)
d.Fallback(...)
```

## Generic Handler

```go
d.Handle("message", func(ctx context.Context, u *models.Update) {
    // ...
})
```

## Middleware

```go
d.Use(mw1, mw2)
```

## Dispatch

```go
d.Dispatch(ctx, update)
```

## Stop

```go
d.Stop()
```

## Stats

```go
dispatched, dropped, panics := d.Stats()
```

## توصیه

Recovery را در ابتدای chain قرار دهید تا panic در Handler باعث از بین رفتن پردازش نشود.
