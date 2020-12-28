package cerrors

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (

	// ErrInParse is signifies an error during parsing of data
	// for exaple while converting proto to json or
	// while converting model to proto
	ErrInParse = errors.New("parsing error")

	// ErrInValidation signifies given input failed some kind of validation
	ErrInValidation = errors.New("validation failed")

	// ErrInvalid signifies invalid input/arg in request
	ErrInvalid = errors.New("invalid")

	// ErrNotFound error when record not found
	ErrNotFound = fmt.Errorf("record Not Found")

	// ErrUnableToMarshalJSON error when json payload corrupt
	ErrUnableToMarshalJSON = fmt.Errorf("json payload corrupt")

	// ErrUpdateFailed error when update fails
	ErrUpdateFailed = fmt.Errorf("db update error")

	// ErrInsertFailed error when insert fails
	ErrInsertFailed = fmt.Errorf("db insert error")

	// ErrDeleteFailed error when delete fails
	ErrDeleteFailed = fmt.Errorf("db delete error")

	// ErrBadParams error when bad params passed in
	ErrBadParams = fmt.Errorf("bad params error")
)

// GrpcHandler checks for specific types of errors and returns
// appropriate grpc errors with the required error codes
func GrpcHandler(err error) error {
	if errors.Is(err, ErrInValidation) {
		return status.Errorf(codes.InvalidArgument, err.Error())
	} else if errors.Is(err, ErrInParse) {
		return status.Errorf(codes.InvalidArgument, err.Error())
	} else if errors.Is(err, ErrNotFound) {
		return status.Errorf(codes.NotFound, err.Error())
	} else if errors.Is(err, ErrInvalid) {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	return status.Errorf(codes.Internal, err.Error())
}
