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

// responseWriter wraps http.ResponseWriter to capture status and size
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int64
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}

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

	// Wrap handler with logging middleware
	loggingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log incoming request
		logger.Info("Incoming request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Create response wrapper to capture status
		rw := &responseWriter{ResponseWriter: w}

		// Process request
		handler.ServeHTTP(rw, r)

		// Log completion
		duration := time.Since(start)
		logger.Access(
			r.Method,
			r.URL.Path,
			rw.status,
			duration,
			rw.size,
			r.RemoteAddr,
		)
	})

	// Create optimized server
	server := &http.Server{
		Addr:              ":8080",
		Handler:           loggingHandler,
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
