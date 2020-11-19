package client

import "context"

type Runner interface {
	GetIdentifier() string
	Build(context.Context) error
	Invoke() (string, error)
	Clear()
}

type RunnerOpts func(Runner)
