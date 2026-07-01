package profile_service

import (
	"context"
	"profile-service/infra/models/domain"
)

type Repository interface {
	GetProfile(ctx context.Context, userID string) (*domain.Profile, error)
	UpsertProfile(ctx context.Context, profile *domain.Profile) (*domain.Profile, error)
	IsProfileComplete(ctx context.Context, userID string) (bool, error)

	CreateManagementCompany(ctx context.Context, company *domain.ManagementCompany) (*domain.ManagementCompany, error)
	ListManagementCompanies(ctx context.Context) ([]*domain.ManagementCompany, error)

	CreateHouse(ctx context.Context, house *domain.House) (*domain.House, error)
	ListHouses(ctx context.Context, ukID string) ([]*domain.House, error)
}
