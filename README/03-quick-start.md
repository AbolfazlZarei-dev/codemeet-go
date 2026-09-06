# شروع سریع

## Echo Bot

```go
package main

import (
    "context"
    "log"
    "os"

    codemeet "github.com/AbolfazlZarei-dev/codemeet-go"
    "github.com/AbolfazlZarei-dev/codemeet-go/models"
    "github.com/AbolfazlZarei-dev/codemeet-go/polling"
)

func main() {
    token := os.Getenv("CODEMEET_BOT_TOKEN")
    bot, err := codemeet.New(token)
    if err != nil { log.Fatal(err) }
    defer bot.Close()

    bot.OnCommand("start", func(ctx context.Context, msg *models.Message) {
        _, err := bot.Send(ctx, msg.Chat.ID, "سلام 👋")
        if err != nil { log.Println(err) }
    })

    bot.OnMessage(func(ctx context.Context, msg *models.Message) {
        if msg.Text == "" { return }
        _, err := bot.Reply(ctx, msg, "دریافت شد: "+msg.Text)
        if err != nil { log.Println(err) }
    })

    ctx := context.Background()
    if err := bot.StartPolling(ctx, polling.DefaultConfig()); err != nil {
        log.Fatal(err)
    }
}
```

## معماری اجرا

```text
New()
  ↓
Register handlers
  ↓
Use middleware
  ↓
StartPolling / StartWebhook
  ↓
Update
  ↓
Dispatcher
  ↓
Handler
  ↓
API response
```

## Callback

```go
bot.OnCallback(func(ctx context.Context, cq *models.CallbackQuery) {
    _ = bot.AnswerCallback(ctx, cq.ID, "انجام شد", false)
})
```

برای پروژه واقعی، Recovery و Rate Limit را نیز در نظر بگیرید.
