package profile_service

import (
	"context"

	"github.com/rs/zerolog"

	"api-gateway/internal/models/domain"
)

type ProfileService struct {
	client ProfileClient
	logger *zerolog.Logger
}

func New(client ProfileClient, logger *zerolog.Logger) *ProfileService {
	return &ProfileService{client: client, logger: logger}
}

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	p, err := s.client.GetProfile(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("get profile failed")

		return nil, err
	}

	return p, nil
}

func (s *ProfileService) UpsertProfile(ctx context.Context, req *domain.UpsertProfileRequest) (*domain.Profile, error) {
	p, err := s.client.UpsertProfile(ctx, req)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", req.UserID).Msg("upsert profile failed")

		return nil, err
	}

	s.logger.Info().Str("user_id", req.UserID).Msg("profile upserted")

	return p, nil
}

func (s *ProfileService) IsProfileComplete(ctx context.Context, userID string) (bool, error) {
	return s.client.IsProfileComplete(ctx, userID)
}

func (s *ProfileService) CreateManagementCompany(ctx context.Context, req *domain.CreateManagementCompanyRequest) (*domain.ManagementCompany, error) {
	c, err := s.client.CreateManagementCompany(ctx, req)
	if err != nil {
		s.logger.Error().Err(err).Str("name", req.Name).Msg("create management company failed")

		return nil, err
	}

	s.logger.Info().Str("id", c.ID).Msg("management company created")

	return c, nil
}

func (s *ProfileService) ListManagementCompanies(ctx context.Context) ([]*domain.ManagementCompany, error) {
	return s.client.ListManagementCompanies(ctx)
}

func (s *ProfileService) CreateHouse(ctx context.Context, req *domain.CreateHouseRequest) (*domain.House, error) {
	h, err := s.client.CreateHouse(ctx, req)
	if err != nil {
		s.logger.Error().Err(err).Str("uk_id", req.UKID).Msg("create house failed")

		return nil, err
	}

	s.logger.Info().Str("id", h.ID).Msg("house created")

	return h, nil
}

func (s *ProfileService) ListHouses(ctx context.Context, ukID string) ([]*domain.House, error) {
	return s.client.ListHouses(ctx, ukID)
}
