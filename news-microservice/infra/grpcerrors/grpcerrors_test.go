package grpcerrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	infraerrors "news-service/infra/errors"
)

func TestToGrpcError(t *testing.T) {
	for _, tc := range []struct {
		e error
		c codes.Code
	}{
		{infraerrors.ErrNewsNotFound, codes.NotFound}, {infraerrors.ErrTitleRequired, codes.InvalidArgument}, {infraerrors.ErrTitleTooShort, codes.InvalidArgument},
		{infraerrors.ErrTitleTooLong, codes.InvalidArgument}, {infraerrors.ErrContentRequired, codes.InvalidArgument}, {infraerrors.ErrContentTooShort, codes.InvalidArgument},
		{infraerrors.ErrHouseIDRequired, codes.InvalidArgument}, {infraerrors.ErrNewsIDRequired, codes.InvalidArgument}, {errors.New("x"), codes.Internal},
	} {
		assert.Equal(t, tc.c, status.Code(ToGrpcError(fmt.Errorf("wrap: %w", tc.e))))
	}
}
