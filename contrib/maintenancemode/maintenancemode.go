package maintenancemode

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

// ActionFunc تابعی است که هنگام فعال بودن حالت تعمیراتی
// برای ارسال پیام یا انجام عملیات موردنظر روی کاربر اجرا می‌شود.
type ActionFunc func(ctx context.Context, userID string, update *models.Update)

// Config تنظیمات حالت تعمیراتی را مشخص می‌کند.
type Config struct {
	// AdminIDs شناسه کاربرانی هستند که حتی در حالت تعمیراتی
	// اجازه عبور از میدل‌ور را دارند.
	AdminIDs []string

	// IsEnabled مشخص می‌کند که حالت تعمیراتی هنگام ساخت
	// نمونه به صورت پیش‌فرض فعال باشد یا خیر.
	IsEnabled bool

	// MaintenanceMsg پیام پیش‌فرض حالت تعمیراتی است.
	// این فیلد در اختیار Action قرار می‌گیرد و خود میدل‌ور
	// مستقیماً پیام ارسال نمی‌کند.
	MaintenanceMsg string

	// NotifyCooldown فاصله زمانی بین دو اخطار برای یک کاربر
	// بر حسب ثانیه است.
	NotifyCooldown int64

	// Action تابعی است که برای اطلاع‌رسانی به کاربر اجرا می‌شود.
	Action ActionFunc

	// CleanupInterval فاصله زمانی پاکسازی اطلاعات کاربران
	// منقضی‌شده را مشخص می‌کند.
	//
	// اگر مقدار آن صفر یا منفی باشد، مقدار پیش‌فرض استفاده می‌شود.
	CleanupInterval time.Duration
}

// DefaultConfig تنظیمات پیش‌فرض حالت تعمیراتی را برمی‌گرداند.
func DefaultConfig() Config {
	return Config{
		AdminIDs:        nil,
		IsEnabled:       false,
		MaintenanceMsg:  "🛠️ ربات در حال به‌روزرسانی و نگهداری است.\nلطفاً چند دقیقه دیگر مجدداً تلاش کنید.",
		NotifyCooldown:  300, // ۵ دقیقه
		Action:          nil,
		CleanupInterval: 10 * time.Minute,
	}
}

// MaintenanceMode مدیریت حالت تعمیراتی ربات را بر عهده دارد.
type MaintenanceMode struct {
	cfg Config

	// وضعیت فعال بودن حالت تعمیراتی.
	// استفاده از atomic باعث می‌شود بررسی وضعیت بدون Mutex انجام شود.
	enabled atomic.Bool

	// فهرست ادمین‌ها.
	//
	// این Map فقط هنگام ساخت نمونه نوشته می‌شود و بعد از آن
	// فقط خوانده می‌شود؛ بنابراین نیازی به Mutex ندارد.
	admins map[string]struct{}

	// زمان آخرین اخطار هر کاربر.
	//
	// Map ساده به همراه Mutex در این سناریو معمولاً سبک‌تر و
	// قابل‌کنترل‌تر از sync.Map است، چون ساختار داده مشخص و
	// عملیات بسیار ساده است.
	lastNotified map[string]int64

	// قفل مربوط به lastNotified.
	notifyMu sync.Mutex

	// کانال توقف Goroutine پاکسازی.
	stopCh chan struct{}

	// برای جلوگیری از چند بار close شدن stopCh.
	stopOnce sync.Once
}

// New یک نمونه جدید از حالت تعمیراتی ایجاد می‌کند.
func New(cfg Config) *MaintenanceMode {
	// مقادیر نامعتبر را به مقدار مناسب تبدیل می‌کنیم.
	if cfg.NotifyCooldown < 0 {
		cfg.NotifyCooldown = 0
	}

	if cfg.CleanupInterval <= 0 {
		cfg.CleanupInterval = 10 * time.Minute
	}

	// اگر پیام تعمیراتی خالی باشد، پیام پیش‌فرض استفاده می‌شود.
	if cfg.MaintenanceMsg == "" {
		cfg.MaintenanceMsg = DefaultConfig().MaintenanceMsg
	}

	// ساخت Map ادمین‌ها برای دسترسی O(1).
	admins := make(map[string]struct{}, len(cfg.AdminIDs))

	for _, adminID := range cfg.AdminIDs {
		if adminID == "" {
			continue
		}

		admins[adminID] = struct{}{}
	}

	// یک کپی از Config ایجاد می‌کنیم تا تغییرات بیرونی روی
	// Slice مربوط به AdminIDs روی نمونه تأثیر نگذارد.
	cfg.AdminIDs = nil

	mm := &MaintenanceMode{
		cfg:          cfg,
		admins:       admins,
		lastNotified: make(map[string]int64),
		stopCh:       make(chan struct{}),
	}

	mm.enabled.Store(cfg.IsEnabled)

	// فقط یک Goroutine سبک برای پاکسازی اطلاعات منقضی‌شده اجرا می‌شود.
	go mm.cleanupLoop()

	return mm
}

// cleanupLoop اطلاعات قدیمی کاربران را به صورت دوره‌ای
// از حافظه حذف می‌کند.
//
// این کار باعث می‌شود اگر تعداد کاربران بسیار زیاد شود،
// Map مربوط به cooldown برای همیشه رشد نکند.
func (mm *MaintenanceMode) cleanupLoop() {
	ticker := time.NewTicker(mm.cfg.CleanupInterval)
	defer ticker.Stop()

	// مدت زمانی که پس از آن اطلاعات کاربر قدیمی محسوب می‌شود.
	//
	// حداقل یک ساعت در نظر گرفته شده تا اطلاعات زودتر از
	// زمان موردنیاز حذف نشوند.
	cleanupAfter := int64(time.Hour / time.Second)

	for {
		select {
		case <-ticker.C:
			now := time.Now().Unix()

			mm.notifyMu.Lock()

			for userID, lastTime := range mm.lastNotified {
				if now-lastTime > cleanupAfter {
					delete(mm.lastNotified, userID)
				}
			}

			mm.notifyMu.Unlock()

		case <-mm.stopCh:
			return
		}
	}
}

// SetEnabled حالت تعمیراتی را در زمان اجرا روشن یا خاموش می‌کند.
//
// این عملیات کاملاً اتمیک است و برای استفاده همزمان از چند Goroutine
// ایمن است.
func (mm *MaintenanceMode) SetEnabled(enabled bool) {
	if mm == nil {
		return
	}

	mm.enabled.Store(enabled)
}

// IsEnabled وضعیت فعلی حالت تعمیراتی را برمی‌گرداند.
func (mm *MaintenanceMode) IsEnabled() bool {
	if mm == nil {
		return false
	}

	return mm.enabled.Load()
}

// Stop عملیات پاکسازی پس‌زمینه را متوقف می‌کند.
//
// این تابع را می‌توان چندین بار صدا زد و در این حالت
// هیچ Panicای رخ نخواهد داد.
func (mm *MaintenanceMode) Stop() {
	if mm == nil {
		return
	}

	mm.stopOnce.Do(func() {
		close(mm.stopCh)
	})
}

// isAdmin بررسی می‌کند که آیا کاربر ادمین است یا خیر.
//
// دسترسی به Map فقط خواندنی است و بعد از New دیگر تغییری در آن ایجاد نمی‌شود.
func (mm *MaintenanceMode) isAdmin(userID string) bool {
	if userID == "" {
		return false
	}

	_, ok := mm.admins[userID]
	return ok
}

// getUserID شناسه کاربر را از Update استخراج می‌کند.
func getUserID(u *models.Update) string {
	if u == nil {
		return ""
	}

	if u.Message != nil && u.Message.From != nil {
		return u.Message.From.ID
	}

	if u.CallbackQuery != nil && u.CallbackQuery.From != nil {
		return u.CallbackQuery.From.ID
	}

	return ""
}

// shouldNotify بررسی می‌کند که آیا زمان ارسال اخطار برای کاربر
// فرا رسیده است یا خیر.
//
// این تابع علاوه بر بررسی cooldown، زمان جدید را نیز به صورت
// اتمیک ثبت می‌کند.
//
// این موضوع بسیار مهم است؛ چون اگر چند درخواست همزمان برای
// یک کاربر برسد، فقط یکی از آن‌ها اجازه ارسال پیام خواهد داشت.
func (mm *MaintenanceMode) shouldNotify(userID string, now int64) bool {
	mm.notifyMu.Lock()
	defer mm.notifyMu.Unlock()

	lastTime, exists := mm.lastNotified[userID]

	if exists {
		if now-lastTime < mm.cfg.NotifyCooldown {
			return false
		}
	}

	// زمان جدید را همان لحظه ثبت می‌کنیم.
	//
	// بنابراین درخواست‌های همزمان نمی‌توانند چند بار
	// از cooldown عبور کنند.
	mm.lastNotified[userID] = now

	return true
}

// Middleware میدل‌ور حالت تعمیراتی را ایجاد می‌کند.
func (mm *MaintenanceMode) Middleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {

			// اگر نمونه معتبر نباشد، درخواست را بدون ایجاد خطا عبور می‌دهیم.
			if mm == nil {
				if next != nil {
					next(ctx, u)
				}
				return
			}

			// اگر حالت تعمیراتی خاموش باشد، هیچ پردازش اضافی انجام نمی‌دهیم.
			//
			// این سریع‌ترین مسیر اجرای Middleware است.
			if !mm.enabled.Load() {
				if next != nil {
					next(ctx, u)
				}
				return
			}

			// استخراج شناسه کاربر.
			userID := getUserID(u)

			// اگر شناسه کاربر مشخص نباشد، نمی‌توانیم برای او
			// پیام تعمیراتی ارسال کنیم.
			if userID == "" {
				return
			}

			// ادمین‌ها در حالت تعمیراتی نیز اجازه عبور دارند.
			if mm.isAdmin(userID) {
				if next != nil {
					next(ctx, u)
				}
				return
			}

			// زمان فعلی را فقط یک بار محاسبه می‌کنیم.
			now := time.Now().Unix()

			// بررسی و ثبت cooldown به صورت اتمیک.
			if !mm.shouldNotify(userID, now) {
				return
			}

			// اجرای Action برای اطلاع‌رسانی به کاربر.
			//
			// cooldown قبل از اجرای Action ثبت شده است تا اگر
			// Action زمان‌بر بود یا چند درخواست همزمان رسید،
			// پیام دوباره ارسال نشود.
			if mm.cfg.Action != nil {
				mm.cfg.Action(ctx, userID, u)
			}
		}
	}
}
