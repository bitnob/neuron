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
