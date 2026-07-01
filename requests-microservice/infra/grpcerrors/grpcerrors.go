package grpcerrors

import (
	"errors"

	infra_errors "requests-service/infra/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errCodeMap = map[error]codes.Code{
	infra_errors.ErrRequestNotFound: codes.NotFound,

	infra_errors.ErrRequestIDRequired:       codes.InvalidArgument,
	infra_errors.ErrUserIDRequired:          codes.InvalidArgument,
	infra_errors.ErrTitleRequired:           codes.InvalidArgument,
	infra_errors.ErrTitleTooShort:           codes.InvalidArgument,
	infra_errors.ErrTitleTooLong:            codes.InvalidArgument,
	infra_errors.ErrDescriptionRequired:     codes.InvalidArgument,
	infra_errors.ErrDescriptionTooShort:     codes.InvalidArgument,
	infra_errors.ErrTypeRequired:            codes.InvalidArgument,
	infra_errors.ErrTypeInvalid:             codes.InvalidArgument,
	infra_errors.ErrStatusRequired:          codes.InvalidArgument,
	infra_errors.ErrStatusInvalid:           codes.InvalidArgument,
	infra_errors.ErrCommentContentRequired:  codes.InvalidArgument,
	infra_errors.ErrCommentContentTooLong:   codes.InvalidArgument,
	infra_errors.ErrCommentRequestIDMissing: codes.InvalidArgument,
}

func ToGrpcError(err error) error {
	for target, code := range errCodeMap {
		if errors.Is(err, target) {
			return status.Error(code, err.Error())
		}
	}
	return status.Error(codes.Internal, infra_errors.ErrInternal.Error())
}
