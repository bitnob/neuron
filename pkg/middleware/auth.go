package middleware

import (
	"context"
	"net/http"
	"neuron/pkg/router"
	"strings"
)

// Use router types instead of local types
type (
	HandlerFunc    = router.HandlerFunc
	MiddlewareFunc = router.MiddlewareFunc
	Context        = router.Context
)

// Add Error struct
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AuthConfig struct {
	TokenType      string
	HeaderName     string
	QueryParam     string
	ContextKey     string
	SkipPaths      []string
	TokenValidator func(string) (interface{}, error)
}

func extractToken(r *http.Request, config AuthConfig) string {
	// Check header
	if config.HeaderName != "" {
		if token := r.Header.Get(config.HeaderName); token != "" {
			return strings.TrimPrefix(token, config.TokenType+" ")
		}
	}

	// Check query parameter
	if config.QueryParam != "" {
		if token := r.URL.Query().Get(config.QueryParam); token != "" {
			return token
		}
	}

	return ""
}

func NewAuthMiddleware(config AuthConfig) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			// Skip authentication for specified paths
			for _, path := range config.SkipPaths {
				if strings.HasPrefix(c.Request.URL.Path, path) {
					return next(c)
				}
			}

			token := extractToken(c.Request, config)
			if token == "" {
				return c.JSON(http.StatusUnauthorized, Error{
					Code:    "UNAUTHORIZED",
					Message: "Authentication required",
				})
			}

			// Validate token
			claims, err := config.TokenValidator(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, Error{
					Code:    "INVALID_TOKEN",
					Message: "Invalid authentication token",
				})
			}

			// Set claims in context
			ctx := context.WithValue(c.Request.Context(), config.ContextKey, claims)
			c.Request = c.Request.WithContext(ctx)

			return next(c)
		}
	}
}
