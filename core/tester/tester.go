package tester

import (
	"context"

	"github.com/thejasn/tester/pkg/log"

	"github.com/thejasn/tester/core/client"
)

type Executor func() (string, string, error)

func GrpcExecutor(ctx context.Context, cc client.Runner, opts ...client.RunnerOpts) Executor {
	return func() (string, string, error) {
		for _, opt := range opts {
			opt(cc)
		}
		err := cc.Build(ctx)
		if err != nil {
			return "", "", err
		}
		jsonStr, err := cc.Invoke()
		if err != nil {
			return "", "", err
		}
		log.GetLogger(ctx).Debugf("Response: %+v", jsonStr)
		cc.Clear()
		return cc.GetIdentifier(), jsonStr, nil
	}
}

func RestExecutor(ctx context.Context, cc client.Runner, opts ...client.RunnerOpts) Executor {
	return func() (string, string, error) {
		for _, opt := range opts {
			opt(cc)
		}
		err := cc.Build(ctx)
		if err != nil {
			return "", "", err
		}
		jsonStr, err := cc.Invoke()
		if err != nil {
			return "", "", err
		}
		log.GetLogger(ctx).Debugf("Response: %+v", jsonStr)
		cc.Clear()
		return cc.GetIdentifier(), jsonStr, nil
	}
}
