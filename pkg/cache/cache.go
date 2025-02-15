// pkg/cache/cache.go
package cache

import (
	"context"
	"fmt"
	"time"
)

// Cache defines the interface for all cache implementations
type Cache interface {
	// Basic operations
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error

	// Bulk operations
	GetMany(ctx context.Context, keys []string) (map[string]interface{}, error)
	SetMany(ctx context.Context, items map[string]interface{}, ttl time.Duration) error
	DeleteMany(ctx context.Context, keys []string) error

	// Advanced features
	Remember(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error)) (interface{}, error)
	Tags(tags ...string) TaggedCache
	WithPrefix(prefix string) Cache
	Increment(ctx context.Context, key string, value int64) (int64, error)
	Decrement(ctx context.Context, key string, value int64) (int64, error)
}

// TaggedCache interface for tag-based caching
type TaggedCache interface {
	Cache
	Flush(ctx context.Context) error
}

// Options for cache configuration
type Options struct {
	Prefix          string
	DefaultTTL      time.Duration
	MaxEntries      int
	MaxMemory       int64
	Compression     bool
	SerializeFunc   func(interface{}) ([]byte, error)
	DeserializeFunc func([]byte) (interface{}, error)
}

// CacheError represents cache-specific errors
type CacheError struct {
	Op  string
	Key string
	Err error
}

func (e *CacheError) Error() string {
	return fmt.Sprintf("cache %s failed for key '%s': %v", e.Op, e.Key, e.Err)
}

// Factory for creating cache instances
type Factory struct {
	providers map[string]func(Options) (Cache, error)
}

func NewFactory() *Factory {
	return &Factory{
		providers: make(map[string]func(Options) (Cache, error)),
	}
}
