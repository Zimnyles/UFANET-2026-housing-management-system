package auth_client

import (
	"context"

	authpb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/auth/langs/go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClientConnection interface {
	Close() error
}

type AuthServiceClient interface {
	Register(ctx context.Context, in *authpb.RegisterRequest, opts ...grpc.CallOption) (*authpb.AuthResponse, error)
	Login(ctx context.Context, in *authpb.LoginRequest, opts ...grpc.CallOption) (*authpb.AuthResponse, error)
	Refresh(ctx context.Context, in *authpb.RefreshRequest, opts ...grpc.CallOption) (*authpb.RefreshResponse, error)
	Logout(ctx context.Context, in *authpb.LogoutRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}
