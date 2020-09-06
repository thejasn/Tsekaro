package tester

import (
	"github.com/thejasn/tester/core/client"
)

type Executor func() (string, string, error)

func GrpcExecutor(cc client.Runner, opts ...client.RunnerOpts) Executor {
	return func() (string, string, error) {
		for _, opt := range opts {
			opt(cc)
		}
		err := cc.Build()
		if err != nil {
			return "", "", err
		}
		jsonStr, err := cc.Invoke()
		if err != nil {
			return "", "", err
		}
		cc.Clear()
		return cc.GetIdentifier(), jsonStr, nil
	}
}

func RestExecutor(cc client.Runner, opts ...client.RunnerOpts) Executor {
	return func() (string, string, error) {
		for _, opt := range opts {
			opt(cc)
		}
		err := cc.Build()
		if err != nil {
			return "", "", err
		}
		jsonStr, err := cc.Invoke()
		if err != nil {
			return "", "", err
		}
		cc.Clear()
		return cc.GetIdentifier(), jsonStr, nil
	}
}
