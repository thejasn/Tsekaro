package stream

import (
	"encoding/json"
	"strings"

	"github.com/thejasn/tester/core/util/flatmap"
	"github.com/tidwall/gjson"
)

type Context interface {
	Mapper(int, string) string
	Store(int, interface{})
	Get(int) interface{}
}

type InMemoryContext struct {
	ctx map[int]interface{}
}

func NewInMemoryContext() *InMemoryContext {
	return &InMemoryContext{
		ctx: make(map[int]interface{}),
	}
}

func (c InMemoryContext) Mapper(actionID int, input string) string {
	if gjson.Valid(input) {
		if m, ok := gjson.Parse(input).Value().(map[string]interface{}); ok {
			newMap := map[string]interface{}{
				"CONST": m,
			}
			fm := flatmap.Flatten(newMap)
			src, err := json.Marshal(c.ctx)
			if err != nil {
				panic(err)
			}
			for k, v := range fm {
				path := strings.TrimPrefix(v, "$")
				newVal := gjson.Get(string(src), strings.Join([]string{string(actionID), path}, ".")).String()
				fm[k] = newVal
			}
			if result, ok := flatmap.Expand(fm, "CONST").(map[string]interface{}); ok {
				result, err := json.Marshal(result)
				if err != nil {
					panic(err)
				}
				return string(result)
			}
		}
	}
	return input
}

func (c *InMemoryContext) Store(k int, v interface{}) {
	c.ctx[k] = v
}

func (c InMemoryContext) Get(k int) interface{} {
	return c.ctx[k]
}
