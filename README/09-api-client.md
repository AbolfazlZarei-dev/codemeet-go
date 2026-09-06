# API Client

Package: `api`

## ساخت

```go
client := api.NewClient(
    "https://botapi.codemeet.chat",
    token,
    log,
)
```

## Request

```go
resp, err := client.Request(ctx, "getMe")
```

## RequestWithParams

```go
resp, err := client.RequestWithParams(
    ctx,
    "sendMessage",
    map[string]interface{}{
        "chat_id": chatID,
        "text": "سلام",
    },
)
```

## Multipart

```go
resp, err := client.RequestWithMultipart(
    ctx,
    "sendDocument",
    map[string]string{"chat_id": chatID},
    map[string]string{"document": "./file.pdf"},
)
```

`RequestWithForm` نیز به Multipart متصل است.

## DownloadFile

```go
r, err := client.DownloadFile(ctx, filePath)
if err != nil { return err }
defer r.Close()
```

## HTTP Client و Timeout

```go
client.SetHTTPClient(customHTTP)
client.SetTimeout(30 * time.Second)
```

در `codemeet.New` نیز:

```go
codemeet.WithHTTPClient(customHTTP)
codemeet.WithTimeout(30*time.Second)
```

## Request ID

```go
ctx = api.WithRequestID(ctx, "job-123")
id := api.GetRequestID(ctx)
```

ID به Header `X-Request-ID` منتقل می‌شود.

## Response

```go
type Response struct {
    Ok          bool
    Result      json.RawMessage
    ErrorCode   int
    Description string
    Parameters  interface{}
    HTTPStatus  int
}
```

## Decode

```go
var user models.User
if err := resp.Decode(&user); err != nil { ... }
```

## AsBool

```go
value, err := resp.AsBool()
```

## retry_after

```go
seconds := resp.ParametersAsRetryAfter()
```

## Stats

```go
s := client.StatsSnapshot()
```

شامل `Requests`, `SuccessCount`, `ErrorCount`, `BytesIn`, `BytesOut`, `AvgLatency` است.

## Circuit Breaker

```go
breaker := client.Breaker()
state := breaker.State()
```

Stateها:

```text
0 Closed
1 Open
2 Half-Open
```

در Half-Open فقط یک Probe همزمان اجازه عبور دارد.
