package broadcast

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/AbolfazlZarei-dev/codemeet-go/ratelimit"
)

type Config struct {
	SendAction func(ctx context.Context, chatID, text string) error
	// تابعی که لیست آیدی‌ها را بر اساس نوع (users, groups, channels, all) برمی‌گرداند
	GetTargetsAction func(ctx context.Context, targetType string) ([]string, error)
}

type Broadcaster struct {
	cfg     Config
	limiter *ratelimit.Limiter
}

func New(cfg Config) *Broadcaster {
	return &Broadcaster{
		cfg:     cfg,
		limiter: ratelimit.New(25), // 25 پیام در ثانیه
	}
}

// Send پیام را به دسته‌بندی مشخص شده ارسال می‌کند
// targetType می‌تواند باشد: "users", "groups", "channels", "all"
func (b *Broadcaster) Send(ctx context.Context, targetType, text string) (successCount, failCount int64, err error) {
	targets, err := b.cfg.GetTargetsAction(ctx, targetType)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get targets: %w", err)
	}

	var success, fail int64
	var wg sync.WaitGroup

	for _, chatID := range targets {
		select {
		case <-ctx.Done():
			break
		default:
		}

		wg.Add(1)
		go func(id string) {
			defer wg.Done()

			if err := b.limiter.Wait(ctx); err != nil {
				atomic.AddInt64(&fail, 1)
				return
			}

			if err := b.cfg.SendAction(ctx, id, text); err != nil {
				atomic.AddInt64(&fail, 1)
			} else {
				atomic.AddInt64(&success, 1)
			}
		}(chatID)
	}

	wg.Wait()
	return success, fail, nil
}
