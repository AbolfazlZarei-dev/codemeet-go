# 09 — Package `codemeet`

Package اصلی facade سطح بالا است.

## ثابت‌ها

```go
Version       = "1.0.0"
Author        = "Abolfazl Zarei"
GitHubProfile = "github.com/AbolfazlZarei-dev"
GitHubRepo    = "github.com/AbolfazlZarei-dev/codemeet-go"
```

## Bot

Bot اجزای اصلی runtime را در خود نگه می‌دارد:

- API client
- Dispatcher
- Rate Limiter
- Retry Policy
- Cache
- Logger
- Methods
- Bot statistics
- Dashboard state

## Constructors

```go
New(token string, opts ...Option) (*Bot, error)
```

## Options

```go
WithBaseURL
WithHTTPClient
WithTimeout
WithRetry
WithRateLimit
WithRateLimitBurst
WithCache
WithShardedCache
WithLogger
WithDashboardAuth
WithMiddleware
WithoutLogger
```

## Accessors

```go
API()
Dispatcher()
Cache()
Logger()
RateLimiter()
RetryPolicy()
Token()
BaseURL()
RunMode()
Stats()
Uptime()
```

## Runtime

```go
StartPolling(ctx, cfg)
StartWebhook(ctx, cfg)
StartDashboard(ctx, addr)
```

## API shortcuts

```go
SetWebhook()
DeleteWebhook()
GetMe()
ResetMe()
```

## Handlers

```go
OnCommand()
OnMessage()
OnCallback()
OnText()
OnRegex()
Fallback()
Use()
```

## Sending

```go
Send()
SendHTML()
SendWithKeyboard()
Reply()
AnswerCallback()
```

## Lifecycle

```go
Close()
HealthCheck()
```

`Close` shutdown را متمرکز می‌کند و منابع runtime را آزاد می‌کند.
