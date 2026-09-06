# Update Methods

## Updates

```go
updates := bot.API().Updates()
```

## Get

```go
items, err := updates.Get(ctx, &methods.GetUpdatesParams{
    Offset: 0,
    Limit: 100,
    Timeout: 20,
    AllowedUpdates: []string{"message", "callback_query"},
})
```

Endpoint:

```text
getUpdates
```

## GetUpdatesParams

```go
type GetUpdatesParams struct {
    Offset int
    Limit int
    Timeout int
    AllowedUpdates []string
}
```

`polling.Poller` همین لایه را برای loop و offset مدیریت می‌کند.
