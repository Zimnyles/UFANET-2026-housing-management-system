package profile_service

import (
	"context"
	"fmt"
	"profile-service/infra/models/domain"
	"profile-service/infra/models/dto"

	"github.com/rs/zerolog"
)

type ProfileService struct {
	repo   Repository
	logger *zerolog.Logger
}

func New(repo Repository, logger *zerolog.Logger) *ProfileService {
	return &ProfileService{repo: repo, logger: logger}
}

// ─── profile ──────────────────────────────────────────────────────────────────

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*dto.Profile, error) {
	p, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("get profile failed")
		return nil, err
	}
	return domain.ProfileToDTO(p), nil
}

func (s *ProfileService) UpsertProfile(ctx context.Context, req *dto.UpsertProfileRequest) (*dto.Profile, error) {
	domReq := domain.UpsertProfileRequestFromDTO(req)
	profile := &domain.Profile{
		UserID:    domReq.UserID,
		FullName:  domReq.FullName,
		Phone:     domReq.Phone,
		Apartment: domReq.Apartment,
		HouseID:   domReq.HouseID,
	}
	p, err := s.repo.UpsertProfile(ctx, profile)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", req.UserID).Msg("upsert profile failed")
		return nil, fmt.Errorf("upsert profile: %w", err)
	}
	s.logger.Info().Str("user_id", req.UserID).Msg("profile upserted")
	return domain.ProfileToDTO(p), nil
}

func (s *ProfileService) IsProfileComplete(ctx context.Context, userID string) (bool, error) {
	complete, err := s.repo.IsProfileComplete(ctx, userID)
	if err != nil {
		return false, err
	}
	return complete, nil
}

// ─── management companies ─────────────────────────────────────────────────────

func (s *ProfileService) CreateManagementCompany(ctx context.Context, req *dto.CreateManagementCompanyRequest) (*dto.ManagementCompany, error) {
	c, err := s.repo.CreateManagementCompany(ctx, req.Name)
	if err != nil {
		s.logger.Error().Err(err).Str("name", req.Name).Msg("create management company failed")
		return nil, fmt.Errorf("create management company: %w", err)
	}
	s.logger.Info().Str("id", c.ID).Str("name", c.Name).Msg("management company created")
	return domain.ManagementCompanyToDTO(c), nil
}

func (s *ProfileService) ListManagementCompanies(ctx context.Context) ([]*dto.ManagementCompany, error) {
	companies, err := s.repo.ListManagementCompanies(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.ManagementCompany, 0, len(companies))
	for _, c := range companies {
		result = append(result, domain.ManagementCompanyToDTO(c))
	}
	return result, nil
}

// ─── houses ───────────────────────────────────────────────────────────────────

func (s *ProfileService) CreateHouse(ctx context.Context, req *dto.CreateHouseRequest) (*dto.House, error) {
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
	return domain.HouseToDTO(h), nil
}

func (s *ProfileService) ListHouses(ctx context.Context, req *dto.ListHousesRequest) ([]*dto.House, error) {
	houses, err := s.repo.ListHouses(ctx, req.UKID)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.House, 0, len(houses))
	for _, h := range houses {
		result = append(result, domain.HouseToDTO(h))
	}
	return result, nil
}
