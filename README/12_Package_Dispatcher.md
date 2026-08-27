# 12 — Package `dispatcher`

Dispatcher قلب event routing کتابخانه است.

## Dispatcher

```go
New(workerCount int)
```

Bot نسخه‌ی `1.0.0` را با Worker Pool داخلی می‌سازد.

## Handler registration

```go
Handle(typeName, handler)
OnMessage(handler)
OnCallback(handler)
OnCommand(command, handler)
OnText(text, handler)
OnRegex(pattern, handler)
Fallback(handler)
```

## Middleware

```go
Use(middlewares...)
```

Middlewareها به زنجیره‌ی handler اضافه می‌شوند.

## Dispatch

```go
Dispatch(ctx, update)
```

Update وارد queue محدودشده می‌شود و workerها آن را پردازش می‌کنند.

## Backpressure

Dispatcher دارای bounded task queue است تا رشد نامحدود مصرف حافظه در شرایط فشار ایجاد نشود.

## Statistics

```go
Stats()
```

آمار شامل:

- total dispatched
- total dropped
- total panics

است.

## Stop

```go
Stop()
```

Workerها را متوقف می‌کند و queue runtime را می‌بندد.
