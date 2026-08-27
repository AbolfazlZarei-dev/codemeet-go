# 02 — Bot API و معماری

## Endpoint

تمام درخواست‌های Bot API از الگوی زیر استفاده می‌کنند:

```text
https://botapi.codemeet.chat/bot{token}/{method}
```

روش‌های API در پروژه به گروه‌های منطقی تقسیم شده‌اند:

```text
Bot
Chat
Messages
Media
Updates
Webhook
```

## معماری

```text
Application
    │
    ▼
codemeet.Bot
    ├── Dispatcher
    │     ├── Handlers
    │     └── Middleware
    │
    ├── methods.Methods
    │     ├── BotMethods
    │     ├── ChatMethods
    │     ├── Messages
    │     ├── Media
    │     ├── UpdatesMethods
    │     └── WebhookMethods
    │
    └── api.Client
          ├── HTTP Transport
          ├── Retry-aware methods
          ├── Rate Limit
          ├── Circuit Breaker
          ├── Statistics
          └── Streaming uploads
```

## API Client

`api.Client` مسئول transport است. Client دارای HTTP Client، Transport، Token، Logger، Statistics و Circuit Breaker است.

## Circuit Breaker

Circuit Breaker سه حالت دارد:

- `closed`
- `open`
- `half-open`

در Half-Open مکانیزم Single Probe وجود دارد تا چند درخواست همزمان برای تست recovery ارسال نشوند.

## Statistics

API Statistics شامل:

- Requests
- SuccessCount
- ErrorCount
- BytesIn
- BytesOut
- AvgLatency

است.

## Request ID

```go
ctx = api.WithRequestID(ctx, "request-123")
```

Request ID از Context خوانده می‌شود و در درخواست HTTP قابل استفاده است.

## Streaming

Multipart upload با streaming و `io.Pipe` انجام می‌شود. این طراحی برای فایل‌های بزرگ مناسب‌تر از ساخت یک buffer کامل در حافظه است.

## Connection Pool

HTTP Transport برای connection reuse تنظیم شده و HTTP/2 را فعال می‌کند.

## Response

`api.Response` امکانات:

- `Decode`
- `AsBool`
- `ParametersAsRetryAfter`

را فراهم می‌کند.
