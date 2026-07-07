package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	requestspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/requests/langs/go"
	"math"
	infraerrors "requests-service/infra/errors"
	"requests-service/infra/models/domain"
	"strings"
	"testing"
	"time"
)

func TestValidators(t *testing.T) {
	for _, tc := range []struct {
		r *requestspb.CreateRequestRequest
		e error
	}{{&requestspb.CreateRequestRequest{}, infraerrors.ErrTitleRequired}, {&requestspb.CreateRequestRequest{Title: "x"}, infraerrors.ErrTitleTooShort}, {&requestspb.CreateRequestRequest{Title: strings.Repeat("x", 256)}, infraerrors.ErrTitleTooLong}, {&requestspb.CreateRequestRequest{Title: "abc"}, infraerrors.ErrDescriptionRequired}, {&requestspb.CreateRequestRequest{Title: "abc", Description: "x"}, infraerrors.ErrDescriptionTooShort}, {&requestspb.CreateRequestRequest{Title: "abc", Description: "long enough"}, infraerrors.ErrTypeRequired}, {&requestspb.CreateRequestRequest{Title: "abc", Description: "long enough", Type: "x"}, infraerrors.ErrTypeInvalid}, {&requestspb.CreateRequestRequest{Title: "abc", Description: "long enough", Type: "plumber"}, infraerrors.ErrUserIDRequired}} {
		require.ErrorIs(t, ValidateCreateRequest(tc.r), tc.e)
	}
	require.NoError(t, ValidateCreateRequest(&requestspb.CreateRequestRequest{Title: "abc", Description: "long enough", Type: "plumber", UserId: "u"}))
	require.ErrorIs(t, ValidateGetRequest(&requestspb.GetRequestRequest{}), infraerrors.ErrRequestIDRequired)
	require.NoError(t, ValidateGetRequest(&requestspb.GetRequestRequest{Id: "r"}))
	for _, tc := range []struct {
		r *requestspb.UpdateStatusRequest
		e error
	}{{&requestspb.UpdateStatusRequest{}, infraerrors.ErrRequestIDRequired}, {&requestspb.UpdateStatusRequest{Id: "r"}, infraerrors.ErrStatusRequired}, {&requestspb.UpdateStatusRequest{Id: "r", Status: "x"}, infraerrors.ErrStatusInvalid}, {&requestspb.UpdateStatusRequest{Id: "r", Status: "done"}, infraerrors.ErrUserIDRequired}} {
		require.ErrorIs(t, ValidateUpdateStatus(tc.r), tc.e)
	}
	require.NoError(t, ValidateUpdateStatus(&requestspb.UpdateStatusRequest{Id: "r", Status: "done", UserId: "u"}))
	for _, tc := range []struct {
		r *requestspb.AddCommentRequest
		e error
	}{{&requestspb.AddCommentRequest{}, infraerrors.ErrCommentRequestIDMissing}, {&requestspb.AddCommentRequest{RequestId: "r"}, infraerrors.ErrUserIDRequired}, {&requestspb.AddCommentRequest{RequestId: "r", UserId: "u"}, infraerrors.ErrCommentContentRequired}, {&requestspb.AddCommentRequest{RequestId: "r", UserId: "u", Content: strings.Repeat("x", 1001)}, infraerrors.ErrCommentContentTooLong}} {
		require.ErrorIs(t, ValidateAddComment(tc.r), tc.e)
	}
	require.NoError(t, ValidateAddComment(&requestspb.AddCommentRequest{RequestId: "r", UserId: "u", Content: "x"}))
	require.ErrorIs(t, ValidateGetComments(&requestspb.GetCommentsRequest{}), infraerrors.ErrCommentRequestIDMissing)
	require.NoError(t, ValidateGetComments(&requestspb.GetCommentsRequest{RequestId: "r"}))
	assert.True(t, validType("electrician"))
	assert.True(t, validStatus("open"))
	assert.True(t, validStatus("in_progress"))
	assert.True(t, validStatus("cancelled"))
}
func TestConverters(t *testing.T) {
	assert.Equal(t, &domain.CreateRequestRequest{Title: "t", Description: "d", Type: "p", UserID: "u"}, protoToCreateRequest(&requestspb.CreateRequestRequest{Title: "t", Description: "d", Type: "p", UserId: "u"}))
	assert.Equal(t, "u", protoToGetRequests(&requestspb.GetRequestsRequest{UserId: "u"}).UserID)
	assert.Equal(t, "r", protoToUpdateStatus(&requestspb.UpdateStatusRequest{Id: "r"}).ID)
	assert.Equal(t, "r", protoToAddComment(&requestspb.AddCommentRequest{RequestId: "r"}).RequestID)
	now := time.Now()
	assert.Equal(t, "r", domainToProtoRequest(&domain.MaintenanceRequest{ID: "r", CreatedAt: now, UpdatedAt: now}).Id)
	assert.Equal(t, "c", domainToProtoComment(&domain.Comment{ID: "c", CreatedAt: now}).Id)
	assert.Equal(t, int32(0), totalToInt32(-1))
	assert.Equal(t, int32(2), totalToInt32(2))
	assert.Equal(t, int32(math.MaxInt32), totalToInt32(math.MaxInt64))
}
