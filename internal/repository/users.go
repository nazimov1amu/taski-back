package repository

import (
	"context"
	"database/sql"
	"log"
	db "taski_backend/internal/db/queries"
	"taski_backend/internal/models"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Get(ctx context.Context, id string) (models.UserResponse, error) {
	row := r.db.QueryRowContext(ctx, db.UserQueries.Get, id)
	return scanUserResponse(row)
}

func (r *UsersRepository) Create(ctx context.Context, user models.CreateUserRequest) (models.UserResponse, error) {
	row := r.db.QueryRowContext(ctx, db.UserQueries.Create, user.Email, user.Username, user.Password)
	return scanUserResponse(row)
}

func (r *UsersRepository) Update(ctx context.Context, user models.UpdateUserRequest) (models.UserResponse, error) {
	row := r.db.QueryRowContext(ctx, db.UserQueries.Update, user.ID, user.Email, user.Username, user.Password)
	return scanUserResponse(row)
}

func (r *UsersRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, db.UserQueries.Delete, id)
	if err != nil {
		log.Println("error deleting user:", err)
		return err
	}
	return nil
}

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (models.UserDB, error) {
	row := r.db.QueryRowContext(ctx, db.UserQueries.GetByEmail, email)
	return scanUserDB(row)
}

func (r *UsersRepository) CreateRefreshToken(ctx context.Context, userID string, token string) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token)
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(ctx, query, userID, token)
	if err != nil {
		log.Println("error creating refresh token:", err)
		return err
	}

	return nil
}

func (r *UsersRepository) GetRefreshToken(ctx context.Context, token string) (models.RefreshTokenDB, error) {
	query := `
		SELECT user_id, expires_at FROM refresh_tokens
		WHERE token = $1
	`
	var refreshToken models.RefreshTokenDB
	err := r.db.QueryRowContext(ctx, query, token).Scan(&refreshToken.UserID, &refreshToken.ExpiresAt)
	if err != nil {
		log.Println("error getting refresh token:", err)
		return models.RefreshTokenDB{}, err
	}
	return refreshToken, nil
}

func scanUserResponse(row *sql.Row) (models.UserResponse, error) {
	var user models.UserResponse
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		log.Println("error scanning user response:", err)
		return models.UserResponse{}, err
	}
	return user, nil
}

func scanUserDB(row *sql.Row) (models.UserDB, error) {
	var user models.UserDB
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		log.Println("error scanning user db:", err)
		return models.UserDB{}, err
	}
	return user, nil
}