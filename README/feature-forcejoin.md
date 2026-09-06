# ForceJoin

برای الزام عضویت در کانال‌های مشخص.

```go
fj := forcejoin.New(forcejoin.Config{
    RequiredChannels: []string{"@channel1", "@channel2"},
    CacheTTL: 10*time.Minute,
    CheckMembership: checkMembership,
    NotJoinedAction: notify,
})
bot.Use(fj.Middleware())
```

## Config

```text
RequiredChannels
CacheTTL
AdminIDs
CheckMembership
NotJoinedAction
```

پیش‌فرض `CacheTTL = 10m` است.

## Helpers

```go
fj.ClearUserCache(userID)
fj.Stats()
```

برای کاهش API callها نتیجه membership cache می‌شود.
