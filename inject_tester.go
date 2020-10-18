//+build wireinject

package main

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/google/wire"
	flowrepo "github.com/thejasn/tester/domain/flow/repo"
	testcaserepo "github.com/thejasn/tester/domain/testcase/repo"
	"github.com/thejasn/tester/service"
	"github.com/thejasn/tester/transport/http"
	"github.com/thejasn/tester/transport/http/handler"
	"gorm.io/gorm"
)

func injectTester(ctx context.Context, r *chi.Mux, db *gorm.DB) http.Router {
	wire.Build(
		flowrepo.NewFlowRepo,
		testcaserepo.NewTestcaseRepo,
		service.NewFlowSvc,
		service.NewTestcaseSvc,
		wire.Struct(new(handler.Set), "*"),
		handler.NewFlowHandler,
		handler.NewTestcaseHandler,
		http.NewRouter,
	)
	return http.Router{}
}
