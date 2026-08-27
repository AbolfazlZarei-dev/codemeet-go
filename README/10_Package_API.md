# 10 — Package `api`

## Client

`api.Client` لایه‌ی transport است.

### State

- Base URL
- Token
- HTTP Client
- Transport
- Logger
- Statistics
- Circuit Breaker

## Constructor

```go
NewClient(baseURL, token, logger)
```

## Configuration

```go
SetHTTPClient()
SetTimeout()
UserAgent()
```

## Requests

```go
RequestWithParams()
Request()
RequestWithMultipart()
RequestWithForm()
```

## Files

```go
DownloadFile()
```

## Statistics

```go
Stats()
StatsSnapshot()
```

Snapshot شامل:

```text
Requests
SuccessCount
ErrorCount
BytesIn
BytesOut
AvgLatency
```

## Circuit Breaker

```go
NewCircuitBreaker()
Allow()
RecordSuccess()
RecordFailure()
State()
```

## Request ID

```go
WithRequestID(ctx, id)
GetRequestID(ctx)
```

## Response

`api.Response`:

```go
Decode()
AsBool()
ParametersAsRetryAfter()
```

## طراحی عملکردی

- `sync.Pool` برای buffer/encoder
- JSON streaming decode
- multipart streaming
- connection pooling
- atomic statistics
- circuit breaker
- محدودسازی response
