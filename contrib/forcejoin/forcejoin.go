package forcejoin

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

// Config تنظیمات سیستم عضویت اجباری
type Config struct {
	// لیست آیدی چت‌هایی که کاربر باید عضو آن‌ها باشد (مثلا @channelid)
	RequiredChannels []string

	// مدت زمان کش کردن وضعیت عضویت کاربر (برای کاهش درخواست به API)
	CacheTTL time.Duration

	// آیدی کاربرانی که حتی اگر عضو نباشند هم ربات برایشان کار می‌کند (ادمین‌ها)
	AdminIDs []string

	// تابع بررسی عضویت توسط ربات (GetChatMember)
	CheckMembership func(ctx context.Context, userID, chatID string) (bool, error)

	// اکشنی که هنگام عدم عضویت اجرا می‌شود (مثلا ارسال پیام با دکمه شیشهی جوین)
	NotJoinedAction func(ctx context.Context, userID, chatID string)
}

// DefaultConfig تنظیمات پیش‌فرض
func DefaultConfig() Config {
	return Config{
		CacheTTL: 10 * time.Minute, // هر ۱۰ دقیقه یکبار عضویت کاربر چک می‌شود
	}
}

type userCache struct {
	isMember bool
	expires  time.Time
}

type stats struct {
	checks  atomic.Int64
	blocked atomic.Int64
	passed  atomic.Int64
}

// ForceJoin ساختار اصلی پکیج
type ForceJoin struct {
	cfg   Config
	cache sync.Map // ذخیره وضعیت کاربران برای کاهش API Calls
	stats stats
}

// New ساخت یک نمونه جدید
func New(cfg Config) *ForceJoin {
	if cfg.CacheTTL <= 0 {
		cfg.CacheTTL = 10 * time.Minute
	}

	fj := &ForceJoin{cfg: cfg}

	// راه‌اندازی پاکسازی دوره‌ای حافظه
	go fj.cleanupLoop()

	return fj
}

// Middleware میدل‌ور عضویت اجباری
func (fj *ForceJoin) Middleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			if u == nil {
				next(ctx, u)
				return
			}

			// استخراج آیدی کاربر و چت
			var userID, chatID string
			if u.Message != nil && u.Message.From != nil {
				userID = u.Message.From.ID
				chatID = u.Message.Chat.ID
			} else if u.CallbackQuery != nil && u.CallbackQuery.From != nil {
				userID = u.CallbackQuery.From.ID
				if u.CallbackQuery.Message != nil && u.CallbackQuery.Message.Chat != nil {
					chatID = u.CallbackQuery.Message.Chat.ID
				}
			}

			if userID == "" || chatID == "" {
				next(ctx, u)
				return
			}

			// عبور ادمین‌ها
			for _, adminID := range fj.cfg.AdminIDs {
				if userID == adminID {
					next(ctx, u)
					return
				}
			}

			// اگر کانالی تعریف نشده بود، عبور کن
			if len(fj.cfg.RequiredChannels) == 0 {
				next(ctx, u)
				return
			}

			fj.stats.checks.Add(1)

			// ۱. بررسی کش
			if val, ok := fj.cache.Load(userID); ok {
				uc := val.(*userCache)
				if time.Now().Before(uc.expires) {
					if uc.isMember {
						fj.stats.passed.Add(1)
						next(ctx, u)
					} else {
						fj.stats.blocked.Add(1)
						fj.triggerNotJoined(ctx, userID, chatID)
					}
					return
				}
				fj.cache.Delete(userID) // کش منقضی شده
			}

			// ۲. بررسی عضویت در کانال‌ها (API Call)
			isMember := true
			for _, channelID := range fj.cfg.RequiredChannels {
				joined, err := fj.cfg.CheckMembership(ctx, userID, channelID)
				if err != nil || !joined {
					isMember = false
					break
				}
			}

			// ۳. ذخیره در کش
			fj.cache.Store(userID, &userCache{
				isMember: isMember,
				expires:  time.Now().Add(fj.cfg.CacheTTL),
			})

			// ۴. اجرای اکشن مناسب
			if isMember {
				fj.stats.passed.Add(1)
				next(ctx, u)
			} else {
				fj.stats.blocked.Add(1)
				fj.triggerNotJoined(ctx, userID, chatID)
			}
		}
	}
}

// triggerNotJoined اجرای اکشن عدم عضویت
func (fj *ForceJoin) triggerNotJoined(ctx context.Context, userID, chatID string) {
	if fj.cfg.NotJoinedAction != nil {
		fj.cfg.NotJoinedAction(ctx, userID, chatID)
	}
}

// ClearUserCache پاک کردن کش یک کاربر خاص
// این متد باید صدا زده شود وقتی کاربر روی دکمه "عضو شدم" کلیک می‌کند
func (fj *ForceJoin) ClearUserCache(userID string) {
	fj.cache.Delete(userID)
}

// cleanupLoop پاکسازی حافظه از کاربران غیرفعال
func (fj *ForceJoin) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for now := range ticker.C {
		fj.cache.Range(func(key, value any) bool {
			uc := value.(*userCache)
			if now.After(uc.expires) {
				fj.cache.Delete(key)
			}
			return true
		})
	}
}

// Stats آمار سیستم
func (fj *ForceJoin) Stats() map[string]int64 {
	return map[string]int64{
		"total_checks": fj.stats.checks.Load(),
		"blocked":      fj.stats.blocked.Load(),
		"passed":       fj.stats.passed.Load(),
	}
}
