package grpcerrors

import (
	infra_errors "auth-service/infra/errors"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errCodeMap = map[error]codes.Code{
	infra_errors.ErrUserNotFound:       codes.NotFound,
	infra_errors.ErrEmailAlreadyExists: codes.AlreadyExists,
	infra_errors.ErrInvalidCredentials: codes.Unauthenticated,
	infra_errors.ErrInvalidToken:       codes.Unauthenticated,
	infra_errors.ErrTokenExpired:       codes.Unauthenticated,
	infra_errors.ErrTokenNotFound:      codes.Unauthenticated,

	infra_errors.ErrNameRequired:     codes.InvalidArgument,
	infra_errors.ErrNameTooShort:     codes.InvalidArgument,
	infra_errors.ErrNameTooLong:      codes.InvalidArgument,
	infra_errors.ErrEmailRequired:    codes.InvalidArgument,
	infra_errors.ErrEmailInvalid:     codes.InvalidArgument,
	infra_errors.ErrPasswordRequired: codes.InvalidArgument,
	infra_errors.ErrPasswordTooShort: codes.InvalidArgument,
	infra_errors.ErrPasswordTooLong:  codes.InvalidArgument,
	infra_errors.ErrRefreshRequired:  codes.InvalidArgument,
}

func ToGrpcError(err error) error {
	for target, code := range errCodeMap {
		if errors.Is(err, target) {
			return status.Error(code, err.Error())
		}
	}
	return status.Error(codes.Internal, infra_errors.ErrInternal.Error())
}
