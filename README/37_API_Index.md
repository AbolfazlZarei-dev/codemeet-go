# 38 — Package Map و API Index

## Core

| Package | نقش |
|---|---|
| `codemeet` | facade و runtime اصلی |
| `api` | HTTP transport و response |
| `methods` | متدهای Bot API |
| `models` | مدل‌های typed |
| `dispatcher` | event routing |
| `middleware` | middlewareهای عمومی |
| `cache` | cache و sharded cache |
| `ratelimit` | request/concurrency limiting |
| `retry` | retry policy |
| `logger` | logging |
| `errors` | typed errors |
| `polling` | long polling |
| `webhook` | webhook server |

## Methods

### Bot

`GetMe`, `SetName`, `GetName`, `GetNameWithLang`, `SetDescription`, `SetDescriptionWithLang`, `GetDescription`, `GetDescriptionWithLang`, `SetShortDescription`, `SetShortDescriptionWithLang`, `GetShortDescription`, `GetShortDescriptionWithLang`, `SetCommands`, `GetCommands`, `DeleteCommands`, `LogOut`, `Close`.

### Chat

`GetChat`, `GetChatMember`, `GetChatAdministrators`, `GetChatMemberCount`, `PinMessage`, `UnpinMessage`, `UnpinAllMessages`, `BanChatMember`, `UnbanChatMember`, `RestrictChatMember`, `PromoteChatMember`, `SetChatAdministratorCustomTitle`, `SetChatTitle`, `SetChatDescription`, `SetChatPhoto`, `DeleteChatPhoto`, `SetChatPermissions`, `LeaveChat`, `ExportChatInviteLink`, `CreateChatInviteLink`, `EditChatInviteLink`, `RevokeChatInviteLink`.

### Messages

`Send`, `SendText`, `SendHTML`, `SendMarkdown`, `SendWithKeyboard`, `Forward`, `Copy`, `EditText`, `EditTextInline`, `EditCaption`, `EditReplyMarkup`, `Delete`, `DeleteMessages`, `SendChatAction`, `AnswerCallback`, `AnswerCallbackSimple`.

### Media

`SendPhoto`, `SendPhotoWithParams`, `SendVideo`, `SendDocument`, `SendVoice`, `SendAudio`, `SendAnimation`, `SendSticker`, `SendVideoNote`, `SendMediaGroup`, `SendLocation`, `SendVenue`, `SendContact`, `SendPoll`, `SendDice`, `GetStickerSet`, `UploadStickerFile`, `GetFile`, `DownloadFile`, `SetMessageReaction`.

### Updates/Webhook

`Get`, `Set`, `GetInfo`, `Delete`, `DeleteWithDrop`.

## Contrib

- `antilink`
- `antispam`
- `forcejoin`
- `gatekeeper`
- `maintenancemode`
- `pagination`
- `profanityfilter`
- `vpndetector`
- `warnsystem`

## Models

- Bot
- Callback
- Chat
- Keyboard
- Message/Media
- Update
- User
- Webhook

## Documentation rule

این index عمداً نام APIهای قابل مشاهده در سورس را جمع می‌کند و برای packageهایی که سورس کامل آن‌ها در dump موجود نبوده، ادعای signature دقیق نمی‌کند.
