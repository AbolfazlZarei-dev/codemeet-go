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

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

var jsonEncoderPool = sync.Pool{
	New: func() interface{} {
		return json.NewEncoder(nil)
	},
}

// countReader برای شمارش دقیق بایت‌های خوانده شده
type countReader struct {
	io.Reader
	n int64
}

func (r *countReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.n += int64(n)
	return n, err
}

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

type Stats struct {
	Requests     int64
	SuccessCount int64
	ErrorCount   int64
	BytesIn      int64
	BytesOut     int64
	totalLatency int64
}

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

func (s *Stats) AvgLatency() time.Duration {
	r := atomic.LoadInt64(&s.Requests)
	if r == 0 {
		return 0
	}
	return time.Duration(atomic.LoadInt64(&s.totalLatency) / r)
}

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

type StatsSnapshot struct {
	Requests     int64
	SuccessCount int64
	ErrorCount   int64
	BytesIn      int64
	BytesOut     int64
	AvgLatency   time.Duration
}

// CircuitBreaker پیاده‌سازی بهینه با Single Probe در Half-Open
type CircuitBreaker struct {
	failureThreshold int32
	resetTimeout     time.Duration
	failures         atomic.Int32
	lastFailure      atomic.Int64
	state            atomic.Int32 // 0=closed, 1=open, 2=half-open
	probeInFlight    atomic.Bool  // فقط یک پروب مجاز است
}

const (
	stateClosed   int32 = 0
	stateOpen     int32 = 1
	stateHalfOpen int32 = 2
)

func NewCircuitBreaker(threshold int, reset time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: int32(threshold),
		resetTimeout:     reset,
	}
}

func (c *CircuitBreaker) Allow() bool {
	state := c.state.Load()
	if state == stateClosed {
		return true
	}
	if state == stateOpen {
		lastFail := time.Unix(0, c.lastFailure.Load())
		if time.Since(lastFail) > c.resetTimeout {
			if c.state.CompareAndSwap(stateOpen, stateHalfOpen) {
				// فقط یک درخواست Probe اجازه عبور دارد
				if c.probeInFlight.CompareAndSwap(false, true) {
					return true
				}
			}
		}
	}
	return false
}

func (c *CircuitBreaker) RecordSuccess() {
	if c.state.Load() == stateHalfOpen {
		c.probeInFlight.Store(false)
	}
	c.failures.Store(0)
	c.state.Store(stateClosed)
}

func (c *CircuitBreaker) RecordFailure() {
	if c.state.Load() == stateHalfOpen {
		c.probeInFlight.Store(false)
	}
	f := c.failures.Add(1)
	c.lastFailure.Store(time.Now().UnixNano())
	if f >= c.failureThreshold {
		c.state.Store(stateOpen)
	}
}

func (c *CircuitBreaker) State() int {
	return int(c.state.Load())
}

func NewClient(baseURL, token string, log *logger.Logger) *Client {
	transport := &http.Transport{
		MaxIdleConns:          1000,
		MaxIdleConnsPerHost:   500,
		MaxConnsPerHost:       500,
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 90 * time.Second,
		ForceAttemptHTTP2:     true,
		DisableCompression:    false,
		Proxy:                 http.ProxyFromEnvironment,
		WriteBufferSize:       128 * 1024,
		ReadBufferSize:        128 * 1024,
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

func (c *Client) SetHTTPClient(hc *http.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.httpClient = hc
	if tr, ok := hc.Transport.(*http.Transport); ok {
		c.transport = tr
	}
}

func (c *Client) SetTimeout(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.httpClient.Timeout = d
}

func (c *Client) Stats() *Stats                { return c.stats }
func (c *Client) StatsSnapshot() StatsSnapshot { return c.stats.Snapshot() }
func (c *Client) Breaker() *CircuitBreaker     { return c.breaker }

func (c *Client) UserAgent() string {
	return fmt.Sprintf("codemeet-go/1.0.0 (%s/%s) go/%s",
		runtime.GOOS, runtime.GOARCH, runtime.Version())
}

func (c *Client) buildURL(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", c.baseURL, c.token, method)
}

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

	if resp.StatusCode >= 500 || resp.StatusCode == 429 {
		c.breaker.RecordFailure()
	} else {
		c.breaker.RecordSuccess()
	}

	return resp, nil
}

// RequestWithParams بهینه‌شده با Streaming Decode
func (c *Client) RequestWithParams(ctx context.Context, method string, params interface{}) (*Response, error) {
	url := c.buildURL(method)

	var body []byte
	if params != nil {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()
		enc := json.NewEncoder(buf)
		if err := enc.Encode(params); err != nil {
			bufPool.Put(buf)
			return nil, fmt.Errorf("failed to marshal params: %w", err)
		}
		body = buf.Bytes()
		bodyCopy := make([]byte, len(body))
		copy(bodyCopy, body)
		bufPool.Put(buf)
		body = bodyCopy
	}

	start := time.Now()

	resp, err := c.doRequest(ctx, http.MethodPost, url, body)
	if err != nil {
		c.stats.Record(time.Since(start), false, 0, int64(len(body)))
		return nil, err
	}
	defer resp.Body.Close()

	// محدود کردن به 10 مگابایت و استفاده از Streaming Decode برای کاهش RAM
	limited := io.LimitReader(resp.Body, 10<<20)
	cr := &countReader{Reader: limited}

	var apiResp Response
	if err := json.NewDecoder(cr).Decode(&apiResp); err != nil {
		c.stats.Record(time.Since(start), false, cr.n, int64(len(body)))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	apiResp.HTTPStatus = resp.StatusCode

	if !apiResp.Ok {
		params := map[string]interface{}{}
		if apiResp.Parameters != nil {
			if m, ok := apiResp.Parameters.(map[string]interface{}); ok {
				params = m
			}
		}
		c.stats.Record(time.Since(start), false, cr.n, int64(len(body)))
		return &apiResp, errors.ParseError(resp.StatusCode, apiResp.Description, params)
	}

	c.stats.Record(time.Since(start), true, cr.n, int64(len(body)))
	return &apiResp, nil
}

func (c *Client) Request(ctx context.Context, method string) (*Response, error) {
	return c.RequestWithParams(ctx, method, nil)
}

// RequestWithMultipart بهینه‌شده با io.Pipe و Streaming Decode
func (c *Client) RequestWithMultipart(ctx context.Context, method string, fields map[string]string, files map[string]string) (*Response, error) {
	url := c.buildURL(method)

	if len(files) == 0 {
		return c.RequestWithParams(ctx, method, fields)
	}

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	contentType := writer.FormDataContentType()

	go func() {
		var err error
		defer pw.Close()

		for k, v := range fields {
			if err = writer.WriteField(k, v); err != nil {
				pw.CloseWithError(err)
				return
			}
		}

		for fieldName, filePath := range files {
			var f *os.File
			f, err = os.Open(filePath)
			if err != nil {
				pw.CloseWithError(err)
				return
			}

			var fi os.FileInfo
			fi, err = f.Stat()
			if err != nil {
				f.Close()
				pw.CloseWithError(err)
				return
			}

			h := textproto.MIMEHeader{}
			h.Set("Content-Disposition",
				fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fi.Name()))
			h.Set("Content-Type", detectContentType(filePath))

			var part io.Writer
			part, err = writer.CreatePart(h)
			if err != nil {
				f.Close()
				pw.CloseWithError(err)
				return
			}

			_, err = io.Copy(part, f)
			f.Close()
			if err != nil {
				pw.CloseWithError(err)
				return
			}
		}

		if err = writer.Close(); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	start := time.Now()

	if !c.breaker.Allow() {
		return nil, errors.NewNetworkError(fmt.Errorf("circuit breaker open"))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, pr)
	if err != nil {
		c.breaker.RecordFailure()
		return nil, fmt.Errorf("failed to create multipart request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.breaker.RecordFailure()
		c.stats.Record(time.Since(start), false, 0, 0)
		return nil, errors.NewNetworkError(err)
	}

	if resp.StatusCode >= 500 || resp.StatusCode == 429 {
		c.breaker.RecordFailure()
	} else {
		c.breaker.RecordSuccess()
	}

	defer resp.Body.Close()

	limited := io.LimitReader(resp.Body, 10<<20)
	cr := &countReader{Reader: limited}

	var apiResp Response
	if err := json.NewDecoder(cr).Decode(&apiResp); err != nil {
		c.stats.Record(time.Since(start), false, cr.n, 0)
		return nil, fmt.Errorf("failed to decode multipart response: %w", err)
	}
	apiResp.HTTPStatus = resp.StatusCode

	if !apiResp.Ok {
		params := map[string]interface{}{}
		if apiResp.Parameters != nil {
			if m, ok := apiResp.Parameters.(map[string]interface{}); ok {
				params = m
			}
		}
		c.stats.Record(time.Since(start), false, cr.n, 0)
		return &apiResp, errors.ParseError(resp.StatusCode, apiResp.Description, params)
	}

	c.stats.Record(time.Since(start), true, cr.n, 0)
	return &apiResp, nil
}

func (c *Client) RequestWithForm(ctx context.Context, method string, form map[string]string, files map[string]string) (*Response, error) {
	return c.RequestWithMultipart(ctx, method, form, files)
}

func (c *Client) Close() error {
	c.transport.CloseIdleConnections()
	return nil
}

type requestIDKey struct{}

func WithRequestID(ctx context.Context, id string) context.Context {
	if id == "" {
		id = uuid.NewString()
	}
	return context.WithValue(ctx, requestIDKey{}, id)
}

func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey{}).(string); ok {
		return v
	}
	return ""
}

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
