package main

import (
	"context"

	"github.com/thejasn/tester/core/asserter"
	grpcClient "github.com/thejasn/tester/core/client/grpc"
	restClient "github.com/thejasn/tester/core/client/rest"
	"github.com/thejasn/tester/core/flow"
	"github.com/thejasn/tester/core/tester"
	"github.com/thejasn/tester/pkg/log"
)

func main() {
	ctx := log.WithLogger(context.Background(), log.L)
	gc := grpcClient.NewConfig(ctx, "something", "localhost", "50051")
	rc := restClient.NewRestConfig(ctx, "http://localhost:8080")

	l := flow.NewLinearFlow()
	l.Execute("0", tester.GrpcExecutor(gc,
		grpcClient.WithRequest(`{"name": "thejas"}`),
		grpcClient.WithMethod("helloworld.Greeter/SayHello")), asserter.Assertion{
		Expected: "Hello thejas",
		Actual:   "message",
		Operator: asserter.Equal})

	l.Execute("1", tester.GrpcExecutor(gc,
		grpcClient.WithRequest(l.Ctx.Mapper("0", `{"name": "$message"}`)),
		grpcClient.WithMethod("helloworld.Greeter/SayHello")), asserter.Assertion{
		Expected: "Hello Hello thejas",
		Actual:   "message",
		Operator: asserter.Equal})

	l.Execute("2", tester.GrpcExecutor(gc,
		grpcClient.WithRequest(l.Ctx.Mapper("1", `{"name": "$message"}`)),
		grpcClient.WithMethod("helloworld.Greeter/SayHello")), asserter.Assertion{
		Expected: "Hello Hello Hello thejas",
		Actual:   "message",
		Operator: asserter.Equal})

	l.Execute("3", tester.RestExecutor(rc,
		restClient.WithMethod("GET"),
		restClient.WithHeaders(map[string]string{"test": "value"}),
		restClient.WithUriPath("/hello")), asserter.Assertion{
		Expected: "Hello World!",
		Actual:   "message",
		Operator: asserter.Equal,
	})

	l.Execute("4", tester.RestExecutor(rc,
		restClient.WithMethod("POST"),
		restClient.WithHeaders(map[string]string{"test": "value"}),
		restClient.WithUriPath("/hello/say"),
		restClient.WithBody(l.Ctx.Mapper("3", `{"message": "$message"}`))), asserter.Assertion{
		Expected: "Hello Hello World!!",
		Actual:   "message",
		Operator: asserter.Equal,
	})
}
