import (
	"context"
	"runtime"
	"sync"
)

// EngineConfig holds the configuration for the Neuron engine
type EngineConfig struct {
	// Core settings
	MaxProcs         int
	WorkerPoolSize   int
	QueueSize        int
	GracefulShutdown bool
	ShutdownTimeout  time.Duration

	// Performance settings
	EnableCompression bool
	CacheEnabled      bool
	CacheSize         int
	PoolSize          int

	// Security settings
	EnableCSRF          bool
	EnableSecureSession bool
	EnableRateLimit     bool
}

// Engine is the core of the Neuron framework
type Engine struct {
	config   *EngineConfig
	modules  *ModuleRegistry
	router   *Router
	pool     *WorkerPool
	cache    Cache
	metrics  *MetricsCollector
	shutdown chan struct{}
	wg       sync.WaitGroup
}

// ModuleRegistry manages framework modules
type ModuleRegistry struct {
	modules map[string]Module
	mu      sync.RWMutex
}

// Module interface for extensible components
type Module interface {
	Init(ctx context.Context) error
	Name() string
	Shutdown(ctx context.Context) error
}

// WorkerPool manages goroutines for request handling
type WorkerPool struct {
	workers chan *Worker
	queue   chan Job
	size    int
}

// Performance optimizations
func (e *Engine) optimize() {
	// Set GOMAXPROCS
	if e.config.MaxProcs > 0 {
		runtime.GOMAXPROCS(e.config.MaxProcs)
	}

	// Initialize worker pool
	e.pool = NewWorkerPool(e.config.WorkerPoolSize, e.config.QueueSize)

	// Enable compression
	if e.config.EnableCompression {
		e.Use(middleware.Compression())
	}

	// Initialize cache
	if e.config.CacheEnabled {
		e.cache = NewCache(e.config.CacheSize)
	}
}