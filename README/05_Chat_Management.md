# 05 — Chats، Groups و Channels

`ChatMethods` تمام عملیات اصلی مدیریت Chat را پوشش می‌دهد.

## اطلاعات Chat

```go
chat, err := bot.API().Chat().GetChat(ctx, chatID)
```

## Member

```go
member, err := bot.API().Chat().GetChatMember(
    ctx,
    chatID,
    userID,
)
```

## Admins

```go
admins, err := bot.API().Chat().GetChatAdministrators(ctx, chatID)
```

## Member Count

```go
count, err := bot.API().Chat().GetChatMemberCount(ctx, chatID)
```

## Pin

- `PinMessage`
- `UnpinMessage`
- `UnpinAllMessages`

## مدیریت عضو

- `BanChatMember`
- `UnbanChatMember`
- `RestrictChatMember`
- `PromoteChatMember`
- `SetChatAdministratorCustomTitle`

## تنظیمات

- `SetChatTitle`
- `SetChatDescription`
- `SetChatPhoto`
- `DeleteChatPhoto`
- `SetChatPermissions`
- `LeaveChat`

## Invite Link

- `ExportChatInviteLink`
- `CreateChatInviteLink`
- `EditChatInviteLink`
- `RevokeChatInviteLink`

## Chat type helpers

```go
chat.IsPrivate()
chat.IsGroup()
chat.IsSupergroup()
chat.IsChannel()
```

عملیات مدیریتی نیازمند مجوز مناسب Bot در Chat است.
