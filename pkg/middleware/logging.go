package middleware

import (
	"net/http"
	"strings"
	"time"
)

// Logger interface for middleware logging
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

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
			recorder := newResponseRecorder(c.Response)
			c.Response = recorder

			err := next(c)

			// Calculate duration
			duration := time.Since(start)

			// Log response
			logResponse(recorder, duration, err, config)

			return err
		}
	}
}

func logRequest(r *http.Request, config LogConfig) {
	fields := []interface{}{
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	}

	if config.LogHeaders {
		fields = append(fields, "headers", r.Header)
	}

	config.Logger.Info("Request", fields...)
}

func logResponse(recorder *responseRecorder, duration time.Duration, err error, config LogConfig) {
	fields := []interface{}{
		"status", recorder.Status(),
		"duration", duration.String(),
	}

	if err != nil {
		fields = append(fields, "error", err.Error())
		config.Logger.Error("Response", fields...)
		return
	}

	if duration > config.SlowThreshold {
		config.Logger.Warn("Slow Response", fields...)
		return
	}

	config.Logger.Info("Response", fields...)
}
