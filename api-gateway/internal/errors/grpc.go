package app_errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FromGRPC(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		if errors.Is(err, ErrClientDial) {
			return ErrClientDial
		}

		return ErrInternal
	}

	switch st.Code() {
	case codes.Unauthenticated:
		return ErrUnauthorized
	case codes.PermissionDenied:
		return ErrForbidden
	case codes.NotFound:
		return ErrNotFound
	case codes.AlreadyExists, codes.InvalidArgument:
		return ErrBadRequest
	case codes.Unavailable:
		return ErrClientDial
	case codes.DeadlineExceeded, codes.Canceled:
		return ErrRequestTimeout
	default:
		return ErrInternal
	}
}
