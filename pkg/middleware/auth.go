package middleware

import (
	"context"
	"net/http"
	"strings"
)

type AuthConfig struct {
	TokenType      string
	HeaderName     string
	QueryParam     string
	ContextKey     string
	SkipPaths      []string
	TokenValidator func(string) (interface{}, error)
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
				return c.Status(http.StatusUnauthorized).JSON(Error{
					Code:    "UNAUTHORIZED",
					Message: "Authentication required",
				})
			}

			// Validate token
			claims, err := config.TokenValidator(token)
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(Error{
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
