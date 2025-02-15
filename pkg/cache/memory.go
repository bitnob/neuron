package cache

import (
	"context"
	"sync"
	"time"
)

type MemoryCache struct {
	mu       sync.RWMutex
	items    map[string]*Item
	maxItems int
	options  Options
}

type Item struct {
	Value      interface{}
	Expiration int64
	Tags       []string
}

func NewMemoryCache(opts Options) *MemoryCache {
	cache := &MemoryCache{
		items:    make(map[string]*Item),
		maxItems: opts.MaxEntries,
		options:  opts,
	}

	// Start cleanup routine
	go cache.cleanupLoop()

	return cache
}

func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, ErrKeyNotFound
	}

	if item.Expiration > 0 && item.Expiration < time.Now().UnixNano() {
		return nil, ErrKeyExpired
	}

	return item.Value, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.maxItems > 0 && len(c.items) >= c.maxItems {
		c.evict()
	}

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	c.items[key] = &Item{
		Value:      value,
		Expiration: exp,
	}

	return nil
}
