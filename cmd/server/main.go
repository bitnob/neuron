package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"neuron/pkg/logger"
	"neuron/pkg/router"
	"neuron/pkg/server"
)

func init() {
	// Optimize GC
	debug.SetGCPercent(500)                      // Less frequent GC
	debug.SetMemoryLimit(4 * 1024 * 1024 * 1024) // 4GB memory limit
}

func main() {
	// Create context that will be canceled on interrupt
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create router
	r := router.New()

	// Ensure logger is initialized
	if r.Logger == nil {
		r.Logger = logger.New()
	}

	// Log server startup
	r.Logger.Info("Starting server...")

	// Configure routes
	setupRoutes(r)

	// Create optimized server
	srv, ln := server.NewServer(r, r.Logger)

	// Start server in goroutine
	go func() {
		r.Logger.Info("Server listening on :8080")
		if err := srv.Serve(ln); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	// Shutdown server gracefully
	log.Println("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}

// setupRoutes configures the API routes
func setupRoutes(r *router.Router) {
	r.GET("/", func(c *router.Context) error {
		return c.String(200, "Hello World!")
	})
}
