# معماری

## اجزای اصلی

```text
codemeet.Bot
 ├── api.Client
 ├── methods.Methods
 ├── dispatcher.Dispatcher
 ├── polling.Poller
 ├── webhook.Server
 ├── retry.Policy
 ├── ratelimit.Limiter
 ├── cache.Cache
 └── logger.Logger
```

## Request Pipeline

```text
Handler
  ↓
Methods
  ↓
Retry / RateLimit
  ↓
api.Client
  ↓
CircuitBreaker
  ↓
HTTP
  ↓
Response Decode
  ↓
APIError یا Result
```

## Concurrency

در بخش‌های مختلف از `atomic`, `sync.RWMutex`, `sync.Map`, `singleflight` و workerها استفاده شده است.

## Lifecycle

```text
New
 ↓
configure
 ↓
register handlers
 ↓
StartPolling / StartWebhook
 ↓
Stop/Context cancel
 ↓
Close resources
```

## Separation of concerns

- `models`: داده
- `methods`: API wrapper
- `dispatcher`: routing
- `polling/webhook`: ingress
- `middleware`: policy
- `api`: transport
- `errors`: error model
- `logger`: observability
- `cache/database`: state
