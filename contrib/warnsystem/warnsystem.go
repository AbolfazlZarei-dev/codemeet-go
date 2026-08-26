package warnsystem

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

const (
	shardCount      = 64
	stateTTL        = 24 * time.Hour // اخطارها پس از ۲۴ ساعت از حافظه پاک می‌شوند
	cleanupInterval = 1 * time.Hour
)

// Config تنظیمات سیستم اخطار
type Config struct {
	// حداکثر اخطار مجاز قبل از اعمال مجازات
	MaxWarnings int

	// تابع فراخوانی شده هنگام دریافت اخطار جدید
	WarnAction func(ctx context.Context, chatID, userID string, current, max int)

	// تابع فراخوانی شده هنگام رسیدن به سقف اخطارها (مثلا اخراج کاربر)
	MaxWarnAction func(ctx context.Context, chatID, userID string)
}

// DefaultConfig تنظیمات پیش‌فرض
func DefaultConfig() Config {
	return Config{
		MaxWarnings:   3,
		WarnAction:    nil,
		MaxWarnAction: nil,
	}
}

type userWarns struct {
	count      atomic.Int32
	lastActive atomic.Int64
}

type stats struct {
	totalWarns    atomic.Int64
	totalMaxWarns atomic.Int64
}

// WarnSystem ساختار اصلی پکیج
type WarnSystem struct {
	cfg          Config
	shardedWarns [shardCount]sync.Map // استفاده از sync.Map برای خواندن بدون قفل
	stats        stats
	ctx          context.Context
	cancel       context.CancelFunc
}

// New ساخت یک نمونه جدید
func New(cfg Config) *WarnSystem {
	if cfg.MaxWarnings <= 0 {
		cfg.MaxWarnings = 3
	}

	ctx, cancel := context.WithCancel(context.Background())
	ws := &WarnSystem{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	go ws.cleanupLoop()
	return ws
}

// getShard پیدا کردن shard مربوط به کاربر برای ذخیره سازی
func (ws *WarnSystem) getShard(chatID, userID string) int {
	var h uint32 = 2166136261
	key := chatID + userID
	for i := 0; i < len(key); i++ {
		h ^= uint32(key[i])
		h *= 16777619
	}
	return int(h) & (shardCount - 1)
}

// AddWarning اضافه کردن یک اخطار به کاربر
func (ws *WarnSystem) AddWarning(ctx context.Context, chatID, userID string) int {
	shard := ws.getShard(chatID, userID)
	shardMap := &ws.shardedWarns[shard]

	key := chatID + "_" + userID

	val, _ := shardMap.LoadOrStore(key, &userWarns{})
	uw := val.(*userWarns)
	uw.lastActive.Store(time.Now().Unix())

	current := uw.count.Add(1)
	ws.stats.totalWarns.Add(1)

	if int(current) >= ws.cfg.MaxWarnings {
		ws.stats.totalMaxWarns.Add(1)
		if ws.cfg.MaxWarnAction != nil {
			ws.cfg.MaxWarnAction(ctx, chatID, userID)
		}
	} else {
		if ws.cfg.WarnAction != nil {
			ws.cfg.WarnAction(ctx, chatID, userID, int(current), ws.cfg.MaxWarnings)
		}
	}

	return int(current)
}

// ResetWarnings صفر کردن اخطارهای یک کاربر
func (ws *WarnSystem) ResetWarnings(chatID, userID string) bool {
	shard := ws.getShard(chatID, userID)
	shardMap := &ws.shardedWarns[shard]
	key := chatID + "_" + userID

	if _, ok := shardMap.LoadAndDelete(key); ok {
		return true
	}
	return false
}

// GetWarnings دریافت تعداد اخطارهای فعلی کاربر
func (ws *WarnSystem) GetWarnings(chatID, userID string) int {
	shard := ws.getShard(chatID, userID)
	shardMap := &ws.shardedWarns[shard]
	key := chatID + "_" + userID

	val, ok := shardMap.Load(key)
	if !ok {
		return 0
	}
	uw := val.(*userWarns)
	return int(uw.count.Load())
}

// cleanupLoop پاکسازی حافظه از اخطارهای قدیمی و غیرفعال
func (ws *WarnSystem) cleanupLoop() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now().Unix()
			for i := 0; i < shardCount; i++ {
				ws.shardedWarns[i].Range(func(key, value any) bool {
					uw := value.(*userWarns)
					if now-uw.lastActive.Load() > int64(stateTTL.Seconds()) {
						ws.shardedWarns[i].Delete(key)
					}
					return true
				})
			}
		case <-ws.ctx.Done():
			return
		}
	}
}

// Stats آمار سیستم
func (ws *WarnSystem) Stats() map[string]int64 {
	return map[string]int64{
		"total_warns":     ws.stats.totalWarns.Load(),
		"total_max_warns": ws.stats.totalMaxWarns.Load(),
	}
}

// Close بستن منابع
func (ws *WarnSystem) Close() {
	ws.cancel()
}
