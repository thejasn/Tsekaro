package flow

import (
	"encoding/json"
	"strconv"

	"github.com/thejasn/tester/core/asserter"
	"github.com/thejasn/tester/core/tester"
	"github.com/tidwall/gjson"
)

type Engine interface {
	Execute(tester.Executor) Engine
}

type Linear struct {
	currentKey int
	Ctx        Context
}

func NewLinearFlow() Linear {
	return Linear{
		currentKey: 0,
		Ctx:        NewInMemoryContext(),
	}
}

func (l *Linear) Execute(actionId string, fn tester.Executor, actions ...asserter.Assertion) *Linear {
	_, content, _ := fn()
	var dest map[string]interface{}
	err := json.Unmarshal([]byte(content), &dest)
	if err != nil {
		panic(err)
	}
	l.Ctx.Store(strconv.Itoa(l.currentKey), dest)
	for _, a := range actions {
		src, err := json.Marshal(l.Ctx.Get(actionId))
		if err != nil {
			panic(err)
		}
		a.Actual = gjson.Get(string(src), a.Actual.(string)).Value()
		ok, msg := a.Assert()
		if !ok {
			panic(msg)
		}
	}
	l.currentKey++
	return l
}
