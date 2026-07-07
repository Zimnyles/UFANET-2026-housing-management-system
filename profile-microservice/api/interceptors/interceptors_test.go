package interceptors

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	profilepb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/profile/langs/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	infraerrors "profile-service/infra/errors"
	"testing"
	"time"
)

func TestInterceptors(t *testing.T) {
	r, e := ErrorMapping()(context.Background(), nil, nil, func(context.Context, interface{}) (interface{}, error) { return "ok", nil })
	require.NoError(t, e)
	assert.Equal(t, "ok", r)
	_, e = ErrorMapping()(context.Background(), nil, nil, func(context.Context, interface{}) (interface{}, error) { return nil, infraerrors.ErrProfileNotFound })
	assert.Equal(t, codes.NotFound, status.Code(e))
	_, e = Timeout(time.Second)(context.Background(), nil, nil, func(c context.Context, _ interface{}) (interface{}, error) {
		_, ok := c.Deadline()
		assert.True(t, ok)
		return nil, nil
	})
	require.NoError(t, e)
}
func TestValidation(t *testing.T) {
	for _, r := range []interface{}{&profilepb.GetProfileRequest{UserId: "u"}, &profilepb.UpsertProfileRequest{UserId: "u", FullName: "nn", Phone: "p", Apartment: "a", HouseId: "h"}, &profilepb.CreateManagementCompanyRequest{Name: "c"}, &profilepb.CreateHouseRequest{Name: "h", Address: "a", UkId: "u"}, struct{}{}} {
		_, e := Validation()(context.Background(), r, nil, func(context.Context, interface{}) (interface{}, error) { return nil, errors.New("called") })
		assert.EqualError(t, e, "called")
	}
	_, e := Validation()(context.Background(), &profilepb.GetProfileRequest{}, nil, func(context.Context, interface{}) (interface{}, error) { return nil, nil })
	require.ErrorIs(t, e, infraerrors.ErrUserIDRequired)
}
