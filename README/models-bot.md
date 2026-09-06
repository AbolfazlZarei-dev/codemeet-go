# Bot Models

Package: `models`

## `BotName`

```go
type BotName struct {
    Name string `json:"name"`
}
```

## `BotDescription`

```go
type BotDescription struct {
    Description string `json:"description"`
}
```

## `BotShortDescription`

```go
type BotShortDescription struct {
    ShortDescription string `json:"short_description"`
}
```

## `BotCommand`

```go
type BotCommand struct {
    Command string `json:"command"`
    Description string `json:"description"`
}
```

## `SetCommandsRequest`

```go
type SetCommandsRequest struct {
    Commands []BotCommand `json:"commands"`
    LanguageCode string `json:"language_code,omitempty"`
    Scope *BotCommandScope `json:"scope,omitempty"`
}
```

## `BotCommandScope`

```go
type BotCommandScope struct {
    Type string `json:"type"`
    ChatID string `json:"chat_id,omitempty"`
    UserID string `json:"user_id,omitempty"`
}
```

### Scope constants

```text
default
all_private_chats
all_group_chats
all_chat_administrators
chat
chat_administrators
chat_member
```
