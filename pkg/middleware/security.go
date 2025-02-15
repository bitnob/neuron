package middleware

import "fmt"

type SecurityConfig struct {
	HSTS                  bool
	HSTSMaxAge            int
	HSTSIncludeSubdomains bool
	FrameOptions          string
	ContentTypeOptions    string
	XSSProtection         string
	CSPDirectives         map[string][]string
}

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
