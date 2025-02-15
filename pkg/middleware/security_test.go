package middleware

import (
	"net/http"
	"net/http/httptest"
	"neuron/pkg/router"
	"testing"
)

func TestSecurityMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		config     SecurityConfig
		wantHeader map[string]string
	}{
		{
			name: "HSTS enabled",
			config: SecurityConfig{
				HSTS:                  true,
				HSTSMaxAge:            31536000,
				HSTSIncludeSubdomains: true,
			},
			wantHeader: map[string]string{
				"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
			},
		},
		{
			name: "Content Security Policy",
			config: SecurityConfig{
				CSPDirectives: map[string][]string{
					"default-src": {"'self'"},
					"script-src":  {"'self'", "'unsafe-inline'"},
				},
			},
			wantHeader: map[string]string{
				"Content-Security-Policy": "default-src 'self'; script-src 'self' 'unsafe-inline'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := NewSecurityMiddleware(tt.config)
			handler := middleware(func(c *Context) error {
				return c.String(http.StatusOK, "OK")
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := router.NewContext(req, rec)

			err := handler(c)
			if err != nil {
				t.Errorf("Handler returned error: %v", err)
			}

			for k, v := range tt.wantHeader {
				if got := rec.Header().Get(k); got != v {
					t.Errorf("Header[%s] = %v, want %v", k, got, v)
				}
			}
		})
	}
}
