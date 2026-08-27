# Models و Data Types

پکیج `models` ساختارهای JSON ورودی و خروجی Bot API را به typeهای Go تبدیل می‌کند.

## User

فیلدهای اصلی:

- `ID string`
- `IsBot bool`
- `FirstName string`
- `Username string`
- `LastName string`
- `LanguageCode string`
- `IsPremium bool`
- `CanJoinGroups bool`
- `CanReadAllGroupMessages bool`
- `SupportsInlineQueries bool`

Helperها:

```go
user.FullName()
user.Mention()
user.HTMLMention()
user.IsUser()
```

## Chat

Chat می‌تواند نوع‌های `private`، `group` یا `channel` داشته باشد و اطلاعاتی مانند ID، title، username و description را نگهداری می‌کند.

## Message

فیلدهای اصلی:

- `message_id`
- `date`
- `chat`
- `from`
- `text`
- `caption`
- `entities`
- `reply_to_message`
- `reply_markup`

## MessageEntity

برای قالب‌بندی متن استفاده می‌شود.

Typeهای تعریف‌شده در مستندات Bot API شامل مواردی مانند:

`bold`, `italic`, `underline`, `strikethrough`, `spoiler`, `code`, `pre`, `text_link`, `mention`, `hashtag`, `bot_command`

هستند.

## Update

Update شناسه‌ی ترتیبی `update_id` دارد و می‌تواند انواع مختلف رویداد را در خود داشته باشد.

مثال:

```go
if update.Message != nil {
    fmt.Println(update.Message.Text)
}

if update.CallbackQuery != nil {
    fmt.Println(update.CallbackQuery.Data)
}
```

## CallbackQuery

فیلدهای اصلی:

- `id`
- `from`
- `message`
- `data`

## Webhook Models

`SetWebhookRequest` شامل:

- `URL`
- `SecretToken`
- `DropPendingUpdates`
- `MaxConnections`
- `AllowedUpdates`
- `IPAddress`

است.

`WebhookInfo` نیز URL، وضعیت certificate، تعداد Updateهای pending، خطای آخر و اطلاعات connection را نگهداری می‌کند.

## Effective Helpers

Update helperهایی برای استخراج User، Message و Chat موثر دارد. این موضوع برای handlerهایی که چند نوع Update را دریافت می‌کنند کاربردی است.
