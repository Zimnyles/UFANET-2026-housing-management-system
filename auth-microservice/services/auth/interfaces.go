package auth

import (
	"context"

	"auth-service/infra/models/domain"
)

type Repository interface {
	IsActiveAdminCode(ctx context.Context, code string) (bool, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type JWTManager interface {
	GenerateAccess(userID, role string) (string, error)
	GenerateRefresh(userID, role string) (string, error)
	ParseRefresh(tokenStr string) (*domain.TokenClaims, error)
}

type Hasher interface {
	Hash(password string) (string, error)
	Check(password, hash string) bool
}
