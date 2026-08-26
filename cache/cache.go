package cache

import (
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type Cache struct {
	mu    sync.RWMutex
	items map[string]*item
	ttl   time.Duration
	stop  chan struct{}
	sf    singleflight.Group
}

type item struct {
	value     interface{}
	expiresAt time.Time
	ttl       time.Duration
}

func New(ttl time.Duration) *Cache {
	return newCacheInternal(ttl, true)
}

// newCacheInternal ساخت کش با قابلیت کنترل راه‌اندازی Goroutine پاکسازی
func newCacheInternal(ttl time.Duration, startCleanup bool) *Cache {
	c := &Cache{
		items: make(map[string]*item),
		ttl:   ttl,
		stop:  make(chan struct{}),
	}
	if startCleanup {
		go c.cleanupLoop()
	}
	return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if !it.expiresAt.IsZero() && time.Now().After(it.expiresAt) {
		c.mu.Lock()
		if it, ok := c.items[key]; ok && !it.expiresAt.IsZero() && time.Now().After(it.expiresAt) {
			delete(c.items, key)
		}
		c.mu.Unlock()
		return nil, false
	}
	return it.value, true
}

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

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &item{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
		ttl:       c.ttl,
	}
}

func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
		ttl:       ttl,
	}
}

func (c *Cache) SetForever(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &item{
		value:     value,
		expiresAt: time.Time{},
		ttl:       0,
	}
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*item)
}

func (c *Cache) GetOrSet(key string, fn func() interface{}) interface{} {
	if v, ok := c.Get(key); ok {
		return v
	}

	v, _, _ := c.sf.Do(key, func() (interface{}, error) {
		if v, ok := c.Get(key); ok {
			return v, nil
		}
		val := fn()
		c.Set(key, val)
		return val, nil
	})
	return v
}

func (c *Cache) GetOrSetWithTTL(key string, ttl time.Duration, fn func() interface{}) interface{} {
	if v, ok := c.Get(key); ok {
		return v
	}

	v, _, _ := c.sf.Do(key, func() (interface{}, error) {
		if v, ok := c.Get(key); ok {
			return v, nil
		}
		val := fn()
		c.SetWithTTL(key, val, ttl)
		return val, nil
	})
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
		if !v.expiresAt.IsZero() && now.After(v.expiresAt) {
			delete(c.items, k)
		}
	}
}

func (c *Cache) Close() {
	select {
	case <-c.stop:
	default:
		close(c.stop)
	}
}

// ShardedCache با Central Scheduler برای کاهش Goroutineها
type ShardedCache struct {
	shards []*Cache
	stop   chan struct{}
}

func NewSharded(shards int, ttl time.Duration) *ShardedCache {
	if shards <= 0 {
		shards = 32
	}
	s := &ShardedCache{
		shards: make([]*Cache, shards),
		stop:   make(chan struct{}),
	}
	// ساخت Shard ها بدون Goroutine اختصاصی
	for i := range s.shards {
		s.shards[i] = newCacheInternal(ttl, false)
	}
	// یک Goroutine مرکزی برای پاکسازی همه Shardها
	go s.cleanupLoop()
	return s
}

func (s *ShardedCache) getShard(key string) *Cache {
	h := fnv32(key)
	return s.shards[int(h)%len(s.shards)]
}

func (s *ShardedCache) Get(key string) (interface{}, bool) { return s.getShard(key).Get(key) }
func (s *ShardedCache) Set(key string, value interface{})  { s.getShard(key).Set(key, value) }
func (s *ShardedCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	s.getShard(key).SetWithTTL(key, value, ttl)
}
func (s *ShardedCache) SetForever(key string, value interface{}) {
	s.getShard(key).SetForever(key, value)
}
func (s *ShardedCache) Delete(key string) { s.getShard(key).Delete(key) }
func (s *ShardedCache) GetOrSet(key string, fn func() interface{}) interface{} {
	return s.getShard(key).GetOrSet(key, fn)
}
func (s *ShardedCache) GetOrSetWithTTL(key string, ttl time.Duration, fn func() interface{}) interface{} {
	return s.getShard(key).GetOrSetWithTTL(key, ttl, fn)
}
func (s *ShardedCache) Len() int {
	total := 0
	for _, c := range s.shards {
		total += c.Len()
	}
	return total
}
func (s *ShardedCache) Keys() []string {
	var keys []string
	for _, c := range s.shards {
		keys = append(keys, c.Keys()...)
	}
	return keys
}
func (s *ShardedCache) Clear() {
	for _, c := range s.shards {
		c.Clear()
	}
}

// Central Scheduler برای پاکسازی دوره‌ای
func (s *ShardedCache) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for _, c := range s.shards {
				c.cleanup()
			}
		case <-s.stop:
			return
		}
	}
}

func (s *ShardedCache) Close() {
	select {
	case <-s.stop:
	default:
		close(s.stop)
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
