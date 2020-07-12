package asserter

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type Operations interface {
	Assert() (bool, string)
}

type Assertion struct {
	Expected interface{}
	Actual   interface{}
	Operator string
}

var (
	Equal = "equal"
)

func (a Assertion) Assert() (bool, string) {
	switch a.Operator {
	case Equal:
		if cmp.Equal(a.Expected, a.Actual) {
			return true, ""
		} else {
			return false, fmt.Sprintf("expected %v but found %v", a.Expected, a.Actual)
		}
	default:
		return false, "invalid operator"
	}
}
