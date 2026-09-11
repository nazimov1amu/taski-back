package repository

import (
	"context"
	"database/sql"
	"log"
	"taski_backend/internal/models"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Get(ctx context.Context, id string) (models.UserResponse, error) {
	query := `
		SELECT id, username FROM users
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanUserResponse(row)
}

func (r *UsersRepository) Create(ctx context.Context, user models.CreateUserRequest) (string, error) {
	query := `
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Username, user.Password).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *UsersRepository) Update(ctx context.Context, user models.UpdateUserRequest) (models.UserResponse, error) {
	query := `
		UPDATE users
		SET email = $2, username = $3, password = $4
		WHERE id = $1
		RETURNING id, username
	`
	row := r.db.QueryRowContext(ctx, query, user.ID, user.Email, user.Username, user.Password)
	return scanUserResponse(row)
}

func (r *UsersRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Println("error deleting user:", err)
		return err
	}
	return nil
}

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (models.UserDB, error) {
	query := `
		SELECT id, username, password FROM users
		WHERE email = $1
	`
	row := r.db.QueryRowContext(ctx, query, email)
	return scanUserDB(row)
}

func (r *UsersRepository) CreateRefreshToken(ctx context.Context, userID string, token string) error {
	query := `
		INSERT INTO user_tokens (user_id, token)
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
		SELECT user_id, expires_at FROM user_tokens
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

func (r *UsersRepository) CreateCode(ctx context.Context, userID string, code string, codeChallenge string) error {
	query := `
		INSERT INTO user_codes (user_id, code, code_challenge)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, userID, code, codeChallenge)
	if err != nil {
		log.Println("error creating code:", err)
		return err
	}
	return nil
}

func (r *UsersRepository) DeleteCode(ctx context.Context, code string) error {
	query := `
		DELETE FROM user_codes
		WHERE code = $1
	`
	_, err := r.db.ExecContext(ctx, query, code)
	if err != nil {
		log.Println("error deleting code:", err)
		return err
	}
	return nil
}

func (r *UsersRepository) GetCode(ctx context.Context, code string) (models.CodeDB, error) {
	query := `
		SELECT user_id, code_challenge, expires_at FROM user_codes
		WHERE code = $1
	`
	var codeDB models.CodeDB
	err := r.db.QueryRowContext(ctx, query, code).Scan(&codeDB.UserID, &codeDB.CodeChallenge, &codeDB.ExpiresAt)
	if err != nil {
		log.Println("error getting code:", err)
		return models.CodeDB{}, err
	}
	return codeDB, nil
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
