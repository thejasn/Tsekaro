package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/thejasn/tester/core/asserter"
	"github.com/thejasn/tester/core/client/grpc"
	"github.com/thejasn/tester/core/client/rest"
	"github.com/thejasn/tester/core/stream"
	"github.com/thejasn/tester/core/tester"
	"github.com/thejasn/tester/domain/flow/model"
	"github.com/thejasn/tester/domain/flow/repo"
	tmodel "github.com/thejasn/tester/domain/testcase/model"
	trepo "github.com/thejasn/tester/domain/testcase/repo"
)

type Flow interface {
	GetAll(ctx context.Context, page, pagesize int64, order string) ([]*model.Flow, int64, error)
	Get(context.Context, int) (model.Flow, error)
	Add(context.Context, *model.Flow) (*model.Flow, int64, error)
	Update(context.Context, int, *model.Flow) (*model.Flow, int64, error)
	Delete(context.Context, int) (int64, error)
	Execute(context.Context, int) (bool, error)
}

func NewFlowSvc(r repo.Flow, t trepo.Testcase) Flow {
	return flow{
		repo:  r,
		trepo: t,
	}
}

type flow struct {
	repo  repo.Flow
	trepo trepo.Testcase
}

func (f flow) GetAll(ctx context.Context, page, pagesize int64, order string) ([]*model.Flow, int64, error) {
	return f.repo.GetAll(ctx, page, pagesize, order)
}

func (f flow) Get(ctx context.Context, id int) (model.Flow, error) {
	return f.repo.Get(ctx, id)
}

func (f flow) Add(ctx context.Context, m *model.Flow) (*model.Flow, int64, error) {
	return f.repo.Add(ctx, m)
}

func (f flow) Update(ctx context.Context, id int, m *model.Flow) (*model.Flow, int64, error) {
	return f.repo.Update(ctx, id, m)
}

func (f flow) Delete(ctx context.Context, id int) (int64, error) {
	return f.repo.Delete(ctx, id)
}

func (f flow) Execute(ctx context.Context, id int) (bool, error) {
	fl, err := f.repo.Get(ctx, id)
	if err != nil {
		return false, fmt.Errorf("could not execute as flow %w", err)
	}

	tests, _, err := f.trepo.GetAllOrderedWhere(ctx, map[string]interface{}{
		"flow_id": fl.ID,
	})
	if err != nil {
		return false, fmt.Errorf("could not find flows for id: %d as %w", fl.ID, err)
	}

	l := stream.NewLinearFlow()
	for _, test := range tests {
		tc, err := f.trepo.Get(ctx, test.ID)
		if err != nil {
			return false, fmt.Errorf("could not execute as testcase %w", err)
		}

		var expected tmodel.Result
		err = json.Unmarshal(tc.Expected, &expected)
		if err != nil {
			return false, fmt.Errorf("corrupt data stored for 'expected' in testcase")
		}

		switch tc.API {
		case "REST":
			cfg := rest.NewRestConfig(ctx, tc.Scheme+"://"+tc.Host+":"+strconv.Itoa(tc.Port))
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
			l.Execute(tc.TestCaseID, tester.GrpcExecutor(ctx, cfg,
				grpc.WithRequest(tc.Body.String),
				grpc.WithMethod(tc.Path)), asserter.Assertion{
				Expected: expected.Data,
				Actual:   tc.Actual.String,
				Operator: tc.Operation})
		}
	}
	return true, nil
}
