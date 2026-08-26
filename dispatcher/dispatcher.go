package dispatcher

import (
	"context"
	"log"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

type HandlerFunc func(ctx context.Context, update *models.Update)
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Dispatcher struct {
	mu              sync.RWMutex
	handlers        map[string]HandlerFunc
	commandHandlers map[string]HandlerFunc
	regexHandlers   []regexHandler
	textHandlers    []textHandler
	middlewares     []MiddlewareFunc
	fallback        HandlerFunc

	workerPool chan struct{}
	wg         sync.WaitGroup
	stopCh     chan struct{}
	closed     atomic.Bool

	// آمار
	stats dispatchStats
}

type dispatchStats struct {
	totalDispatched atomic.Int64
	totalDropped    atomic.Int64
	totalPanics     atomic.Int64
}

type regexHandler struct {
	pattern *regexp.Regexp
	handler HandlerFunc
}

type textHandler struct {
	text    string
	handler HandlerFunc
}

func New(maxWorkers int) *Dispatcher {
	if maxWorkers <= 0 {
		maxWorkers = 200
	}
	return &Dispatcher{
		handlers:        make(map[string]HandlerFunc),
		commandHandlers: make(map[string]HandlerFunc),
		workerPool:      make(chan struct{}, maxWorkers),
		stopCh:          make(chan struct{}),
	}
}

func (d *Dispatcher) Use(mw ...MiddlewareFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.middlewares = append(d.middlewares, mw...)
}

func (d *Dispatcher) Handle(updateType string, h HandlerFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[updateType] = h
}

func (d *Dispatcher) OnMessage(h func(ctx context.Context, msg *models.Message)) {
	d.Handle("message", func(ctx context.Context, u *models.Update) {
		if u.Message != nil {
			h(ctx, u.Message)
		}
	})
}

func (d *Dispatcher) OnCallback(h func(ctx context.Context, cq *models.CallbackQuery)) {
	d.Handle("callback_query", func(ctx context.Context, u *models.Update) {
		if u.CallbackQuery != nil {
			h(ctx, u.CallbackQuery)
		}
	})
}

func (d *Dispatcher) OnCommand(cmd string, h func(ctx context.Context, msg *models.Message)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	cmd = strings.TrimPrefix(cmd, "/")
	d.commandHandlers[cmd] = func(ctx context.Context, u *models.Update) {
		if u.Message != nil {
			h(ctx, u.Message)
		}
	}
}

func (d *Dispatcher) OnText(text string, h func(ctx context.Context, msg *models.Message)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.textHandlers = append(d.textHandlers, textHandler{
		text: text,
		handler: func(ctx context.Context, u *models.Update) {
			if u.Message != nil {
				h(ctx, u.Message)
			}
		},
	})
}

func (d *Dispatcher) Fallback(h HandlerFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.fallback = h
}

func (d *Dispatcher) Dispatch(ctx context.Context, update *models.Update) {
	if d.closed.Load() {
		return
	}

	handler := d.matchHandler(update)
	if handler == nil {
		return
	}

	handler = d.applyMiddlewares(handler)
	d.stats.totalDispatched.Add(1)

	select {
	case d.workerPool <- struct{}{}:
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			defer func() { <-d.workerPool }()
			defer func() {
				if r := recover(); r != nil {
					d.stats.totalPanics.Add(1)
					log.Printf("PANIC recovered in dispatcher: %v\n%s", r, debugStack())
				}
			}()
			handler(ctx, update)
		}()
	default:
		// worker pool پر است — پردازش sync در goroutine جدید بدون محدودیت
		// این از drop کردن آپدیت‌ها جلوگیری می‌کند
		d.stats.totalDropped.Add(1)
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					d.stats.totalPanics.Add(1)
					log.Printf("PANIC (overflow): %v", r)
				}
			}()
			handler(ctx, update)
		}()
	}
}

func debugStack() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

func (d *Dispatcher) matchHandler(update *models.Update) HandlerFunc {
	updateType := update.Type()

	d.mu.RLock()
	defer d.mu.RUnlock()

	if updateType == "message" && update.Message != nil && update.Message.Text != "" {
		text := update.Message.Text
		if strings.HasPrefix(text, "/") {
			parts := strings.Fields(text)
			if len(parts) > 0 {
				cmd := strings.TrimPrefix(parts[0], "/")
				if idx := strings.Index(cmd, "@"); idx != -1 {
					cmd = cmd[:idx]
				}
				if h, ok := d.commandHandlers[cmd]; ok {
					return h
				}
			}
		}

		for _, th := range d.textHandlers {
			if th.text == text {
				return th.handler
			}
		}

		for _, rh := range d.regexHandlers {
			if rh.pattern.MatchString(text) {
				return rh.handler
			}
		}
	}

	if h, ok := d.handlers[updateType]; ok {
		return h
	}

	if d.fallback != nil {
		return d.fallback
	}

	return nil
}

func (d *Dispatcher) applyMiddlewares(handler HandlerFunc) HandlerFunc {
	d.mu.RLock()
	mws := make([]MiddlewareFunc, len(d.middlewares))
	copy(mws, d.middlewares)
	d.mu.RUnlock()

	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler
}

func (d *Dispatcher) Stop() {
	if !d.closed.CompareAndSwap(false, true) {
		return
	}
	close(d.stopCh)
	d.wg.Wait()
}

func (d *Dispatcher) OnRegex(pattern string, h func(ctx context.Context, msg *models.Message, matches []string)) error {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.regexHandlers = append(d.regexHandlers, regexHandler{
		pattern: compiled,
		handler: func(ctx context.Context, u *models.Update) {
			if u.Message != nil && u.Message.Text != "" {
				matches := compiled.FindStringSubmatch(u.Message.Text)
				if matches != nil {
					h(ctx, u.Message, matches)
				}
			}
		},
	})
	return nil
}

// Stats آمار dispatcher
func (d *Dispatcher) Stats() (dispatched, dropped, panics int64) {
	return d.stats.totalDispatched.Load(),
		d.stats.totalDropped.Load(),
		d.stats.totalPanics.Load()
}
