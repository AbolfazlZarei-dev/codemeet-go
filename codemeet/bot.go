package codemeet

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/api"
	"github.com/AbolfazlZarei-dev/codemeet-go/cache"
	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/errors"
	"github.com/AbolfazlZarei-dev/codemeet-go/logger"
	"github.com/AbolfazlZarei-dev/codemeet-go/methods"
	"github.com/AbolfazlZarei-dev/codemeet-go/middleware"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
	"github.com/AbolfazlZarei-dev/codemeet-go/polling"
	"github.com/AbolfazlZarei-dev/codemeet-go/ratelimit"
	"github.com/AbolfazlZarei-dev/codemeet-go/retry"
	"github.com/AbolfazlZarei-dev/codemeet-go/webhook"
	"github.com/AbolfazlZarei-dev/codemeet-go/ws"
)

// Version نسخه کتابخانه
const Version = "1.0.0"

// Author سازنده کتابخانه
const Author = "Abolfazl Zarei"

// GitHubProfile آدرس گیت‌هاب سازنده
const GitHubProfile = "github.com/AbolfazlZarei-dev"

// GitHubRepo آدرس مخزن کتابخانه
const GitHubRepo = "github.com/AbolfazlZarei-dev/codemeet"

// Cache اینترفیس کش برای پشتیبانی از انواع کش‌ها
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
	Close()
}

// DashboardWriter بافر لاگ‌ها برای داشبورد
type DashboardWriter struct {
	mu   sync.Mutex
	logs []string
}

// Write پاکسازی کدهای ANSI قبل از ذخیره در بافر داشبورد
func (d *DashboardWriter) Write(p []byte) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// حذف کدهای رنگی (ANSI escape codes) برای نمایش تمیز در پنل وب
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	clean := re.ReplaceAll(p, []byte{})

	s := strings.TrimSpace(string(clean))
	if s != "" {
		d.logs = append(d.logs, s)
		if len(d.logs) > 100 { // نگه‌داشتن آخرین ۱۰۰ لاگ
			d.logs = d.logs[len(d.logs)-100:]
		}
	}
	return len(p), nil
}

func (d *DashboardWriter) GetLogs() []string {
	d.mu.Lock()
	defer d.mu.Unlock()
	cp := make([]string, len(d.logs))
	copy(cp, d.logs)
	return cp
}

// Bot ربات اصلی کدمیت با تمام قابلیت‌ها
type Bot struct {
	token       string
	baseURL     string
	api         *api.Client
	dispatcher  *dispatcher.Dispatcher
	rateLimit   *ratelimit.Limiter
	retry       *retry.Policy
	cache       Cache
	logger      *logger.Logger
	wsHub       *ws.Hub
	methods     *methods.Methods
	middlewares []middleware.Middleware

	mu             sync.RWMutex
	me             *models.User
	meFetched      bool
	runMode        string
	activeFeatures []string
	dashWriter     *DashboardWriter
}

// Option تابع پیکربندی
type Option func(*Bot)

func WithBaseURL(url string) Option {
	return func(b *Bot) { b.baseURL = url }
}

func WithHTTPClient(c *http.Client) Option {
	return func(b *Bot) { b.api.SetHTTPClient(c) }
}

func WithTimeout(d time.Duration) Option {
	return func(b *Bot) { b.api.SetTimeout(d) }
}

func WithRetry(p *retry.Policy) Option {
	return func(b *Bot) {
		b.retry = p
		b.activeFeatures = append(b.activeFeatures, "Retry Policy")
	}
}

func WithRateLimit(rps int) Option {
	return func(b *Bot) {
		b.rateLimit = ratelimit.New(rps)
		b.activeFeatures = append(b.activeFeatures, "Rate Limiter")
	}
}

func WithRateLimitBurst(rps, burst int) Option {
	return func(b *Bot) {
		b.rateLimit = ratelimit.NewWithBurst(rps, burst)
		b.activeFeatures = append(b.activeFeatures, "Rate Limiter (Burst)")
	}
}

func WithCache(ttl time.Duration) Option {
	return func(b *Bot) {
		b.cache = cache.New(ttl)
		b.activeFeatures = append(b.activeFeatures, "Cache")
	}
}

func WithShardedCache(shards int, ttl time.Duration) Option {
	return func(b *Bot) {
		b.cache = cache.NewSharded(shards, ttl)
		b.activeFeatures = append(b.activeFeatures, "Sharded Cache")
	}
}

func WithLogger(l *logger.Logger) Option {
	return func(b *Bot) { b.logger = l }
}

// WithMiddleware افزودن میان‌افزارها
func WithMiddleware(mws ...middleware.Middleware) Option {
	return func(b *Bot) {
		b.middlewares = append(b.middlewares, mws...)
		b.activeFeatures = append(b.activeFeatures, fmt.Sprintf("Middlewares (%d)", len(mws)))
		for _, mw := range mws {
			mw := mw
			b.dispatcher.Use(func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
				return func(ctx context.Context, u *models.Update) {
					wrappedNext := func(ctx context.Context, u *models.Update) {
						next(ctx, u)
					}
					mw(wrappedNext)(ctx, u)
				}
			})
		}
	}
}

// New ساخت ربات جدید
func New(token string, opts ...Option) (*Bot, error) {
	if token == "" {
		return nil, errors.NewValidationError("token", "token is required")
	}

	b := &Bot{
		token:       token,
		baseURL:     "https://botapi.codemeet.chat",
		dispatcher:  dispatcher.New(100),
		retry:       retry.DefaultPolicy(),
		rateLimit:   ratelimit.New(30),
		logger:      logger.New(logger.LevelInfo),
		middlewares: []middleware.Middleware{},
		activeFeatures: []string{
			"Dispatcher",
		},
	}

	b.api = api.NewClient(b.baseURL, token, b.logger)
	b.wsHub = ws.NewHub(b.api, b.logger)

	for _, opt := range opts {
		opt(b)
	}

	b.methods = methods.New(b.api, b.retry, b.rateLimit)
	return b, nil
}

// API دسترسی به متدهای API
func (b *Bot) API() *methods.Methods { return b.methods }

// Dispatcher دسترسی به Update Router
func (b *Bot) Dispatcher() *dispatcher.Dispatcher { return b.dispatcher }

// StartPolling شروع Long Polling
func (b *Bot) StartPolling(ctx context.Context, cfg polling.Config) error {
	b.runMode = "Long Polling"
	b.printStartupBanner(ctx)
	p := polling.New(b.api, b.dispatcher, b.logger, cfg)
	return p.Start(ctx)
}

// StartWebhook راه‌اندازی سرور Webhook
func (b *Bot) StartWebhook(ctx context.Context, cfg webhook.Config) error {
	b.runMode = "Webhook"
	b.printStartupBanner(ctx)
	wh := webhook.New(b.api, b.dispatcher, b.logger, cfg)
	return wh.Start(ctx)
}

// printStartupBanner چاپ اطلاعات سازنده و وضعیت ربات در ترمینال
func (b *Bot) printStartupBanner(ctx context.Context) {
	var botName, botUser string
	me, err := b.API().Bot().GetMe(ctx)
	if err == nil && me != nil {
		botName = me.FullName()
		botUser = "@" + me.Username
	} else {
		botName = "Unknown (Check Token/Network)"
		botUser = "Unknown"
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println(" CodeMeet Go Bot Started Successfully!")
	fmt.Println("--------------------------------------------------")
	fmt.Printf(" Author       : %s\n", Author)
	fmt.Printf(" GitHub       : %s\n", GitHubProfile)
	fmt.Printf(" Repository   : %s\n", GitHubRepo)
	fmt.Printf(" Version      : %s\n", Version)
	fmt.Println("--------------------------------------------------")
	fmt.Printf(" Bot Name     : %s\n", botName)
	fmt.Printf(" Bot Username : %s\n", botUser)
	fmt.Println("--------------------------------------------------")
	fmt.Printf(" Run Mode     : %s\n", b.runMode)
	fmt.Println(" Active Features:")
	for _, f := range b.activeFeatures {
		fmt.Printf("  - %s\n", f)
	}
	fmt.Println("--------------------------------------------------")
	fmt.Println(" Waiting for incoming messages...")
	fmt.Println("--------------------------------------------------")
}

// StartDashboard راه‌اندازی داشبورد وب برای مانیتورینگ ربات
func (b *Bot) StartDashboard(ctx context.Context, addr string) error {
	b.dashWriter = &DashboardWriter{}

	// اتصال خروجی لاگر به بافر داشبورد برای نمایش لاگ‌های زنده
	if b.logger != nil {
		b.logger.SetOutput(io.MultiWriter(os.Stdout, b.dashWriter))
		b.logger.SetLevel(logger.LevelInfo) // اضافه شد: تنظیم سطح لاگر برای داشبورد
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, dashboardHTML)
	})

	mux.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"author":   Author,
			"github":   GitHubProfile,
			"repo":     GitHubRepo,
			"version":  Version,
			"runMode":  b.RunMode(),
			"features": b.activeFeatures,
		})
	})

	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(b.Stats())
	})

	mux.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(b.dashWriter.GetLogs())
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	fmt.Printf(" Dashboard is running at http://localhost%s\n", addr)
	return srv.ListenAndServe()
}

// RunMode متدی برای دریافت وضعیت فعلی اجرای ربات
func (b *Bot) RunMode() string {
	if b.runMode == "" {
		return "Stopped"
	}
	return b.runMode
}

// SetWebhook تنظیم Webhook در سرور کدمیت
func (b *Bot) SetWebhook(ctx context.Context, url, secretToken string) error {
	return b.API().Webhook().Set(ctx, &models.SetWebhookRequest{
		URL:         url,
		SecretToken: secretToken,
	})
}

// DeleteWebhook حذف Webhook
func (b *Bot) DeleteWebhook(ctx context.Context) error {
	return b.API().Webhook().Delete(ctx)
}

// GetMe دریافت اطلاعات ربات (با کش)
func (b *Bot) GetMe(ctx context.Context) (*models.User, error) {
	if b.cache != nil {
		if v, ok := b.cache.Get("me"); ok {
			if u, ok := v.(*models.User); ok {
				return u, nil
			}
		}
	}

	b.mu.RLock()
	if b.meFetched && b.me != nil {
		me := b.me
		b.mu.RUnlock()
		return me, nil
	}
	b.mu.RUnlock()

	b.mu.Lock()
	if b.meFetched && b.me != nil {
		me := b.me
		b.mu.Unlock()
		return me, nil
	}
	b.mu.Unlock()

	me, err := b.API().Bot().GetMe(ctx)
	if err != nil {
		return nil, err
	}

	b.mu.Lock()
	b.me = me
	b.meFetched = true
	b.mu.Unlock()

	if b.cache != nil {
		b.cache.Set("me", me)
	}
	return me, nil
}

// ResetMe پاک کردن cache مربوط به GetMe
func (b *Bot) ResetMe() {
	b.mu.Lock()
	b.me = nil
	b.meFetched = false
	b.mu.Unlock()
	if b.cache != nil {
		b.cache.Delete("me")
	}
}

// WS دسترسی به WebSocket Hub
func (b *Bot) WS() *ws.Hub { return b.wsHub }

// Cache دسترسی به کش
func (b *Bot) Cache() Cache { return b.cache }

// Logger دسترسی به لاگر
func (b *Bot) Logger() *logger.Logger { return b.logger }

// RateLimiter دسترسی به rate limiter
func (b *Bot) RateLimiter() *ratelimit.Limiter { return b.rateLimit }

// RetryPolicy دسترسی به retry policy
func (b *Bot) RetryPolicy() *retry.Policy { return b.retry }

// Helper methods
func (b *Bot) OnCommand(cmd string, h func(ctx context.Context, msg *models.Message)) {
	b.dispatcher.OnCommand(cmd, h)
}

func (b *Bot) OnMessage(h func(ctx context.Context, msg *models.Message)) {
	b.dispatcher.OnMessage(h)
}

func (b *Bot) OnCallback(h func(ctx context.Context, cq *models.CallbackQuery)) {
	b.dispatcher.OnCallback(h)
}

func (b *Bot) OnText(text string, h func(ctx context.Context, msg *models.Message)) {
	b.dispatcher.OnText(text, h)
}

func (b *Bot) OnRegex(pattern string, h func(ctx context.Context, msg *models.Message, matches []string)) {
	b.dispatcher.OnRegex(pattern, h)
}

func (b *Bot) Fallback(h func(ctx context.Context, u *models.Update)) {
	b.dispatcher.Fallback(func(ctx context.Context, u *models.Update) {
		h(ctx, u)
	})
}

func (b *Bot) Use(mw ...dispatcher.MiddlewareFunc) {
	b.dispatcher.Use(mw...)
}

func (b *Bot) Send(ctx context.Context, chatID, text string) (*models.Message, error) {
	return b.API().Messages().SendText(ctx, chatID, text)
}

func (b *Bot) SendHTML(ctx context.Context, chatID, text string) (*models.Message, error) {
	return b.API().Messages().SendHTML(ctx, chatID, text)
}

func (b *Bot) SendWithKeyboard(ctx context.Context, chatID, text string, markup interface{}) (*models.Message, error) {
	return b.API().Messages().SendWithKeyboard(ctx, chatID, text, markup)
}

func (b *Bot) Reply(ctx context.Context, msg *models.Message, text string) (*models.Message, error) {
	return b.API().Messages().Send(ctx, &methods.SendMessageRequest{
		ChatID:           msg.Chat.ID,
		Text:             text,
		ReplyToMessageID: msg.MessageID,
	})
}

func (b *Bot) AnswerCallback(ctx context.Context, callbackID, text string, showAlert bool) error {
	return b.API().Messages().AnswerCallbackSimple(ctx, callbackID, text, showAlert)
}

// Close بستن اتصالات
func (b *Bot) Close() error {
	var errs []error
	if b.wsHub != nil {
		b.wsHub.Close()
	}
	if b.cache != nil {
		b.cache.Close()
	}
	if b.dispatcher != nil {
		b.dispatcher.Stop()
	}
	if b.rateLimit != nil {
		b.rateLimit.Close()
	}
	if b.api != nil {
		if err := b.api.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if b.logger != nil {
		b.logger.Close()
	}
	if len(errs) > 0 {
		return errors.NewValidationError("close", "errors during close")
	}
	return nil
}

// HealthCheck بررسی سلامت بات
func (b *Bot) HealthCheck(ctx context.Context) error {
	_, err := b.API().Bot().GetMe(ctx)
	return err
}

// Stats آمار بات
func (b *Bot) Stats() api.StatsSnapshot {
	return b.api.StatsSnapshot()
}
