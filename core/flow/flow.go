package flow

import (
	"sync"

	"github.com/thejasn/tester/core/asserter"
	"github.com/thejasn/tester/core/tester"
)

type Engine interface {
	Execute(tester.Executor) Engine
	PostExecute(...asserter.Assertion) Engine
}

type Linear struct {
	ctx sync.Map
}

func (l *Linear) Execute(fn tester.Executor) Engine {
	key, content := fn.Execute()
	val, ok := l.ctx.Load(key)
	if ok {
		for k, v := range content {
			val.(map[string]interface{})[k] = v
		}
		l.ctx.Store(key, val)
	} else {
		l.ctx.Store(key, content)
	}
	return l
}

func (l *Linear) PostExecute(actions ...asserter.Assertion) Engine {
	for _, a := range actions {
		ok, msg := a.Assert()
		if !ok {
			panic(msg)
		}
	}
	return l
}
