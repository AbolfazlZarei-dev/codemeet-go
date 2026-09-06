# Chat Methods

دسترسی:

```go
chat := bot.API().Chat()
```

> وجود wrapper به معنی فعال بودن Endpoint سمت سرور نیست.

| Method | Endpoint |
|---|---|
| GetChat | getChat |
| GetChatMember | getChatMember |
| GetChatAdministrators | getChatAdministrators |
| GetChatMemberCount | getChatMemberCount |
| PinMessage | pinChatMessage |
| UnpinMessage | unpinChatMessage |
| UnpinAllMessages | unpinAllChatMessages |
| BanChatMember | banChatMember |
| UnbanChatMember | unbanChatMember |
| RestrictChatMember | restrictChatMember |
| PromoteChatMember | promoteChatMember |
| SetChatAdministratorCustomTitle | setChatAdministratorCustomTitle |
| SetChatTitle | setChatTitle |
| SetChatDescription | setChatDescription |
| SetChatPhoto | setChatPhoto |
| DeleteChatPhoto | deleteChatPhoto |
| SetChatPermissions | setChatPermissions |
| LeaveChat | leaveChat |
| ExportChatInviteLink | exportChatInviteLink |
| CreateChatInviteLink | createChatInviteLink |
| EditChatInviteLink | editChatInviteLink |
| RevokeChatInviteLink | revokeChatInviteLink |

## مثال

```go
c, err := chat.GetChat(ctx, chatID)
member, err := chat.GetChatMember(ctx, chatID, userID)
admins, err := chat.GetChatAdministrators(ctx, chatID)
count, err := chat.GetChatMemberCount(ctx, chatID)

err = chat.PinMessage(ctx, chatID, messageID, false)
err = chat.SetChatTitle(ctx, chatID, "عنوان جدید")
```

## Restrict

```go
perms := &models.ChatPermissions{
    CanSendMessages: false,
}
err := chat.RestrictChatMember(ctx, chatID, userID, perms, 0)
```

## Promote

```go
rights := &models.ChatAdministratorRights{
    CanDeleteMessages: true,
    CanRestrictMembers: true,
}
err := chat.PromoteChatMember(ctx, chatID, userID, rights)
```

## Invite

```go
req := &models.CreateChatInviteLinkRequest{
    ChatID: chatID,
    Name: "main",
}
link, err := chat.CreateChatInviteLink(ctx, req)
```

عملیات مدیریتی نیازمند Permission مناسب Bot هستند.
