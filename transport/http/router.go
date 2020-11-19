package http

import (
	"context"
	"net/http"

	"github.com/thejasn/tester/pkg/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thejasn/tester/transport/http/handler"
)

// Router encapsulates the http router and the service layer object
type Router struct {
	mux     *chi.Mux
	handler handler.Set
}

// NewRouter is a provider/constructor for building a Router instance
func NewRouter(r *chi.Mux, h handler.Set) Router {
	log.InitInterceptor()
	return Router{
		mux:     r,
		handler: h,
	}
}

// Route returns the build router, after binding all handler from service layer
func (r Router) Route(ctx context.Context) *chi.Mux {
	r.mux.Use(
		middleware.RealIP,
		middleware.AllowContentType("application/json"),
		middleware.StripSlashes,
		middleware.Recoverer,
		log.RequestID,
		log.Logger(ctx, log.DefaultCodeToLevel),
	)

	r.mux.Route("/v1", func(m chi.Router) {
		m.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		m.Group(r.handler.Flow.ConfigFlowsRouter)
		m.Group(r.handler.Testcase.ConfigTestcasesRouter)
	})
	log.GetLogger(ctx).Info("Registering handlers")
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.GetLogger(ctx).Infof("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r.mux, walkFunc); err != nil {
		log.GetLogger(ctx).Errorf("Failed to walk routes %s\n", err.Error())
	}
	return r.mux
}
