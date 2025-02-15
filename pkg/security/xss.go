package security

import (
	"net/http"
	"regexp"
	"strings"
)

type XSSConfig struct {
	EnableCSP           bool
	CSPDirectives       map[string][]string
	EnableXFrameOptions bool
	XFrameOptions       string
	EnableXSSProtection bool
	XSSProtection       string
}

type XSSProtector struct {
	config   XSSConfig
	patterns []*regexp.Regexp
}

func NewXSSProtector(config XSSConfig) *XSSProtector {
	if config.CSPDirectives == nil {
		config.CSPDirectives = map[string][]string{
			"default-src": {"'self'"},
			"script-src":  {"'self'"},
			"style-src":   {"'self'"},
			"img-src":     {"'self'"},
		}
	}

	return &XSSProtector{
		config:   config,
		patterns: compileXSSPatterns(),
	}
}

func (x *XSSProtector) ApplyHeaders(w http.ResponseWriter) {
	if x.config.EnableCSP {
		w.Header().Set("Content-Security-Policy", x.buildCSPHeader())
	}

	if x.config.EnableXFrameOptions {
		w.Header().Set("X-Frame-Options", x.config.XFrameOptions)
	}

	if x.config.EnableXSSProtection {
		w.Header().Set("X-XSS-Protection", x.config.XSSProtection)
	}
}

func compileXSSPatterns() []*regexp.Regexp {
	patterns := []string{
		`<script[^>]*>.*?</script>`,
		`javascript:.*`,
		`onload=.*`,
		`onerror=.*`,
	}

	compiled := make([]*regexp.Regexp, len(patterns))
	for i, p := range patterns {
		compiled[i] = regexp.MustCompile(p)
	}
	return compiled
}

func (x *XSSProtector) buildCSPHeader() string {
	var policies []string
	for directive, sources := range x.config.CSPDirectives {
		if len(sources) > 0 {
			policies = append(policies, directive+" "+strings.Join(sources, " "))
		}
	}
	return strings.Join(policies, "; ")
}
