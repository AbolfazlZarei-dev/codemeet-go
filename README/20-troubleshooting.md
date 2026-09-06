# Troubleshooting

## 400

ورودی، JSON، پارامتر یا markup را بررسی کنید.

## 401

Token را بررسی کنید.

## 403

موارد رایج:

```text
the user must start the bot first
bot must be an administrator
```

کاربر باید Bot را Start کند یا Bot باید Permission مناسب داشته باشد.

## 404

نام Method/URL یا وجود پیام را بررسی کنید؛ همچنین احتمال Endpoint غیرفعال را در نظر بگیرید.

## 409

Conflict بین Polling و Webhook یکی از سناریوهای رایج است.

```go
_ = bot.DeleteWebhook(ctx)
```

## 413

Payload بیش از حد مجاز است. در Webhook `MaxBodySize` را نیز بررسی کنید.

## 429

`retry_after` را رعایت کنید.

## 5xx

از Retry Policy استفاده کنید.

## Markdown/HTML

خطای `can't parse entities` معمولاً به syntax نادرست markup مربوط است.

## Webhook

```text
1. /health
2. HTTPS
3. مسیر درست
4. Secret Header
5. JSON body
6. Reverse Proxy
7. logs
```

## اصل Endpoint

اگر یک Method وجود دارد ولی Server `404 Not Found` می‌دهد، صرفاً از روی وجود Method نتیجه نگیرید که SDK خراب است؛ availability Endpoint را روی Bot API مقصد بررسی کنید.
