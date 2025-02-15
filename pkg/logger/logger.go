package logger

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	accessLog *log.Logger
	errorLog  *log.Logger
	pool      sync.Pool
}

type logEntry struct {
	Method    string
	Path      string
	Status    int
	Latency   time.Duration
	Size      int64
	UserAgent string
	IP        string
	Time      time.Time
}

func New() *Logger {
	return &Logger{
		accessLog: log.New(os.Stdout, "[ACCESS] ", log.LstdFlags),
		errorLog:  log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		pool: sync.Pool{
			New: func() interface{} {
				return &logEntry{}
			},
		},
	}
}

func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get log entry from pool
		entry := l.pool.Get().(*logEntry)
		defer l.pool.Put(entry)

		// Reset entry
		*entry = logEntry{
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IP:        getClientIP(r),
			Time:      start,
		}

		// Create response wrapper
		rw := &responseWriter{ResponseWriter: w}

		// Process request
		next.ServeHTTP(rw, r)

		// Update entry with response info
		entry.Status = rw.status
		entry.Size = rw.size
		entry.Latency = time.Since(start)

		// Log entry with color based on status
		statusColor := getStatusColor(entry.Status)
		l.accessLog.Printf("%s %s %s%d%s %v %d bytes [%s] %s",
			entry.Method,
			entry.Path,
			statusColor,
			entry.Status,
			"\033[0m",
			entry.Latency,
			entry.Size,
			entry.IP,
			entry.UserAgent,
		)
	})
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLog.Printf(format, v...)
}

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
	// Set default status if not set
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}

// getClientIP gets the real client IP from headers or RemoteAddr
func getClientIP(r *http.Request) string {
	// Check X-Real-IP header
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Check X-Forwarded-For header
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0]
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// getStatusColor returns ANSI color code based on status code
func getStatusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "\033[32m" // Green
	case status >= 300 && status < 400:
		return "\033[34m" // Blue
	case status >= 400 && status < 500:
		return "\033[33m" // Yellow
	default:
		return "\033[31m" // Red
	}
}

// Add new method
func (l *Logger) Info(format string, v ...interface{}) {
	l.accessLog.Printf(format, v...)
}

func (l *Logger) Access(method, path string, status int, latency time.Duration, size int64, ip string) {
	statusColor := getStatusColor(status)
	l.accessLog.Printf("%s %s %s%d%s %v %d bytes [%s]",
		method,
		path,
		statusColor,
		status,
		"\033[0m",
		latency,
		size,
		ip,
	)
}

// Add new method for fasthttp logging
func (l *Logger) AccessFastHTTP(method, path []byte, status int, latency time.Duration, size int64, ip string) {
	statusColor := getStatusColor(status)
	l.accessLog.Printf("%s %s %s%d%s %v %d bytes [%s]",
		string(method),
		string(path),
		statusColor,
		status,
		"\033[0m",
		latency,
		size,
		ip,
	)
}
