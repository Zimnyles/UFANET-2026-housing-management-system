package server

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authpb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/auth/langs/go"

	infraerrors "auth-service/infra/errors"
	"auth-service/infra/models/domain"
)

func TestValidators(t *testing.T) {
	long := strings.Repeat("x", 73)
	for _, tc := range []struct {
		name string
		run  func() error
		want error
	}{
		{"register email required", func() error { return ValidateRegisterRequest(&authpb.RegisterRequest{}) }, infraerrors.ErrEmailRequired},
		{"register invalid email", func() error { return ValidateRegisterRequest(&authpb.RegisterRequest{Email: "x"}) }, infraerrors.ErrEmailInvalid},
		{"register password required", func() error { return ValidateRegisterRequest(&authpb.RegisterRequest{Email: "a@b.ru"}) }, infraerrors.ErrPasswordRequired},
		{"register short password", func() error {
			return ValidateRegisterRequest(&authpb.RegisterRequest{Email: "a@b.ru", Password: "123"})
		}, infraerrors.ErrPasswordTooShort},
		{"register long password", func() error { return ValidateRegisterRequest(&authpb.RegisterRequest{Email: "a@b.ru", Password: long}) }, infraerrors.ErrPasswordTooLong},
		{"register long admin code", func() error {
			return ValidateRegisterRequest(&authpb.RegisterRequest{Email: "a@b.ru", Password: "12345678", AdminCode: long})
		}, infraerrors.ErrInvalidAdminCode},
		{"login email required", func() error { return ValidateLoginRequest(&authpb.LoginRequest{}) }, infraerrors.ErrEmailRequired},
		{"login invalid email", func() error { return ValidateLoginRequest(&authpb.LoginRequest{Email: "x"}) }, infraerrors.ErrEmailInvalid},
		{"login password required", func() error { return ValidateLoginRequest(&authpb.LoginRequest{Email: "a@b.ru"}) }, infraerrors.ErrPasswordRequired},
		{"refresh required", func() error { return ValidateRefreshRequest(&authpb.RefreshRequest{}) }, infraerrors.ErrRefreshRequired},
		{"logout required", func() error { return ValidateLogoutRequest(&authpb.LogoutRequest{}) }, infraerrors.ErrRefreshRequired},
	} {
		t.Run(tc.name, func(t *testing.T) { require.ErrorIs(t, tc.run(), tc.want) })
	}
	require.NoError(t, ValidateRegisterRequest(&authpb.RegisterRequest{Email: "a@b.ru", Password: "12345678"}))
	require.NoError(t, ValidateLoginRequest(&authpb.LoginRequest{Email: "a@b.ru", Password: "x"}))
	require.NoError(t, ValidateRefreshRequest(&authpb.RefreshRequest{RefreshToken: "r"}))
	require.NoError(t, ValidateLogoutRequest(&authpb.LogoutRequest{RefreshToken: "r"}))
}

func TestConverters(t *testing.T) {
	assert.Equal(t, &domain.RegisterRequest{Email: "a", Password: "p", AdminCode: "c"}, ConvertFromProtoToRegisterRequest(&authpb.RegisterRequest{Email: "a", Password: "p", AdminCode: "c"}))
	assert.Equal(t, &domain.LoginRequest{Email: "a", Password: "p"}, ConvertFromProtoToLoginRequest(&authpb.LoginRequest{Email: "a", Password: "p"}))
	assert.Equal(t, &domain.RefreshRequest{RefreshToken: "r"}, ConvertFromProtoToRefreshRequest(&authpb.RefreshRequest{RefreshToken: "r"}))
	assert.Equal(t, &domain.LogoutRequest{RefreshToken: "r"}, ConvertFromProtoToLogoutRequest(&authpb.LogoutRequest{RefreshToken: "r"}))
	assert.Equal(t, &authpb.AuthResponse{UserId: "u", AccessToken: "a", RefreshToken: "r"}, ConvertFromDomainToAuthResponse(&domain.AuthResult{UserID: "u", AccessToken: "a", RefreshToken: "r"}))
	assert.Equal(t, &authpb.RefreshResponse{AccessToken: "a"}, ConvertFromDomainToRefreshResponse(&domain.RefreshResult{AccessToken: "a"}))
}

func TestAuthServer(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()
	ctrl := gomock.NewController(t)
	service := NewMockAuthService(ctrl)
	server := NewAuthServer(service, &logger)
	service.EXPECT().Register(ctx, &domain.RegisterRequest{Email: "a", Password: "p"}).Return(&domain.AuthResult{UserID: "u"}, nil)
	got, err := server.Register(ctx, &authpb.RegisterRequest{Email: "a", Password: "p"})
	require.NoError(t, err)
	assert.Equal(t, "u", got.UserId)
	service.EXPECT().Login(ctx, gomock.Any()).Return(nil, errors.New("login"))
	_, err = server.Login(ctx, &authpb.LoginRequest{})
	assert.EqualError(t, err, "login")
	service.EXPECT().Refresh(ctx, &domain.RefreshRequest{RefreshToken: "r"}).Return(&domain.RefreshResult{AccessToken: "a"}, nil)
	refresh, err := server.Refresh(ctx, &authpb.RefreshRequest{RefreshToken: "r"})
	require.NoError(t, err)
	assert.Equal(t, "a", refresh.AccessToken)
	service.EXPECT().Logout(ctx, gomock.Any()).Return(nil)
	empty, err := server.Logout(ctx, &authpb.LogoutRequest{RefreshToken: "r"})
	require.NoError(t, err)
	require.NotNil(t, empty)

	ctrl = gomock.NewController(t)
	service = NewMockAuthService(ctrl)
	server = NewAuthServer(service, &logger)
	service.EXPECT().Register(ctx, gomock.Any()).Return(nil, errors.New("register"))
	_, err = server.Register(ctx, &authpb.RegisterRequest{})
	assert.EqualError(t, err, "register")
	service.EXPECT().Login(ctx, gomock.Any()).Return(&domain.AuthResult{AccessToken: "a"}, nil)
	login, err := server.Login(ctx, &authpb.LoginRequest{})
	require.NoError(t, err)
	assert.Equal(t, "a", login.AccessToken)
	service.EXPECT().Refresh(ctx, gomock.Any()).Return(nil, errors.New("refresh"))
	_, err = server.Refresh(ctx, &authpb.RefreshRequest{})
	assert.EqualError(t, err, "refresh")
	service.EXPECT().Logout(ctx, gomock.Any()).Return(errors.New("logout"))
	_, err = server.Logout(ctx, &authpb.LogoutRequest{})
	assert.EqualError(t, err, "logout")
}
