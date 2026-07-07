package interceptors

import (
	"context"

	profilepb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/profile/langs/go"
	"google.golang.org/grpc"

	"profile-service/api/server"
)

func Validation() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := validate(req); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func validate(req interface{}) error {
	switch r := req.(type) {
	case *profilepb.GetProfileRequest:
		return server.ValidateGetProfileRequest(r)
	case *profilepb.UpsertProfileRequest:
		return server.ValidateUpsertProfileRequest(r)
	case *profilepb.CreateManagementCompanyRequest:
		return server.ValidateCreateManagementCompanyRequest(r)
	case *profilepb.CreateHouseRequest:
		return server.ValidateCreateHouseRequest(r)
	}

	return nil
}
