package router

import (
	"sync"
	"time"

	"neuron/pkg/logger"

	"github.com/valyala/fasthttp"
)

type FastRouter struct {
	routes      map[string][]FastRoute
	contextPool sync.Pool
	Logger      *logger.Logger
	routeTrie   *fastNode
}

type FastRoute struct {
	Method  string
	Path    string
	Handler FastHandlerFunc
}

type FastHandlerFunc func(*FastContext) error

type FastContext struct {
	RequestCtx *fasthttp.RequestCtx
	Params     []Param
	store      map[string]interface{}
}

func NewFastRouter() *FastRouter {
	r := &FastRouter{
		routes:    make(map[string][]FastRoute, 32),
		routeTrie: &fastNode{children: make([]*fastNode, 0, 8)},
		Logger:    logger.New(),
	}

	r.contextPool = sync.Pool{
		New: func() interface{} {
			return &FastContext{
				store:  make(map[string]interface{}, 8),
				Params: make([]Param, 0, 8),
			}
		},
	}

	return r
}

func (r *FastRouter) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		// Log incoming request
		r.Logger.Info("[%s] Incoming request: %s %s",
			ctx.RemoteIP().String(),
			ctx.Method(),
			ctx.Path())

		// Get context from pool
		c := r.contextPool.Get().(*FastContext)
		c.Reset(ctx)
		defer r.contextPool.Put(c)

		// Find handler
		method := string(ctx.Method())
		path := string(ctx.Path())

		if handler, params := r.routeTrie.find(method, path); handler != nil {
			c.Params = params
			if err := handler(c); err != nil {
				r.Logger.Error("Handler error: %v", err)
				ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			}
		} else {
			r.Logger.Error("No route found for %s %s", method, path)
			ctx.NotFound()
		}

		// Log request completion
		duration := time.Since(start)
		r.Logger.Access(
			method,
			path,
			ctx.Response.StatusCode(),
			duration,
			int64(ctx.Response.Header.ContentLength()),
			ctx.RemoteIP().String(),
		)
	}
}

func (r *FastRouter) GET(path string, handler FastHandlerFunc) {
	r.Handle("GET", path, handler)
}

func (r *FastRouter) POST(path string, handler FastHandlerFunc) {
	r.Handle("POST", path, handler)
}

func (r *FastRouter) Handle(method, path string, handler FastHandlerFunc) {
	r.routeTrie.insert(method, path, handler)
	r.routes[method] = append(r.routes[method], FastRoute{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

func (c *FastContext) Reset(ctx *fasthttp.RequestCtx) {
	c.RequestCtx = ctx
	c.store = make(map[string]interface{}, 8)
	if cap(c.Params) > 32 {
		c.Params = c.Params[:0]
	} else {
		c.Params = make([]Param, 0, 8)
	}
}
