# 18 — Package `errors`

Package errors خطاها را typed و قابل تشخیص می‌کند.

## APIError

برای خطاهای Bot API:

- Code
- Description
- Parameters
- RetryAfter
- HTTP status information

## ValidationError

خطاهای ورودی برنامه:

```go
errors.NewValidationError("token", "token is required")
```

## NetworkError

برای خطاهای شبکه.

## MultiError

برای جمع‌آوری چند خطا:

```go
Add(err)
HasError()
```

## Helpers

```go
IsRetryable(err)
AsAPIError(err)
AsNetworkError(err)
AsValidationError(err)
ParseError(...)
```

## Retryable

خطاهای قابل retry از مسیر `retry.Policy` دوباره اجرا می‌شوند.

429 و برخی خطاهای شبکه/5xx در مسیر retry اهمیت ویژه دارند.
