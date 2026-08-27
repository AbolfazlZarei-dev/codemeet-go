# 11 — Package `methods`

`methods` facade typed روی Bot API است.

## Methods

```go
New(client, retryPolicy, rateLimiter)
```

Accessorها:

```go
Messages()
Media()
Bot()
Chat()
Webhook()
Updates()
Client()
```

## دسته‌بندی

| API | Package type |
|---|---|
| Bot | `BotMethods` |
| Chat | `ChatMethods` |
| Messages | `Messages` |
| Media | `Media` |
| Updates | `UpdatesMethods` |
| Webhook | `WebhookMethods` |

تمام عملیات سطح methods از لایه‌ی API client استفاده می‌کنند و مسیر Retry/Rate Limit را از facade دریافت می‌کنند.
