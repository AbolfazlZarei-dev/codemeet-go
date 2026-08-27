# مدیریت Chat، گروه و کانال

لایه‌ی `ChatMethods` برای کار با گفتگوها، گروه‌ها و کانال‌ها است.

## دریافت اطلاعات Chat

```go
chat, err := bot.API().Chat().GetChat(ctx, chatID)
```

مدل Chat اطلاعاتی مانند ID، type، title، username و description را پوشش می‌دهد.

## اطلاعات عضو

```go
member, err := bot.API().Chat().GetChatMember(
    ctx,
    chatID,
    userID,
)
```

## مدیران

```go
admins, err := bot.API().Chat().GetChatAdministrators(ctx, chatID)
```

## تعداد اعضا

```go
count, err := bot.API().Chat().GetChatMemberCount(ctx, chatID)
```

## Pin

```go
err := bot.API().Chat().PinMessage(
    ctx,
    chatID,
    messageID,
    false,
)
```

همچنین:

- `UnpinMessage`
- `UnpinAllMessages`

وجود دارند.

## مدیریت اعضا

کتابخانه متدهای مدیریت عضو مانند:

- `BanChatMember`
- `UnbanChatMember`
- `RestrictChatMember`

را ارائه می‌کند.

دسترسی واقعی این عملیات به سطح دسترسی Bot در Chat وابسته است.

## تنظیمات Chat

متدهای موجود شامل:

- `SetChatTitle`
- `SetChatDescription`
- `SetChatPhoto`
- `DeleteChatPhoto`
- `SetChatPermissions`
- `LeaveChat`

است.

## لینک دعوت

امکانات مدیریت Invite Link شامل:

- `ExportChatInviteLink`
- `CreateChatInviteLink`
- `EditChatInviteLink`
- `RevokeChatInviteLink`

است.

## نکته

برای عملیات مدیریتی گروه/کانال، Bot باید عضو Chat باشد و برای عملیات مدیریتی لازم، مجوز مناسب داشته باشد.
