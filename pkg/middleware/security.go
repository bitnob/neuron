package middleware

import "fmt"

// Package middleware provides a collection of middleware components for web applications.

// SecurityConfig defines the configuration options for the security middleware.
type SecurityConfig struct {
	// HSTS enables HTTP Strict Transport Security
	HSTS bool

	// HSTSMaxAge sets the max-age directive for HSTS in seconds
	HSTSMaxAge int

	// HSTSIncludeSubdomains enables includeSubDomains directive for HSTS
	HSTSIncludeSubdomains bool

	// FrameOptions sets the X-Frame-Options header
	// Valid values: "DENY", "SAMEORIGIN", or "ALLOW-FROM uri"
	FrameOptions string

	// ContentTypeOptions sets the X-Content-Type-Options header
	// Recommended value: "nosniff"
	ContentTypeOptions string

	// XSSProtection sets the X-XSS-Protection header
	// Recommended value: "1; mode=block"
	XSSProtection string

	// CSPDirectives defines Content Security Policy directives
	// Example: {"default-src": ["'self'"], "script-src": ["'self'", "trusted-scripts.com"]}
	CSPDirectives map[string][]string
}

// NewSecurityMiddleware creates a new middleware that adds security headers to responses.
//
// Example usage:
//
//	app.Use(middleware.NewSecurityMiddleware(middleware.SecurityConfig{
//		HSTS: true,
//		HSTSMaxAge: 31536000,
//		FrameOptions: "DENY",
//		CSPDirectives: map[string][]string{
//			"default-src": {"'self'"},
//			"script-src":  {"'self'", "trusted-scripts.com"},
//		},
//	}))
func NewSecurityMiddleware(config SecurityConfig) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			// Set security headers
			if config.HSTS {
				value := fmt.Sprintf("max-age=%d", config.HSTSMaxAge)
				if config.HSTSIncludeSubdomains {
					value += "; includeSubDomains"
				}
				c.Response().Header().Set("Strict-Transport-Security", value)
			}

			c.Response().Header().Set("X-Frame-Options", config.FrameOptions)
			c.Response().Header().Set("X-Content-Type-Options", config.ContentTypeOptions)
			c.Response().Header().Set("X-XSS-Protection", config.XSSProtection)

			if len(config.CSPDirectives) > 0 {
				c.Response().Header().Set("Content-Security-Policy", buildCSPHeader(config.CSPDirectives))
			}

			return next(c)
		}
	}
}
