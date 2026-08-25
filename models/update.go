package models

type Update struct {
	UpdateID          int                `json:"update_id"`
	Message           *Message           `json:"message,omitempty"`
	EditedMessage     *Message           `json:"edited_message,omitempty"`
	CallbackQuery     *CallbackQuery     `json:"callback_query,omitempty"`
	ChannelPost       *Message           `json:"channel_post,omitempty"`
	EditedChannelPost *Message           `json:"edited_channel_post,omitempty"`
	MyChatMember      *ChatMemberUpdated `json:"my_chat_member,omitempty"`
	ChatMember        *ChatMemberUpdated `json:"chat_member,omitempty"`
	ChatJoinRequest   *ChatJoinRequest   `json:"chat_join_request,omitempty"`
}

// ChatMemberUpdated تغییر عضویت
type ChatMemberUpdated struct {
	Chat          *Chat       `json:"chat"`
	From          *User       `json:"from"`
	Date          int64       `json:"date"`
	OldChatMember *ChatMember `json:"old_chat_member"`
	NewChatMember *ChatMember `json:"new_chat_member"`
}

// ChatJoinRequest درخواست عضویت
type ChatJoinRequest struct {
	Chat       *Chat           `json:"chat"`
	From       *User           `json:"from"`
	Date       int64           `json:"date"`
	Bio        string          `json:"bio,omitempty"`
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}

func (u *Update) Type() string {
	if u.Message != nil {
		return "message"
	}
	if u.EditedMessage != nil {
		return "edited_message"
	}
	if u.CallbackQuery != nil {
		return "callback_query"
	}
	if u.ChannelPost != nil {
		return "channel_post"
	}
	if u.EditedChannelPost != nil {
		return "edited_channel_post"
	}
	if u.MyChatMember != nil {
		return "my_chat_member"
	}
	if u.ChatMember != nil {
		return "chat_member"
	}
	if u.ChatJoinRequest != nil {
		return "chat_join_request"
	}
	return "unknown"
}

// EffectiveMessage پیام موثر (message یا edited_message یا channel_post)
func (u *Update) EffectiveMessage() *Message {
	if u.Message != nil {
		return u.Message
	}
	if u.EditedMessage != nil {
		return u.EditedMessage
	}
	if u.ChannelPost != nil {
		return u.ChannelPost
	}
	if u.EditedChannelPost != nil {
		return u.EditedChannelPost
	}
	return nil
}

// EffectiveUser کاربر موثر
func (u *Update) EffectiveUser() *User {
	if u.Message != nil && u.Message.From != nil {
		return u.Message.From
	}
	if u.CallbackQuery != nil {
		return u.CallbackQuery.From
	}
	if u.EditedMessage != nil && u.EditedMessage.From != nil {
		return u.EditedMessage.From
	}
	if u.MyChatMember != nil {
		return u.MyChatMember.From
	}
	if u.ChatJoinRequest != nil {
		return u.ChatJoinRequest.From
	}
	return nil
}

// EffectiveChat چت موثر
func (u *Update) EffectiveChat() *Chat {
	if m := u.EffectiveMessage(); m != nil {
		return m.Chat
	}
	if u.CallbackQuery != nil && u.CallbackQuery.Message != nil {
		return u.CallbackQuery.Message.Chat
	}
	if u.MyChatMember != nil {
		return u.MyChatMember.Chat
	}
	if u.ChatJoinRequest != nil {
		return u.ChatJoinRequest.Chat
	}
	return nil
}
