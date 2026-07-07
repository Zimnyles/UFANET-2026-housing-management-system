package grpcerrors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	infra_errors "profile-service/infra/errors"
)

var errCodeMap = map[error]codes.Code{
	infra_errors.ErrProfileNotFound: codes.NotFound,

	infra_errors.ErrUserIDRequired:       codes.InvalidArgument,
	infra_errors.ErrFullNameRequired:     codes.InvalidArgument,
	infra_errors.ErrFullNameTooShort:     codes.InvalidArgument,
	infra_errors.ErrFullNameTooLong:      codes.InvalidArgument,
	infra_errors.ErrPhoneRequired:        codes.InvalidArgument,
	infra_errors.ErrPhoneTooLong:         codes.InvalidArgument,
	infra_errors.ErrApartmentRequired:    codes.InvalidArgument,
	infra_errors.ErrApartmentTooLong:     codes.InvalidArgument,
	infra_errors.ErrHouseIDRequired:      codes.InvalidArgument,
	infra_errors.ErrHouseIDInvalid:       codes.InvalidArgument,
	infra_errors.ErrCompanyNameRequired:  codes.InvalidArgument,
	infra_errors.ErrCompanyNameTooLong:   codes.InvalidArgument,
	infra_errors.ErrHouseNameRequired:    codes.InvalidArgument,
	infra_errors.ErrHouseAddressRequired: codes.InvalidArgument,
	infra_errors.ErrUKIDRequired:         codes.InvalidArgument,
}

func ToGrpcError(err error) error {
	for target, code := range errCodeMap {
		if errors.Is(err, target) {
			return status.Error(code, err.Error())
		}
	}

	return status.Error(codes.Internal, infra_errors.ErrInternal.Error())
}
