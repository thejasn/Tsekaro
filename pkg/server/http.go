package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/thejasn/tester/pkg/log"
)

type httpServerBuilder struct {
	addr    string
	handler http.Handler
}

type httpServer struct {
	srv *http.Server
}

type httpOption func(*httpServerBuilder)

// WithHTTPAddr provides an override option for the http server
func WithHTTPAddr(addr string, port uint) httpOption {
	return func(h *httpServerBuilder) {
		h.addr = fmt.Sprintf("%s:%d", addr, port)
	}
}

// WithHTTPHandler is used to override the handler for the http server
func WithHTTPHandler(handler http.Handler) httpOption {
	return func(h *httpServerBuilder) {
		h.handler = handler
	}
}

// BuildHttp initializes a http server with the provided set of httpOptions
func BuildHttp(opts ...httpOption) Server {

	sb := &httpServerBuilder{}

	for _, opt := range opts {
		opt(sb)
	}

	return &httpServer{
		srv: &http.Server{
			Addr:    sb.addr,
			Handler: sb.handler,
		},
	}
}

// Start starts the http server
func (h *httpServer) Start(ctx context.Context) {
	go func() {
		if err := h.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.GetLogger(ctx).Fatalf("listen:%+s\n", err)
		}
	}()
	log.GetLogger(ctx).Infof("started http server with addr: %v", h.srv.Addr)
}

// CleanUp is responsible for graceful shutdown of http server
func (h *httpServer) CleanUp(ctx context.Context) error {
	log.GetLogger(ctx).Infof("http server stopping")

	newCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := h.srv.Shutdown(newCtx); err != nil {
		return fmt.Errorf("http server shutdown was not graceful:%w", err)
	}
	log.GetLogger(newCtx).Infof("http server gracefully shutdown")
	return nil
}
