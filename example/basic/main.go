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

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
