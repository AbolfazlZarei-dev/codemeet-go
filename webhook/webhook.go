package webhook

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/api"
	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/logger"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

// Config تنظیمات Webhook
type Config struct {
	ListenAddr        string
	Path              string
	SecretToken       string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	MaxHeaderBytes    int
	MaxBodySize       int64
	HTTPS             bool
	CertFile          string
	KeyFile           string
}

// DefaultConfig تنظیمات پیش‌فرض
func DefaultConfig() Config {
	return Config{
		ListenAddr:        ":8443",
		Path:              "/webhook",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20,  // 1MB
		MaxBodySize:       10 << 20, // 10MB
	}
}

// Server سرور Webhook
type Server struct {
	api        *api.Client
	dispatcher *dispatcher.Dispatcher
	logger     *logger.Logger
	cfg        Config
	srv        *http.Server
	requests   int64
	errors     int64
}

// New ساخت سرور
func New(c *api.Client, d *dispatcher.Dispatcher, log *logger.Logger, cfg Config) *Server {
	s := &Server{api: c, dispatcher: d, logger: log, cfg: cfg}
	mux := http.NewServeMux()
	mux.HandleFunc(cfg.Path, s.handleWebhook)
	mux.HandleFunc("/health", s.healthCheck)
	mux.HandleFunc("/metrics", s.metrics)
	s.srv = &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           mux,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}
	return s
}

func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&s.requests, 1)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		atomic.AddInt64(&s.errors, 1)
		return
	}

	// اعتبارسنجی Secret Token
	if s.cfg.SecretToken != "" {
		token := r.Header.Get("X-CodeMeet-Bot-Api-Secret-Token")
		if subtle.ConstantTimeCompare([]byte(token), []byte(s.cfg.SecretToken)) != 1 {
			if s.logger != nil {
				s.logger.Warn("invalid secret token", "ip", r.RemoteAddr)
			}
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			atomic.AddInt64(&s.errors, 1)
			return
		}
	}

	// محدود کردن سایز body
	body := http.MaxBytesReader(w, r.Body, s.cfg.MaxBodySize)
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("failed to read body", "error", err)
		}
		http.Error(w, "body too large or unreadable", http.StatusRequestEntityTooLarge)
		atomic.AddInt64(&s.errors, 1)
		return
	}

	var update models.Update
	if err := json.Unmarshal(data, &update); err != nil {
		if s.logger != nil {
			s.logger.Error("failed to decode update", "error", err)
		}
		http.Error(w, "bad request", http.StatusBadRequest)
		atomic.AddInt64(&s.errors, 1)
		return
	}

	s.dispatcher.Dispatch(r.Context(), &update)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func (s *Server) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"requests": atomic.LoadInt64(&s.requests),
		"errors":   atomic.LoadInt64(&s.errors),
	})
}

// Start شروع سرور
func (s *Server) Start(ctx context.Context) error {
	if s.logger != nil {
		s.logger.Info("starting webhook server", "addr", s.cfg.ListenAddr, "path", s.cfg.Path)
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.srv.Shutdown(shutdownCtx)
	}()

	if s.cfg.HTTPS {
		return s.srv.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
	}
	return s.srv.ListenAndServe()
}

// Stats آمار سرور
func (s *Server) Stats() (requests, errors int64) {
	return atomic.LoadInt64(&s.requests), atomic.LoadInt64(&s.errors)
}
