package domain

import (
	"time"
)

type TokenClaims struct {
	UserID string
	Role   string
}

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}

type RefreshToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type RegisterRequest struct {
	Name     string
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}

type RefreshRequest struct {
	RefreshToken string
}

type LogoutRequest struct {
	RefreshToken string
}
