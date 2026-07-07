package repository

import (
	"context"

	"auth-service/infra/models/domain"
)

type AuthRepository interface {
	Migrate(ctx context.Context) error
	IsActiveAdminCode(ctx context.Context, code string) (bool, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
