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

	// todo: fasthttp is more trouble than it's worth, honestly
	// net/http is plenty fast and well optimized
	r := router.NewFastRouter()

	// Configure routes
	setupRoutes(r)

	// Create optimized server
	srv := server.NewFastServer(r.Handler(), r.Logger)

	// Start server in goroutine
	go func() {
		// todo: this fatal here will only kill the goroutine and not propagate
		// to the main routine which means the server will be alive but not be
		// able to receive any http requests. ideally we want to capture this
		// error here and use a channel to send it back to the main routine
		if err := srv.ListenAndServe(":8080"); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// now we can use a select {} here to listen for both the error channel and
	// context.Done(), that way which ever one happens first will trigger the
	// shutdown happy path
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
