package cache

import (
	"sync"
	"time"
)

// Cache کش در حافظه با TTL — thread-safe و بهینه
type Cache struct {
	mu    sync.RWMutex
	items map[string]*item
	ttl   time.Duration
	stop  chan struct{}
}

type item struct {
	value     interface{}
	expiresAt time.Time
	ttl       time.Duration // برای TTL اختصاصی
}

// New ساخت کش جدید
func New(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]*item),
		ttl:   ttl,
		stop:  make(chan struct{}),
	}
	go c.cleanupLoop()
	return c
}

// Get دریافت مقدار
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().After(it.expiresAt) {
		// Lazy deletion
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return nil, false
	}
	return it.value, true
}

// GetTyped دریافت مقدار با تایپ مشخص
func GetTyped[T any](c *Cache, key string) (T, bool) {
	v, ok := c.Get(key)
	if !ok {
		var zero T
		return zero, false
	}
	t, ok := v.(T)
	if !ok {
		var zero T
		return zero, false
	}
	return t, true
}

// Set ذخیره مقدار با TTL پیش‌فرض
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &item{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
		ttl:       c.ttl,
	}
}

// SetWithTTL ذخیره با TTL اختصاصی
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
		ttl:       ttl,
	}
}

// SetForever ذخیره بدون انقضا
func (c *Cache) SetForever(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &item{
		value:     value,
		expiresAt: time.Now().AddDate(100, 0, 0),
		ttl:       0,
	}
}

// Delete حذف مقدار
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Len تعداد آیتم‌ها
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Keys لیست کلیدها
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

// Clear پاک کردن کل کش
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*item)
}

// GetOrSet دریافت یا تنظیم با تابع
func (c *Cache) GetOrSet(key string, fn func() interface{}) interface{} {
	if v, ok := c.Get(key); ok {
		return v
	}
	v := fn()
	c.Set(key, v)
	return v
}

// GetOrSetWithTTL دریافت یا تنظیم با TTL اختصاصی
func (c *Cache) GetOrSetWithTTL(key string, ttl time.Duration, fn func() interface{}) interface{} {
	if v, ok := c.Get(key); ok {
		return v
	}
	v := fn()
	c.SetWithTTL(key, v, ttl)
	return v
}

func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stop:
			return
		}
	}
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, v := range c.items {
		if now.After(v.expiresAt) {
			delete(c.items, k)
		}
	}
}

// Close توقف cleanup
func (c *Cache) Close() {
	select {
	case <-c.stop:
		// already closed
	default:
		close(c.stop)
	}
}

// ShardedCache کش چند-بخشی برای کاهش lock contention
type ShardedCache struct {
	shards []*Cache
}

// NewSharded ساخت کش چند-بخشی برای سرعت بالا
func NewSharded(shards int, ttl time.Duration) *ShardedCache {
	if shards <= 0 {
		shards = 32
	}
	s := &ShardedCache{shards: make([]*Cache, shards)}
	for i := range s.shards {
		s.shards[i] = New(ttl)
	}
	return s
}

func (s *ShardedCache) getShard(key string) *Cache {
	h := fnv32(key)
	return s.shards[int(h)%len(s.shards)]
}

func (s *ShardedCache) Get(key string) (interface{}, bool) { return s.getShard(key).Get(key) }
func (s *ShardedCache) Set(key string, value interface{})  { s.getShard(key).Set(key, value) }
func (s *ShardedCache) Delete(key string)                  { s.getShard(key).Delete(key) }
func (s *ShardedCache) Close() {
	for _, c := range s.shards {
		c.Close()
	}
}

func fnv32(key string) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		h *= 16777619
		h ^= uint32(key[i])
	}
	return h
}
