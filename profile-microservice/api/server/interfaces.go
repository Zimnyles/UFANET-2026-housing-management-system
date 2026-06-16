package server

import (
	"context"
	"profile-service/infra/models/dto"
)

type ProfileService interface {
	GetProfile(ctx context.Context, userID string) (*dto.Profile, error)
	UpsertProfile(ctx context.Context, req *dto.UpsertProfileRequest) (*dto.Profile, error)
	IsProfileComplete(ctx context.Context, userID string) (bool, error)

	CreateManagementCompany(ctx context.Context, req *dto.CreateManagementCompanyRequest) (*dto.ManagementCompany, error)
	ListManagementCompanies(ctx context.Context) ([]*dto.ManagementCompany, error)

	CreateHouse(ctx context.Context, req *dto.CreateHouseRequest) (*dto.House, error)
	ListHouses(ctx context.Context, req *dto.ListHousesRequest) ([]*dto.House, error)
}
