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
				return c.Status(http.StatusTooManyRequests).JSON(Error{
					Code:    "RATE_LIMIT_EXCEEDED",
					Message: "Too many requests",
				})
			}
			return next(c)
		}
	}
}
