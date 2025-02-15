package main

import (
	"log"
	"runtime"

	// "neuron"
	"neuron/pkg/middleware"
)

func main() {
	// Create a new Neuron engine
	app := neuron.New(&neuron.EngineConfig{
		MaxProcs:       runtime.NumCPU(),
		WorkerPoolSize: 100,
		QueueSize:      1000,
	})

	// Add middleware
	app.Use(middleware.Logger())
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
	app.GET("/", func(c *neuron.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Welcome to Neuron!",
		})
	})

	// Start the server
	log.Fatal(app.Start())
}
