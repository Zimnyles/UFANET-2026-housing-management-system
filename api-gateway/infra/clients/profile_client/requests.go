package profile_client

import (
	"context"

	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/domain"
)

func (c *ProfileClient) GetProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	resp, err := c.client.GetProfile(ctx, toProtoGetProfileRequest(userID))
	if err != nil {
		return nil, app_errors.FromGRPC(err)
	}

	return toDomainProfile(resp.GetProfile()), nil
}

func (c *ProfileClient) UpsertProfile(ctx context.Context, req *domain.UpsertProfileRequest) (*domain.Profile, error) {
	resp, err := c.client.UpsertProfile(ctx, toProtoUpsertProfileRequest(req))
	if err != nil {
		return nil, app_errors.FromGRPC(err)
	}

	return toDomainProfile(resp.GetProfile()), nil
}

func (c *ProfileClient) IsProfileComplete(ctx context.Context, userID string) (bool, error) {
	resp, err := c.client.IsProfileComplete(ctx, toProtoIsProfileCompleteRequest(userID))
	if err != nil {
		return false, app_errors.FromGRPC(err)
	}

	return resp.GetComplete(), nil
}

func (c *ProfileClient) CreateManagementCompany(ctx context.Context, req *domain.CreateManagementCompanyRequest) (*domain.ManagementCompany, error) {
	resp, err := c.client.CreateManagementCompany(ctx, toProtoCreateManagementCompanyRequest(req))
	if err != nil {
		return nil, app_errors.FromGRPC(err)
	}

	return toDomainManagementCompany(resp.GetCompany()), nil
}

func (c *ProfileClient) ListManagementCompanies(ctx context.Context) ([]*domain.ManagementCompany, error) {
	resp, err := c.client.ListManagementCompanies(ctx, toProtoListManagementCompaniesRequest())
	if err != nil {
		return nil, app_errors.FromGRPC(err)
	}

	result := make([]*domain.ManagementCompany, 0, len(resp.GetCompanies()))
	for _, c := range resp.GetCompanies() {
		result = append(result, toDomainManagementCompany(c))
	}

	return result, nil
}

func (c *ProfileClient) CreateHouse(ctx context.Context, req *domain.CreateHouseRequest) (*domain.House, error) {
	resp, err := c.client.CreateHouse(ctx, toProtoCreateHouseRequest(req))
	if err != nil {
		return nil, app_errors.FromGRPC(err)
	}

	return toDomainHouse(resp.GetHouse()), nil
}

func (c *ProfileClient) ListHouses(ctx context.Context, ukID string) ([]*domain.House, error) {
	resp, err := c.client.ListHouses(ctx, toProtoListHousesRequest(ukID))
	if err != nil {
		return nil, app_errors.FromGRPC(err)
	}

	result := make([]*domain.House, 0, len(resp.GetHouses()))
	for _, h := range resp.GetHouses() {
		result = append(result, toDomainHouse(h))
	}

	return result, nil
}
