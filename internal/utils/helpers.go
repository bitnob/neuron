package utils

import (
	"math/rand"
	"net"
	"strings"
)

// StringHelpers provides string manipulation utilities
type StringHelpers struct{}

// Slugify creates a URL-friendly slug from a string
func (h *StringHelpers) Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, s)
}

// RandomString generates a random string of given length
func (h *StringHelpers) RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// NetworkHelpers provides network-related utilities
type NetworkHelpers struct{}

// GetLocalIP returns the non-loopback local IP of the host
func (h *NetworkHelpers) GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
