package models

type User struct {
	ID                      string `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	Username                string `json:"username,omitempty"`
	LastName                string `json:"last_name,omitempty"`
	LanguageCode            string `json:"language_code,omitempty"`
	IsPremium               bool   `json:"is_premium,omitempty"`
	AddedToAttachmentMenu   bool   `json:"added_to_attachment_menu,omitempty"`
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`
}

// IsUser آیا کاربر عادی است (نه ربات)
func (u *User) IsUser() bool { return !u.IsBot }

// FullName نام کامل کاربر
func (u *User) FullName() string {
	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.FirstName
}

// Mention منشن کاربر
func (u *User) Mention() string {
	if u.Username != "" {
		return "@" + u.Username
	}
	return u.FullName()
}

// HTMLMention منشن کاربر با لینک HTML
func (u *User) HTMLMention() string {
	if u.Username != "" {
		return `<a href="https://codemeet.chat/` + u.Username + `">` + u.FullName() + "</a>"
	}
	return u.FullName()
}
