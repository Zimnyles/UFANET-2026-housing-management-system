package profile_handler

import (
	"context"

	"api-gateway/internal/models/domain"
)

type ProfileService interface {
	GetProfile(ctx context.Context, userID string) (*domain.Profile, error)
	UpsertProfile(ctx context.Context, req *domain.UpsertProfileRequest) (*domain.Profile, error)

	CreateManagementCompany(ctx context.Context, req *domain.CreateManagementCompanyRequest) (*domain.ManagementCompany, error)
	ListManagementCompanies(ctx context.Context) ([]*domain.ManagementCompany, error)

	CreateHouse(ctx context.Context, req *domain.CreateHouseRequest) (*domain.House, error)
	ListHouses(ctx context.Context, ukID string) ([]*domain.House, error)
}
