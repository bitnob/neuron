package benchmark

import (
	"net/http/httptest"
	neuron "neuron/pkg"
	"neuron/pkg/router"
	"testing"
)

func BenchmarkRouting(b *testing.B) {
	app := neuron.New(neuron.DefaultConfig())
	r := app.Router()

	// Add test routes
	r.GET("/", func(c *router.Context) error {
		return c.String(200, "Hello World")
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
		}
	})
}

func BenchmarkMiddleware(b *testing.B) {
	// Similar benchmark for middleware chain
}

func BenchmarkJSONSerialization(b *testing.B) {
	// Benchmark JSON handling
}
