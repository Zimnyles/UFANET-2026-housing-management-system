package domain

import (
	"time"
)

type TokenClaims struct {
	UserID string
	Role   string
}

type User struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password;not null"`
	Role         string    `gorm:"not null;default:user"`
	CreatedAt    time.Time `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}

type RefreshToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type RegisterRequest struct {
	Email     string
	Password  string
	AdminCode string
}

type AuthResult struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

type RefreshResult struct {
	AccessToken string
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
