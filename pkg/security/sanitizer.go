package security

import (
	"html"
)

type SanitizerConfig struct {
	AllowedTags      []string
	AllowedAttrs     map[string][]string
	AllowedProtocols []string
	MaxLength        int
}

type Sanitizer struct {
	config SanitizerConfig
}

func NewSanitizer(config SanitizerConfig) *Sanitizer {
	if config.AllowedTags == nil {
		config.AllowedTags = []string{"b", "i", "u", "p", "br", "a"}
	}
	if config.AllowedProtocols == nil {
		config.AllowedProtocols = []string{"http", "https", "mailto"}
	}
	return &Sanitizer{config: config}
}

func (s *Sanitizer) SanitizeString(input string) string {
	if s.config.MaxLength > 0 && len(input) > s.config.MaxLength {
		input = input[:s.config.MaxLength]
	}
	return html.EscapeString(input)
}
