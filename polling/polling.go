package polling

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/api"
	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/errors"
	"github.com/AbolfazlZarei-dev/codemeet-go/logger"
	"github.com/AbolfazlZarei-dev/codemeet-go/methods"
	"github.com/AbolfazlZarei-dev/codemeet-go/ratelimit"
	"github.com/AbolfazlZarei-dev/codemeet-go/retry"
)

type Config struct {
	Timeout            int
	PollInterval       time.Duration
	Limit              int
	AllowedUpdates     []string
	BufferSize         int
	DeleteWebhookFirst bool
	MaxRetries         int
}

func DefaultConfig() Config {
	return Config{
		Timeout:            10,
		PollInterval:       2 * time.Second,
		Limit:              100,
		BufferSize:         1000,
		DeleteWebhookFirst: true,
		MaxRetries:         5,
	}
}

type Poller struct {
	api        *api.Client
	dispatcher *dispatcher.Dispatcher
	logger     *logger.Logger
	cfg        Config
	offset     int
	updates    *methods.UpdatesMethods
	webhook    *methods.WebhookMethods
	retries    int
}

func New(c *api.Client, d *dispatcher.Dispatcher, log *logger.Logger, cfg Config) *Poller {
	defaultRetry := retry.DefaultPolicy()
	defaultLimit := ratelimit.New(30)
	m := methods.New(c, defaultRetry, defaultLimit)

	return &Poller{
		api:        c,
		dispatcher: d,
		logger:     log,
		cfg:        cfg,
		updates:    m.Updates(),
		webhook:    m.Webhook(),
	}
}

func (p *Poller) Start(ctx context.Context) error {
	if p.cfg.DeleteWebhookFirst {
		if err := p.webhook.DeleteWithDrop(ctx, true); err != nil {
			if p.logger != nil {
				p.logger.Warn("failed to delete existing webhook", "error", err)
			}
		} else if p.logger != nil {
			p.logger.Info("deleted existing webhook before polling")
		}
	}

	if p.logger != nil {
		p.logger.Info("starting long polling", "timeout", p.cfg.Timeout)
	}

	for {
		select {
		case <-ctx.Done():
			if p.logger != nil {
				p.logger.Info("polling stopped")
			}
			return ctx.Err()
		default:
		}

		updates, err := p.updates.Get(ctx, &methods.GetUpdatesParams{
			Offset:         p.offset,
			Limit:          p.cfg.Limit,
			Timeout:        p.cfg.Timeout,
			AllowedUpdates: p.cfg.AllowedUpdates,
		})
		if err != nil {
			p.retries++
			if p.logger != nil {
				p.logger.Error("polling error", "error", err, "retry", p.retries)
			}

			if apiErr, ok := errors.AsAPIError(err); ok && apiErr.Code == errors.CodeTooManyRequests {
				wait := time.Duration(apiErr.RetryAfter) * time.Second
				if wait < 2*time.Second {
					wait = 2 * time.Second
				}
				wait += 500 * time.Millisecond // buffer

				if p.logger != nil {
					p.logger.Warn("rate limit hit (429), waiting...", "wait", wait)
				}
				select {
				case <-time.After(wait):
				case <-ctx.Done():
					return ctx.Err()
				}
				p.retries = 0
				continue
			}

			backoff := time.Duration(p.retries) * time.Second
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
			}
			jitter := time.Duration(rand.Int63n(int64(backoff/2) + 1))
			select {
			case <-time.After(backoff + jitter):
			case <-ctx.Done():
				return ctx.Err()
			}

			if p.retries > p.cfg.MaxRetries {
				return fmt.Errorf("max retries (%d) exceeded: %w", p.cfg.MaxRetries, err)
			}
			continue
		}

		p.retries = 0

		for _, update := range updates {
			p.offset = update.UpdateID + 1
			p.dispatcher.Dispatch(ctx, &update)
		}

		if len(updates) == 0 {
			select {
			case <-time.After(p.cfg.PollInterval):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func (p *Poller) Offset() int  { return p.offset }
func (p *Poller) ResetOffset() { p.offset = 0 }
