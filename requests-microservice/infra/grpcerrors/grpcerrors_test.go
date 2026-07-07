package grpcerrors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	infraerrors "requests-service/infra/errors"
	"testing"
)

func TestMapping(t *testing.T) {
	assert.Equal(t, codes.NotFound, status.Code(ToGrpcError(infraerrors.ErrRequestNotFound)))
	for _, e := range []error{infraerrors.ErrRequestIDRequired, infraerrors.ErrUserIDRequired, infraerrors.ErrTitleRequired, infraerrors.ErrTitleTooShort, infraerrors.ErrTitleTooLong, infraerrors.ErrDescriptionRequired, infraerrors.ErrDescriptionTooShort, infraerrors.ErrTypeRequired, infraerrors.ErrTypeInvalid, infraerrors.ErrStatusRequired, infraerrors.ErrStatusInvalid, infraerrors.ErrCommentContentRequired, infraerrors.ErrCommentContentTooLong, infraerrors.ErrCommentRequestIDMissing} {
		assert.Equal(t, codes.InvalidArgument, status.Code(ToGrpcError(e)))
	}
	assert.Equal(t, codes.Internal, status.Code(ToGrpcError(errors.New("x"))))
}
