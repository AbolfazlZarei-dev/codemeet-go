# Update Models

Package: `models`

## `CallbackQuery`

```go
type CallbackQuery struct {
    ID string `json:"id"`
    From *User `json:"from"`
    Message *Message `json:"message,omitempty"`
    InlineMessageID string `json:"inline_message_id,omitempty"`
    Data string `json:"data,omitempty"`
    GameShortName string `json:"game_short_name,omitempty"`
}
```

## `AnswerCallbackRequest`

```go
type AnswerCallbackRequest struct {
    CallbackQueryID string `json:"callback_query_id"`
    Text string `json:"text,omitempty"`
    ShowAlert bool `json:"show_alert,omitempty"`
    URL string `json:"url,omitempty"`
    CacheTime int `json:"cache_time,omitempty"`
}
```

## `Update`

```go
type Update struct {
    UpdateID int `json:"update_id"`
    Message *Message `json:"message,omitempty"`
    EditedMessage *Message `json:"edited_message,omitempty"`
    CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
    ChannelPost *Message `json:"channel_post,omitempty"`
    EditedChannelPost *Message `json:"edited_channel_post,omitempty"`
    MyChatMember *ChatMemberUpdated `json:"my_chat_member,omitempty"`
    ChatMember *ChatMemberUpdated `json:"chat_member,omitempty"`
    ChatJoinRequest *ChatJoinRequest `json:"chat_join_request,omitempty"`
}
```

## `ChatMemberUpdated`

```go
type ChatMemberUpdated struct {
    Chat *Chat `json:"chat"`
    From *User `json:"from"`
    Date int64 `json:"date"`
    OldChatMember *ChatMember `json:"old_chat_member"`
    NewChatMember *ChatMember `json:"new_chat_member"`
}
```

## `ChatJoinRequest`

```go
type ChatJoinRequest struct {
    Chat *Chat `json:"chat"`
    From *User `json:"from"`
    Date int64 `json:"date"`
    Bio string `json:"bio,omitempty"`
    InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}
```

## Helperهای Update

```go
u.Type()
u.EffectiveMessage()
u.EffectiveUser()
u.EffectiveChat()
```

`CallbackQuery.Data` داده callback دکمه را حمل می‌کند.
