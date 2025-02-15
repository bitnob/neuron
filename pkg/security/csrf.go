package security

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

type CSRFConfig struct {
	TokenLength    int
	CookieName     string
	HeaderName     string
	CookieMaxAge   int
	CookiePath     string
	CookieDomain   string
	CookieSecure   bool
	CookieHTTPOnly bool
	TrustedOrigins []string
}

type CSRFProtector struct {
	config CSRFConfig
	tokens map[string]time.Time
	mu     sync.RWMutex
}

func NewCSRFProtector(config CSRFConfig) *CSRFProtector {
	if config.TokenLength == 0 {
		config.TokenLength = 32
	}
	if config.CookieName == "" {
		config.CookieName = "_csrf"
	}
	if config.HeaderName == "" {
		config.HeaderName = "X-CSRF-Token"
	}

	return &CSRFProtector{
		config: config,
		tokens: make(map[string]time.Time),
	}
}

func (c *CSRFProtector) GenerateToken() (string, error) {
	b := make([]byte, c.config.TokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)

	c.mu.Lock()
	c.tokens[token] = time.Now().Add(time.Duration(c.config.CookieMaxAge) * time.Second)
	c.mu.Unlock()

	return token, nil
}
