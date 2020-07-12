package tester

type Executor interface {
	Execute() (string, map[string]interface{})
}
