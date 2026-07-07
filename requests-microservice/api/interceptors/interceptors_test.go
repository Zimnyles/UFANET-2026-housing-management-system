package interceptors

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	requestspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/requests/langs/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	infraerrors "requests-service/infra/errors"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	r, e := ErrorMapping()(context.Background(), nil, nil, func(context.Context, interface{}) (interface{}, error) { return "ok", nil })
	require.NoError(t, e)
	assert.Equal(t, "ok", r)
	_, e = ErrorMapping()(context.Background(), nil, nil, func(context.Context, interface{}) (interface{}, error) { return nil, infraerrors.ErrRequestNotFound })
	assert.Equal(t, codes.NotFound, status.Code(e))
	_, e = Timeout(time.Second)(context.Background(), nil, nil, func(c context.Context, _ interface{}) (interface{}, error) {
		_, ok := c.Deadline()
		assert.True(t, ok)
		return nil, nil
	})
	require.NoError(t, e)
	for _, v := range []interface{}{&requestspb.CreateRequestRequest{Title: "abc", Description: "long enough", Type: "plumber", UserId: "u"}, &requestspb.GetRequestRequest{Id: "r"}, &requestspb.UpdateStatusRequest{Id: "r", Status: "done", UserId: "u"}, &requestspb.AddCommentRequest{RequestId: "r", UserId: "u", Content: "x"}, &requestspb.GetCommentsRequest{RequestId: "r"}, struct{}{}} {
		_, e = Validation()(context.Background(), v, nil, func(context.Context, interface{}) (interface{}, error) { return nil, errors.New("called") })
		assert.EqualError(t, e, "called")
	}
	_, e = Validation()(context.Background(), &requestspb.GetRequestRequest{}, nil, func(context.Context, interface{}) (interface{}, error) { return nil, nil })
	require.ErrorIs(t, e, infraerrors.ErrRequestIDRequired)
}
