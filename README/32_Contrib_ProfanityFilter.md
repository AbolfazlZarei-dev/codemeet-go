# 33 — `contrib/profanityfilter`

فیلتر کلمات نامناسب با normalize کردن متن و تشخیص کلمات ممنوعه کار می‌کند.

## Config

```go
type Config struct {
    BannedWords []string
    DeleteMessage bool
    WarnUser bool
    WarnText string
    Action func(ctx context.Context, userID, chatID string, messageID int, reason string)
}
```

## DefaultConfig

لیست پیش‌فرض شامل چند کلمه‌ی نمونه‌ی انگلیسی است و `DeleteMessage` و `WarnUser` فعال هستند.

برای production بهتر است لیست را صریحاً با policy پروژه تنظیم کنید.

## Normalize

- lowercase
- تبدیل Leetspeak
- تبدیل punctuation به space

نمونه‌های Leetspeak که در replacer پوشش داده می‌شوند:

```text
@ -> a
4 -> a
3 -> e
1 -> i
! -> i
0 -> o
$ -> s
5 -> s
7 -> t
```

## تشخیص

ابتدا tokenها جداگانه بررسی می‌شوند تا false positive کاهش یابد.

سپس متن بدون فاصله نیز بررسی می‌شود تا مواردی مانند جدا کردن حروف با فاصله/کاراکتر تشخیص داده شود.

## Middleware

```go
pf := profanityfilter.New(profanityfilter.Config{
    BannedWords: []string{"spamword"},
    DeleteMessage: true,
    WarnUser: true,
})

bot.Use(pf.Middleware())
```

## API

```text
DefaultConfig
New
Middleware
```
