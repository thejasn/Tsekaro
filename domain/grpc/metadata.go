package grpc

import (
	"github.com/thejasn/tester/cerrors"
)

type Metadata struct {
	Method  string
	Service string
}

func ParseMetadata(headers map[string][]string) (Metadata, error) {
	var m Metadata
	if v, ok := headers["Method"]; ok {
		m.Method = v[0]
	} else {
		return Metadata{}, cerrors.ErrInValidation
	}
	if v, ok := headers["Service"]; ok {
		m.Service = v[0]
	} else {
		return Metadata{}, cerrors.ErrInValidation
	}
	return m, nil
}
