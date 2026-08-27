# 13 — Package `middleware`

Middleware برای cross-cutting concerns طراحی شده است.

## Recovery

```go
middleware.Recovery(logger)
```

panic handler را recover می‌کند و stack trace و Update ID را log می‌کند.

## Logging

```go
middleware.Logging(logger)
```

شروع و پایان پردازش Update و duration را ثبت می‌کند.

## RateLimit

```go
middleware.RateLimit(perUser, window)
```

Rate limit را per-user اعمال می‌کند و از 64 shard برای state استفاده می‌کند.

## Metrics

```go
NewMetricsCounter()
Inc()
Snapshot()
Metrics()
```

## Timeout

```go
middleware.Timeout(duration)
```

## Access control

- `BotOnly`
- `UserOnly`
- `AdminOnly`

## Lists

- `Blacklist`
- `Whitelist`

Middlewareها `dispatcher.MiddlewareFunc` هستند و با `Bot.WithMiddleware` یا `Dispatcher.Use` ثبت می‌شوند.
