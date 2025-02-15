package router

import (
	"fmt"
	"net/http"
	"sync"
)

// HandlerFunc defines a function to serve HTTP requests
type HandlerFunc func(*Context) error

// MiddlewareFunc defines middleware function
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Route represents a route with its handler and configuration
type Route struct {
	Method     string
	Path       string
	Handler    HandlerFunc
	Middleware []MiddlewareFunc
	match      func(string) ([]Param, bool)
}

// Router handles HTTP routing
type Router struct {
	routes      map[string][]Route
	middleware  []MiddlewareFunc
	groups      []*RouteGroup
	notFound    http.HandlerFunc
	pool        sync.Pool
	contextPool sync.Pool
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Debug logging
	fmt.Printf("Handling request: %s %s\n", req.Method, req.URL.Path)
	fmt.Printf("Available routes: %+v\n", r.routes)

	// Get a context from the pool
	c := r.pool.Get().(*Context)
	c.Reset(w, req)

	// Find and execute the matching route
	found := false
	for _, route := range r.routes[req.Method] {
		fmt.Printf("Checking route: %s %s\n", route.Method, route.Path)
		if params, ok := route.match(req.URL.Path); ok {
			c.Params = params
			found = true

			// Build middleware chain
			handler := route.Handler
			for i := len(r.middleware) - 1; i >= 0; i-- {
				handler = r.middleware[i](handler)
			}

			// Execute the handler
			if err := handler(c); err != nil {
				// TODO: Implement error handling
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			break
		}
	}

	// Handle 404 Not Found
	if !found {
		if r.notFound != nil {
			r.notFound(w, req)
		} else {
			http.NotFound(w, req)
		}
	}

	// Put the context back in the pool
	r.pool.Put(c)
}

// Handle registers a new route
func (r *Router) Handle(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	if r == nil {
		r = New()
	}
	if r.routes == nil {
		r.routes = make(map[string][]Route)
	}

	route := Route{
		Method:     method,
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	}

	// Add route matching function
	route.match = func(reqPath string) ([]Param, bool) {
		if reqPath == path {
			return nil, true
		}
		return nil, false
	}

	r.routes[method] = append(r.routes[method], route)
}

// Use adds middleware to the router
func (r *Router) Use(middleware ...MiddlewareFunc) {
	if r == nil {
		r = New()
	}
	if r.middleware == nil {
		r.middleware = make([]MiddlewareFunc, 0)
	}
	if len(middleware) > 0 {
		r.middleware = append(r.middleware, middleware...)
	}
}

// New creates a new router instance
func New() *Router {
	r := &Router{
		routes:     make(map[string][]Route),
		middleware: make([]MiddlewareFunc, 0),
		groups:     make([]*RouteGroup, 0),
	}

	r.pool = sync.Pool{
		New: func() interface{} {
			return &Context{
				store: make(map[string]interface{}),
			}
		},
	}

	r.contextPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]interface{})
		},
	}

	return r
}

// Context represents the request context
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Params   []Param
	store    map[string]interface{}
}

// Reset resets the context for reuse
func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
	c.Request = r
	c.Response = w
	c.Params = c.Params[:0]
	for k := range c.store {
		delete(c.store, k)
	}
}

// GET registers a GET route
func (r *Router) GET(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.Handle("GET", path, handler, middleware...)
}

// POST registers a POST route
func (r *Router) POST(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.Handle("POST", path, handler, middleware...)
}

// PUT registers a PUT route
func (r *Router) PUT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.Handle("PUT", path, handler, middleware...)
}

// DELETE registers a DELETE route
func (r *Router) DELETE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.Handle("DELETE", path, handler, middleware...)
}
