package auth

import "auth-service/infra/models/domain"

func registerRequestToUser(req *domain.RegisterRequest, passwordHash, role string) *domain.User {
	return &domain.User{Email: req.Email, PasswordHash: passwordHash, Role: role}
}

func tokensToAuthResult(userID, accessToken, refreshToken string) *domain.AuthResult {
	return &domain.AuthResult{UserID: userID, AccessToken: accessToken, RefreshToken: refreshToken}
}

func accessTokenToRefreshResult(accessToken string) *domain.RefreshResult {
	return &domain.RefreshResult{AccessToken: accessToken}
}
