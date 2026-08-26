package retry

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/errors"
)

// Policy سیاست Retry
type Policy struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	Jitter       bool
	MaxTotalTime time.Duration // حداکثر زمان کل retry
}

// DefaultPolicy سیاست پیش‌فرض
func DefaultPolicy() *Policy {
	return &Policy{
		MaxAttempts:  3,
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
		MaxTotalTime: 60 * time.Second,
	}
}

// AggressivePolicy سیاست تهاجمی برای عملیات‌های مهم
func AggressivePolicy() *Policy {
	return &Policy{
		MaxAttempts:  5,
		InitialDelay: 200 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   1.5,
		Jitter:       true,
		MaxTotalTime: 30 * time.Second,
	}
}

// ConservativePolicy سیاست محتاطانه
func ConservativePolicy() *Policy {
	return &Policy{
		MaxAttempts:  2,
		InitialDelay: 1 * time.Second,
		MaxDelay:     20 * time.Second,
		Multiplier:   2.0,
		Jitter:       false,
		MaxTotalTime: 120 * time.Second,
	}
}

// Do اجرای تابع با Retry
func (p *Policy) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	var lastErr error
	startTime := time.Now()

	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		// بررسی max total time
		if p.MaxTotalTime > 0 && time.Since(startTime) > p.MaxTotalTime {
			if lastErr == nil {
				lastErr = errors.NewValidationError("retry", "max total time exceeded")
			}
			return lastErr
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}
		lastErr = err

		if !errors.IsRetryable(err) {
			return err
		}

		// اگر خطای 429 بود، از retry_after استفاده کن
		if apiErr, ok := errors.AsAPIError(err); ok && apiErr.RetryAfter > 0 {
			wait := time.Duration(apiErr.RetryAfter) * time.Second
			// اضافه کردن مقدار کمی buffer
			wait += 500 * time.Millisecond
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				return ctx.Err()
			}
			continue
		}

		if attempt == p.MaxAttempts {
			break
		}

		delay := p.calcDelay(attempt)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return lastErr
}

func (p *Policy) calcDelay(attempt int) time.Duration {
	delay := float64(p.InitialDelay) * math.Pow(p.Multiplier, float64(attempt-1))
	if delay > float64(p.MaxDelay) {
		delay = float64(p.MaxDelay)
	}
	if p.Jitter && delay > 0 {
		// Full jitter — بین 0 تا delay
		jitterAmount := rand.Float64() * delay * 0.25
		delay = delay + jitterAmount - (delay * 0.125)
	}
	return time.Duration(delay)
}
