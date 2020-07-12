package grpc

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thejasn/tester/cerrors"
)

func TestParseHeaders(t *testing.T) {

	cases := []struct {
		input  map[string][]string
		result Metadata
		err    error
	}{
		{
			input: map[string][]string{
				"Method":  {"com.example.com"},
				"Service": {"com.example.Service"},
			},
			result: Metadata{
				Method:  "com.example.com",
				Service: "com.example.Service",
			},
			err: nil,
		},
		{
			input: map[string][]string{
				"Service": {"com.example.Service"},
			},
			result: Metadata{},
			err:    cerrors.ErrInValidation,
		},
		{
			input: map[string][]string{
				"Method": {"com.example.com"},
			},
			result: Metadata{},
			err:    cerrors.ErrInValidation,
		},
		{
			input:  map[string][]string{},
			result: Metadata{},
			err:    cerrors.ErrInValidation,
		},
		{
			input: map[string][]string{
				"Method":       {"com.example.com"},
				"Service":      {"com.example.Service"},
				"otherheaders": {"com.example.com"},
			},
			result: Metadata{
				Method:  "com.example.com",
				Service: "com.example.Service",
			},
			err: nil,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r, err := ParseMetadata(c.input)
			if err != nil {
				if !errors.Is(err, c.err) {
					t.Errorf("expected %#v but got %#v", c.err, err)
				}
			}
			assert.Equal(t, c.result, r)
		})
	}
}
