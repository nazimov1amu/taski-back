package models

import "time"

type UserDB struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	ID       string `json:"id"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	CodeChallenge string `json:"code_challenge"`
}

type LoginResponse struct {
	Token string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Code string `json:"code,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Token string `json:"access_token"`
}

type RefreshTokenDB struct {
	UserID string `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type CodeRequest struct {
	Code string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
}

type CodeDB struct {
	UserID string `json:"id"`
	CodeChallenge string `json:"code_challenge"`
	ExpiresAt time.Time `json:"expires_at"`
}

