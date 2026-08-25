package methods

import (
	"context"

	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

type WebhookMethods struct {
	parent *Messages
}

func newWebhookMethods(m *Messages) *WebhookMethods { return &WebhookMethods{parent: m} }

func (w *WebhookMethods) Set(ctx context.Context, req *models.SetWebhookRequest) error {
	return w.parent.doWithRetry(ctx, "setWebhook", req, nil)
}

func (w *WebhookMethods) GetInfo(ctx context.Context) (*models.WebhookInfo, error) {
	var info models.WebhookInfo
	err := w.parent.doWithRetry(ctx, "getWebhookInfo", nil, &info)
	return &info, err
}

func (w *WebhookMethods) Delete(ctx context.Context) error {
	return w.parent.doWithRetry(ctx, "deleteWebhook", map[string]bool{"drop_pending_updates": false}, nil)
}

// DeleteWithDrop حذف وب‌هوک همراه با پاک کردن آپدیت‌های معلق
func (w *WebhookMethods) DeleteWithDrop(ctx context.Context, drop bool) error {
	return w.parent.doWithRetry(ctx, "deleteWebhook", map[string]bool{"drop_pending_updates": drop}, nil)
}
