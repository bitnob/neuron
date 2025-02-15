package middleware

import (
	"net/http"
	"time"
)

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
			recorder := newResponseRecorder(c.Response())
			c.Response().Writer = recorder

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
