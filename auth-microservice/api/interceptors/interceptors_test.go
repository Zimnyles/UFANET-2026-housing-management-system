package interceptors

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authpb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/auth/langs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	infraerrors "auth-service/infra/errors"
)

func TestErrorMapping(t *testing.T) {
	resp, err := ErrorMapping()(context.Background(), "req", nil, func(context.Context, interface{}) (interface{}, error) { return "ok", nil })
	require.NoError(t, err)
	assert.Equal(t, "ok", resp)
	resp, err = ErrorMapping()(context.Background(), nil, nil, func(context.Context, interface{}) (interface{}, error) { return nil, infraerrors.ErrUserNotFound })
	assert.Nil(t, resp)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestTimeout(t *testing.T) {
	resp, err := Timeout(time.Second)(context.Background(), nil, nil, func(ctx context.Context, _ interface{}) (interface{}, error) {
		deadline, ok := ctx.Deadline()
		assert.True(t, ok)
		assert.WithinDuration(t, time.Now().Add(time.Second), deadline, 100*time.Millisecond)
		return "ok", nil
	})
	require.NoError(t, err)
	assert.Equal(t, "ok", resp)
}

func TestValidation(t *testing.T) {
	valid := []interface{}{
		&authpb.RegisterRequest{Email: "a@b.ru", Password: "12345678"}, &authpb.LoginRequest{Email: "a@b.ru", Password: "p"},
		&authpb.RefreshRequest{RefreshToken: "r"}, &authpb.LogoutRequest{RefreshToken: "r"}, struct{}{},
	}
	for _, req := range valid {
		called := false
		_, err := Validation()(context.Background(), req, &grpc.UnaryServerInfo{}, func(context.Context, interface{}) (interface{}, error) {
			called = true
			return nil, errors.New("handler")
		})
		assert.EqualError(t, err, "handler")
		assert.True(t, called)
	}
	called := false
	_, err := Validation()(context.Background(), &authpb.RegisterRequest{}, nil, func(context.Context, interface{}) (interface{}, error) {
		called = true
		return nil, nil
	})
	require.ErrorIs(t, err, infraerrors.ErrEmailRequired)
	assert.False(t, called)
}
