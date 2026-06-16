package interceptors

import (
	"auth-service/api/server"
	"context"

	authpb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/auth/langs/go"
	"google.golang.org/grpc"
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
	case *authpb.RegisterRequest:
		return server.ValidateRegisterRequest(r)
	case *authpb.LoginRequest:
		return server.ValidateLoginRequest(r)
	case *authpb.RefreshRequest:
		return server.ValidateRefreshRequest(r)
	case *authpb.LogoutRequest:
		return server.ValidateLogoutRequest(r)
	}
	return nil
}
