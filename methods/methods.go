package methods

import (
	"github.com/AbolfazlZarei-dev/codemeet-go/api"
	"github.com/AbolfazlZarei-dev/codemeet-go/ratelimit"
	"github.com/AbolfazlZarei-dev/codemeet-go/retry"
)

// Methods تمام متدهای API در یک جا
type Methods struct {
	api      *api.Client
	messages *Messages
	media    *Media
	bot      *BotMethods
	chat     *ChatMethods
	webhook  *WebhookMethods
	updates  *UpdatesMethods
}

func New(c *api.Client, r *retry.Policy, l *ratelimit.Limiter) *Methods {
	m := newMessages(c, r, l)
	return &Methods{
		api:      c,
		messages: m,
		media:    newMedia(m),
		bot:      newBotMethods(m),
		chat:     newChatMethods(m),
		webhook:  newWebhookMethods(m),
		updates:  newUpdatesMethods(m),
	}
}

func (m *Methods) Messages() *Messages      { return m.messages }
func (m *Methods) Media() *Media            { return m.media }
func (m *Methods) Bot() *BotMethods         { return m.bot }
func (m *Methods) Chat() *ChatMethods       { return m.chat }
func (m *Methods) Webhook() *WebhookMethods { return m.webhook }
func (m *Methods) Updates() *UpdatesMethods { return m.updates }

// Client دسترسی مستقیم به API Client (پیشرفته)
func (m *Methods) Client() *api.Client { return m.api }
