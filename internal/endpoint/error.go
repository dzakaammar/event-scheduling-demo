package endpoint

import (
	"errors"

	"github.com/dzakaammar/event-scheduling-example/internal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapErrToStatusCode(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, internal.ErrInvalidRequest) ||
		errors.Is(err, internal.ErrInvalidTimezone) ||
		errors.Is(err, internal.ErrValidationFailed) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
