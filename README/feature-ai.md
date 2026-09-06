# AI Client

Package: `ai`

## Config

```go
type Config struct {
    APIKey string
    BaseURL string
}
```

BaseURL پیش‌فرض:

```text
https://openrouter.ai/api/v1
```

## Message / Request

```go
type Message struct {
    Role string
    Content string
}

type ChatRequest struct {
    Model string
    Messages []Message
    Stream bool
}
```

## ساخت

```go
client := ai.New(ai.Config{
    APIKey: os.Getenv("OPENROUTER_API_KEY"),
})
```

## Ask

```go
answer, err := client.Ask(ctx, ai.ChatRequest{
    Model: "model-name",
    Messages: []ai.Message{
        {Role: "user", Content: "سلام"},
    },
})
```

## Streaming

```go
err := client.StreamChat(
    ctx,
    request,
    func(content string) {
        fmt.Print(content)
    },
)
```

API Key را داخل source قرار ندهید.
