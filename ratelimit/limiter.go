package ratelimit

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var ErrLimiterClosed = errors.New("rate limiter is closed")

// Limiter توکن باکت بهینه با atomic
type Limiter struct {
	tokens  chan struct{}
	rate    int
	burst   int
	refill  time.Duration
	stop    chan struct{}
	stopped atomic.Bool
	total   atomic.Int64
	dropped atomic.Int64
}

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
	for i := 0; i < rate; i++ {
		l.tokens <- struct{}{}
	}
	go l.refillLoop()
	return l
}

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
				// bucket پر است
			}
		case <-l.stop:
			return
		}
	}
}

func (l *Limiter) Wait(ctx context.Context) error {
	if l.stopped.Load() {
		return ErrLimiterClosed
	}
	select {
	case <-l.tokens:
		l.total.Add(1)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *Limiter) TryWait() bool {
	if l.stopped.Load() {
		return false
	}
	select {
	case <-l.tokens:
		return true
	default:
		l.dropped.Add(1)
		return false
	}
}

func (l *Limiter) WaitTimeout(timeout time.Duration) bool {
	if l.stopped.Load() {
		return false
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-l.tokens:
		return true
	case <-timer.C:
		l.dropped.Add(1)
		return false
	}
}

func (l *Limiter) Available() int { return len(l.tokens) }
func (l *Limiter) Rate() int      { return l.rate }
func (l *Limiter) Total() int64   { return l.total.Load() }
func (l *Limiter) Dropped() int64 { return l.dropped.Load() }

func (l *Limiter) Close() {
	if l.stopped.CompareAndSwap(false, true) {
		close(l.stop)
	}
}

// ConcurrencyLimiter برای محدود کردن درخواست‌های همزمان (Semaphore)
// این کلاس Rate Limit را از Concurrency Limit جدا می‌کند
type ConcurrencyLimiter struct {
	sem chan struct{}
}

func NewConcurrencyLimiter(maxConcurrent int) *ConcurrencyLimiter {
	if maxConcurrent <= 0 {
		maxConcurrent = 50
	}
	return &ConcurrencyLimiter{
		sem: make(chan struct{}, maxConcurrent),
	}
}

func (c *ConcurrencyLimiter) Acquire(ctx context.Context) error {
	select {
	case c.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *ConcurrencyLimiter) Release() {
	<-c.sem
}
