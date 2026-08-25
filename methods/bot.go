package methods

import (
	"context"

	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

type BotMethods struct {
	parent *Messages
}

func newBotMethods(m *Messages) *BotMethods { return &BotMethods{parent: m} }

func (b *BotMethods) GetMe(ctx context.Context) (*models.User, error) {
	var u models.User
	err := b.parent.doWithRetry(ctx, "getMe", nil, &u)
	return &u, err
}

func (b *BotMethods) SetName(ctx context.Context, name string) error {
	return b.parent.doWithRetry(ctx, "setMyName", map[string]string{"name": name}, nil)
}

func (b *BotMethods) GetName(ctx context.Context) (string, error) {
	var r models.BotName
	err := b.parent.doWithRetry(ctx, "getMyName", nil, &r)
	return r.Name, err
}

func (b *BotMethods) GetNameWithLang(ctx context.Context, lang string) (string, error) {
	var r models.BotName
	err := b.parent.doWithRetry(ctx, "getMyName", map[string]string{"language_code": lang}, &r)
	return r.Name, err
}

func (b *BotMethods) SetDescription(ctx context.Context, desc string) error {
	return b.parent.doWithRetry(ctx, "setMyDescription", map[string]string{"description": desc}, nil)
}

func (b *BotMethods) SetDescriptionWithLang(ctx context.Context, desc, lang string) error {
	return b.parent.doWithRetry(ctx, "setMyDescription", map[string]string{
		"description":   desc,
		"language_code": lang,
	}, nil)
}

func (b *BotMethods) GetDescription(ctx context.Context) (string, error) {
	var r models.BotDescription
	err := b.parent.doWithRetry(ctx, "getMyDescription", nil, &r)
	return r.Description, err
}

func (b *BotMethods) GetDescriptionWithLang(ctx context.Context, lang string) (string, error) {
	var r models.BotDescription
	err := b.parent.doWithRetry(ctx, "getMyDescription", map[string]string{"language_code": lang}, &r)
	return r.Description, err
}

func (b *BotMethods) SetShortDescription(ctx context.Context, desc string) error {
	return b.parent.doWithRetry(ctx, "setMyShortDescription", map[string]string{"short_description": desc}, nil)
}

func (b *BotMethods) SetShortDescriptionWithLang(ctx context.Context, desc, lang string) error {
	return b.parent.doWithRetry(ctx, "setMyShortDescription", map[string]string{
		"short_description": desc,
		"language_code":     lang,
	}, nil)
}

func (b *BotMethods) GetShortDescription(ctx context.Context) (string, error) {
	var r models.BotShortDescription
	err := b.parent.doWithRetry(ctx, "getMyShortDescription", nil, &r)
	return r.ShortDescription, err
}

func (b *BotMethods) GetShortDescriptionWithLang(ctx context.Context, lang string) (string, error) {
	var r models.BotShortDescription
	err := b.parent.doWithRetry(ctx, "getMyShortDescription", map[string]string{"language_code": lang}, &r)
	return r.ShortDescription, err
}

func (b *BotMethods) SetCommands(ctx context.Context, cmds []models.BotCommand, lang string) error {
	return b.parent.doWithRetry(ctx, "setMyCommands", models.SetCommandsRequest{Commands: cmds, LanguageCode: lang}, nil)
}

func (b *BotMethods) GetCommands(ctx context.Context, lang string) ([]models.BotCommand, error) {
	var cmds []models.BotCommand
	err := b.parent.doWithRetry(ctx, "getMyCommands", map[string]string{"language_code": lang}, &cmds)
	return cmds, err
}

func (b *BotMethods) DeleteCommands(ctx context.Context, lang string) error {
	return b.parent.doWithRetry(ctx, "deleteMyCommands", map[string]string{"language_code": lang}, nil)
}

// LogOut خروج ربات از سرور کدمیت (برای انتقال به سرور محلی)
func (b *BotMethods) LogOut(ctx context.Context) (bool, error) {
	resp, err := b.parent.api.RequestWithParams(ctx, "logOut", nil)
	if err != nil {
		return false, err
	}
	return resp.AsBool()
}

// Close بستن سرور بات (پیش از انتقال)
func (b *BotMethods) Close(ctx context.Context) (bool, error) {
	resp, err := b.parent.api.RequestWithParams(ctx, "close", nil)
	if err != nil {
		return false, err
	}
	return resp.AsBool()
}
