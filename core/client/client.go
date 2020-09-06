package client

type Runner interface {
	GetIdentifier() string
	Build() error
	Invoke() (string, error)
	Clear()
}

type RunnerOpts func(Runner)
