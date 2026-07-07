package grpcerrors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	infraerrors "profile-service/infra/errors"
	"testing"
)

func TestToGrpcError(t *testing.T) {
	assert.Equal(t, codes.NotFound, status.Code(ToGrpcError(infraerrors.ErrProfileNotFound)))
	for _, e := range []error{infraerrors.ErrUserIDRequired, infraerrors.ErrFullNameRequired, infraerrors.ErrFullNameTooShort, infraerrors.ErrFullNameTooLong, infraerrors.ErrPhoneRequired, infraerrors.ErrPhoneTooLong, infraerrors.ErrApartmentRequired, infraerrors.ErrApartmentTooLong, infraerrors.ErrHouseIDRequired, infraerrors.ErrHouseIDInvalid, infraerrors.ErrCompanyNameRequired, infraerrors.ErrCompanyNameTooLong, infraerrors.ErrHouseNameRequired, infraerrors.ErrHouseAddressRequired, infraerrors.ErrUKIDRequired} {
		assert.Equal(t, codes.InvalidArgument, status.Code(ToGrpcError(e)))
	}
	assert.Equal(t, codes.Internal, status.Code(ToGrpcError(errors.New("x"))))
}
