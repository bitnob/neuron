package cache

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExpired  = errors.New("key expired")
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

func (c *MemoryCache) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now().UnixNano()
		for key, item := range c.items {
			if item.Expiration > 0 && item.Expiration < now {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// evict removes the oldest item from the cache
func (c *MemoryCache) evict() {
	var oldestKey string
	var oldestTime int64 = math.MaxInt64

	for key, item := range c.items {
		if item.Expiration < oldestTime {
			oldestTime = item.Expiration
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}

// Clear removes all items from the cache
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*Item)
	return nil
}

func (c *MemoryCache) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return 0, ErrKeyNotFound
	}

	current, ok := item.Value.(int64)
	if !ok {
		return 0, errors.New("value is not an integer")
	}

	current -= value
	item.Value = current
	return current, nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *MemoryCache) DeleteMany(ctx context.Context, keys []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		delete(c.items, key)
	}
	return nil
}

func (c *MemoryCache) GetMany(ctx context.Context, keys []string) (map[string]interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	results := make(map[string]interface{}, len(keys))
	for _, key := range keys {
		if item, found := c.items[key]; found {
			if item.Expiration == 0 || item.Expiration > time.Now().UnixNano() {
				results[key] = item.Value
			}
		}
	}
	return results, nil
}

// Increment atomically increments a numeric value
func (c *MemoryCache) Increment(ctx context.Context, key string, value int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return 0, ErrKeyNotFound
	}

	current, ok := item.Value.(int64)
	if !ok {
		return 0, errors.New("value is not an integer")
	}

	current += value
	item.Value = current
	return current, nil
}

// SetMany sets multiple key-value pairs
func (c *MemoryCache) SetMany(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	for key, value := range items {
		if c.maxItems > 0 && len(c.items) >= c.maxItems {
			c.evict()
		}
		c.items[key] = &Item{
			Value:      value,
			Expiration: exp,
		}
	}
	return nil
}

// Tags returns all items with the given tags
func (c *MemoryCache) Tags(tags ...string) TaggedCache {
	return &MemoryTaggedCache{
		cache: c,
		tags:  tags,
	}
}

// MemoryTaggedCache implements TaggedCache for memory cache
type MemoryTaggedCache struct {
	cache *MemoryCache
	tags  []string
}

func (t *MemoryTaggedCache) Get(ctx context.Context, key string) (interface{}, error) {
	item, err := t.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	// Check if item has all required tags
	if cacheItem, ok := t.cache.items[key]; ok {
		if !t.hasAllTags(cacheItem.Tags) {
			return nil, ErrKeyNotFound
		}
	}

	return item, nil
}

func (t *MemoryTaggedCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	t.cache.mu.Lock()
	defer t.cache.mu.Unlock()

	if t.cache.maxItems > 0 && len(t.cache.items) >= t.cache.maxItems {
		t.cache.evict()
	}

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	t.cache.items[key] = &Item{
		Value:      value,
		Expiration: exp,
		Tags:       t.tags,
	}

	return nil
}

func (t *MemoryTaggedCache) hasAllTags(itemTags []string) bool {
	if len(t.tags) == 0 {
		return true
	}
	tagMap := make(map[string]bool)
	for _, tag := range itemTags {
		tagMap[tag] = true
	}
	for _, tag := range t.tags {
		if !tagMap[tag] {
			return false
		}
	}
	return true
}

func (t *MemoryTaggedCache) Clear(ctx context.Context) error {
	t.cache.mu.Lock()
	defer t.cache.mu.Unlock()

	// Remove all items that have all the required tags
	for key, item := range t.cache.items {
		if t.hasAllTags(item.Tags) {
			delete(t.cache.items, key)
		}
	}
	return nil
}

func (t *MemoryTaggedCache) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	if item, ok := t.cache.items[key]; ok {
		if !t.hasAllTags(item.Tags) {
			return 0, ErrKeyNotFound
		}
	}
	return t.cache.Decrement(ctx, key, value)
}

func (t *MemoryTaggedCache) Delete(ctx context.Context, key string) error {
	if item, ok := t.cache.items[key]; ok {
		if !t.hasAllTags(item.Tags) {
			return ErrKeyNotFound
		}
	}
	return t.cache.Delete(ctx, key)
}

func (t *MemoryTaggedCache) DeleteMany(ctx context.Context, keys []string) error {
	t.cache.mu.Lock()
	defer t.cache.mu.Unlock()

	for _, key := range keys {
		if item, ok := t.cache.items[key]; ok {
			if t.hasAllTags(item.Tags) {
				delete(t.cache.items, key)
			}
		}
	}
	return nil
}

func (t *MemoryTaggedCache) Flush(ctx context.Context) error {
	return t.Clear(ctx)
}

func (t *MemoryTaggedCache) GetMany(ctx context.Context, keys []string) (map[string]interface{}, error) {
	t.cache.mu.RLock()
	defer t.cache.mu.RUnlock()

	results := make(map[string]interface{})
	for _, key := range keys {
		if item, found := t.cache.items[key]; found {
			if item.Expiration == 0 || item.Expiration > time.Now().UnixNano() {
				if t.hasAllTags(item.Tags) {
					results[key] = item.Value
				}
			}
		}
	}
	return results, nil
}

func (t *MemoryTaggedCache) Increment(ctx context.Context, key string, value int64) (int64, error) {
	if item, ok := t.cache.items[key]; ok {
		if !t.hasAllTags(item.Tags) {
			return 0, ErrKeyNotFound
		}
	}
	return t.cache.Increment(ctx, key, value)
}

func (t *MemoryTaggedCache) Remember(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	// Try to get from cache first
	if value, err := t.Get(ctx, key); err == nil {
		return value, nil
	}

	// Generate value using callback
	value, err := callback()
	if err != nil {
		return nil, err
	}

	// Store in cache
	if err := t.Set(ctx, key, value, ttl); err != nil {
		return nil, err
	}

	return value, nil
}

func (t *MemoryTaggedCache) SetMany(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	t.cache.mu.Lock()
	defer t.cache.mu.Unlock()

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	for key, value := range items {
		if t.cache.maxItems > 0 && len(t.cache.items) >= t.cache.maxItems {
			t.cache.evict()
		}
		t.cache.items[key] = &Item{
			Value:      value,
			Expiration: exp,
			Tags:       t.tags,
		}
	}
	return nil
}

func (t *MemoryTaggedCache) Tags(tags ...string) TaggedCache {
	return &MemoryTaggedCache{
		cache: t.cache,
		tags:  append(t.tags, tags...),
	}
}

func (t *MemoryTaggedCache) WithPrefix(prefix string) Cache {
	t.cache.options.Prefix = prefix
	return t.cache
}

func (c *MemoryCache) Remember(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	// Try to get from cache first
	if value, err := c.Get(ctx, key); err == nil {
		return value, nil
	}

	// Generate value using callback
	value, err := callback()
	if err != nil {
		return nil, err
	}

	// Store in cache
	if err := c.Set(ctx, key, value, ttl); err != nil {
		return nil, err
	}

	return value, nil
}

func (c *MemoryCache) WithPrefix(prefix string) Cache {
	c.options.Prefix = prefix
	return c
}
