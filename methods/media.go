package methods

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

// Media متدهای رسانه
type Media struct {
	parent *Messages
}

func newMedia(m *Messages) *Media { return &Media{parent: m} }

// SendPhotoRequest پارامترهای ارسال عکس
type SendPhotoRequest struct {
	ChatID              string           `json:"chat_id"`
	Photo               string           `json:"photo,omitempty"`
	Caption             string           `json:"caption,omitempty"`
	ParseMode           models.ParseMode `json:"parse_mode,omitempty"`
	ReplyToMessageID    int              `json:"reply_to_message_id,omitempty"`
	DisableNotification bool             `json:"disable_notification,omitempty"`
	ReplyMarkup         interface{}      `json:"reply_markup,omitempty"`
}

// SendPhoto ارسال تصویر (مسیر فایل یا media_id)
func (m *Media) SendPhoto(ctx context.Context, chatID, photo, caption string) (*models.Message, error) {
	return m.SendPhotoWithParams(ctx, &SendPhotoRequest{
		ChatID:  chatID,
		Photo:   photo,
		Caption: caption,
	})
}

// SendPhotoWithParams ارسال عکس با پارامترهای کامل
func (m *Media) SendPhotoWithParams(ctx context.Context, req *SendPhotoRequest) (*models.Message, error) {
	var result models.Message
	var err error

	if req.Photo != "" {
		if _, e := os.Stat(req.Photo); e == nil {
			// آپلود فایل محلی از طریق Multipart
			fields := map[string]string{
				"chat_id": req.ChatID,
				"caption": req.Caption,
			}
			if req.ParseMode != "" {
				fields["parse_mode"] = string(req.ParseMode)
			}
			if req.ReplyToMessageID != 0 {
				fields["reply_to_message_id"] = strconv.Itoa(req.ReplyToMessageID)
			}
			if req.DisableNotification {
				fields["disable_notification"] = "true"
			}
			if req.ReplyMarkup != nil {
				markupBytes, _ := json.Marshal(req.ReplyMarkup)
				fields["reply_markup"] = string(markupBytes)
			}

			files := map[string]string{"photo": req.Photo}
			err = m.parent.doWithRetryMultipart(ctx, "sendPhoto", fields, files, &result)
		} else {
			// ارسال با media_id از طریق JSON
			err = m.parent.doWithRetry(ctx, "sendPhoto", req, &result)
		}
	}
	return &result, err
}

// SendVideo ارسال ویدیو
func (m *Media) SendVideo(ctx context.Context, chatID, video, caption string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendVideo", "video", chatID, video, caption, 0, false, nil)
}

// SendDocument ارسال فایل
func (m *Media) SendDocument(ctx context.Context, chatID, doc, caption string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendDocument", "document", chatID, doc, caption, 0, false, nil)
}

// SendVoice ارسال وویس
func (m *Media) SendVoice(ctx context.Context, chatID, voice, caption string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendVoice", "voice", chatID, voice, caption, 0, false, nil)
}

// SendAudio ارسال صوت
func (m *Media) SendAudio(ctx context.Context, chatID, audio, caption string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendAudio", "audio", chatID, audio, caption, 0, false, nil)
}

// SendAnimation ارسال گیف/انیمیشن
func (m *Media) SendAnimation(ctx context.Context, chatID, animation, caption string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendAnimation", "animation", chatID, animation, caption, 0, false, nil)
}

// SendSticker ارسال استیکر
func (m *Media) SendSticker(ctx context.Context, chatID, sticker string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendSticker", "sticker", chatID, sticker, "", 0, false, nil)
}

// SendVideoNote ارسال ویژری ویدیو گرد
func (m *Media) SendVideoNote(ctx context.Context, chatID, videoNote string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendVideoNote", "video_note", chatID, videoNote, "", 0, false, nil)
}

// sendMediaFile متد کمکی برای ارسال مدیا
func (m *Media) sendMediaFile(ctx context.Context, method, fieldName, chatID, filePath, caption string, replyToID int, disableNotif bool, replyMarkup interface{}) (*models.Message, error) {
	var result models.Message
	var err error

	if _, e := os.Stat(filePath); e == nil {
		// آپلود فایل محلی از طریق Multipart
		fields := map[string]string{
			"chat_id": chatID,
			"caption": caption,
		}
		if replyToID != 0 {
			fields["reply_to_message_id"] = strconv.Itoa(replyToID)
		}
		if disableNotif {
			fields["disable_notification"] = "true"
		}
		if replyMarkup != nil {
			markupBytes, _ := json.Marshal(replyMarkup)
			fields["reply_markup"] = string(markupBytes)
		}

		files := map[string]string{fieldName: filePath}
		err = m.parent.doWithRetryMultipart(ctx, method, fields, files, &result)
	} else {
		// فایل محلی نیست — استفاده از media_id
		params := map[string]interface{}{
			"chat_id": chatID,
			fieldName: filePath,
		}
		if caption != "" {
			params["caption"] = caption
		}
		if replyToID != 0 {
			params["reply_to_message_id"] = replyToID
		}
		if disableNotif {
			params["disable_notification"] = true
		}
		if replyMarkup != nil {
			params["reply_markup"] = replyMarkup
		}
		err = m.parent.doWithRetry(ctx, method, params, &result)
	}
	return &result, err
}

// SendMediaGroupRequest پارامترهای ارسال گروه مدیا
type SendMediaGroupRequest struct {
	ChatID              string              `json:"chat_id"`
	Media               []models.InputMedia `json:"media"`
	DisableNotification bool                `json:"disable_notification,omitempty"`
	ReplyToMessageID    int                 `json:"reply_to_message_id,omitempty"`
}

// SendMediaGroup ارسال گروهی چند مدیا
func (m *Media) SendMediaGroup(ctx context.Context, chatID string, media []models.InputMedia) ([]models.Message, error) {
	var messages []models.Message
	err := m.parent.doWithRetry(ctx, "sendMediaGroup", SendMediaGroupRequest{
		ChatID: chatID,
		Media:  media,
	}, &messages)
	return messages, err
}

// SendLocationRequest پارامترهای ارسال موقعیت
type SendLocationRequest struct {
	ChatID               string      `json:"chat_id"`
	Latitude             float64     `json:"latitude"`
	Longitude            float64     `json:"longitude"`
	LivePeriod           int         `json:"live_period,omitempty"`
	Heading              int         `json:"heading,omitempty"`
	ProximityAlertRadius int         `json:"proximity_alert_radius,omitempty"`
	DisableNotification  bool        `json:"disable_notification,omitempty"`
	ReplyToMessageID     int         `json:"reply_to_message_id,omitempty"`
	ReplyMarkup          interface{} `json:"reply_markup,omitempty"`
}

// SendLocation ارسال موقعیت مکانی
func (m *Media) SendLocation(ctx context.Context, chatID string, lat, lng float64) (*models.Message, error) {
	var result models.Message
	err := m.parent.doWithRetry(ctx, "sendLocation", &SendLocationRequest{
		ChatID: chatID, Latitude: lat, Longitude: lng,
	}, &result)
	return &result, err
}

// SendVenueRequest پارامترهای ارسال مکان
type SendVenueRequest struct {
	ChatID              string      `json:"chat_id"`
	Latitude            float64     `json:"latitude"`
	Longitude           float64     `json:"longitude"`
	Title               string      `json:"title"`
	Address             string      `json:"address"`
	FoursquareID        string      `json:"foursquare_id,omitempty"`
	DisableNotification bool        `json:"disable_notification,omitempty"`
	ReplyToMessageID    int         `json:"reply_to_message_id,omitempty"`
	ReplyMarkup         interface{} `json:"reply_markup,omitempty"`
}

// SendVenue ارسال مکان روی نقشه
func (m *Media) SendVenue(ctx context.Context, chatID string, lat, lng float64, title, address string) (*models.Message, error) {
	var result models.Message
	err := m.parent.doWithRetry(ctx, "sendVenue", &SendVenueRequest{
		ChatID: chatID, Latitude: lat, Longitude: lng, Title: title, Address: address,
	}, &result)
	return &result, err
}

// SendContactRequest پارامترهای ارسال مخاطب
type SendContactRequest struct {
	ChatID              string      `json:"chat_id"`
	PhoneNumber         string      `json:"phone_number"`
	FirstName           string      `json:"first_name"`
	LastName            string      `json:"last_name,omitempty"`
	VCard               string      `json:"vcard,omitempty"`
	DisableNotification bool        `json:"disable_notification,omitempty"`
	ReplyToMessageID    int         `json:"reply_to_message_id,omitempty"`
	ReplyMarkup         interface{} `json:"reply_markup,omitempty"`
}

// SendContact ارسال مخاطب
func (m *Media) SendContact(ctx context.Context, chatID, phone, firstName, lastName string) (*models.Message, error) {
	var result models.Message
	err := m.parent.doWithRetry(ctx, "sendContact", &SendContactRequest{
		ChatID: chatID, PhoneNumber: phone, FirstName: firstName, LastName: lastName,
	}, &result)
	return &result, err
}

// SendPollRequest پارامترهای ارسال نظرسنجی
type SendPollRequest struct {
	ChatID                string      `json:"chat_id"`
	Question              string      `json:"question"`
	Options               []string    `json:"options"`
	IsAnonymous           bool        `json:"is_anonymous,omitempty"`
	Type                  string      `json:"type,omitempty"`
	AllowsMultipleAnswers bool        `json:"allows_multiple_answers,omitempty"`
	CorrectOptionID       int         `json:"correct_option_id,omitempty"`
	Explanation           string      `json:"explanation,omitempty"`
	OpenPeriod            int         `json:"open_period,omitempty"`
	CloseDate             int64       `json:"close_date,omitempty"`
	IsClosed              bool        `json:"is_closed,omitempty"`
	DisableNotification   bool        `json:"disable_notification,omitempty"`
	ReplyToMessageID      int         `json:"reply_to_message_id,omitempty"`
	ReplyMarkup           interface{} `json:"reply_markup,omitempty"`
}

// SendPoll ارسال نظرسنجی
func (m *Media) SendPoll(ctx context.Context, chatID, question string, options []string) (*models.Message, error) {
	var result models.Message
	err := m.parent.doWithRetry(ctx, "sendPoll", &SendPollRequest{
		ChatID: chatID, Question: question, Options: options,
	}, &result)
	return &result, err
}

// SendDice ارسال تاس
func (m *Media) SendDice(ctx context.Context, chatID, emoji string) (*models.Message, error) {
	var result models.Message
	err := m.parent.doWithRetry(ctx, "sendDice", map[string]string{
		"chat_id": chatID, "emoji": emoji,
	}, &result)
	return &result, err
}

// --- استیکرها ---

// GetStickerSet دریافت مجموعه استیکر
func (m *Media) GetStickerSet(ctx context.Context, name string) (*models.StickerSet, error) {
	var ss models.StickerSet
	err := m.parent.doWithRetry(ctx, "getStickerSet", map[string]string{"name": name}, &ss)
	return &ss, err
}

// UploadStickerFile آپلود فایل استیکر
func (m *Media) UploadStickerFile(ctx context.Context, userID, stickerPath string) (*models.File, error) {
	var f models.File
	err := m.parent.doWithRetryMultipart(ctx, "uploadStickerFile",
		map[string]string{"user_id": userID},
		map[string]string{"sticker": stickerPath}, &f)
	return &f, err
}

// --- فایل‌ها ---

// GetFile دریافت اطلاعات فایل
func (m *Media) GetFile(ctx context.Context, fileID string) (*models.File, error) {
	var f models.File
	err := m.parent.doWithRetry(ctx, "getFile", map[string]string{"file_id": fileID}, &f)
	return &f, err
}

// DownloadFile دانلود فایل از سرور کدمیت
func (m *Media) DownloadFile(ctx context.Context, filePath string) ([]byte, error) {
	// استفاده از متد DownloadFile کلاینت API برای دریافت Stream
	reader, err := m.parent.api.DownloadFile(ctx, filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// --- مدیریت واکنش‌ها ---

// SetMessageReaction تنظیم واکنش به پیام
func (m *Media) SetMessageReaction(ctx context.Context, chatID string, messageID int, reaction []models.ReactionType) error {
	return m.parent.doWithRetry(ctx, "setMessageReaction", map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"reaction":   reaction,
	}, nil)
}
