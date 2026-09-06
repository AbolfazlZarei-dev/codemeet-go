# Database

Package: `database`

## ساخت

```go
db, err := database.New("./data/bot.db")
if err != nil { log.Fatal(err) }
defer db.Close()
```

## SQL

```go
sqlDB := db.SQL()
```

## Group

```go
err := db.SaveGroup(ctx, chatID, title, groupType)
```

## KV با TTL

```go
err := db.SetKV(ctx, "session", "user:1", value, 30*time.Minute)
value, err := db.GetKV(ctx, "session", "user:1")
err = db.DeleteKV(ctx, "session", "user:1")
err = db.CleanupExpired(ctx)
```

## User

```go
err := db.SaveUser(ctx, userID, chatID, firstName, username, isBot)
u, err := db.GetUser(ctx, userID)
```

## Ban

```go
_ = db.BanUser(ctx, userID)
_ = db.UnbanUser(ctx, userID)

banned, err := db.IsBanned(ctx, userID)
```

## User Model

```go
type User struct {
    ID string
    FirstName string
    Username string
    IsBot bool
    JoinedAt int64
    Banned bool
}
```

Database برای state پایدار مناسب است؛ Cache حافظه‌ای را جایگزین persistence نکنید.
