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
}

type LoginResponse struct {
	Token string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenDB struct {
	UserID string `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Token string `json:"access_token"`
}