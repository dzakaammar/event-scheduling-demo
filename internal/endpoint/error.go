package endpoint

import (
	"github.com/dzakaammar/event-scheduling-example/internal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapErrToStatusCode(err error) error {
	switch err {
	case nil:
		return nil
	case internal.ErrInvalidRequest, internal.ErrInvalidTimezone, internal.ErrValidationFailed:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
