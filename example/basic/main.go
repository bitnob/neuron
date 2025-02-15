package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	neuron "neuron/pkg"
	"neuron/pkg/middleware"
	"neuron/pkg/router"
)

type loggerAdapter struct {
	*log.Logger
}

func (l *loggerAdapter) Debug(msg string, fields ...interface{}) {
	l.Printf("[DEBUG] "+msg, fields...)
}

func (l *loggerAdapter) Info(msg string, fields ...interface{}) {
	l.Printf("[INFO] "+msg, fields...)
}

func (l *loggerAdapter) Error(msg string, fields ...interface{}) {
	l.Printf("[ERROR] "+msg, fields...)
}

func (l *loggerAdapter) Warn(msg string, fields ...interface{}) {
	l.Printf("[WARN] "+msg, fields...)
}

func main() {
	// Wait for interrupt signal
	listenForShutdown := make(chan os.Signal, 1)

	// todo: we probably want to make this a context with signal.NotifyContext()
	// so we can pass it to neuron.New() or something. That way we can listen
	// for context.Done() as our cancellation request from a failing service
	// deep within our app and kill the application in a graceful manner
	signal.Notify(listenForShutdown, syscall.SIGINT, syscall.SIGTERM)

	// Create a new Neuron engine with default configuration
	config := neuron.DefaultConfig()
	config.Host = "0.0.0.0" // Listen on all interfaces
	config.Port = 8080

	app := neuron.New(config)

	// Get the underlying router
	r := app.Router()

	// Add middleware
	app.Use(middleware.NewLoggingMiddleware(middleware.LogConfig{
		Logger:        &loggerAdapter{Logger: log.Default()},
		SlowThreshold: 5 * time.Millisecond, // Only log as slow if request takes longer than 5ms
	}))
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())

	// Add security headers
	app.Use(middleware.NewSecurityMiddleware(middleware.SecurityConfig{
		HSTS:          true,
		HSTSMaxAge:    31536000,
		FrameOptions:  "DENY",
		XSSProtection: "1; mode=block",
	}))

	// Add routes
	r.GET("/", func(c *router.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Welcome to Neuron!",
		})
	})

	// Start the server
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

	// this becomes superfluous as we need to use select {} on the interrupt
	// signal and app.Start() error at the same time. So using the
	// signal.NotifyContext() feels more appropriate
	<-listenForShutdown

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
