# 25 — UpdatesMethods و WebhookMethods

## Updates

```go
type GetUpdatesParams struct {
    Offset int
    Limit int
    Timeout int
    AllowedUpdates []string
}
```

متد:

```go
Get(ctx, params)
```

به `getUpdates` متصل است.

## Webhook

```text
Set
GetInfo
Delete
DeleteWithDrop
```

## مثال

```go
req := &models.SetWebhookRequest{
    URL: "https://example.com/webhook",
    SecretToken: "secret",
    AllowedUpdates: []string{"message", "callback_query"},
}

err := bot.API().Webhook().Set(ctx, req)
```

## Delete with drop

```go
err := bot.API().Webhook().DeleteWithDrop(ctx, true)
```

در این حالت Updateهای pending نیز طبق پارامتر API حذف می‌شوند.
