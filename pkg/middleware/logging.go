package middleware

import (
	"strings"
	"time"
)

type LogConfig struct {
	Logger        Logger
	SkipPaths     []string
	LogHeaders    bool
	LogBody       bool
	SlowThreshold time.Duration
}

func NewLoggingMiddleware(config LogConfig) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			start := time.Now()
			path := c.Request.URL.Path

			// Skip logging for specified paths
			for _, skip := range config.SkipPaths {
				if strings.HasPrefix(path, skip) {
					return next(c)
				}
			}

			// Log request
			logRequest(c.Request, config)

			// Create response recorder
			recorder := newResponseRecorder(c.Response())
			c.Response().Writer = recorder

			err := next(c)

			// Calculate duration
			duration := time.Since(start)

			// Log response
			logResponse(recorder, duration, err, config)

			return err
		}
	}
}
