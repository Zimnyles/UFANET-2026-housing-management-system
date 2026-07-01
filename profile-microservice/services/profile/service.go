package profile_service

import (
	"context"
	"fmt"
	"profile-service/infra/models/domain"

	"github.com/rs/zerolog"
)

type ProfileService struct {
	repo   Repository
	logger *zerolog.Logger
}

func New(repo Repository, logger *zerolog.Logger) *ProfileService {
	return &ProfileService{repo: repo, logger: logger}
}

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	p, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("get profile failed")
		return nil, err
	}
	return p, nil
}

func (s *ProfileService) UpsertProfile(ctx context.Context, req *domain.UpsertProfileRequest) (*domain.Profile, error) {
	profile := &domain.Profile{
		UserID:    req.UserID,
		FullName:  req.FullName,
		Phone:     req.Phone,
		Apartment: req.Apartment,
		HouseID:   req.HouseID,
	}
	p, err := s.repo.UpsertProfile(ctx, profile)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", req.UserID).Msg("upsert profile failed")
		return nil, fmt.Errorf("upsert profile: %w", err)
	}
	s.logger.Info().Str("user_id", req.UserID).Msg("profile upserted")
	return p, nil
}

func (s *ProfileService) IsProfileComplete(ctx context.Context, userID string) (bool, error) {
	complete, err := s.repo.IsProfileComplete(ctx, userID)
	if err != nil {
		return false, err
	}
	return complete, nil
}

func (s *ProfileService) CreateManagementCompany(ctx context.Context, req *domain.CreateManagementCompanyRequest) (*domain.ManagementCompany, error) {
	c, err := s.repo.CreateManagementCompany(ctx, &domain.ManagementCompany{Name: req.Name})
	if err != nil {
		s.logger.Error().Err(err).Str("name", req.Name).Msg("create management company failed")
		return nil, fmt.Errorf("create management company: %w", err)
	}
	s.logger.Info().Str("id", c.ID).Str("name", c.Name).Msg("management company created")
	return c, nil
}

func (s *ProfileService) ListManagementCompanies(ctx context.Context) ([]*domain.ManagementCompany, error) {
	return s.repo.ListManagementCompanies(ctx)
}

func (s *ProfileService) CreateHouse(ctx context.Context, req *domain.CreateHouseRequest) (*domain.House, error) {
	h, err := s.repo.CreateHouse(ctx, &domain.House{
		Name:    req.Name,
		Address: req.Address,
		UKID:    req.UKID,
	})
	if err != nil {
		s.logger.Error().Err(err).Str("uk_id", req.UKID).Msg("create house failed")
		return nil, fmt.Errorf("create house: %w", err)
	}
	s.logger.Info().Str("id", h.ID).Str("name", h.Name).Msg("house created")
	return h, nil
}

func (s *ProfileService) ListHouses(ctx context.Context, req *domain.ListHousesRequest) ([]*domain.House, error) {
	return s.repo.ListHouses(ctx, req.UKID)
}
