# Keyboard Models

Package: `models`

## `InlineKeyboardMarkup`

```go
type InlineKeyboardMarkup struct {
    InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}
```

## `InlineKeyboardButton`

```go
type InlineKeyboardButton struct {
    Text string `json:"text"`
    CallbackData string `json:"callback_data,omitempty"`
    URL string `json:"url,omitempty"`
    SwitchInlineQuery string `json:"switch_inline_query,omitempty"`
    SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"`
    WebApp *WebAppInfo `json:"web_app,omitempty"`
    LoginURL *LoginURL `json:"login_url,omitempty"`
}
```

## `WebAppInfo`

```go
type WebAppInfo struct {
    URL string `json:"url"`
}
```

## `LoginURL`

```go
type LoginURL struct {
    URL string `json:"url"`
    ForwardText string `json:"forward_text,omitempty"`
    BotUsername string `json:"bot_username,omitempty"`
    RequestWriteAccess bool `json:"request_write_access,omitempty"`
}
```

## `ReplyKeyboardMarkup`

```go
type ReplyKeyboardMarkup struct {
    Keyboard [][]KeyboardButton `json:"keyboard"`
    ResizeKeyboard bool `json:"resize_keyboard,omitempty"`
    OneTimeKeyboard bool `json:"one_time_keyboard,omitempty"`
    IsPersistent bool `json:"is_persistent,omitempty"`
    Selective bool `json:"selective,omitempty"`
    InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
}
```

## `KeyboardButton`

```go
type KeyboardButton struct {
    Text string `json:"text"`
    RequestContact bool `json:"request_contact,omitempty"`
    RequestLocation bool `json:"request_location,omitempty"`
    RequestPoll *KeyboardButtonPollType `json:"request_poll,omitempty"`
    WebApp *WebAppInfo `json:"web_app,omitempty"`
}
```

## `KeyboardButtonPollType`

```go
type KeyboardButtonPollType struct {
    Type string `json:"type,omitempty"`
}
```

## `ReplyKeyboardRemove`

```go
type ReplyKeyboardRemove struct {
    RemoveKeyboard bool `json:"remove_keyboard"`
    Selective bool `json:"selective,omitempty"`
}
```

## `ForceReply`

```go
type ForceReply struct {
    ForceReply bool `json:"force_reply"`
    Selective bool `json:"selective,omitempty"`
    InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
}
```

## Builder helpers

```go
models.NewInlineKeyboard(...)
models.InlineRow(...)
models.Btn(...)
models.URLBtn(...)
models.WebAppBtn(...)
models.SwitchInlineBtn(...)

models.NewReplyKeyboard(...)
models.ReplyRow(...)
models.KBtn(...)
models.ContactBtn(...)
models.LocationBtn(...)
```

Inline keyboard به صورت `[][]InlineKeyboardButton` مدل می‌شود؛ یعنی آرایه‌ای از سطرها.
