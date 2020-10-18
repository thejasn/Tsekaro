package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/thejasn/tester/core/client/grpc"

	"github.com/thejasn/tester/core/asserter"
	"github.com/thejasn/tester/core/client/rest"
	"github.com/thejasn/tester/core/stream"
	"github.com/thejasn/tester/core/tester"
	"github.com/thejasn/tester/domain/testcase/model"
	"github.com/thejasn/tester/domain/testcase/repo"
)

type Testcase interface {
	GetAll(context.Context, int64, int64, string) ([]*model.Testcase, int64, error)
	Get(context.Context, int) (model.Testcase, error)
	Add(context.Context, *model.Testcase) (*model.Testcase, int64, error)
	Update(context.Context, int, *model.Testcase) (*model.Testcase, int64, error)
	Delete(context.Context, int) (int64, error)
	Execute(context.Context, int) (bool, error)
}

func NewTestcaseSvc(r repo.Testcase) Testcase {
	return testcase{
		r: r,
	}
}

type testcase struct {
	r repo.Testcase
}

func (t testcase) GetAll(ctx context.Context, page int64, pagesize int64, order string) ([]*model.Testcase, int64, error) {
	return t.r.GetAll(ctx, page, pagesize, order)
}

func (t testcase) Get(ctx context.Context, id int) (model.Testcase, error) {
	return t.r.Get(ctx, id)
}

func (t testcase) Add(ctx context.Context, ts *model.Testcase) (*model.Testcase, int64, error) {
	return t.r.Add(ctx, ts)
}

func (t testcase) Update(ctx context.Context, id int, tc *model.Testcase) (*model.Testcase, int64, error) {
	return t.r.Update(ctx, id, tc)
}

func (t testcase) Delete(ctx context.Context, id int) (int64, error) {
	return t.r.Delete(ctx, id)
}

func (t testcase) Execute(ctx context.Context, id int) (bool, error) {
	tc, err := t.r.Get(ctx, id)
	if err != nil {
		return false, fmt.Errorf("could not execute as testcase %w", err)
	}

	var expected model.Result
	err = json.Unmarshal(tc.Expected, &expected)
	if err != nil {
		return false, fmt.Errorf("corrupt data stored for 'expected' in testcase")
	}

	switch tc.API {
	case "REST":
		cfg := rest.NewRestConfig(ctx, tc.Scheme+"://"+tc.Host+":"+strconv.Itoa(tc.Port))
		l := stream.NewLinearFlow()
		l.Execute(tc.TestCaseID, tester.RestExecutor(ctx, cfg,
			rest.WithMethod(tc.Method.String),
			rest.WithBody(tc.Body.String),
			rest.WithUriPath(tc.Path)), asserter.Assertion{
			Expected: expected.Data,
			Actual:   tc.Actual.String,
			Operator: tc.Operation,
		})
	case "GRPC":
		cfg := grpc.NewConfig(ctx, "something", tc.Host, strconv.Itoa(tc.Port))
		l := stream.NewLinearFlow()
		l.Execute(tc.TestCaseID, tester.GrpcExecutor(ctx, cfg,
			grpc.WithRequest(tc.Body.String),
			grpc.WithMethod(tc.Path)), asserter.Assertion{
			Expected: expected.Data,
			Actual:   tc.Actual.String,
			Operator: tc.Operation})
	}
	return true, nil
}
