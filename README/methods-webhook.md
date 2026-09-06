# Webhook Methods

```go
w := bot.API().Webhook()
```

| Method | Endpoint |
|---|---|
| Set | setWebhook |
| GetInfo | getWebhookInfo |
| Delete | deleteWebhook |
| DeleteWithDrop | deleteWebhook |

## نمونه

```go
err := w.Set(ctx, &models.SetWebhookRequest{
    URL: "https://example.com/webhook",
    SecretToken: "secret",
})

info, err := w.GetInfo(ctx)

err = w.Delete(ctx)
err = w.DeleteWithDrop(ctx, true)
```

> فعال بودن Endpoint را در Bot API فعلی سرویس مقصد بررسی کنید.
