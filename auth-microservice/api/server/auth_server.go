package server

import (
	"context"

	authpb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/auth/langs/go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	service AuthService
	logger  *zerolog.Logger
}

func NewAuthServer(service AuthService, logger *zerolog.Logger) *AuthServer {
	return &AuthServer{service: service, logger: logger}
}

func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	result, err := s.service.Register(ctx, ConvertFromProtoToRegisterRequest(req))
	if err != nil {
		return nil, err
	}
	return ConvertFromDomainToAuthResponse(result), nil
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	result, err := s.service.Login(ctx, ConvertFromProtoToLoginRequest(req))
	if err != nil {
		return nil, err
	}
	return ConvertFromDomainToAuthResponse(result), nil
}

func (s *AuthServer) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	result, err := s.service.Refresh(ctx, ConvertFromProtoToRefreshRequest(req))
	if err != nil {
		return nil, err
	}
	return ConvertFromDomainToRefreshResponse(result), nil
}

func (s *AuthServer) Logout(ctx context.Context, req *authpb.LogoutRequest) (*emptypb.Empty, error) {
	if err := s.service.Logout(ctx, ConvertFromProtoToLogoutRequest(req)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
