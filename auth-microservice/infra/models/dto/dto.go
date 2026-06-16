package dto

type RegisterRequest struct {
	Name      string
	Email     string
	Password  string
	AdminCode string
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

type AuthResult struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

type RefreshResult struct {
	AccessToken string
}
