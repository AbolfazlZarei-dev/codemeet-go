package middleware

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/logger"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

// Middleware نوع میدل‌ور
type Middleware = func(next HandlerFunc) HandlerFunc

// HandlerFunc امضای هندلر
type HandlerFunc func(ctx context.Context, update *models.Update)

// Recovery میدل‌ور بازیابی panic
func Recovery(log *logger.Logger) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			defer func() {
				if r := recover(); r != nil {
					if log != nil {
						log.Error("panic recovered",
							"error", fmt.Sprintf("%v", r),
							"stack", string(debug.Stack()),
							"update_id", u.UpdateID,
						)
					}
				}
			}()
			next(ctx, u)
		}
	}
}

// Logging میدل‌ور لاگ‌گیری
func Logging(log *logger.Logger) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			start := time.Now()
			if log != nil {
				log.Info("processing update",
					"update_id", u.UpdateID,
					"type", u.Type(),
				)
			}
			next(ctx, u)
			if log != nil {
				log.Info("update processed",
					"update_id", u.UpdateID,
					"duration", time.Since(start),
				)
			}
		}
	}
}

// RateLimit میدل‌ور محدودیت نرخ هر کاربر
// بهینه با sharded map برای کاهش lock contention
func RateLimit(perUser int, window time.Duration) Middleware {
	type counter struct {
		count int
		reset time.Time
	}
	const shards = 64
	var (
		shardedUsers [shards]map[string]*counter
		shardedMu    [shards]sync.Mutex
	)
	for i := range shardedUsers {
		shardedUsers[i] = make(map[string]*counter)
	}

	hash := func(s string) int {
		h := uint32(2166136261)
		for i := 0; i < len(s); i++ {
			h *= 16777619
			h ^= uint32(s[i])
		}
		return int(h) % shards
	}

	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			var userID string
			if u.Message != nil && u.Message.From != nil {
				userID = u.Message.From.ID
			} else if u.CallbackQuery != nil {
				userID = u.CallbackQuery.From.ID
			}
			if userID == "" {
				next(ctx, u)
				return
			}

			idx := hash(userID)
			shardedMu[idx].Lock()
			c, ok := shardedUsers[idx][userID]
			if !ok || time.Now().After(c.reset) {
				c = &counter{count: 0, reset: time.Now().Add(window)}
				shardedUsers[idx][userID] = c
			}
			c.count++
			allowed := c.count <= perUser
			shardedMu[idx].Unlock()

			if !allowed {
				return
			}
			next(ctx, u)
		}
	}
}

// Metrics میدل‌ور متریک
func Metrics(counter *MetricsCounter) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			counter.Inc(u.Type())
			next(ctx, u)
		}
	}
}

// MetricsCounter شمارنده متریک — بهینه با atomic
type MetricsCounter struct {
	counts sync.Map // map[string]*int64
}

func NewMetricsCounter() *MetricsCounter {
	return &MetricsCounter{}
}

// Inc افزایش شمارنده با استفاده از atomic استاندارد
func (m *MetricsCounter) Inc(updateType string) {
	v, _ := m.counts.LoadOrStore(updateType, new(int64))
	atomic.AddInt64(v.(*int64), 1)
}

// Snapshot گرفتن کپی آمار با استفاده از atomic استاندارد
func (m *MetricsCounter) Snapshot() map[string]int64 {
	out := make(map[string]int64)
	m.counts.Range(func(k, v interface{}) bool {
		out[k.(string)] = atomic.LoadInt64(v.(*int64))
		return true
	})
	return out
}

// Timeout میدل‌ور تایم‌اوت برای هر هندلر
func Timeout(d time.Duration) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			ctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()
			done := make(chan struct{})
			go func() {
				defer close(done)
				next(ctx, u)
			}()
			select {
			case <-done:
			case <-ctx.Done():
			}
		}
	}
}

// BotOnly فقط بات‌ها اجازه عبور داشته باشند
func BotOnly() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			if u.Message != nil && u.Message.From != nil && u.Message.From.IsBot {
				next(ctx, u)
			}
		}
	}
}

// UserOnly فقط کاربران عادی
func UserOnly() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			if u.Message != nil && u.Message.From != nil && !u.Message.From.IsBot {
				next(ctx, u)
			} else if u.CallbackQuery != nil && !u.CallbackQuery.From.IsBot {
				next(ctx, u)
			}
		}
	}
}

// AdminOnly فقط ادمین‌ها (با تابع چک)
func AdminOnly(isAdmin func(userID string) bool) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			var userID string
			if u.Message != nil && u.Message.From != nil {
				userID = u.Message.From.ID
			} else if u.CallbackQuery != nil {
				userID = u.CallbackQuery.From.ID
			}
			if userID != "" && isAdmin(userID) {
				next(ctx, u)
			}
		}
	}
}

// Blacklist میدل‌ور لیست سیاه
func Blacklist(blacklist func(userID string) bool) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			var userID string
			if u.Message != nil && u.Message.From != nil {
				userID = u.Message.From.ID
			} else if u.CallbackQuery != nil {
				userID = u.CallbackQuery.From.ID
			}
			if userID != "" && blacklist(userID) {
				return
			}
			next(ctx, u)
		}
	}
}

// Whitelist میدل‌ور لیست سفید
func Whitelist(whitelist func(userID string) bool) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			var userID string
			if u.Message != nil && u.Message.From != nil {
				userID = u.Message.From.ID
			} else if u.CallbackQuery != nil {
				userID = u.CallbackQuery.From.ID
			}
			if userID != "" && !whitelist(userID) {
				return
			}
			next(ctx, u)
		}
	}
}
