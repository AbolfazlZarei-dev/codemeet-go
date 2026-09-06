# Errors

Package: `errors`

## کدها

```text
400 BadRequest
401 Unauthorized
403 Forbidden
404 NotFound
405 MethodNotAllowed
409 Conflict
413 RequestTooLarge
429 TooManyRequests
500 Internal
501 NotImplemented
502 BadGateway
503 ServiceUnavailable
504 GatewayTimeout
```

## APIError

```go
apiErr, ok := errors.AsAPIError(err)
if ok {
    fmt.Println(apiErr.Code)
    fmt.Println(apiErr.Description)
    fmt.Println(apiErr.RetryAfter)
}
```

## سایر خطاها

```text
ValidationError
NetworkError
MultiError
```

## Retry

```go
if errors.IsRetryable(err) {
    // retry
}
```

## MultiError

```go
var me errors.MultiError
me.Add(err1)
me.Add(err2)

if me.HasError() { ... }
```

## ParseError

API Client پاسخ ناموفق را parse می‌کند و خطای ساختاریافته برمی‌گرداند.

## 429

`ParametersAsRetryAfter()` و `APIError.RetryAfter` برای رعایت rate-limit کاربرد دارند.
