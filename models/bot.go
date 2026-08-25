package models

type BotName struct {
	Name string `json:"name"`
}

type BotDescription struct {
	Description string `json:"description"`
}

type BotShortDescription struct {
	ShortDescription string `json:"short_description"`
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type SetCommandsRequest struct {
	Commands     []BotCommand     `json:"commands"`
	LanguageCode string           `json:"language_code,omitempty"`
	Scope        *BotCommandScope `json:"scope,omitempty"`
}

// BotCommandScope — محدوده دستورات
type BotCommandScope struct {
	Type   string `json:"type"`
	ChatID string `json:"chat_id,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

const (
	BotCommandScopeDefault               = "default"
	BotCommandScopeAllPrivateChats       = "all_private_chats"
	BotCommandScopeAllGroupChats         = "all_group_chats"
	BotCommandScopeAllChatAdministrators = "all_chat_administrators"
	BotCommandScopeChat                  = "chat"
	BotCommandScopeChatAdministrators    = "chat_administrators"
	BotCommandScopeChatMember            = "chat_member"
)
