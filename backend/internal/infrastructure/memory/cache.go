package memory

import (
	"sync"
	"sync/atomic"
	"time"
)

type cacheItem struct {
	value      interface{}
	expiration int64
	lastAccess int64
}

// MemoryCache is a high-performance thread-safe in-memory cache.
type MemoryCache struct {
	mu          sync.RWMutex
	items       map[string]*cacheItem
	hits        int64
	misses      int64
	evictions   int64
	defaultTTL  time.Duration
	stopCleanup chan struct{}
}

// NewMemoryCache creates a new in-memory cache instance.
func NewMemoryCache(defaultTTL time.Duration, cleanupInterval time.Duration) *MemoryCache {
	c := &MemoryCache{
		items:       make(map[string]*cacheItem),
		defaultTTL:  defaultTTL,
		stopCleanup: make(chan struct{}),
	}

	go c.startCleanupLoop(cleanupInterval)
	return c
}

// Set stores a key-value pair with TTL.
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	if ttl <= 0 {
		ttl = c.defaultTTL
	}

	now := time.Now().UnixNano()
	exp := now + ttl.Nanoseconds()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &cacheItem{
		value:      value,
		expiration: exp,
		lastAccess: now,
	}
}

// Get retrieves a key from cache.
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, exists := c.items[key]
	if !exists {
		c.mu.RUnlock()
		atomic.AddInt64(&c.misses, 1)
		return nil, false
	}

	now := time.Now().UnixNano()
	if item.expiration > 0 && now > item.expiration {
		c.mu.RUnlock()
		c.Delete(key)
		atomic.AddInt64(&c.misses, 1)
		atomic.AddInt64(&c.evictions, 1)
		return nil, false
	}

	atomic.StoreInt64(&item.lastAccess, now)
	val := item.value
	c.mu.RUnlock()

	atomic.AddInt64(&c.hits, 1)
	return val, true
}

// Delete removes a key.
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear flushes all items.
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*cacheItem)
}

// Stats returns cache metrics.
func (c *MemoryCache) Stats() map[string]interface{} {
	c.mu.RLock()
	itemCount := len(c.items)
	c.mu.RUnlock()

	hits := atomic.LoadInt64(&c.hits)
	misses := atomic.LoadInt64(&c.misses)
	evictions := atomic.LoadInt64(&c.evictions)
	totalOps := hits + misses
	hitRatio := 0.0
	if totalOps > 0 {
		hitRatio = float64(hits) / float64(totalOps) * 100.0
	}

	return map[string]interface{}{
		"itemCount": itemCount,
		"hits":      hits,
		"misses":    misses,
		"evictions": evictions,
		"hitRatio":  hitRatio,
	}
}

// Close stops the cleanup worker.
func (c *MemoryCache) Close() {
	select {
	case <-c.stopCleanup:
		return
	default:
		close(c.stopCleanup)
	}
}

func (c *MemoryCache) startCleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCleanup:
			return
		case <-ticker.C:
			c.evictExpired()
		}
	}
}

func (c *MemoryCache) evictExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, item := range c.items {
		if item.expiration > 0 && now > item.expiration {
			delete(c.items, k)
			atomic.AddInt64(&c.evictions, 1)
		}
	}
}
