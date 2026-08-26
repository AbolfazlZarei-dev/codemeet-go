package codemeet

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/api"
	"github.com/AbolfazlZarei-dev/codemeet-go/cache"
	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/errors"
	"github.com/AbolfazlZarei-dev/codemeet-go/logger"
	"github.com/AbolfazlZarei-dev/codemeet-go/methods"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
	"github.com/AbolfazlZarei-dev/codemeet-go/polling"
	"github.com/AbolfazlZarei-dev/codemeet-go/ratelimit"
	"github.com/AbolfazlZarei-dev/codemeet-go/retry"
	"github.com/AbolfazlZarei-dev/codemeet-go/webhook"
	"golang.org/x/sync/singleflight"
)

const (
	Version       = "1.0.0"
	Author        = "Abolfazl Zarei"
	GitHubProfile = "github.com/AbolfazlZarei-dev"
	GitHubRepo    = "github.com/AbolfazlZarei-dev/codemeet-go"
)

// Cache interface برای کش
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
	Close()
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// DashboardWriter برای داشبورد
type DashboardWriter struct {
	mu   sync.Mutex
	logs []string
}

func (d *DashboardWriter) Write(p []byte) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	clean := ansiRegex.ReplaceAll(p, []byte{})
	s := strings.TrimSpace(string(clean))
	if s != "" {
		d.logs = append(d.logs, s)
		if len(d.logs) > 200 {
			d.logs = d.logs[len(d.logs)-200:]
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

// Bot ساختار اصلی ربات
type Bot struct {
	token      string
	baseURL    string
	api        *api.Client
	dispatcher *dispatcher.Dispatcher
	rateLimit  *ratelimit.Limiter
	retry      *retry.Policy
	cache      Cache
	logger     *logger.Logger
	methods    *methods.Methods

	mu             sync.RWMutex
	me             *models.User
	meFetched      atomic.Bool
	runMode        string
	activeFeatures []string
	dashWriter     *DashboardWriter

	// احراز هویت داشبورد
	dashUser     string
	dashPass     string
	dashSessions sync.Map

	getMeSF singleflight.Group

	// آمار داخلی ربات
	stats botStats
}

type botStats struct {
	UpdatesProcessed atomic.Int64
	CommandsExecuted atomic.Int64
	MessagesSent     atomic.Int64
	ErrorsCount      atomic.Int64
	StartTime        time.Time
}

type Option func(*Bot)

func WithBaseURL(url string) Option {
	return func(b *Bot) {
		b.baseURL = url
		b.logger.Debug("Base URL changed", "url", url)
	}
}

func WithHTTPClient(c *http.Client) Option {
	return func(b *Bot) {
		b.api.SetHTTPClient(c)
		b.logger.Debug("Custom HTTP Client set")
	}
}

func WithTimeout(d time.Duration) Option {
	return func(b *Bot) {
		b.api.SetTimeout(d)
		b.logger.Debug("HTTP Timeout set", "timeout", d.String())
	}
}

func WithRetry(p *retry.Policy) Option {
	return func(b *Bot) {
		b.retry = p
		b.activeFeatures = append(b.activeFeatures, "Retry Policy")
		b.logger.Info("Retry policy configured", "max_attempts", p.MaxAttempts)
	}
}

func WithRateLimit(rps int) Option {
	return func(b *Bot) {
		b.rateLimit = ratelimit.New(rps)
		b.activeFeatures = append(b.activeFeatures, "Rate Limiter")
		b.logger.Info("Rate limiter configured", "rps", rps)
	}
}

func WithRateLimitBurst(rps, burst int) Option {
	return func(b *Bot) {
		b.rateLimit = ratelimit.NewWithBurst(rps, burst)
		b.activeFeatures = append(b.activeFeatures, "Rate Limiter (Burst)")
		b.logger.Info("Rate limiter with burst configured", "rps", rps, "burst", burst)
	}
}

func WithCache(ttl time.Duration) Option {
	return func(b *Bot) {
		b.cache = cache.New(ttl)
		b.activeFeatures = append(b.activeFeatures, "Cache")
		b.logger.Info("Cache initialized", "ttl", ttl.String())
	}
}

func WithShardedCache(shards int, ttl time.Duration) Option {
	return func(b *Bot) {
		b.cache = cache.NewSharded(shards, ttl)
		b.activeFeatures = append(b.activeFeatures, "Sharded Cache")
		b.logger.Info("Sharded cache initialized", "shards", shards, "ttl", ttl.String())
	}
}

func WithLogger(l *logger.Logger) Option {
	return func(b *Bot) {
		b.logger = l
		b.logger.Debug("Custom logger attached")
	}
}

// WithDashboardAuth فعال‌سازی صفحه لاگین برای داشبورد
func WithDashboardAuth(user, pass string) Option {
	return func(b *Bot) {
		b.dashUser = user
		b.dashPass = pass
		b.activeFeatures = append(b.activeFeatures, "Dashboard Auth")
		b.logger.Info("Dashboard authentication enabled", "user", user)
	}
}

// WithMiddleware برای اضافه کردن میدل‌ورها به صورت مستقل
func WithMiddleware(mws ...dispatcher.MiddlewareFunc) Option {
	return func(b *Bot) {
		b.activeFeatures = append(b.activeFeatures, fmt.Sprintf("Middlewares (%d)", len(mws)))
		b.dispatcher.Use(mws...)
		b.logger.Info("Middlewares registered", "count", len(mws))
	}
}

// New ساخت ربات جدید
func New(token string, opts ...Option) (*Bot, error) {
	if token == "" {
		return nil, errors.NewValidationError("token", "token is required")
	}

	b := &Bot{
		token:      token,
		baseURL:    "https://botapi.codemeet.chat",
		dispatcher: dispatcher.New(200),
		retry:      retry.DefaultPolicy(),
		rateLimit:  ratelimit.New(30),
		logger:     logger.New(logger.LevelInfo),
		activeFeatures: []string{
			"Dispatcher",
			"Worker Pool",
		},
		stats: botStats{StartTime: time.Now()},
	}

	b.logger.Info("Initializing CodeMeet Bot...")

	b.api = api.NewClient(b.baseURL, token, b.logger)

	for _, opt := range opts {
		opt(b)
	}

	b.methods = methods.New(b.api, b.retry, b.rateLimit)
	b.logger.Info("Bot initialized successfully")
	return b, nil
}

// Getters
func (b *Bot) API() *methods.Methods              { return b.methods }
func (b *Bot) Dispatcher() *dispatcher.Dispatcher { return b.dispatcher }
func (b *Bot) Cache() Cache                       { return b.cache }
func (b *Bot) Logger() *logger.Logger             { return b.logger }
func (b *Bot) RateLimiter() *ratelimit.Limiter    { return b.rateLimit }
func (b *Bot) RetryPolicy() *retry.Policy         { return b.retry }
func (b *Bot) Token() string                      { return b.token }
func (b *Bot) BaseURL() string                    { return b.baseURL }

// StartPolling شروع با Long Polling
func (b *Bot) StartPolling(ctx context.Context, cfg polling.Config) error {
	b.setRunMode("Long Polling")
	b.printStartupBanner(ctx)
	b.logger.Info("Starting Long Polling...", "timeout", cfg.Timeout, "limit", cfg.Limit)
	p := polling.New(b.api, b.dispatcher, b.logger, cfg)
	return p.Start(ctx)
}

// StartWebhook شروع با Webhook
func (b *Bot) StartWebhook(ctx context.Context, cfg webhook.Config) error {
	b.setRunMode("Webhook")
	b.printStartupBanner(ctx)
	b.logger.Info("Starting Webhook Server...", "addr", cfg.ListenAddr)
	wh := webhook.New(b.api, b.dispatcher, b.logger, cfg)
	return wh.Start(ctx)
}

func (b *Bot) setRunMode(mode string) {
	b.mu.Lock()
	b.runMode = mode
	b.mu.Unlock()
	b.logger.Debug("Run mode set", "mode", mode)
}

func (b *Bot) printStartupBanner(ctx context.Context) {
	var botName, botUser string

	bannerCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	b.logger.Debug("Fetching bot info for banner...")
	me, err := b.API().Bot().GetMe(bannerCtx)
	if err == nil && me != nil {
		botName = me.FullName()
		botUser = "@" + me.Username
	} else {
		botName = "Unknown (Check Token/Network)"
		botUser = "Unknown"
		if err != nil {
			b.logger.Error("Failed to fetch bot info for banner", "error", err)
		}
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

	b.mu.RLock()
	mode := b.runMode
	features := b.activeFeatures
	b.mu.RUnlock()

	fmt.Printf(" Run Mode     : %s\n", mode)
	fmt.Println(" Active Features:")
	for _, f := range features {
		fmt.Printf("  - %s\n", f)
	}
	fmt.Println("--------------------------------------------------")
	fmt.Println(" Waiting for incoming messages...")
	fmt.Println("--------------------------------------------------")
}

// authMiddleware بررسی اینکه آیا کاربر لاگین کرده است یا خیر
func (b *Bot) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if b.dashUser == "" && b.dashPass == "" {
			next.ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		c, err := r.Cookie("cm_dash_session")
		if err != nil || !b.isValidSession(c.Value) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (b *Bot) isValidSession(token string) bool {
	val, ok := b.dashSessions.Load(token)
	if !ok {
		return false
	}
	if time.Now().After(val.(time.Time)) {
		b.dashSessions.Delete(token)
		return false
	}
	return true
}

func (b *Bot) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		errMsg := ""
		if r.URL.Query().Get("err") == "1" {
			errMsg = `<div style="color: #ef4444; margin-bottom: 15px; font-weight: bold;">Invalid username or password!</div>`
		}
		html := strings.Replace(loginHTML, "{{ERR_MSG}}", errMsg, 1)
		fmt.Fprint(w, html)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		user := r.FormValue("username")
		pass := r.FormValue("password")

		if user == b.dashUser && pass == b.dashPass {
			token := generateToken()
			b.dashSessions.Store(token, time.Now().Add(24*time.Hour))
			http.SetCookie(w, &http.Cookie{
				Name:     "cm_dash_session",
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().Add(24 * time.Hour),
			})
			b.logger.Info("Dashboard login successful", "user", user, "ip", r.RemoteAddr)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		b.logger.Warn("Failed dashboard login attempt", "ip", r.RemoteAddr)
		http.Redirect(w, r, "/login?err=1", http.StatusSeeOther)
	}
}

func (b *Bot) handleLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("cm_dash_session")
	if err == nil {
		b.dashSessions.Delete(c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "cm_dash_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
	b.logger.Info("Dashboard user logged out", "ip", r.RemoteAddr)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// StartDashboard داشبورد مانیتورینگ
func (b *Bot) StartDashboard(ctx context.Context, addr string) error {
	b.dashWriter = &DashboardWriter{}
	if b.logger != nil {
		b.logger.SetOutput(io.MultiWriter(os.Stdout, b.dashWriter))
	}

	b.logger.Info("Starting Web Dashboard...", "addr", addr)

	mux := http.NewServeMux()

	if b.dashUser != "" && b.dashPass != "" {
		mux.HandleFunc("/login", b.handleLogin)
		mux.HandleFunc("/logout", b.handleLogout)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, dashboardHTML)
	})

	mux.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b.mu.RLock()
		info := map[string]interface{}{
			"author":       Author,
			"github":       GitHubProfile,
			"repo":         GitHubRepo,
			"version":      Version,
			"runMode":      b.runMode,
			"features":     b.activeFeatures,
			"uptime":       time.Since(b.stats.StartTime).String(),
			"logs_enabled": b.logger.IsEnabled(),
		}
		b.mu.RUnlock()
		json.NewEncoder(w).Encode(info)
	})

	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		s := b.api.StatsSnapshot()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"api":               s,
			"updates_processed": b.stats.UpdatesProcessed.Load(),
			"commands_executed": b.stats.CommandsExecuted.Load(),
			"messages_sent":     b.stats.MessagesSent.Load(),
			"errors":            b.stats.ErrorsCount.Load(),
			"uptime":            time.Since(b.stats.StartTime).String(),
		})
	})

	mux.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(b.dashWriter.GetLogs())
	})

	// Endpoint جدید برای استاپ و استارت لاگ‌ها
	mux.HandleFunc("/api/logs/toggle", func(w http.ResponseWriter, r *http.Request) {
		if b.logger != nil {
			newState := !b.logger.IsEnabled()
			b.logger.SetEnabled(newState)
			b.logger.Info("Logging state toggled", "new_state", newState)
			// چون لاگ بالا ممکنه بعد از خاموش شدن پرینت نشه، ما اینجا مستقیم می‌فرستیم
			if !newState {
				b.dashWriter.Write([]byte("[DASHBOARD] Logging has been STOPPED by user.\n"))
			} else {
				b.dashWriter.Write([]byte("[DASHBOARD] Logging has been RESUMED by user.\n"))
			}
			json.NewEncoder(w).Encode(map[string]bool{"enabled": newState})
			return
		}
		http.Error(w, "logger not initialized", http.StatusInternalServerError)
	})

	handler := b.authMiddleware(mux)

	srv := &http.Server{Addr: addr, Handler: handler}
	go func() {
		<-ctx.Done()
		b.logger.Info("Dashboard shutting down...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	fmt.Printf(" Dashboard is running at http://localhost%s\n", addr)
	return srv.ListenAndServe()
}

func (b *Bot) RunMode() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.runMode == "" {
		return "Stopped"
	}
	return b.runMode
}

func (b *Bot) SetWebhook(ctx context.Context, url, secretToken string) error {
	b.logger.Info("Setting webhook", "url", url)
	return b.API().Webhook().Set(ctx, &models.SetWebhookRequest{
		URL:         url,
		SecretToken: secretToken,
	})
}

func (b *Bot) DeleteWebhook(ctx context.Context) error {
	b.logger.Info("Deleting webhook")
	return b.API().Webhook().Delete(ctx)
}

// GetMe با کش و singleflight
func (b *Bot) GetMe(ctx context.Context) (*models.User, error) {
	b.logger.Debug("GetMe called")
	if b.cache != nil {
		if v, ok := b.cache.Get("me"); ok {
			if u, ok := v.(*models.User); ok {
				b.logger.Debug("GetMe returned from cache")
				return u, nil
			}
		}
	}
	v, err, _ := b.getMeSF.Do("getMe", func() (interface{}, error) {
		b.logger.Debug("Fetching bot info from API...")
		return b.API().Bot().GetMe(ctx)
	})
	if err != nil {
		b.logger.Error("GetMe failed", "error", err)
		return nil, err
	}
	me := v.(*models.User)
	b.mu.Lock()
	b.me = me
	b.meFetched.Store(true)
	b.mu.Unlock()
	if b.cache != nil {
		b.cache.Set("me", me)
	}
	b.logger.Info("Bot info fetched successfully", "username", me.Username)
	return me, nil
}

func (b *Bot) ResetMe() {
	b.logger.Debug("Resetting bot info cache")
	b.mu.Lock()
	b.me = nil
	b.mu.Unlock()
	b.meFetched.Store(false)
	if b.cache != nil {
		b.cache.Delete("me")
	}
}

// Helper dispatcher methods
func (b *Bot) OnCommand(cmd string, h func(ctx context.Context, msg *models.Message)) {
	b.logger.Debug("Command handler registered", "command", cmd)
	b.dispatcher.OnCommand(cmd, func(ctx context.Context, msg *models.Message) {
		b.stats.CommandsExecuted.Add(1)
		b.stats.UpdatesProcessed.Add(1)
		b.logger.Info("Command executed", "cmd", cmd, "user_id", msg.From.ID)
		h(ctx, msg)
	})
}

func (b *Bot) OnMessage(h func(ctx context.Context, msg *models.Message)) {
	b.dispatcher.OnMessage(func(ctx context.Context, msg *models.Message) {
		b.stats.UpdatesProcessed.Add(1)
		b.logger.Debug("Message received", "chat_id", msg.Chat.ID, "text", msg.Text)
		h(ctx, msg)
	})
}

func (b *Bot) OnCallback(h func(ctx context.Context, cq *models.CallbackQuery)) {
	b.dispatcher.OnCallback(func(ctx context.Context, cq *models.CallbackQuery) {
		b.stats.UpdatesProcessed.Add(1)
		b.logger.Debug("Callback query received", "data", cq.Data, "user_id", cq.From.ID)
		h(ctx, cq)
	})
}

func (b *Bot) OnText(text string, h func(ctx context.Context, msg *models.Message)) {
	b.logger.Debug("Text handler registered", "text", text)
	b.dispatcher.OnText(text, h)
}

func (b *Bot) OnRegex(pattern string, h func(ctx context.Context, msg *models.Message, matches []string)) {
	b.logger.Debug("Regex handler registered", "pattern", pattern)
	b.dispatcher.OnRegex(pattern, h)
}

func (b *Bot) Fallback(h func(ctx context.Context, u *models.Update)) {
	b.logger.Debug("Fallback handler registered")
	b.dispatcher.Fallback(func(ctx context.Context, u *models.Update) {
		b.stats.UpdatesProcessed.Add(1)
		b.logger.Debug("Fallback handler executed", "update_id", u.UpdateID)
		h(ctx, u)
	})
}

func (b *Bot) Use(mw ...dispatcher.MiddlewareFunc) {
	b.logger.Debug("Registering middlewares", "count", len(mw))
	b.dispatcher.Use(mw...)
}

// Message helpers
func (b *Bot) Send(ctx context.Context, chatID, text string) (*models.Message, error) {
	b.logger.Info("Sending message", "chat_id", chatID, "text_length", len(text))
	m, err := b.API().Messages().SendText(ctx, chatID, text)
	if err == nil {
		b.stats.MessagesSent.Add(1)
		b.logger.Debug("Message sent successfully", "message_id", m.MessageID)
	} else {
		b.stats.ErrorsCount.Add(1)
		b.logger.Error("Failed to send message", "chat_id", chatID, "error", err)
	}
	return m, err
}

func (b *Bot) SendHTML(ctx context.Context, chatID, text string) (*models.Message, error) {
	b.logger.Info("Sending HTML message", "chat_id", chatID)
	m, err := b.API().Messages().SendHTML(ctx, chatID, text)
	if err == nil {
		b.stats.MessagesSent.Add(1)
	} else {
		b.stats.ErrorsCount.Add(1)
		b.logger.Error("Failed to send HTML message", "error", err)
	}
	return m, err
}

func (b *Bot) SendWithKeyboard(ctx context.Context, chatID, text string, markup interface{}) (*models.Message, error) {
	b.logger.Info("Sending message with keyboard", "chat_id", chatID)
	m, err := b.API().Messages().SendWithKeyboard(ctx, chatID, text, markup)
	if err == nil {
		b.stats.MessagesSent.Add(1)
	} else {
		b.stats.ErrorsCount.Add(1)
		b.logger.Error("Failed to send message with keyboard", "error", err)
	}
	return m, err
}

func (b *Bot) Reply(ctx context.Context, msg *models.Message, text string) (*models.Message, error) {
	b.logger.Info("Replying to message", "chat_id", msg.Chat.ID, "reply_to", msg.MessageID)
	m, err := b.API().Messages().Send(ctx, &methods.SendMessageRequest{
		ChatID:           msg.Chat.ID,
		Text:             text,
		ReplyToMessageID: msg.MessageID,
	})
	if err == nil {
		b.stats.MessagesSent.Add(1)
	} else {
		b.stats.ErrorsCount.Add(1)
		b.logger.Error("Failed to reply", "error", err)
	}
	return m, err
}

func (b *Bot) AnswerCallback(ctx context.Context, callbackID, text string, showAlert bool) error {
	b.logger.Debug("Answering callback query", "callback_id", callbackID)
	return b.API().Messages().AnswerCallbackSimple(ctx, callbackID, text, showAlert)
}

// Close بستن منابع
func (b *Bot) Close() error {
	b.logger.Info("Shutting down bot gracefully...")
	var errs []error
	if b.cache != nil {
		b.cache.Close()
		b.logger.Debug("Cache closed")
	}
	if b.dispatcher != nil {
		b.dispatcher.Stop()
		b.logger.Debug("Dispatcher stopped")
	}
	if b.rateLimit != nil {
		b.rateLimit.Close()
		b.logger.Debug("Rate limiter closed")
	}
	if b.api != nil {
		if err := b.api.Close(); err != nil {
			errs = append(errs, err)
		} else {
			b.logger.Debug("API client closed")
		}
	}
	if b.logger != nil {
		b.logger.Debug("Logger closed")
		b.logger.Close()
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors during close: %v", errs)
	}
	return nil
}

// HealthCheck بررسی سلامت
func (b *Bot) HealthCheck(ctx context.Context) error {
	b.logger.Debug("Performing health check...")
	_, err := b.API().Bot().GetMe(ctx)
	if err != nil {
		b.logger.Error("Health check failed", "error", err)
	}
	return err
}

// Stats آمار
func (b *Bot) Stats() api.StatsSnapshot {
	return b.api.StatsSnapshot()
}

// Uptime زمان فعالیت
func (b *Bot) Uptime() time.Duration {
	return time.Since(b.stats.StartTime)
}

// WithoutLogger برای غیرفعال کردن کامل لاگ‌ها در ترمینال
func WithoutLogger() Option {
	return func(b *Bot) {
		l := logger.New(logger.LevelFatal)
		l.SetOutput(io.Discard)
		b.logger = l
	}
}
