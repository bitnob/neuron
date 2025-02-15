// pkg/router/router.go
package router

import (
	"context"
	"net/http"
	"sync"
)

type Router struct {
	routes      map[string][]Route
	middleware  []MiddlewareFunc
	groups      []*RouteGroup
	notFound    http.HandlerFunc
	pool        sync.Pool
	contextPool sync.Pool
}

type Route struct {
	Method      string
	Path        string
	Handler     HandlerFunc
	Middleware  []MiddlewareFunc
	Params      []Param
}

type HandlerFunc func(*Context) error
type MiddlewareFunc func(HandlerFunc) HandlerFunc

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

func (r *Router) Handle(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	route := Route{
		Method:     method,
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	}

	r.routes[method] = append(r.routes[method], route)
}

// pkg/router/group.go
package router

type RouteGroup struct {
	prefix     string
	router     *Router
	middleware []MiddlewareFunc
}

func (r *Router) Group(prefix string, middleware ...MiddlewareFunc) *RouteGroup {
	group := &RouteGroup{
		prefix:     prefix,
		router:     r,
		middleware: middleware,
	}
	r.groups = append(r.groups, group)
	return group
}

func (g *RouteGroup) Handle(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	// Combine group middleware with route middleware
	finalMiddleware := make([]MiddlewareFunc, 0, len(g.middleware)+len(middleware))
	finalMiddleware = append(finalMiddleware, g.middleware...)
	finalMiddleware = append(finalMiddleware, middleware...)

	fullPath := g.prefix + path
	g.router.Handle(method, fullPath, handler, finalMiddleware...)
}