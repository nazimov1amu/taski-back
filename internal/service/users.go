package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"os"
	"taski_backend/internal/apperrors"
	"taski_backend/internal/models"
	"taski_backend/internal/repository"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type UsersService struct {
	repo   *repository.UsersRepository
	logger *zap.SugaredLogger
}

func NewUsersService(repo *repository.UsersRepository, logger *zap.SugaredLogger) *UsersService {
	return &UsersService{repo: repo, logger: logger}
}

func (s *UsersService) verifyPassword(password, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		s.logger.Error("error comparing password and hash", "error", err)
		return false
	}
	return match
}

func (s *UsersService) hashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		s.logger.Error("error hashing password", "error", err)
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
		s.logger.Error("error hashing password", "error", err)
		return models.UserResponse{}, apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
	}
	user.Password = hash
	return s.repo.Create(ctx, user)
}

func (s *UsersService) Update(ctx context.Context, user models.UpdateUserRequest) (models.UserResponse, error) {
	if user.Password != "" {
		hash, err := s.hashPassword(user.Password)
		if err != nil {
			s.logger.Error("error hashing password", "error", err)
			return models.UserResponse{}, apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
		}
		user.Password = hash
	}
	return s.repo.Update(ctx, user)
}

func (s *UsersService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *UsersService) CreateTokens(ctx context.Context, userID string) (string, string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}).SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		s.logger.Error("error creating token", "error", err)
		return "", "", apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
	}

	refreshToken := make([]byte, 32)
	_, err = rand.Read(refreshToken)
	if err != nil {
		return "", "", apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
	}
	refreshTokenStr := hex.EncodeToString(refreshToken)

	err = s.repo.CreateRefreshToken(ctx, userID, refreshTokenStr)
	if err != nil {
		s.logger.Error("error creating refresh token", "error", err)
		return "", "", apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
	}

	return token, refreshTokenStr, nil
}

func (s *UsersService) Login(ctx context.Context, login models.LoginRequest) (string, string, string, error) {
	userDB, err := s.repo.GetByEmail(ctx, login.Email)
	if err != nil {
		s.logger.Error("error getting user by email", "error", err)
		return "", "", "", apperrors.NewAppError(apperrors.ErrForbidden, "invalid_email_or_password")
	}
	if !s.verifyPassword(login.Password, userDB.PasswordHash) {
		s.logger.Error("invalid password")
		return "", "", "", apperrors.NewAppError(apperrors.ErrForbidden, "invalid_email_or_password")
	}

	token, refreshToken, err := s.CreateTokens(ctx, userDB.ID)
	if err != nil {
		s.logger.Error("error creating tokens", "error", err)
		return "", "", "", err
	}

	if login.CodeChallenge != "" {
		code := make([]byte, 6)
		_, err = rand.Read(code)
		if err != nil {
			s.logger.Error("error creating code", "error", err)
			return "", "", "", apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
		}
		codeStr := hex.EncodeToString(code)

		err = s.repo.CreateCode(ctx, userDB.ID, codeStr, login.CodeChallenge)
		if err != nil {
			s.logger.Error("error creating code", "error", err)
			return "", "", "", apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
		}
		return token, refreshToken, codeStr, nil
	}

	return token, refreshToken, "", nil
}

func (s *UsersService) VerifyCode(ctx context.Context, request models.CodeRequest) (string, string, error) {
	codeDB, err := s.repo.GetCode(ctx, request.Code)

	if err != nil {
		s.logger.Error("error getting code", "error", err)
		return "", "", apperrors.NewAppError(apperrors.ErrForbidden, "invalid_code")
	}

	if codeDB.ExpiresAt.Before(time.Now()) {
		s.logger.Error("expired code")
		return "", "", apperrors.NewAppError(apperrors.ErrForbidden, "expired_code")
	}

	reqCodeSum := sha256.Sum256([]byte(request.CodeVerifier))
	challenge := base64.RawURLEncoding.EncodeToString(reqCodeSum[:])
	if challenge != codeDB.CodeChallenge {
		s.logger.Error("invalid code challenge")
		return "", "", apperrors.NewAppError(apperrors.ErrForbidden, "invalid_code_challenge")
	}

	jwtToken, refreshToken, err := s.CreateTokens(ctx, codeDB.UserID)
	if err != nil {
		s.logger.Error("error creating tokens", "error", err)
		return "", "", err
	}

	err = s.repo.DeleteCode(ctx, request.Code)
	if err != nil {
		s.logger.Error("error deleting code", "error", err)
		return "", "", apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
	}

	return jwtToken, refreshToken, nil
}

func (s *UsersService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	refreshTokenDB, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.Error("error getting refresh token", "error", err)
		return "", apperrors.NewAppError(apperrors.ErrForbidden, "invalid_refresh_token")
	}

	if refreshTokenDB.ExpiresAt.Before(time.Now()) {
		s.logger.Error("expired refresh token")
		return "", apperrors.NewAppError(apperrors.ErrForbidden, "expired_refresh_token")
	}

	jwtToken, _, err := s.CreateTokens(ctx, refreshTokenDB.UserID)
	if err != nil {
		s.logger.Error("error creating tokens", "error", err)
		return "", err
	}

	return jwtToken, nil
}
