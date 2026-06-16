package profile_client

import (
	"context"

	profilepb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/profile/langs/go"
	"google.golang.org/grpc"
)

type ClientConnection interface {
	Close() error
}

type ProfileServiceClient interface {
	GetProfile(ctx context.Context, in *profilepb.GetProfileRequest, opts ...grpc.CallOption) (*profilepb.ProfileResponse, error)
	UpsertProfile(ctx context.Context, in *profilepb.UpsertProfileRequest, opts ...grpc.CallOption) (*profilepb.ProfileResponse, error)
	IsProfileComplete(ctx context.Context, in *profilepb.IsProfileCompleteRequest, opts ...grpc.CallOption) (*profilepb.IsProfileCompleteResponse, error)

	CreateManagementCompany(ctx context.Context, in *profilepb.CreateManagementCompanyRequest, opts ...grpc.CallOption) (*profilepb.ManagementCompanyResponse, error)
	ListManagementCompanies(ctx context.Context, in *profilepb.ListManagementCompaniesRequest, opts ...grpc.CallOption) (*profilepb.ListManagementCompaniesResponse, error)

	CreateHouse(ctx context.Context, in *profilepb.CreateHouseRequest, opts ...grpc.CallOption) (*profilepb.HouseResponse, error)
	ListHouses(ctx context.Context, in *profilepb.ListHousesRequest, opts ...grpc.CallOption) (*profilepb.ListHousesResponse, error)
}
