package server

import (
	"context"

	profilepb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/profile/langs/go"
	"github.com/rs/zerolog"
)

type ProfileServer struct {
	profilepb.UnimplementedProfileServiceServer
	service ProfileService
	logger  *zerolog.Logger
}

func NewProfileServer(service ProfileService, logger *zerolog.Logger) *ProfileServer {
	return &ProfileServer{service: service, logger: logger}
}

func (s *ProfileServer) GetProfile(ctx context.Context, req *profilepb.GetProfileRequest) (*profilepb.ProfileResponse, error) {
	p, err := s.service.GetProfile(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &profilepb.ProfileResponse{Profile: domainToProtoProfile(p)}, nil
}

func (s *ProfileServer) UpsertProfile(ctx context.Context, req *profilepb.UpsertProfileRequest) (*profilepb.ProfileResponse, error) {
	p, err := s.service.UpsertProfile(ctx, protoToUpsertProfile(req))
	if err != nil {
		return nil, err
	}
	return &profilepb.ProfileResponse{Profile: domainToProtoProfile(p)}, nil
}

func (s *ProfileServer) IsProfileComplete(ctx context.Context, req *profilepb.IsProfileCompleteRequest) (*profilepb.IsProfileCompleteResponse, error) {
	complete, err := s.service.IsProfileComplete(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &profilepb.IsProfileCompleteResponse{Complete: complete}, nil
}

func (s *ProfileServer) CreateManagementCompany(ctx context.Context, req *profilepb.CreateManagementCompanyRequest) (*profilepb.ManagementCompanyResponse, error) {
	c, err := s.service.CreateManagementCompany(ctx, protoToCreateCompany(req))
	if err != nil {
		return nil, err
	}
	return &profilepb.ManagementCompanyResponse{Company: domainToProtoCompany(c)}, nil
}

func (s *ProfileServer) ListManagementCompanies(ctx context.Context, _ *profilepb.ListManagementCompaniesRequest) (*profilepb.ListManagementCompaniesResponse, error) {
	companies, err := s.service.ListManagementCompanies(ctx)
	if err != nil {
		return nil, err
	}
	pb := make([]*profilepb.ManagementCompany, 0, len(companies))
	for _, c := range companies {
		pb = append(pb, domainToProtoCompany(c))
	}
	return &profilepb.ListManagementCompaniesResponse{Companies: pb}, nil
}

func (s *ProfileServer) CreateHouse(ctx context.Context, req *profilepb.CreateHouseRequest) (*profilepb.HouseResponse, error) {
	h, err := s.service.CreateHouse(ctx, protoToCreateHouse(req))
	if err != nil {
		return nil, err
	}
	return &profilepb.HouseResponse{House: domainToProtoHouse(h)}, nil
}

func (s *ProfileServer) ListHouses(ctx context.Context, req *profilepb.ListHousesRequest) (*profilepb.ListHousesResponse, error) {
	houses, err := s.service.ListHouses(ctx, protoToListHouses(req))
	if err != nil {
		return nil, err
	}
	pb := make([]*profilepb.House, 0, len(houses))
	for _, h := range houses {
		pb = append(pb, domainToProtoHouse(h))
	}
	return &profilepb.ListHousesResponse{Houses: pb}, nil
}
