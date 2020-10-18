package main

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/thejasn/tester/pkg/config"
	"github.com/thejasn/tester/pkg/db"
	"github.com/thejasn/tester/pkg/log"
	"github.com/thejasn/tester/pkg/server"
	"golang.org/x/sync/errgroup"
)

func main() {

	ctx := log.WithLogger(context.Background(), log.Init())

	router := injectTester(ctx, chi.NewMux(), db.LoadDatabase(ctx, config.LoadAppConfig()))
	httpServer := server.BuildHttp(
		server.WithHTTPAddr("0.0.0.0", 8080),
		server.WithHTTPHandler(router.Route(ctx)),
	)
	httpServer.Start(ctx)

	<-server.AwaitTermination()

	grace, ctx := errgroup.WithContext(ctx)

	grace.Go(func() error {
		return httpServer.CleanUp(ctx)
	})

	if err := grace.Wait(); err != nil {
		log.GetLogger(ctx).Errorf("application shutdown was not graceful:%+v", err)
	}
}
