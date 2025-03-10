package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type RateConfig struct {
	Limit        rate.Limit
	Burst        int
	StoreType    string // "memory" or "redis"
	KeyFunc      func(*Context) string
	ErrorHandler func(*Context, error) error
}

type RateLimiter struct {
	config   RateConfig
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
}

func NewRateMiddleware(config RateConfig) MiddlewareFunc {
	limiter := &RateLimiter{
		config:   config,
		limiters: make(map[string]*rate.Limiter),
	}

	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			key := config.KeyFunc(c)
			if !limiter.Allow(key) {
				return c.JSON(http.StatusTooManyRequests, Error{
					Code:    "RATE_LIMIT_EXCEEDED",
					Message: "Too many requests",
				})
			}
			return next(c)
		}
	}
}

// Allow checks if a request is allowed based on the rate limit
func (l *RateLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, exists := l.limiters[key]
	if !exists {
		b = rate.NewLimiter(l.config.Limit, l.config.Burst)
		l.limiters[key] = b
	}

	return b.Allow()
}
