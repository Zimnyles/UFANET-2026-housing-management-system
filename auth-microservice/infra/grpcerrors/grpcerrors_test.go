package grpcerrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	infraerrors "auth-service/infra/errors"
)

func TestToGrpcError(t *testing.T) {
	for _, tc := range []struct {
		err  error
		code codes.Code
	}{
		{infraerrors.ErrUserNotFound, codes.NotFound}, {infraerrors.ErrEmailAlreadyExists, codes.AlreadyExists},
		{infraerrors.ErrInvalidCredentials, codes.Unauthenticated}, {infraerrors.ErrInvalidToken, codes.Unauthenticated},
		{infraerrors.ErrTokenExpired, codes.Unauthenticated}, {infraerrors.ErrTokenNotFound, codes.Unauthenticated},
		{infraerrors.ErrEmailRequired, codes.InvalidArgument}, {infraerrors.ErrEmailInvalid, codes.InvalidArgument},
		{infraerrors.ErrPasswordRequired, codes.InvalidArgument}, {infraerrors.ErrPasswordTooShort, codes.InvalidArgument},
		{infraerrors.ErrPasswordTooLong, codes.InvalidArgument}, {infraerrors.ErrRefreshRequired, codes.InvalidArgument},
		{infraerrors.ErrInvalidAdminCode, codes.InvalidArgument}, {errors.New("unknown"), codes.Internal},
	} {
		assert.Equal(t, tc.code, status.Code(ToGrpcError(fmt.Errorf("wrapped: %w", tc.err))))
	}
}
