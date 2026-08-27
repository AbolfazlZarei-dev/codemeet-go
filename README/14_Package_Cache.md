# 14 — Package `cache`

Cache برای داده‌های موقت با TTL طراحی شده است.

## Basic Cache

```go
c := cache.New(5 * time.Minute)
defer c.Close()
```

## Operations

```go
Get()
Set()
SetWithTTL()
SetForever()
Delete()
Len()
Keys()
Clear()
GetOrSet()
GetOrSetWithTTL()
Close()
```

## Example

```go
value, ok := c.Get("user:123")

if !ok {
    value = loadUser()
    c.Set("user:123", value)
}
```

## Singleflight

`GetOrSet` و `GetOrSetWithTTL` از singleflight برای جلوگیری از اجرای موازی تکراری محاسبه‌ی یک key استفاده می‌کنند.

## ShardedCache

```go
sc := cache.NewSharded(32, 5*time.Minute)
defer sc.Close()
```

Key با FNV-1a به shard نگاشت می‌شود.

ویژگی مهم: هر shard goroutine cleanup اختصاصی ندارد؛ یک scheduler مرکزی cleanup همه‌ی shardها را انجام می‌دهد.

این طراحی تعداد goroutineهای پس‌زمینه را کاهش می‌دهد.

## نکته

TTL پیش‌فرض در constructor تعیین می‌شود و `SetWithTTL` امکان TTL اختصاصی را فراهم می‌کند.
