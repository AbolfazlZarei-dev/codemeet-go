package models

// InlineKeyboardMarkup کیبورد شیشه‌ای
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton دکمه شیشه‌ای
type InlineKeyboardButton struct {
	Text                         string      `json:"text"`
	CallbackData                 string      `json:"callback_data,omitempty"`
	URL                          string      `json:"url,omitempty"`
	SwitchInlineQuery            string      `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat string      `json:"switch_inline_query_current_chat,omitempty"`
	WebApp                       *WebAppInfo `json:"web_app,omitempty"`
	LoginURL                     *LoginURL   `json:"login_url,omitempty"`
}

// WebAppInfo اطلاعات Web App
type WebAppInfo struct {
	URL string `json:"url"`
}

// LoginURL اطلاعات لاگین
type LoginURL struct {
	URL                string `json:"url"`
	ForwardText        string `json:"forward_text,omitempty"`
	BotUsername        string `json:"bot_username,omitempty"`
	RequestWriteAccess bool   `json:"request_write_access,omitempty"`
}

// ReplyKeyboardMarkup کیبورد معمولی
type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard        bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard,omitempty"`
	IsPersistent          bool               `json:"is_persistent,omitempty"`
	Selective             bool               `json:"selective,omitempty"`
	InputFieldPlaceholder string             `json:"input_field_placeholder,omitempty"`
}

// KeyboardButton دکمه کیبورد معمولی
type KeyboardButton struct {
	Text            string                  `json:"text"`
	RequestContact  bool                    `json:"request_contact,omitempty"`
	RequestLocation bool                    `json:"request_location,omitempty"`
	RequestPoll     *KeyboardButtonPollType `json:"request_poll,omitempty"`
	WebApp          *WebAppInfo             `json:"web_app,omitempty"`
}

// KeyboardButtonPollType نوع نظرسنجی
type KeyboardButtonPollType struct {
	Type string `json:"type,omitempty"`
}

// ReplyKeyboardRemove حذف کیبورد
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective,omitempty"`
}

// ForceReply پاسخ اجباری
type ForceReply struct {
	ForceReply            bool   `json:"force_reply"`
	Selective             bool   `json:"selective,omitempty"`
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
}

// --- Helpers ---

func NewInlineKeyboard(rows ...[]InlineKeyboardButton) *InlineKeyboardMarkup {
	return &InlineKeyboardMarkup{InlineKeyboard: rows}
}

func InlineRow(btns ...InlineKeyboardButton) []InlineKeyboardButton {
	return btns
}

func Btn(text, callbackData string) InlineKeyboardButton {
	return InlineKeyboardButton{Text: text, CallbackData: callbackData}
}

func URLBtn(text, url string) InlineKeyboardButton {
	return InlineKeyboardButton{Text: text, URL: url}
}

func WebAppBtn(text string, url string) InlineKeyboardButton {
	return InlineKeyboardButton{Text: text, WebApp: &WebAppInfo{URL: url}}
}

func SwitchInlineBtn(text, query string) InlineKeyboardButton {
	return InlineKeyboardButton{Text: text, SwitchInlineQuery: query}
}

func NewReplyKeyboard(rows ...[]KeyboardButton) *ReplyKeyboardMarkup {
	return &ReplyKeyboardMarkup{Keyboard: rows}
}

func ReplyRow(btns ...KeyboardButton) []KeyboardButton {
	return btns
}

func KBtn(text string) KeyboardButton {
	return KeyboardButton{Text: text}
}

func ContactBtn(text string) KeyboardButton {
	return KeyboardButton{Text: text, RequestContact: true}
}

func LocationBtn(text string) KeyboardButton {
	return KeyboardButton{Text: text, RequestLocation: true}
}
