# Webhook Models

Package: `models`

## `SetWebhookRequest`

```go
type SetWebhookRequest struct {
    URL string `json:"url"`
    SecretToken string `json:"secret_token,omitempty"`
    DropPendingUpdates bool `json:"drop_pending_updates,omitempty"`
    MaxConnections int `json:"max_connections,omitempty"`
    AllowedUpdates []string `json:"allowed_updates,omitempty"`
    IPAddress string `json:"ip_address,omitempty"`
}
```

## `WebhookInfo`

```go
type WebhookInfo struct {
    URL string `json:"url"`
    HasCustomCertificate bool `json:"has_custom_certificate"`
    PendingUpdateCount int `json:"pending_update_count"`
    LastError string `json:"last_error,omitempty"`
    LastErrorMessage string `json:"last_error_message,omitempty"`
    LastErrorDate int64 `json:"last_error_date,omitempty"`
    MaxConnections int `json:"max_connections,omitempty"`
    IPAddress string `json:"ip_address,omitempty"`
}
```

## نمونه

```go
req := &models.SetWebhookRequest{
    URL: "https://example.com/webhook",
    SecretToken: "secret",
    AllowedUpdates: []string{"message", "callback_query"},
}
```
