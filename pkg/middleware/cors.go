package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

func NewCORSMiddleware(config CORSConfig) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			origin := c.Request.Header.Get("Origin")

			// Handle preflight requests
			if c.Request.Method == http.MethodOptions {
				c.Response.Header().Set("Access-Control-Allow-Origin", getAllowedOrigin(origin, config.AllowOrigins))
				c.Response.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ","))
				c.Response.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ","))
				c.Response.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))

				if config.AllowCredentials {
					c.Response.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				return c.NoContent(http.StatusNoContent)
			}

			// Set CORS headers for regular requests
			c.Response.Header().Set("Access-Control-Allow-Origin", getAllowedOrigin(origin, config.AllowOrigins))
			if config.AllowCredentials {
				c.Response.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if len(config.ExposeHeaders) > 0 {
				c.Response.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ","))
			}

			return next(c)
		}
	}
}

func getAllowedOrigin(origin string, allowedOrigins []string) string {
	if len(allowedOrigins) == 0 {
		return "*"
	}

	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return origin
		}
	}

	return allowedOrigins[0]
}

func CORS() MiddlewareFunc {
	return NewCORSMiddleware(CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{},
		AllowCredentials: true,
		MaxAge:           86400,
	})
}
