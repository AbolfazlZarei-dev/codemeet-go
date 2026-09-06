# Endpoint Map

این جدول ارتباط wrapperهای اصلی `methods` با نام Endpoint را نشان می‌دهد.

| Go | Endpoint |
|---|---|
| Bot.GetMe | getMe |
| Bot.SetName/GetName | setMyName/getMyName |
| Bot.SetDescription/GetDescription | setMyDescription/getMyDescription |
| Bot.SetShortDescription/GetShortDescription | setMyShortDescription/getMyShortDescription |
| Bot.SetCommands/GetCommands/DeleteCommands | setMyCommands/getMyCommands/deleteMyCommands |
| Chat.GetChat | getChat |
| Chat.GetChatMember | getChatMember |
| Chat.GetChatAdministrators | getChatAdministrators |
| Chat.GetChatMemberCount | getChatMemberCount |
| Chat.PinMessage | pinChatMessage |
| Chat.UnpinMessage | unpinChatMessage |
| Chat.UnpinAllMessages | unpinAllChatMessages |
| Chat.BanChatMember | banChatMember |
| Chat.UnbanChatMember | unbanChatMember |
| Chat.RestrictChatMember | restrictChatMember |
| Chat.PromoteChatMember | promoteChatMember |
| Chat.SetChatAdministratorCustomTitle | setChatAdministratorCustomTitle |
| Chat.SetChatTitle | setChatTitle |
| Chat.SetChatDescription | setChatDescription |
| Chat.SetChatPhoto | setChatPhoto |
| Chat.DeleteChatPhoto | deleteChatPhoto |
| Chat.SetChatPermissions | setChatPermissions |
| Chat.LeaveChat | leaveChat |
| Chat.ExportChatInviteLink | exportChatInviteLink |
| Chat.CreateChatInviteLink | createChatInviteLink |
| Chat.EditChatInviteLink | editChatInviteLink |
| Chat.RevokeChatInviteLink | revokeChatInviteLink |
| Messages.Send | sendMessage |
| Messages.Forward | forwardMessage |
| Messages.Copy | copyMessage |
| Messages.EditText/EditTextInline | editMessageText |
| Messages.EditCaption | editMessageCaption |
| Messages.EditReplyMarkup | editMessageReplyMarkup |
| Messages.Delete | deleteMessage |
| Messages.DeleteMessages | deleteMessages |
| Messages.SendChatAction | sendChatAction |
| Messages.AnswerCallback | answerCallbackQuery |
| Media.SendPhoto | sendPhoto |
| Media.SendVideo | sendVideo |
| Media.SendDocument | sendDocument |
| Media.SendVoice | sendVoice |
| Media.SendAudio | sendAudio |
| Media.SendAnimation | sendAnimation |
| Media.SendSticker | sendSticker |
| Media.SendVideoNote | sendVideoNote |
| Media.SendMediaGroup | sendMediaGroup |
| Media.SendLocation | sendLocation |
| Media.SendVenue | sendVenue |
| Media.SendContact | sendContact |
| Media.SendPoll | sendPoll |
| Media.SendDice | sendDice |
| Media.GetStickerSet | getStickerSet |
| Media.UploadStickerFile | uploadStickerFile |
| Media.GetFile | getFile |
| Media.SetMessageReaction | setMessageReaction |
| Updates.Get | getUpdates |
| Webhook.Set | setWebhook |
| Webhook.GetInfo | getWebhookInfo |
| Webhook.Delete/DeleteWithDrop | deleteWebhook |

> **هشدار:** این map درباره wrapper/Endpoint در SDK است، نه تضمین availability سمت Server. در صورت `404` یا `501`، نسخه فعلی Bot API را بررسی کنید.
