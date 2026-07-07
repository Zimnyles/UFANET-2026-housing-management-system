package interceptors

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	newspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/news/langs/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	infraerrors "news-service/infra/errors"
)

func TestInterceptors(t *testing.T) {
	resp, err := ErrorMapping()(context.Background(), nil, nil, func(context.Context, interface{}) (interface{}, error) { return "ok", nil })
	require.NoError(t, err)
	assert.Equal(t, "ok", resp)
	_, err = ErrorMapping()(context.Background(), nil, nil, func(context.Context, interface{}) (interface{}, error) { return nil, infraerrors.ErrNewsNotFound })
	assert.Equal(t, codes.NotFound, status.Code(err))
	_, err = Timeout(time.Second)(context.Background(), nil, nil, func(ctx context.Context, _ interface{}) (interface{}, error) {
		_, ok := ctx.Deadline()
		assert.True(t, ok)
		return nil, nil
	})
	require.NoError(t, err)
}

func TestValidation(t *testing.T) {
	for _, req := range []interface{}{&newspb.CreateNewsRequest{Title: "abc", Content: "long enough", HouseId: "h"}, &newspb.GetNewsItemRequest{Id: "n"}, struct{}{}} {
		_, err := Validation()(context.Background(), req, nil, func(context.Context, interface{}) (interface{}, error) { return nil, errors.New("called") })
		assert.EqualError(t, err, "called")
	}
	_, err := Validation()(context.Background(), &newspb.CreateNewsRequest{}, nil, func(context.Context, interface{}) (interface{}, error) { return nil, nil })
	require.ErrorIs(t, err, infraerrors.ErrTitleRequired)
}
