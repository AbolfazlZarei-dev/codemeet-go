# 08 — Models

## User

فیلدهای اصلی:

- `ID`
- `IsBot`
- `FirstName`
- `Username`
- `LastName`
- `LanguageCode`
- `IsPremium`
- `AddedToAttachmentMenu`
- `CanJoinGroups`
- `CanReadAllGroupMessages`
- `SupportsInlineQueries`

Helperها:

```go
user.IsUser()
user.FullName()
user.Mention()
user.HTMLMention()
```

## Chat

Chat اطلاعات type و مشخصات گفتگو را نگه می‌دارد.

Helperها:

```go
chat.IsPrivate()
chat.IsGroup()
chat.IsSupergroup()
chat.IsChannel()
```

## Message

Message مدل مرکزی رویدادهای پیام است و mediaهای متعدد را پوشش می‌دهد.

Helperها:

```go
msg.HasMedia()
msg.IsCommand()
msg.CommandName()
msg.CommandArgs()
```

## Media Models

- `PhotoSize`
- `Video`
- `Audio`
- `Document`
- `Animation`
- `Voice`
- `VideoNote`
- `Contact`
- `Location`
- `Venue`
- `Poll`
- `PollOption`
- `Dice`
- `Sticker`
- `StickerSet`
- `File`
- `InputMedia`
- `ReactionType`

## Callback

- `CallbackQuery`
- `AnswerCallbackRequest`

## Update

- `Update`
- `ChatMemberUpdated`
- `ChatJoinRequest`

Update helperها:

```go
update.Type()
update.EffectiveMessage()
update.EffectiveUser()
update.EffectiveChat()
```

## Bot Models

- `BotName`
- `BotDescription`
- `BotShortDescription`
- `BotCommand`
- `SetCommandsRequest`
- `BotCommandScope`

## Webhook Models

`SetWebhookRequest`:

- URL
- SecretToken
- DropPendingUpdates
- MaxConnections
- AllowedUpdates
- IPAddress

`WebhookInfo`:

- URL
- HasCustomCertificate
- PendingUpdateCount
- LastError
- LastErrorMessage
- LastErrorDate
- MaxConnections
- IPAddress
