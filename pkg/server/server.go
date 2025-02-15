package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"runtime"
	"time"

	"neuron/pkg/logger"

	"golang.org/x/net/http2"
)

func NewServer(handler http.Handler, logger *logger.Logger) (*http.Server, net.Listener) {
	// Set GOMAXPROCS to match CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Configure TCP keep-alive listener
	lc := net.ListenConfig{
		KeepAlive: 30 * time.Second,
	}
	ln, err := lc.Listen(context.Background(), "tcp", ":8080")
	if err != nil {
		panic(err)
	}

	// HTTP/2 config
	h2Config := &http2.Server{
		MaxConcurrentStreams: 250,
		MaxReadFrameSize:     1048576,
		IdleTimeout:          10 * time.Second,
	}

	// TLS config for HTTP/2
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
	}

	// Wrap handler with logger middleware
	wrappedHandler := logger.Middleware(handler)

	// Create optimized server
	server := &http.Server{
		Addr:              ":8080",
		Handler:           wrappedHandler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		ReadHeaderTimeout: 2 * time.Second,
		TLSConfig:         tlsConfig,
	}

	// Enable HTTP/2
	http2.ConfigureServer(server, h2Config)

	return server, ln
}
