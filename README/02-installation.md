# نصب

## Go Module

```bash
go get github.com/AbolfazlZarei-dev/codemeet-go@v1.1.0
go mod tidy
```

## Import

```go
import codemeet "github.com/AbolfazlZarei-dev/codemeet-go"
```

## Token

Token را در environment نگه دارید:

```bash
export CODEMEET_BOT_TOKEN="YOUR_TOKEN"
```

Windows PowerShell:

```powershell
$env:CODEMEET_BOT_TOKEN="YOUR_TOKEN"
```

## Build

```bash
go build -o codemeet-bot .
```

Cross compile برای Linux amd64:

```powershell
$env:CGO_ENABLED="0"
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o codemeet-bot .
```

## نسخه

```go
fmt.Println(codemeet.Version) // 1.1.0
```

> Token را commit یا داخل log عمومی چاپ نکنید.
