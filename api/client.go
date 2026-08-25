package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/errors"
	"github.com/AbolfazlZarei-dev/codemeet-go/logger"
	"github.com/google/uuid"
)

// Client کلاینت HTTP با Connection Pooling، Circuit Breaker و بهینه‌سازی
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	logger     *logger.Logger
	transport  *http.Transport
	mu         sync.RWMutex
	stats      *Stats
	breaker    *CircuitBreaker
}

// Stats آمار کلاینت — هم thread-safe و هم lock-free برای شمارش
type Stats struct {
	Requests     int64
	SuccessCount int64
	ErrorCount   int64
	BytesIn      int64
	BytesOut     int64
	totalLatency int64 // nanoseconds
	mu           sync.Mutex
}

// Record ثبت آمار یک درخواست
func (s *Stats) Record(latency time.Duration, success bool, bytesIn, bytesOut int64) {
	atomic.AddInt64(&s.Requests, 1)
	atomic.AddInt64(&s.totalLatency, int64(latency))
	atomic.AddInt64(&s.BytesIn, bytesIn)
	atomic.AddInt64(&s.BytesOut, bytesOut)
	if success {
		atomic.AddInt64(&s.SuccessCount, 1)
	} else {
		atomic.AddInt64(&s.ErrorCount, 1)
	}
}

// AvgLatency میانگین تأخیر
func (s *Stats) AvgLatency() time.Duration {
	r := atomic.LoadInt64(&s.Requests)
	if r == 0 {
		return 0
	}
	return time.Duration(atomic.LoadInt64(&s.totalLatency) / r)
}

// Snapshot گرفتن کپی از آمار
func (s *Stats) Snapshot() StatsSnapshot {
	return StatsSnapshot{
		Requests:     atomic.LoadInt64(&s.Requests),
		SuccessCount: atomic.LoadInt64(&s.SuccessCount),
		ErrorCount:   atomic.LoadInt64(&s.ErrorCount),
		BytesIn:      atomic.LoadInt64(&s.BytesIn),
		BytesOut:     atomic.LoadInt64(&s.BytesOut),
		AvgLatency:   s.AvgLatency(),
	}
}

// StatsSnapshot یک کپی immutable از آمار
type StatsSnapshot struct {
	Requests     int64
	SuccessCount int64
	ErrorCount   int64
	BytesIn      int64
	BytesOut     int64
	AvgLatency   time.Duration
}

// CircuitBreaker الگوی Circuit Breaker برای جلوگیری از ارسال درخواست به سرور آسیب‌دیده
type CircuitBreaker struct {
	mu               sync.Mutex
	failureThreshold int
	resetTimeout     time.Duration
	failures         int
	lastFailure      time.Time
	state            CircuitState
}

type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// NewCircuitBreaker ساخت Circuit Breaker
func NewCircuitBreaker(threshold int, reset time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: threshold,
		resetTimeout:     reset,
		state:            StateClosed,
	}
}

// Allow آیا اجازه ارسال درخواست داریم؟
func (c *CircuitBreaker) Allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch c.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(c.lastFailure) > c.resetTimeout {
			c.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return true
}

// RecordSuccess ثبت موفقیت
func (c *CircuitBreaker) RecordSuccess() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures = 0
	c.state = StateClosed
}

// RecordFailure ثبت شکست
func (c *CircuitBreaker) RecordFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures++
	c.lastFailure = time.Now()
	if c.failures >= c.failureThreshold {
		c.state = StateOpen
	}
}

// State دریافت وضعیت فعلی
func (c *CircuitBreaker) State() CircuitState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

// NewClient ساخت کلاینت جدید با تنظیمات بهینه
func NewClient(baseURL, token string, log *logger.Logger) *Client {
	transport := &http.Transport{
		MaxIdleConns:          500,
		MaxIdleConnsPerHost:   200,
		MaxConnsPerHost:       200,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ForceAttemptHTTP2:     true,
		DisableCompression:    false,
		Proxy:                 http.ProxyFromEnvironment,
		WriteBufferSize:       64 * 1024,
		ReadBufferSize:        64 * 1024,
	}

	hc := &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
	}

	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: hc,
		logger:     log,
		transport:  transport,
		stats:      &Stats{},
		breaker:    NewCircuitBreaker(10, 30*time.Second),
	}
}

// SetHTTPClient تنظیم HTTP Client سفارشی (حفظ breaker)
func (c *Client) SetHTTPClient(hc *http.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.httpClient = hc
	if tr, ok := hc.Transport.(*http.Transport); ok {
		c.transport = tr
	}
}

// SetTimeout تنظیم تایم‌اوت
func (c *Client) SetTimeout(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.httpClient.Timeout = d
}

// Stats دریافت آبجکت آمار
func (c *Client) Stats() *Stats { return c.stats }

// StatsSnapshot گرفتن snapshot از آمار
func (c *Client) StatsSnapshot() StatsSnapshot { return c.stats.Snapshot() }

// Breaker دسترسی به Circuit Breaker
func (c *Client) Breaker() *CircuitBreaker { return c.breaker }

// UserAgent ساخت User-Agent
func (c *Client) UserAgent() string {
	return fmt.Sprintf("codemeet-go/1.0.0 (%s/%s) go/%s",
		runtime.GOOS, runtime.GOARCH, runtime.Version())
}

// buildURL ساخت URL کامل
func (c *Client) buildURL(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", c.baseURL, c.token, method)
}

// DownloadFile دانلود فایل از سرور کدمیت (پیاده‌سازی شده)
func (c *Client) DownloadFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/file/bot%s/%s", c.baseURL, c.token, filePath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.NewNetworkError(err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("download failed with status: %s", resp.Status)
	}

	return resp.Body, nil
}

// doRequest اجرای درخواست — با اندازه گیری و breaker
func (c *Client) doRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
	if !c.breaker.Allow() {
		return nil, errors.NewNetworkError(fmt.Errorf("circuit breaker open"))
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		c.breaker.RecordFailure()
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent())
	req.Header.Set("Accept-Encoding", "gzip")

	if reqID, ok := ctx.Value(requestIDKey{}).(string); ok && reqID != "" {
		req.Header.Set("X-Request-ID", reqID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.breaker.RecordFailure()
		return nil, errors.NewNetworkError(err)
	}
	c.breaker.RecordSuccess()
	return resp, nil
}

// Request GET request (ساده‌ترین حالت)
func (c *Client) Request(ctx context.Context, method string) (*Response, error) {
	return c.RequestWithParams(ctx, method, nil)
}

// RequestWithParams ارسال درخواست با پارامتر JSON
func (c *Client) RequestWithParams(ctx context.Context, method string, params interface{}) (*Response, error) {
	url := c.buildURL(method)

	var body []byte
	var err error
	if params != nil {
		body, err = json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %w", err)
		}
	}

	start := time.Now()

	resp, err := c.doRequest(ctx, http.MethodPost, url, body)
	if err != nil {
		c.stats.Record(time.Since(start), false, 0, int64(len(body)))
		return nil, err
	}
	defer resp.Body.Close()

	limited := io.LimitReader(resp.Body, 100<<20) // 100MB
	data, err := io.ReadAll(limited)
	if err != nil {
		c.stats.Record(time.Since(start), false, int64(len(data)), int64(len(body)))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(data, &apiResp); err != nil {
		c.stats.Record(time.Since(start), false, int64(len(data)), int64(len(body)))
		return nil, fmt.Errorf("failed to decode response: %w (body: %s)", err, truncate(string(data), 500))
	}
	apiResp.HTTPStatus = resp.StatusCode

	if !apiResp.Ok {
		params := map[string]interface{}{}
		if apiResp.Parameters != nil {
			if m, ok := apiResp.Parameters.(map[string]interface{}); ok {
				params = m
			}
		}
		c.stats.Record(time.Since(start), false, int64(len(data)), int64(len(body)))
		return &apiResp, errors.ParseError(resp.StatusCode, apiResp.Description, params)
	}

	c.stats.Record(time.Since(start), true, int64(len(data)), int64(len(body)))
	return &apiResp, nil
}

// RequestWithMultipart ارسال multipart واقعی (برای آپلود فایل)
func (c *Client) RequestWithMultipart(ctx context.Context, method string, fields map[string]string, files map[string]string) (*Response, error) {
	url := c.buildURL(method)

	if len(files) == 0 {
		return c.RequestWithParams(ctx, method, fields)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			return nil, fmt.Errorf("failed to write field %s: %w", k, err)
		}
	}

	for fieldName, filePath := range files {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}

		fi, _ := f.Stat()
		h := textproto.MIMEHeader{}
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fi.Name()))
		h.Set("Content-Type", detectContentType(filePath))

		part, err := writer.CreatePart(h)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("failed to create part: %w", err)
		}
		_, err = io.Copy(part, f)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to copy file %s: %w", filePath, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	start := time.Now()

	if !c.breaker.Allow() {
		return nil, errors.NewNetworkError(fmt.Errorf("circuit breaker open"))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		c.breaker.RecordFailure()
		return nil, fmt.Errorf("failed to create multipart request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.breaker.RecordFailure()
		c.stats.Record(time.Since(start), false, 0, int64(buf.Len()))
		return nil, errors.NewNetworkError(err)
	}
	c.breaker.RecordSuccess()
	defer resp.Body.Close()

	limited := io.LimitReader(resp.Body, 100<<20)
	data, err := io.ReadAll(limited)
	if err != nil {
		c.stats.Record(time.Since(start), false, int64(len(data)), int64(buf.Len()))
		return nil, fmt.Errorf("failed to read multipart response: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(data, &apiResp); err != nil {
		c.stats.Record(time.Since(start), false, int64(len(data)), int64(buf.Len()))
		return nil, fmt.Errorf("failed to decode multipart response: %w (body: %s)", err, truncate(string(data), 500))
	}
	apiResp.HTTPStatus = resp.StatusCode

	if !apiResp.Ok {
		params := map[string]interface{}{}
		if apiResp.Parameters != nil {
			if m, ok := apiResp.Parameters.(map[string]interface{}); ok {
				params = m
			}
		}
		c.stats.Record(time.Since(start), false, int64(len(data)), int64(buf.Len()))
		return &apiResp, errors.ParseError(resp.StatusCode, apiResp.Description, params)
	}

	c.stats.Record(time.Since(start), true, int64(len(data)), int64(buf.Len()))
	return &apiResp, nil
}

// RequestWithForm ارسال multipart (برای آپلود فایل)
func (c *Client) RequestWithForm(ctx context.Context, method string, form map[string]string, files map[string]string) (*Response, error) {
	return c.RequestWithMultipart(ctx, method, form, files)
}

// Close بستن اتصالات
func (c *Client) Close() error {
	c.transport.CloseIdleConnections()
	return nil
}

// RequestID نوع کلید برای context
type requestIDKey struct{}

// WithRequestID افزودن Request ID به context
func WithRequestID(ctx context.Context, id string) context.Context {
	if id == "" {
		id = uuid.NewString()
	}
	return context.WithValue(ctx, requestIDKey{}, id)
}

// GetRequestID گرفتن Request ID از context
func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey{}).(string); ok {
		return v
	}
	return ""
}

// detectContentType تشخیص Content-Type بر اساس پسوند
func detectContentType(path string) string {
	switch lowerExt(path) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mp3":
		return "audio/mpeg"
	case ".ogg":
		return "audio/ogg"
	case ".wav":
		return "audio/wav"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	case ".webm-sticker":
		return "application/octet-stream"
	case ".tgs":
		return "application/x-tgsticker"
	default:
		return "application/octet-stream"
	}
}

func lowerExt(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '.' {
			return strings.ToLower(p[i:])
		}
	}
	return ""
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
