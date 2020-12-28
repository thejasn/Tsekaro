package api

import (
	"database/sql"
)

type CreateTestCase struct {

	// Testcase Name
	FlowID int
	Name   string

	// Assertion Details
	Actual    string
	Expected  []byte
	Operation string

	// Testcase metadata
	Body          []byte
	Headers       []byte
	Host          string
	MappingTestID int
	Method        sql.NullString
	Path          string
	Port          int
	Scheme        string
}
