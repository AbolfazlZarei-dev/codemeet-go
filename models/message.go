package models

import (
	"encoding/json"
	"strings"
)

type Message struct {
	MessageID           int                   `json:"message_id"`
	Date                int64                 `json:"date"`
	Chat                *Chat                 `json:"chat"`
	From                *User                 `json:"from,omitempty"`
	SenderChat          *Chat                 `json:"sender_chat,omitempty"`
	Text                string                `json:"text,omitempty"`
	Caption             string                `json:"caption,omitempty"`
	Entities            []MessageEntity       `json:"entities,omitempty"`
	CaptionEntities     []MessageEntity       `json:"caption_entities,omitempty"`
	ReplyToMessage      *Message              `json:"reply_to_message,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	MediaGroupID        string                `json:"media_group_id,omitempty"`
	AuthorSignature     string                `json:"author_signature,omitempty"`
	EditDate            int64                 `json:"edit_date,omitempty"`
	HasProtectedContent bool                  `json:"has_protected_content,omitempty"`

	NewChatMembers []User `json:"new_chat_members,omitempty"`
	LeftChatMember *User  `json:"left_chat_member,omitempty"`

	Photo     []PhotoSize `json:"photo,omitempty"`
	Video     *Video      `json:"video,omitempty"`
	Audio     *Audio      `json:"audio,omitempty"`
	Document  *Document   `json:"document,omitempty"`
	Animation *Animation  `json:"animation,omitempty"`
	Voice     *Voice      `json:"voice,omitempty"`
	VideoNote *VideoNote  `json:"video_note,omitempty"`
	Sticker   *Sticker    `json:"sticker,omitempty"`
	Contact   *Contact    `json:"contact,omitempty"`
	Location  *Location   `json:"location,omitempty"`
	Venue     *Venue      `json:"venue,omitempty"`
	Poll      *Poll       `json:"poll,omitempty"`
	Dice      *Dice       `json:"dice,omitempty"`
}

// HasMedia آیا پیام مدیا دارد
func (m *Message) HasMedia() bool {
	return m.Photo != nil || m.Video != nil || m.Audio != nil ||
		m.Document != nil || m.Animation != nil || m.Voice != nil ||
		m.VideoNote != nil || m.Sticker != nil || m.Contact != nil ||
		m.Location != nil || m.Venue != nil || m.Poll != nil || m.Dice != nil
}

// IsCommand آیا پیام یک دستور است
func (m *Message) IsCommand() bool {
	for _, e := range m.Entities {
		if e.Type == EntityBotCommand {
			return true
		}
	}
	return false
}

// CommandName نام دستور (بدون / و بدون @botname)
func (m *Message) CommandName() string {
	for _, e := range m.Entities {
		if e.Type == EntityBotCommand {
			if e.Offset == 0 {
				cmd := m.Text[e.Offset : e.Offset+e.Length]
				if len(cmd) > 0 && cmd[0] == '/' {
					cmd = cmd[1:]
				}
				// حذف @botname اگر وجود داشت
				if idx := strings.Index(cmd, "@"); idx != -1 {
					cmd = cmd[:idx]
				}
				return cmd
			}
		}
	}
	return ""
}

// CommandArgs آرگومان‌های دستور
func (m *Message) CommandArgs() string {
	if !m.IsCommand() {
		return ""
	}

	cmdName := "/" + m.CommandName()
	// اگر در متن @botname وجود داشت، آن را نادیده می‌گیریم
	if idx := strings.Index(m.Text, "@"); idx != -1 {
		// پیدا کردن فاصله بعد از @botname
		spaceIdx := strings.Index(m.Text[idx:], " ")
		if spaceIdx != -1 {
			return strings.TrimSpace(m.Text[idx+spaceIdx:])
		}
		return ""
	}

	if len(m.Text) > len(cmdName) {
		return strings.TrimSpace(m.Text[len(cmdName)+1:])
	}
	return ""
}

type MessageEntity struct {
	Type        string `json:"type"`
	Offset      int    `json:"offset"`
	Length      int    `json:"length"`
	URL         string `json:"url,omitempty"`
	User        *User  `json:"user,omitempty"`
	Language    string `json:"language,omitempty"`
	CustomEmoji string `json:"custom_emoji_id,omitempty"`
}

const (
	EntityBold          = "bold"
	EntityItalic        = "italic"
	EntityUnderline     = "underline"
	EntityStrikethrough = "strikethrough"
	EntitySpoiler       = "spoiler"
	EntityCode          = "code"
	EntityPre           = "pre"
	EntityTextLink      = "text_link"
	EntityTextMention   = "text_mention"
	EntityMention       = "mention"
	EntityHashtag       = "hashtag"
	EntityCashtag       = "cashtag"
	EntityBotCommand    = "bot_command"
	EntityURL           = "url"
	EntityEmail         = "email"
	EntityPhoneNumber   = "phone_number"
	EntityBlockquote    = "blockquote"
	EntityCustomEmoji   = "custom_emoji"
)

type ParseMode string

const (
	ParseModeHTML       ParseMode = "HTML"
	ParseModeMarkdown   ParseMode = "Markdown"
	ParseModeMarkdownV2 ParseMode = "MarkdownV2"
)

type ChatAction string

const (
	ActionTyping          ChatAction = "typing"
	ActionUploadPhoto     ChatAction = "upload_photo"
	ActionRecordVideo     ChatAction = "record_video"
	ActionUploadVideo     ChatAction = "upload_video"
	ActionRecordVoice     ChatAction = "record_voice"
	ActionUploadVoice     ChatAction = "upload_voice"
	ActionUploadDoc       ChatAction = "upload_document"
	ActionChooseSticker   ChatAction = "choose_sticker"
	ActionFindLocation    ChatAction = "find_location"
	ActionRecordVideoNote ChatAction = "record_video_note"
	ActionUploadVideoNote ChatAction = "upload_video_note"
)

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int64  `json:"file_size,omitempty"`
}

type Video struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

type Audio struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Duration     int        `json:"duration"`
	Performer    string     `json:"performer,omitempty"`
	Title        string     `json:"title,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
}

type Document struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

type Animation struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

type VideoNote struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Length       int        `json:"length"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	VCard       string `json:"vcard,omitempty"`
}

type Location struct {
	Longitude            float64 `json:"longitude"`
	Latitude             float64 `json:"latitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
}

type Venue struct {
	Location        *Location `json:"location"`
	Title           string    `json:"title"`
	Address         string    `json:"address"`
	FoursquareID    string    `json:"foursquare_id,omitempty"`
	FoursquareType  string    `json:"foursquare_type,omitempty"`
	GooglePlaceID   string    `json:"google_place_id,omitempty"`
	GooglePlaceType string    `json:"google_place_type,omitempty"`
}

type Poll struct {
	ID                    string       `json:"id"`
	Question              string       `json:"question"`
	Options               []PollOption `json:"options"`
	TotalVoterCount       int          `json:"total_voter_count"`
	IsClosed              bool         `json:"is_closed"`
	IsAnonymous           bool         `json:"is_anonymous"`
	Type                  string       `json:"type"`
	AllowsMultipleAnswers bool         `json:"allows_multiple_answers"`
	CorrectOptionID       int          `json:"correct_option_id,omitempty"`
	Explanation           string       `json:"explanation,omitempty"`
}

type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

type Dice struct {
	Emoji string `json:"emoji"`
	Value int    `json:"value"`
}

type Sticker struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Type         string     `json:"type"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	IsAnimated   bool       `json:"is_animated"`
	IsVideo      bool       `json:"is_video"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	Emoji        string     `json:"emoji,omitempty"`
	SetName      string     `json:"set_name,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

type StickerSet struct {
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	StickerType string     `json:"sticker_type"`
	IsAnimated  bool       `json:"is_animated"`
	IsVideo     bool       `json:"is_video"`
	Stickers    []Sticker  `json:"stickers"`
	Thumbnail   *PhotoSize `json:"thumbnail,omitempty"`
}

type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

type InputMedia struct {
	Type            string          `json:"type"`
	Media           string          `json:"media"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       string          `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	Thumbnail       string          `json:"thumbnail,omitempty"`
	Width           int             `json:"width,omitempty"`
	Height          int             `json:"height,omitempty"`
	Duration        int             `json:"duration,omitempty"`
	Streaming       bool            `json:"supports_streaming,omitempty"`
	Title           string          `json:"title,omitempty"`
	Performer       string          `json:"performer,omitempty"`
	HasSpoiler      bool            `json:"has_spoiler,omitempty"`
}

type ReactionType struct {
	Type        string `json:"type"`
	Emoji       string `json:"emoji,omitempty"`
	CustomEmoji string `json:"custom_emoji_id,omitempty"`
}

func (im InputMedia) MarshalJSON() ([]byte, error) {
	type alias InputMedia
	return json.Marshal(alias(im))
}
