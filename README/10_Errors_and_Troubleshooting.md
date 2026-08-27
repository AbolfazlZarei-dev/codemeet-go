# خطاها و عیب‌یابی

## ساختار خطای API

Bot API پاسخ استانداردی شبیه زیر دارد:

```json
{
  "ok": false,
  "error_code": 400,
  "description": "Bad Request: chat not found"
}
```

`api.Response` فیلدهای زیر را نگهداری می‌کند:

- `Ok`
- `Result`
- `ErrorCode`
- `Description`
- `Parameters`
- `HTTPStatus`

## Decode

اگر `Result` وجود داشته باشد:

```go
var user models.User

if err := response.Decode(&user); err != nil {
    log.Println(err)
}
```

برای Resultهای Boolean:

```go
ok, err := response.AsBool()
```

## Retry After

در صورت وجود `parameters.retry_after`:

```go
seconds := response.ParametersAsRetryAfter()
```

## خطاهای رایج

### 400 Bad Request

علت‌های رایج:

- JSON نامعتبر
- chat_id نامعتبر
- پارامتر ناقص
- دکمه با action نامعتبر
- قالب‌بندی متن اشتباه

راه‌حل: request و typeهای `models` را بررسی کنید.

### 403

برای خطاهایی مانند نیاز به Start کردن Bot یا نیاز به Administrator بودن، ابتدا وضعیت دسترسی Bot را بررسی کنید.

### 429 Too Many Requests

Rate Limit فعال شده است.

اقدام:

1. Retry-After را رعایت کنید.
2. Rate Limiter کتابخانه را فعال/تنظیم کنید.
3. تعداد درخواست‌های همزمان را کاهش دهید.

Polling نیز برای 429 مقدار `retry_after` را در نظر می‌گیرد و قبل از ادامه صبر می‌کند.

### 500 و خطاهای شبکه

API Client برای وضعیت‌های 5xx و 429 Circuit Breaker را درگیر می‌کند. لایه‌ی Retry نیز می‌تواند درخواست‌های قابل retry را دوباره اجرا کند.

## Panic در Handler

از Middleware زیر استفاده کنید:

```go
bot, err := codemeet.New(
    token,
    codemeet.WithMiddleware(
        middleware.Recovery(bot.Logger()),
    ),
)
```

اگر Middleware به Bot قبل از ساخت نیاز داشته باشد، می‌توانید آن را از طریق Dispatcher بعد از ساخت Bot ثبت کنید.

## Logging

Logger متدهای:

- `Debug`
- `Info`
- `Warn`
- `Error`
- `Fatal`

دارد.

برای غیرفعال کردن کامل خروجی logger:

```go
codemeet.WithoutLogger()
```

## Shutdown

همیشه:

```go
defer bot.Close()
```

را قرار دهید.

این کار Cache، Dispatcher، Rate Limiter، API Client و Logger را می‌بندد.

## Health Check

```go
if err := bot.HealthCheck(ctx); err != nil {
    log.Println("bot unhealthy:", err)
}
```

Health Check با `getMe` وضعیت ارتباط با API را بررسی می‌کند.
