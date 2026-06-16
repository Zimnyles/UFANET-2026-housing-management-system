package repository

import (
	"auth-service/infra/models/domain"
	"context"
)

type AuthRepository interface {
	Migrate(ctx context.Context) error
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
