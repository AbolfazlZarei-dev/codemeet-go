package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/api"
	"github.com/AbolfazlZarei-dev/codemeet-go/logger"
	"github.com/gorilla/websocket"
)

// Hub مدیریت اتصالات WebSocket
type Hub struct {
	api    *api.Client
	logger *logger.Logger
	mu     sync.RWMutex
	conns  map[string]*wsConn
	subs   map[string][]*subscriber

	reconnectEnabled bool
	reconnectDelay   time.Duration
	maxReconnect     int

	messages int64
	errors   int64
}

type wsConn struct {
	conn       *websocket.Conn
	url        string
	closed     int32
	reconnects int
}

type subscriber struct {
	ch     chan *Event
	closed bool
}

// Event رویداد WebSocket
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// NewHub ساخت Hub
func NewHub(c *api.Client, log *logger.Logger) *Hub {
	return &Hub{
		api:    c,
		logger: log,
		conns:  make(map[string]*wsConn),
		subs:   make(map[string][]*subscriber),

		reconnectEnabled: true,
		reconnectDelay:   5 * time.Second,
		maxReconnect:     10,
	}
}

// SetReconnect تنظیم پارامترهای reconnect
func (h *Hub) SetReconnect(enabled bool, delay time.Duration, maxRetries int) {
	h.reconnectEnabled = enabled
	h.reconnectDelay = delay
	h.maxReconnect = maxRetries
}

// Connect اتصال به سرور WebSocket
func (h *Hub) Connect(ctx context.Context, url string) error {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return err
	}
	w := &wsConn{conn: conn, url: url}
	h.mu.Lock()
	h.conns[url] = w
	h.mu.Unlock()

	go h.readLoop(ctx, w)
	go h.pingLoop(ctx, w)
	return nil
}

func (h *Hub) readLoop(ctx context.Context, w *wsConn) {
	defer func() {
		w.conn.Close()
		atomic.StoreInt32(&w.closed, 1)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, data, err := w.conn.ReadMessage()
		if err != nil {
			atomic.AddInt64(&h.errors, 1)
			if h.logger != nil {
				h.logger.Error("ws read error", "url", w.url, "error", err)
			}

			// اگر context cancel نشده، تلاش برای reconnect
			if ctx.Err() == nil && h.reconnectEnabled && w.reconnects < h.maxReconnect {
				w.reconnects++
				if h.reconnect(w, ctx) {
					w.reconnects = 0
					continue
				}
			}
			return
		}

		atomic.AddInt64(&h.messages, 1)
		var event Event
		if err := json.Unmarshal(data, &event); err != nil {
			continue
		}
		h.publish(&event)
	}
}

func (h *Hub) reconnect(w *wsConn, ctx context.Context) bool {
	time.Sleep(h.reconnectDelay)

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, w.url, nil)
	if err != nil {
		if h.logger != nil {
			h.logger.Error("ws reconnect failed", "url", w.url, "error", err)
		}
		return false
	}

	w.conn = conn
	atomic.StoreInt32(&w.closed, 0)
	if h.logger != nil {
		h.logger.Info("ws reconnected", "url", w.url)
	}
	return true
}

func (h *Hub) pingLoop(ctx context.Context, w *wsConn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.conn.WriteControl(websocket.PingMessage, nil,
				time.Now().Add(5*time.Second)); err != nil {
				return
			}
		}
	}
}

func (h *Hub) publish(e *Event) {
	h.mu.RLock()
	subs := h.subs[e.Type]
	h.mu.RUnlock()
	for _, s := range subs {
		if s.closed {
			continue
		}
		select {
		case s.ch <- e:
		default:
			// channel پر است — دور نریز، سعی کن بعدا
		}
	}
}

// Subscribe اشتراک در رویداد
func (h *Hub) Subscribe(eventType string) <-chan *Event {
	ch := make(chan *Event, 100)
	s := &subscriber{ch: ch}
	h.mu.Lock()
	h.subs[eventType] = append(h.subs[eventType], s)
	h.mu.Unlock()
	return ch
}

// Unsubscribe لغو اشتراک
func (h *Hub) Unsubscribe(eventType string, ch <-chan *Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	subs := h.subs[eventType]
	for i, s := range subs {
		if s.ch == ch {
			s.closed = true
			close(s.ch)
			h.subs[eventType] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

// Publish ارسال رویداد به همه اتصالات
func (h *Hub) Publish(ctx context.Context, e *Event) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, w := range h.conns {
		if err := w.conn.WriteJSON(e); err != nil {
			atomic.AddInt64(&h.errors, 1)
			return err
		}
	}
	return nil
}

// Close بستن همه اتصالات
func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for url, w := range h.conns {
		w.conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		w.conn.Close()
		atomic.StoreInt32(&w.closed, 1)
		delete(h.conns, url)
	}
	// بستن تمام subscriber ها
	for eventType, subs := range h.subs {
		for _, s := range subs {
			if !s.closed {
				close(s.ch)
			}
		}
		delete(h.subs, eventType)
	}
}

// Ping ارسال ping برای بررسی اتصال
func (h *Hub) Ping(url string) error {
	h.mu.RLock()
	w, ok := h.conns[url]
	h.mu.RUnlock()
	if !ok {
		return fmt.Errorf("connection not found for url: %s", url)
	}
	return w.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second))
}

// IsConnected آیا اتصال برقرار است
func (h *Hub) IsConnected(url string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	w, ok := h.conns[url]
	if !ok {
		return false
	}
	return atomic.LoadInt32(&w.closed) == 0
}

// Stats آمار
func (h *Hub) Stats() (messages, errors int64) {
	return atomic.LoadInt64(&h.messages), atomic.LoadInt64(&h.errors)
}

// Send ارسال متن به یک URL خاص
func (h *Hub) Send(ctx context.Context, url string, message interface{}) error {
	h.mu.RLock()
	w, ok := h.conns[url]
	h.mu.RUnlock()
	if !ok {
		return fmt.Errorf("connection not found for url: %s", url)
	}
	return w.conn.WriteJSON(message)
}
