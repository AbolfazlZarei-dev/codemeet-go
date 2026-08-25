package ratelimit

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Limiter محدودیت نرخ با الگوریتم Token Bucket
type Limiter struct {
	mu      sync.Mutex
	tokens  chan struct{}
	rate    int
	burst   int
	refill  time.Duration
	stop    chan struct{}
	stopped int32
	total   int64
}

// New ساخت Rate Limiter جدید
func New(rate int) *Limiter {
	if rate <= 0 {
		rate = 30
	}
	l := &Limiter{
		tokens: make(chan struct{}, rate),
		rate:   rate,
		burst:  rate,
		refill: time.Second / time.Duration(rate),
		stop:   make(chan struct{}),
	}
	// پر کردن اولیه
	for i := 0; i < rate; i++ {
		l.tokens <- struct{}{}
	}
	go l.refillLoop()
	return l
}

// NewWithBurst ساخت limiter با burst (حداکثر تعداد درخواست همزمان)
func NewWithBurst(rate, burst int) *Limiter {
	if rate <= 0 {
		rate = 30
	}
	if burst <= 0 {
		burst = rate
	}
	l := &Limiter{
		tokens: make(chan struct{}, burst),
		rate:   rate,
		burst:  burst,
		refill: time.Second / time.Duration(rate),
		stop:   make(chan struct{}),
	}
	for i := 0; i < burst; i++ {
		l.tokens <- struct{}{}
	}
	go l.refillLoop()
	return l
}

func (l *Limiter) refillLoop() {
	ticker := time.NewTicker(l.refill)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			select {
			case l.tokens <- struct{}{}:
			default:
				// token bucket پر است — رها کن
			}
		case <-l.stop:
			return
		}
	}
}

// Wait صبر تا گرفتن توکن — خطا در صورت cancel شدن context
func (l *Limiter) Wait(ctx context.Context) error {
	if atomic.LoadInt32(&l.stopped) == 1 {
		return nil
	}
	atomic.AddInt64(&l.total, 1)
	select {
	case <-l.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryWait تلاش بدون صبر
func (l *Limiter) TryWait() bool {
	select {
	case <-l.tokens:
		return true
	default:
		return false
	}
}

// WaitTimeout صبر با timeout
func (l *Limiter) WaitTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-l.tokens:
		return true
	case <-timer.C:
		return false
	}
}

// Available تعداد توکن‌های موجود
func (l *Limiter) Available() int {
	return len(l.tokens)
}

// Rate نرخ تنظیم شده
func (l *Limiter) Rate() int { return l.rate }

// Total تعداد کل درخواست‌های Wait شده
func (l *Limiter) Total() int64 { return atomic.LoadInt64(&l.total) }

// Close توقف refiller
func (l *Limiter) Close() {
	if atomic.CompareAndSwapInt32(&l.stopped, 0, 1) {
		close(l.stop)
	}
}
