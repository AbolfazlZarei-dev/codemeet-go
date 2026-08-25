package methods

import (
	"context"

	"github.com/AbolfazlZarei-dev/codemeet-go/api"
	"github.com/AbolfazlZarei-dev/codemeet-go/errors"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
	"github.com/AbolfazlZarei-dev/codemeet-go/ratelimit"
	"github.com/AbolfazlZarei-dev/codemeet-go/retry"
)

// Messages متدهای ارسال و مدیریت پیام
type Messages struct {
	api     *api.Client
	retry   *retry.Policy
	limiter *ratelimit.Limiter
}

func newMessages(c *api.Client, r *retry.Policy, l *ratelimit.Limiter) *Messages {
	return &Messages{api: c, retry: r, limiter: l}
}

// SendMessageRequest درخواست ارسال پیام
type SendMessageRequest struct {
	ChatID              string                 `json:"chat_id"`
	Text                string                 `json:"text"`
	ParseMode           models.ParseMode       `json:"parse_mode,omitempty"`
	Entities            []models.MessageEntity `json:"entities,omitempty"`
	ReplyToMessageID    int                    `json:"reply_to_message_id,omitempty"`
	DisableNotification bool                   `json:"disable_notification,omitempty"`
	ProtectContent      bool                   `json:"protect_content,omitempty"`
	ReplyMarkup         interface{}            `json:"reply_markup,omitempty"`
}

// SendMessage ارسال پیام متنی
func (m *Messages) Send(ctx context.Context, req *SendMessageRequest) (*models.Message, error) {
	var result models.Message
	err := m.doWithRetry(ctx, "sendMessage", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SendText ارسال متن ساده
func (m *Messages) SendText(ctx context.Context, chatID, text string) (*models.Message, error) {
	return m.Send(ctx, &SendMessageRequest{ChatID: chatID, Text: text})
}

// SendHTML ارسال متن با HTML
func (m *Messages) SendHTML(ctx context.Context, chatID, text string) (*models.Message, error) {
	return m.Send(ctx, &SendMessageRequest{
		ChatID: chatID, Text: text, ParseMode: models.ParseModeHTML,
	})
}

// SendMarkdown ارسال متن با MarkdownV2
func (m *Messages) SendMarkdown(ctx context.Context, chatID, text string) (*models.Message, error) {
	return m.Send(ctx, &SendMessageRequest{
		ChatID: chatID, Text: text, ParseMode: models.ParseModeMarkdownV2,
	})
}

// SendWithKeyboard ارسال با کیبورد شیشه‌ای
func (m *Messages) SendWithKeyboard(ctx context.Context, chatID, text string, markup interface{}) (*models.Message, error) {
	return m.Send(ctx, &SendMessageRequest{
		ChatID: chatID, Text: text, ReplyMarkup: markup,
	})
}

// ForwardMessage فوروارد پیام
func (m *Messages) Forward(ctx context.Context, chatID, fromChatID string, messageID int) (*models.Message, error) {
	var result models.Message
	err := m.doWithRetry(ctx, "forwardMessage", map[string]interface{}{
		"chat_id": chatID, "from_chat_id": fromChatID, "message_id": messageID,
	}, &result)
	return &result, err
}

// CopyMessage کپی پیام
func (m *Messages) Copy(ctx context.Context, chatID, fromChatID string, messageID int, caption string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := m.doWithRetry(ctx, "copyMessage", map[string]interface{}{
		"chat_id": chatID, "from_chat_id": fromChatID, "message_id": messageID, "caption": caption,
	}, &result)
	return result, err
}

// EditText ویرایش متن پیام (اصلاح شده با ParseMode)
func (m *Messages) EditText(ctx context.Context, chatID string, messageID int, text string, parseMode models.ParseMode, replyMarkup *models.InlineKeyboardMarkup) error {
	return m.doWithRetry(ctx, "editMessageText", map[string]interface{}{
		"chat_id":      chatID,
		"message_id":   messageID,
		"text":         text,
		"parse_mode":   parseMode,
		"reply_markup": replyMarkup,
	}, nil)
}

// EditTextInline ویرایش پیام inline (اصلاح شده با ParseMode)
func (m *Messages) EditTextInline(ctx context.Context, inlineMessageID, text string, parseMode models.ParseMode, replyMarkup *models.InlineKeyboardMarkup) error {
	return m.doWithRetry(ctx, "editMessageText", map[string]interface{}{
		"inline_message_id": inlineMessageID,
		"text":              text,
		"parse_mode":        parseMode,
		"reply_markup":      replyMarkup,
	}, nil)
}

// EditCaption ویرایش کپشن (اصلاح شده با ParseMode)
func (m *Messages) EditCaption(ctx context.Context, chatID string, messageID int, caption string, parseMode models.ParseMode, replyMarkup *models.InlineKeyboardMarkup) error {
	return m.doWithRetry(ctx, "editMessageCaption", map[string]interface{}{
		"chat_id":      chatID,
		"message_id":   messageID,
		"caption":      caption,
		"parse_mode":   parseMode,
		"reply_markup": replyMarkup,
	}, nil)
}

// EditReplyMarkup ویرایش کیبورد
func (m *Messages) EditReplyMarkup(ctx context.Context, chatID string, messageID int, markup *models.InlineKeyboardMarkup) error {
	return m.doWithRetry(ctx, "editMessageReplyMarkup", map[string]interface{}{
		"chat_id": chatID, "message_id": messageID, "reply_markup": markup,
	}, nil)
}

// Delete حذف پیام
func (m *Messages) Delete(ctx context.Context, chatID string, messageID int) error {
	return m.doWithRetry(ctx, "deleteMessage", map[string]interface{}{
		"chat_id": chatID, "message_id": messageID,
	}, nil)
}

// DeleteMessages حذف چند پیام (batch)
func (m *Messages) DeleteMessages(ctx context.Context, chatID string, messageIDs []int) error {
	return m.doWithRetry(ctx, "deleteMessages", map[string]interface{}{
		"chat_id": chatID, "message_ids": messageIDs,
	}, nil)
}

// SendChatAction ارسال اکشن تایپ و...
func (m *Messages) SendChatAction(ctx context.Context, chatID string, action models.ChatAction) error {
	return m.doWithRetry(ctx, "sendChatAction", map[string]interface{}{
		"chat_id": chatID, "action": action,
	}, nil)
}

// AnswerCallback پاسخ به callback query
func (m *Messages) AnswerCallback(ctx context.Context, req *models.AnswerCallbackRequest) error {
	return m.doWithRetry(ctx, "answerCallbackQuery", req, nil)
}

// AnswerCallbackSimple پاسخ ساده به callback
func (m *Messages) AnswerCallbackSimple(ctx context.Context, callbackID, text string, showAlert bool) error {
	return m.AnswerCallback(ctx, &models.AnswerCallbackRequest{
		CallbackQueryID: callbackID,
		Text:            text,
		ShowAlert:       showAlert,
	})
}

// SetMyDescription تنظیم توضیحات
func (m *Messages) SetMyDescription(ctx context.Context, desc, lang string) error {
	return m.doWithRetry(ctx, "setMyDescription", map[string]string{
		"description":   desc,
		"language_code": lang,
	}, nil)
}

// doWithRetry اجرای درخواست با Retry و Rate Limit
func (m *Messages) doWithRetry(ctx context.Context, method string, params interface{}, result interface{}) error {
	if m.limiter != nil {
		if err := m.limiter.Wait(ctx); err != nil {
			return err
		}
	}

	if m.retry == nil {
		resp, err := m.api.RequestWithParams(ctx, method, params)
		if err != nil {
			return err
		}
		if result != nil {
			return resp.Decode(result)
		}
		return nil
	}

	return m.retry.Do(ctx, func(ctx context.Context) error {
		resp, err := m.api.RequestWithParams(ctx, method, params)
		if err != nil {
			return err
		}
		if result != nil {
			return resp.Decode(result)
		}
		return nil
	})
}

// doWithRetryMultipart اجرای درخواست multipart با Retry
func (m *Messages) doWithRetryMultipart(ctx context.Context, method string, fields map[string]string, files map[string]string, result interface{}) error {
	if m.limiter != nil {
		if err := m.limiter.Wait(ctx); err != nil {
			return err
		}
	}

	doOnce := func(ctx context.Context) error {
		resp, err := m.api.RequestWithMultipart(ctx, method, fields, files)
		if err != nil {
			return err
		}
		if result != nil {
			return resp.Decode(result)
		}
		return nil
	}

	if m.retry == nil {
		return doOnce(ctx)
	}

	return m.retry.Do(ctx, doOnce)
}

// --- متدهای کمکی برای Admin ---

// ParseError parse خطای API
func (m *Messages) ParseError(err error) *errors.APIError {
	if e, ok := errors.AsAPIError(err); ok {
		return e
	}
	return nil
}
