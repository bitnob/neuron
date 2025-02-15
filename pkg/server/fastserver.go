package server

import (
	"runtime"
	"time"

	"neuron/pkg/logger"

	"github.com/valyala/fasthttp"
)

type FastServer struct {
	server *fasthttp.Server
	logger *logger.Logger
}

func NewFastServer(handler fasthttp.RequestHandler, logger *logger.Logger) *FastServer {
	// Set GOMAXPROCS
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create server with optimized settings
	server := &fasthttp.Server{
		Handler:                       handler,
		Name:                          "Neuron",
		ReadTimeout:                   30 * time.Second,
		WriteTimeout:                  30 * time.Second,
		IdleTimeout:                   120 * time.Second,
		MaxRequestBodySize:            1024 * 1024 * 10, // 10MB
		DisableHeaderNamesNormalizing: true,
		NoDefaultServerHeader:         true,
		NoDefaultContentType:          true,
		NoDefaultDate:                 true,
		ReduceMemoryUsage:             true,
		Concurrency:                   runtime.NumCPU() * 10000,
		MaxConnsPerIP:                 50000,
		TCPKeepalive:                  true,
		TCPKeepalivePeriod:            30 * time.Second,
		GetOnly:                       true, // Optimize for GET requests
		DisableKeepalive:              false,
		StreamRequestBody:             true,
		DisablePreParseMultipartForm:  true,
	}

	return &FastServer{
		server: server,
		logger: logger,
	}
}

func (s *FastServer) ListenAndServe(addr string) error {
	s.logger.Info("Fast server starting on %s", addr)
	return s.server.ListenAndServe(addr)
}

func (s *FastServer) Shutdown() error {
	return s.server.Shutdown()
}
