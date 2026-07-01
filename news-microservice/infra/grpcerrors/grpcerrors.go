package grpcerrors

import (
	"errors"

	infra_errors "news-service/infra/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errCodeMap = map[error]codes.Code{
	infra_errors.ErrNewsNotFound: codes.NotFound,

	infra_errors.ErrTitleRequired:   codes.InvalidArgument,
	infra_errors.ErrTitleTooShort:   codes.InvalidArgument,
	infra_errors.ErrTitleTooLong:    codes.InvalidArgument,
	infra_errors.ErrContentRequired: codes.InvalidArgument,
	infra_errors.ErrContentTooShort: codes.InvalidArgument,
	infra_errors.ErrHouseIDRequired: codes.InvalidArgument,
	infra_errors.ErrNewsIDRequired:  codes.InvalidArgument,
}

func ToGrpcError(err error) error {
	for target, code := range errCodeMap {
		if errors.Is(err, target) {
			return status.Error(code, err.Error())
		}
	}
	return status.Error(codes.Internal, infra_errors.ErrInternal.Error())
}
