package domain

type RegisterRequest struct {
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

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthResult struct {
	UserID string
	Tokens TokenPair
}
