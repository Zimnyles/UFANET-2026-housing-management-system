package grpcerrors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	infra_errors "notification-service/infra/errors"
)

var errCodeMap = map[error]codes.Code{
	infra_errors.ErrUserIDRequired:      codes.InvalidArgument,
	infra_errors.ErrDeviceTokenRequired: codes.InvalidArgument,
	infra_errors.ErrPlatformRequired:    codes.InvalidArgument,
}

func ToGRPCError(err error) error {
	for target, code := range errCodeMap {
		if errors.Is(err, target) {
			return status.Error(code, err.Error())
		}
	}

	return status.Error(codes.Internal, infra_errors.ErrInternal.Error())
}
