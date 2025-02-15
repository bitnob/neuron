package middleware

import (
	"context"
	"net/http"
	"time"
)

// Cache defines the interface for caching implementations
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

type CacheConfig struct {
	TTL           time.Duration
	KeyPrefix     string
	IgnoreHeaders []string
	Cache         Cache
}

func NewCacheMiddleware(config CacheConfig) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			// Skip caching for non-GET requests
			if c.Request.Method != http.MethodGet {
				return next(c)
			}

			key := generateCacheKey(c.Request, config)

			// Try to get from cache
			if cached, err := config.Cache.Get(c.Request.Context(), key); err == nil {
				response := cached.([]byte)
				return c.Blob(http.StatusOK, "application/json", response)
			}

			// Create response recorder
			recorder := newResponseRecorder(c.Response)
			c.Response = recorder

			err := next(c)
			if err != nil {
				return err
			}

			// Cache the response
			if recorder.Status() == http.StatusOK {
				config.Cache.Set(c.Request.Context(), key, recorder.Body(), config.TTL)
			}

			return nil
		}
	}
}

func generateCacheKey(r *http.Request, config CacheConfig) string {
	key := config.KeyPrefix + r.URL.Path
	if r.URL.RawQuery != "" {
		key += "?" + r.URL.RawQuery
	}
	return key
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	body   []byte
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) Status() int {
	return r.status
}

func (r *responseRecorder) Body() []byte {
	return r.body
}
