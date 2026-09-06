# Console و Broadcast

## Console

Package `contrib/console` برای commandهای مدیریتی محلی.

```go
c := console.New(...)
c.Register("stats", func(args []string) {
    // ...
})
err := c.ExecuteCommand("stats")
c.Start()
```

## Broadcast

Package `contrib/broadcast` برای ارسال پیام به چند مقصد.

```go
bc := broadcast.New(...)
err := bc.Send(...)
```

امضای دقیق constructor/action را از نسخه نصب‌شده دنبال کنید.

## Bot Options

```go
codemeet.WithConsole(c)
codemeet.WithBroadcaster(bc)
```
