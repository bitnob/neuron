package router

import (
	"net/http"
	"neuron/pkg/logger"
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
	contextPool sync.Pool
	paramPool   sync.Pool
	routeTrie   *node
	Logger      *logger.Logger
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Get context from pool - zero allocation path
	ctx := r.contextPool.Get().(*Context)
	ctx.Reset(w, req)
	defer r.contextPool.Put(ctx)

	// Fast path for static routes - no allocations
	if handler, params := r.routeTrie.find(req.Method, req.URL.Path); handler != nil {
		ctx.Params = params
		if err := handler.(HandlerFunc)(ctx); err != nil {
			r.Logger.Error("Handler error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	r.Logger.Error("No route found for %s %s", req.Method, req.URL.Path)
	http.NotFound(w, req)
}

// Handle registers a new route
func (r *Router) Handle(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	if r.routes == nil {
		r.routes = make(map[string][]Route, 32)
	}

	// Pre-compile route into trie for O(1) lookup
	r.routeTrie.insert(method, path, handler)

	// Store route info for debugging/introspection
	r.routes[method] = append(r.routes[method], Route{
		Method:     method,
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	})
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
		routes:     make(map[string][]Route, 32),
		middleware: make([]MiddlewareFunc, 0, 8),
		groups:     make([]*RouteGroup, 0, 8),
		routeTrie:  &node{children: make([]*node, 0, 8)},
		Logger:     logger.New(),
	}

	r.contextPool = sync.Pool{
		New: func() interface{} {
			return &Context{
				store:  make(map[string]interface{}, 8),
				Params: make([]Param, 0, 8),
			}
		},
	}

	r.paramPool = sync.Pool{
		New: func() interface{} {
			return make([]Param, 0, 8)
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
	// Just create new map - faster than clearing
	c.store = make(map[string]interface{}, 8)
	// Reuse param slice
	if cap(c.Params) > 32 {
		c.Params = c.Params[:0]
	} else {
		c.Params = make([]Param, 0, 8)
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
