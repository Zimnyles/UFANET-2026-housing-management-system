package domain

import "auth-service/infra/models/dto"

func ConvertFromDTOToRegisterRequest(d *dto.RegisterRequest) *RegisterRequest {
	return &RegisterRequest{
		Name:     d.Name,
		Email:    d.Email,
		Password: d.Password,
	}
}

func ConvertFromDTOToLoginRequest(d *dto.LoginRequest) *LoginRequest {
	return &LoginRequest{
		Email:    d.Email,
		Password: d.Password,
	}
}

func ConvertFromDTOToRefreshRequest(d *dto.RefreshRequest) *RefreshRequest {
	return &RefreshRequest{
		RefreshToken: d.RefreshToken,
	}
}

func ConvertFromDTOToLogoutRequest(d *dto.LogoutRequest) *LogoutRequest {
	return &LogoutRequest{
		RefreshToken: d.RefreshToken,
	}
}
