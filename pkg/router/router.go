package router

import (
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
	// Get a context from the pool
	c := r.pool.Get().(*Context)
	c.Reset(w, req)

	// Find and execute the matching route
	found := false
	for _, route := range r.routes[req.Method] {
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
	route := Route{
		Method:     method,
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	}

	if r.routes == nil {
		r.routes = make(map[string][]Route)
	}
	r.routes[method] = append(r.routes[method], route)
}

// Use adds middleware to the router
func (r *Router) Use(middleware ...MiddlewareFunc) {
	r.middleware = append(r.middleware, middleware...)
}

// New creates a new router instance
func New() *Router {
	r := &Router{
		routes: make(map[string][]Route),
	}

	r.pool.New = func() interface{} {
		return &Context{}
	}

	r.contextPool.New = func() interface{} {
		return make(map[string]interface{})
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
