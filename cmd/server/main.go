package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"neuron/pkg/router"
	"neuron/pkg/server"

	"github.com/valyala/fasthttp"
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
	r := router.NewFastRouter()

	// Configure routes
	setupRoutes(r)

	// Create optimized server
	srv := server.NewFastServer(r.Handler(), r.Logger)

	// Start server in goroutine
	go func() {
		if err := srv.ListenAndServe(":8080"); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	// Shutdown server gracefully
	log.Println("Shutting down server...")
	if err := srv.Shutdown(); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}

// setupRoutes configures the API routes
func setupRoutes(r *router.FastRouter) {
	r.GET("/", func(c *router.FastContext) error {
		c.RequestCtx.SetStatusCode(fasthttp.StatusOK)
		c.RequestCtx.SetBodyString("Hello World!")
		return nil
	})
}
