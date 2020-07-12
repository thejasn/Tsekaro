package main

import (
	"context"
	"strings"
	"time"

	"github.com/thejasn/tester/core/asserter"
	"github.com/thejasn/tester/core/flow"

	"github.com/pkg/errors"
	"github.com/thejasn/tester/core/client"
	"github.com/thejasn/tester/pkg/log"

	"google.golang.org/grpc"
)

func main() {
	ctx := log.WithLogger(context.Background(), log.L)
	l := flow.Linear{}
	l.Execute().PostExecute(asserter.Assertion{}, asserter.Assertion{}).
		Execute().PostExecute(asserter.Assertion{}).
		Execute().PostExecute()

	dial := func() *grpc.ClientConn {
		clientBuilder := client.GrpcClientBuilder{}
		dialTime := 10 * time.Second
		ctx, cancel := context.WithTimeout(ctx, dialTime)
		defer cancel()
		clientBuilder.WithContext(ctx)
		cc, err := clientBuilder.GetConn("localhost", "6565")
		if err != nil {
			log.G(ctx).Fatal(errors.Wrapf(err, "Failed to dial target host %q", "localhost:6565"))
		}
		return cc
	}

	rc := client.ReflectClientBuilder{}
	rc.WithClientConn(dial())
	rc.WithAdditionalHeaders([]string{"partnerId:1"})
	rc.WithPayload(strings.NewReader(`{"partnerId": 1}`))
	rc.InvokeRPC("com.gonuclei.masterdata.v2.Masterdata/GetPartner")
}
