package dispatcher

import (
	"context"
	"regexp"
	"strings"
	"sync"

	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

// HandlerFunc امضای استاندارد برای تمامی هندلرها
type HandlerFunc func(ctx context.Context, update *models.Update)

// MiddlewareFunc تابعی که یک هندلر را می‌گیرد و هندلر جدیدی برمی‌گرداند
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

// Dispatcher مسیریاب مرکزی برای پردازش رویدادهای (Updates) دریافتی از کدمیت.
type Dispatcher struct {
	mu              sync.RWMutex
	handlers        map[string]HandlerFunc // هندلرهای مبتنی بر نوع
	commandHandlers map[string]HandlerFunc // هندلرهای مبتنی بر دستور
	regexHandlers   []regexHandler         // هندلرهای مبتنی بر regex
	textHandlers    []textHandler          // هندلرهای مبتنی بر متن دقیق
	middlewares     []MiddlewareFunc       // لایه‌های میان‌افزار
	fallback        HandlerFunc            // هندلر پیش‌فرض
	workerPool      chan struct{}          // کانال برای محدود کردن Goroutineها
	wg              sync.WaitGroup         // برای انتظار اتمام Goroutineها
	stopCh          chan struct{}          // برای توقف کامل دیسپاچر
	closed          bool
}

type regexHandler struct {
	pattern *regexp.Regexp
	handler HandlerFunc
}

type textHandler struct {
	text    string
	handler HandlerFunc
}

// New ساخت یک نمونه جدید از Dispatcher
func New(maxWorkers int) *Dispatcher {
	if maxWorkers <= 0 {
		maxWorkers = 100
	}
	return &Dispatcher{
		handlers:        make(map[string]HandlerFunc),
		commandHandlers: make(map[string]HandlerFunc),
		workerPool:      make(chan struct{}, maxWorkers),
		stopCh:          make(chan struct{}),
	}
}

// Use افزودن میان‌افزارها
func (d *Dispatcher) Use(mw ...MiddlewareFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.middlewares = append(d.middlewares, mw...)
}

// Handle ثبت هندلر برای یک نوع آپدیت
func (d *Dispatcher) Handle(updateType string, h HandlerFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[updateType] = h
}

// OnMessage ثبت هندلر برای پیام‌های متنی و رسانه‌ای
func (d *Dispatcher) OnMessage(h func(ctx context.Context, msg *models.Message)) {
	d.Handle("message", func(ctx context.Context, u *models.Update) {
		if u.Message != nil {
			h(ctx, u.Message)
		}
	})
}

// OnCallback ثبت هندلر برای کلیک روی دکمه‌های شیشه‌ای
func (d *Dispatcher) OnCallback(h func(ctx context.Context, cq *models.CallbackQuery)) {
	d.Handle("callback_query", func(ctx context.Context, u *models.Update) {
		if u.CallbackQuery != nil {
			h(ctx, u.CallbackQuery)
		}
	})
}

// OnCommand ثبت هندلر برای دستورات خاص (مثل /start)
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

// OnText ثبت هندلر برای متن دقیق
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

// Fallback ثبت هندلری که در صورت نبود هندلر اختصاصی اجرا می‌شود
func (d *Dispatcher) Fallback(h HandlerFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.fallback = h
}

// Dispatch پردازش و ارسال رویداد به هندلر مربوطه
func (d *Dispatcher) Dispatch(ctx context.Context, update *models.Update) {
	select {
	case <-d.stopCh:
		return
	default:
	}

	handler := d.matchHandler(update)
	if handler == nil {
		return
	}

	handler = d.applyMiddlewares(handler)

	d.wg.Add(1)
	d.workerPool <- struct{}{}

	go func() {
		defer d.wg.Done()
		defer func() { <-d.workerPool }()
		defer func() {
			if r := recover(); r != nil {
				_ = r
			}
		}()
		handler(ctx, update)
	}()
}

// matchHandler پیدا کردن هندلر مناسب
func (d *Dispatcher) matchHandler(update *models.Update) HandlerFunc {
	updateType := update.Type()

	d.mu.RLock()
	defer d.mu.RUnlock()

	// ۱. بررسی دستورات (فقط برای پیام‌های متنی)
	if updateType == "message" && update.Message != nil && update.Message.Text != "" {
		text := update.Message.Text
		if strings.HasPrefix(text, "/") {
			parts := strings.Fields(text)
			if len(parts) > 0 {
				cmd := strings.TrimPrefix(parts[0], "/")
				// پشتیبانی از فرمت /command@bot_username
				if idx := strings.Index(cmd, "@"); idx != -1 {
					cmd = cmd[:idx]
				}
				if h, ok := d.commandHandlers[cmd]; ok {
					return h
				}
			}
		}

		// ۲. بررسی هندلرهای دقیق متن
		for _, th := range d.textHandlers {
			if th.text == text {
				return th.handler
			}
		}

		// ۳. بررسی regex handlers
		for _, rh := range d.regexHandlers {
			if rh.pattern.MatchString(text) {
				return rh.handler
			}
		}
	}

	// ۴. هندلرهای مبتنی بر نوع
	if h, ok := d.handlers[updateType]; ok {
		return h
	}

	// ۵. fallback
	if d.fallback != nil {
		return d.fallback
	}

	return nil
}

// applyMiddlewares اعمال زنجیره middleware
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

// Stop متدی برای توقف امن
func (d *Dispatcher) Stop() {
	d.mu.Lock()
	if d.closed {
		d.mu.Unlock()
		return
	}
	d.closed = true
	close(d.stopCh)
	d.mu.Unlock()
	d.wg.Wait()
}

// OnRegex ثبت هندلر با regex روی متن پیام
func (d *Dispatcher) OnRegex(pattern string, h func(ctx context.Context, msg *models.Message, matches []string)) {
	compiled := regexp.MustCompile(pattern)
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
}
