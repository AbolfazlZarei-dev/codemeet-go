package methods

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/AbolfazlZarei-dev/codemeet-go/errors"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

type Media struct {
	parent *Messages
}

func newMedia(m *Messages) *Media { return &Media{parent: m} }

type SendPhotoRequest struct {
	ChatID              string           `json:"chat_id"`
	Photo               string           `json:"photo,omitempty"`
	Caption             string           `json:"caption,omitempty"`
	ParseMode           models.ParseMode `json:"parse_mode,omitempty"`
	ReplyToMessageID    int              `json:"reply_to_message_id,omitempty"`
	DisableNotification bool             `json:"disable_notification,omitempty"`
	ReplyMarkup         interface{}      `json:"reply_markup,omitempty"`
}

func (m *Media) SendPhoto(ctx context.Context, chatID, photo, caption string) (*models.Message, error) {
	return m.SendPhotoWithParams(ctx, &SendPhotoRequest{
		ChatID:  chatID,
		Photo:   photo,
		Caption: caption,
	})
}

func (m *Media) SendPhotoWithParams(ctx context.Context, req *SendPhotoRequest) (*models.Message, error) {
	if req.ChatID == "" {
		return nil, errors.NewValidationError("chat_id", "chat_id is required")
	}
	if req.Photo == "" {
		return nil, errors.NewValidationError("photo", "photo is required")
	}

	var result models.Message
	var err error

	if _, e := os.Stat(req.Photo); e == nil {
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
	} else if !os.IsNotExist(e) {
		// اگر خطای دسترسی یا چیز دیگر بود
		return nil, errors.NewValidationError("photo", "invalid file path or permission denied")
	} else {
		err = m.parent.doWithRetry(ctx, "sendPhoto", req, &result)
	}

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *Media) SendVideo(ctx context.Context, chatID, video, caption string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendVideo", "video", chatID, video, caption, 0, false, nil)
}

func (m *Media) SendDocument(ctx context.Context, chatID, doc, caption string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendDocument", "document", chatID, doc, caption, 0, false, nil)
}

func (m *Media) SendVoice(ctx context.Context, chatID, voice, caption string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendVoice", "voice", chatID, voice, caption, 0, false, nil)
}

func (m *Media) SendAudio(ctx context.Context, chatID, audio, caption string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendAudio", "audio", chatID, audio, caption, 0, false, nil)
}

func (m *Media) SendAnimation(ctx context.Context, chatID, animation, caption string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendAnimation", "animation", chatID, animation, caption, 0, false, nil)
}

func (m *Media) SendSticker(ctx context.Context, chatID, sticker string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendSticker", "sticker", chatID, sticker, "", 0, false, nil)
}

func (m *Media) SendVideoNote(ctx context.Context, chatID, videoNote string) (*models.Message, error) {
	return m.sendMediaFile(ctx, "sendVideoNote", "video_note", chatID, videoNote, "", 0, false, nil)
}

func (m *Media) sendMediaFile(ctx context.Context, method, fieldName, chatID, filePath, caption string, replyToID int, disableNotif bool, replyMarkup interface{}) (*models.Message, error) {
	if chatID == "" {
		return nil, errors.NewValidationError("chat_id", "chat_id is required")
	}
	if filePath == "" {
		return nil, errors.NewValidationError(fieldName, fieldName+" is required")
	}

	var result models.Message
	var err error

	if _, e := os.Stat(filePath); e == nil {
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
	} else if !os.IsNotExist(e) {
		return nil, errors.NewValidationError(fieldName, "invalid file path or permission denied")
	} else {
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

	if err != nil {
		return nil, err
	}
	return &result, nil
}

type SendMediaGroupRequest struct {
	ChatID              string              `json:"chat_id"`
	Media               []models.InputMedia `json:"media"`
	DisableNotification bool                `json:"disable_notification,omitempty"`
	ReplyToMessageID    int                 `json:"reply_to_message_id,omitempty"`
}

func (m *Media) SendMediaGroup(ctx context.Context, chatID string, media []models.InputMedia) ([]models.Message, error) {
	var messages []models.Message
	err := m.parent.doWithRetry(ctx, "sendMediaGroup", SendMediaGroupRequest{
		ChatID: chatID,
		Media:  media,
	}, &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

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

func (m *Media) SendLocation(ctx context.Context, chatID string, lat, lng float64) (*models.Message, error) {
	var result models.Message
	err := m.parent.doWithRetry(ctx, "sendLocation", &SendLocationRequest{
		ChatID: chatID, Latitude: lat, Longitude: lng,
	}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

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

func (m *Media) SendVenue(ctx context.Context, chatID string, lat, lng float64, title, address string) (*models.Message, error) {
	var result models.Message
	err := m.parent.doWithRetry(ctx, "sendVenue", &SendVenueRequest{
		ChatID: chatID, Latitude: lat, Longitude: lng, Title: title, Address: address,
	}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

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

func (m *Media) SendContact(ctx context.Context, chatID, phone, firstName, lastName string) (*models.Message, error) {
	var result models.Message
	err := m.parent.doWithRetry(ctx, "sendContact", &SendContactRequest{
		ChatID: chatID, PhoneNumber: phone, FirstName: firstName, LastName: lastName,
	}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

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

func (m *Media) SendPoll(ctx context.Context, chatID, question string, options []string) (*models.Message, error) {
	var result models.Message
	err := m.parent.doWithRetry(ctx, "sendPoll", &SendPollRequest{
		ChatID: chatID, Question: question, Options: options,
	}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *Media) SendDice(ctx context.Context, chatID, emoji string) (*models.Message, error) {
	var result models.Message
	err := m.parent.doWithRetry(ctx, "sendDice", map[string]string{
		"chat_id": chatID, "emoji": emoji,
	}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *Media) GetStickerSet(ctx context.Context, name string) (*models.StickerSet, error) {
	var ss models.StickerSet
	err := m.parent.doWithRetry(ctx, "getStickerSet", map[string]string{"name": name}, &ss)
	if err != nil {
		return nil, err
	}
	return &ss, nil
}

func (m *Media) UploadStickerFile(ctx context.Context, userID, stickerPath string) (*models.File, error) {
	var f models.File
	err := m.parent.doWithRetryMultipart(ctx, "uploadStickerFile",
		map[string]string{"user_id": userID},
		map[string]string{"sticker": stickerPath}, &f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (m *Media) GetFile(ctx context.Context, fileID string) (*models.File, error) {
	var f models.File
	err := m.parent.doWithRetry(ctx, "getFile", map[string]string{"file_id": fileID}, &f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (m *Media) DownloadFile(ctx context.Context, filePath string) ([]byte, error) {
	reader, err := m.parent.api.DownloadFile(ctx, filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// محدودیت 20MB برای جلوگیری از OOM
	limited := io.LimitReader(reader, 20<<20)
	return io.ReadAll(limited)
}

func (m *Media) SetMessageReaction(ctx context.Context, chatID string, messageID int, reaction []models.ReactionType) error {
	return m.parent.doWithRetry(ctx, "setMessageReaction", map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"reaction":   reaction,
	}, nil)
}
