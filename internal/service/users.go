package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"taski_backend/internal/apperrors"
	"taski_backend/internal/models"
	"taski_backend/internal/repository"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
)


type UsersService struct {
	repo *repository.UsersRepository
}

func NewUsersService(repo *repository.UsersRepository) *UsersService {
	return &UsersService{repo: repo}
}

func (s *UsersService) verifyPassword(password, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false
	}
	return match
}

func (s *UsersService) hashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
	}
	return hash, nil
}

func (s *UsersService) Get(ctx context.Context, id string) (models.UserResponse, error) {
	return s.repo.Get(ctx, id)
}

func (s *UsersService) Create(ctx context.Context, user models.CreateUserRequest) (models.UserResponse, error) {
	hash, err := s.hashPassword(user.Password)
	if err != nil {
		return models.UserResponse{}, apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
	}
	user.Password = hash
	return s.repo.Create(ctx, user)
}

func (s *UsersService) Update(ctx context.Context, user models.UpdateUserRequest) (models.UserResponse, error) {
	if user.Password != "" {
		hash, err := s.hashPassword(user.Password)
		if err != nil {
			return models.UserResponse{}, apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
		}
		user.Password = hash
	}
	return s.repo.Update(ctx, user)
}

func (s *UsersService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *UsersService) Login(ctx context.Context, login models.LoginRequest) (string, string, error) {
	userDB, err := s.repo.GetByEmail(ctx, login.Email)
	if err != nil {
		log.Println("error getting user by email:", err)
		return "", "", apperrors.NewAppError(apperrors.ErrInvalidInput, "invalid_email_or_password")
	}
	if !s.verifyPassword(login.Password, userDB.PasswordHash) {
		log.Println("invalid password")
		return "", "", apperrors.NewAppError(apperrors.ErrInvalidInput, "invalid_email_or_password")
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userDB.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}).SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if err != nil {
		log.Println("error creating refresh token:", err)
		return "", "", apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
	}

	refreshToken := make([]byte, 32)
	_, err = rand.Read(refreshToken)
	if err != nil {
		return "", "", apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
	}

	refreshTokenStr := hex.EncodeToString(refreshToken)

	err = s.repo.CreateRefreshToken(ctx, userDB.ID, refreshTokenStr)
	if err != nil {
		log.Println("error creating refresh token:", err)
		return "", "", apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
	}

	return token, refreshTokenStr, nil
}

func (s *UsersService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	refreshTokenDB, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Println("error getting refresh token:", err)
		return "", apperrors.NewAppError(apperrors.ErrInvalidInput, "invalid_refresh_token")
	}

	if refreshTokenDB.ExpiresAt.Before(time.Now()) {
		log.Println("expired refresh token")
		return "", apperrors.NewAppError(apperrors.ErrInvalidInput, "expired_refresh_token")
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": refreshTokenDB.UserID,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	}).SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		log.Println("error creating token:", err)
		return "", apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
	}

	return token, nil
}