package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	infraerrors "auth-service/infra/errors"
	"auth-service/infra/models/domain"
)

type fixture struct {
	repo   *MockRepository
	jwt    *MockJWTManager
	hasher *MockHasher
	svc    *AuthService
}

func newFixture(t *testing.T) fixture {
	ctrl := gomock.NewController(t)
	repo, jwtManager, passwordHasher := NewMockRepository(ctrl), NewMockJWTManager(ctrl), NewMockHasher(ctrl)
	logger := zerolog.Nop()
	return fixture{repo: repo, jwt: jwtManager, hasher: passwordHasher, svc: New(repo, jwtManager, passwordHasher, &logger)}
}

func (f fixture) expectTokens(user *domain.User) {
	f.jwt.EXPECT().GenerateAccess(user.ID, user.Role).Return("access", nil)
	f.jwt.EXPECT().GenerateRefresh(user.ID, user.Role).Return("refresh", nil)
}

func TestAuthServiceRegister(t *testing.T) {
	ctx := context.Background()
	req := &domain.RegisterRequest{Email: "user@example.com", Password: "password"}

	t.Run("user success", func(t *testing.T) {
		f := newFixture(t)
		f.hasher.EXPECT().Hash("password").Return("hash", nil)
		user := &domain.User{ID: "u1", Email: req.Email, PasswordHash: "hash", Role: "user"}
		f.repo.EXPECT().CreateUser(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, got *domain.User) (*domain.User, error) {
			assert.Equal(t, "user", got.Role)
			assert.Equal(t, "hash", got.PasswordHash)
			got.ID = "u1"
			return got, nil
		})
		f.expectTokens(user)
		got, err := f.svc.Register(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, &domain.AuthResult{UserID: "u1", AccessToken: "access", RefreshToken: "refresh"}, got)
	})

	t.Run("admin success", func(t *testing.T) {
		f := newFixture(t)
		adminReq := &domain.RegisterRequest{Email: req.Email, Password: req.Password, AdminCode: "code"}
		f.repo.EXPECT().IsActiveAdminCode(ctx, "code").Return(true, nil)
		f.hasher.EXPECT().Hash("password").Return("hash", nil)
		user := &domain.User{ID: "a1", Email: req.Email, PasswordHash: "hash", Role: "admin"}
		f.repo.EXPECT().CreateUser(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, got *domain.User) (*domain.User, error) {
			assert.Equal(t, "admin", got.Role)
			got.ID = "a1"
			return got, nil
		})
		f.expectTokens(user)
		got, err := f.svc.Register(ctx, adminReq)
		require.NoError(t, err)
		assert.Equal(t, "a1", got.UserID)
	})

	t.Run("admin code errors", func(t *testing.T) {
		for _, tc := range []struct {
			name          string
			active        bool
			repoErr, want error
		}{
			{name: "inactive", want: infraerrors.ErrInvalidAdminCode},
			{name: "repository", repoErr: errors.New("db"), want: errors.New("db")},
		} {
			t.Run(tc.name, func(t *testing.T) {
				f := newFixture(t)
				f.repo.EXPECT().IsActiveAdminCode(ctx, "bad").Return(tc.active, tc.repoErr)
				_, err := f.svc.Register(ctx, &domain.RegisterRequest{AdminCode: "bad"})
				if tc.repoErr != nil {
					assert.EqualError(t, err, "db")
				} else {
					require.ErrorIs(t, err, tc.want)
				}
			})
		}
	})

	t.Run("hash error", func(t *testing.T) {
		f := newFixture(t)
		f.hasher.EXPECT().Hash(req.Password).Return("", errors.New("hash failed"))
		_, err := f.svc.Register(ctx, req)
		assert.EqualError(t, err, "hash password: hash failed")
	})

	t.Run("create error", func(t *testing.T) {
		f := newFixture(t)
		f.hasher.EXPECT().Hash(req.Password).Return("hash", nil)
		f.repo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil, infraerrors.ErrEmailAlreadyExists)
		_, err := f.svc.Register(ctx, req)
		require.ErrorIs(t, err, infraerrors.ErrEmailAlreadyExists)
	})
}

func TestAuthServiceLoginRefreshAndLogout(t *testing.T) {
	ctx := context.Background()
	t.Run("login success", func(t *testing.T) {
		f := newFixture(t)
		user := &domain.User{ID: "u", Email: "a@b.ru", PasswordHash: "hash", Role: "user"}
		f.repo.EXPECT().GetUserByEmail(ctx, user.Email).Return(user, nil)
		f.hasher.EXPECT().Check("password", "hash").Return(true)
		f.expectTokens(user)
		got, err := f.svc.Login(ctx, &domain.LoginRequest{Email: user.Email, Password: "password"})
		require.NoError(t, err)
		assert.Equal(t, "access", got.AccessToken)
	})

	t.Run("login repository and password errors", func(t *testing.T) {
		f := newFixture(t)
		f.repo.EXPECT().GetUserByEmail(ctx, "missing").Return(nil, infraerrors.ErrUserNotFound)
		_, err := f.svc.Login(ctx, &domain.LoginRequest{Email: "missing"})
		require.ErrorIs(t, err, infraerrors.ErrUserNotFound)

		f = newFixture(t)
		f.repo.EXPECT().GetUserByEmail(ctx, "a").Return(&domain.User{PasswordHash: "h"}, nil)
		f.hasher.EXPECT().Check("bad", "h").Return(false)
		_, err = f.svc.Login(ctx, &domain.LoginRequest{Email: "a", Password: "bad"})
		require.ErrorIs(t, err, infraerrors.ErrInvalidCredentials)
	})

	t.Run("token generation errors", func(t *testing.T) {
		f := newFixture(t)
		user := &domain.User{ID: "u", Role: "user"}
		f.repo.EXPECT().GetUserByEmail(ctx, "a").Return(user, nil)
		f.hasher.EXPECT().Check("p", "").Return(true)
		f.jwt.EXPECT().GenerateAccess("u", "user").Return("", errors.New("sign"))
		_, err := f.svc.Login(ctx, &domain.LoginRequest{Email: "a", Password: "p"})
		assert.EqualError(t, err, "generate access token: sign")

		f = newFixture(t)
		f.repo.EXPECT().GetUserByEmail(ctx, "a").Return(user, nil)
		f.hasher.EXPECT().Check("p", "").Return(true)
		f.jwt.EXPECT().GenerateAccess("u", "user").Return("a", nil)
		f.jwt.EXPECT().GenerateRefresh("u", "user").Return("", errors.New("sign"))
		_, err = f.svc.Login(ctx, &domain.LoginRequest{Email: "a", Password: "p"})
		assert.EqualError(t, err, "generate refresh token: sign")
	})

	t.Run("refresh", func(t *testing.T) {
		f := newFixture(t)
		f.jwt.EXPECT().ParseRefresh("r").Return(&domain.TokenClaims{UserID: "u", Role: "admin"}, nil)
		f.jwt.EXPECT().GenerateAccess("u", "admin").Return("new", nil)
		got, err := f.svc.Refresh(ctx, &domain.RefreshRequest{RefreshToken: "r"})
		require.NoError(t, err)
		assert.Equal(t, &domain.RefreshResult{AccessToken: "new"}, got)

		f = newFixture(t)
		f.jwt.EXPECT().ParseRefresh("bad").Return(nil, errors.New("parse"))
		_, err = f.svc.Refresh(ctx, &domain.RefreshRequest{RefreshToken: "bad"})
		require.ErrorIs(t, err, infraerrors.ErrInvalidToken)

		f = newFixture(t)
		f.jwt.EXPECT().ParseRefresh("r").Return(&domain.TokenClaims{UserID: "u", Role: "user"}, nil)
		f.jwt.EXPECT().GenerateAccess("u", "user").Return("", errors.New("sign"))
		_, err = f.svc.Refresh(ctx, &domain.RefreshRequest{RefreshToken: "r"})
		assert.EqualError(t, err, "generate access token: sign")
	})

	f := newFixture(t)
	require.NoError(t, f.svc.Logout(ctx, &domain.LogoutRequest{}))
}
