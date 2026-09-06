# Cache

Package: `cache`

## Cache

```go
c := cache.New(10 * time.Minute)
defer c.Close()
```

متدها:

```text
Get
Set
SetWithTTL
SetForever
Delete
Len
Keys
Clear
GetOrSet
GetOrSetWithTTL
Close
```

## نمونه

```go
c.Set("user:1", user)

v, ok := c.Get("user:1")

user2, ok := cache.GetTyped[*models.User](c, "user:1")
```

## GetOrSet

```go
v := c.GetOrSet("expensive", func() interface{} {
    return load()
})
```

برای جلوگیری از محاسبه همزمان یک key از `singleflight` استفاده می‌شود.

## ShardedCache

```go
sc := cache.NewSharded(32, 10*time.Minute)
defer sc.Close()
```

کلید با FNV بین shardها توزیع می‌شود و یک scheduler مرکزی cleanup را انجام می‌دهد.

## Bot

```go
codemeet.WithCache(10*time.Minute)
codemeet.WithShardedCache(32, 10*time.Minute)
```

## نکته

Cache حافظه‌ای است؛ برای persistence از Database استفاده کنید.
