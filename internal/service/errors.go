package service

import (
	"database/sql"
	"errors"

	"taski_backend/internal/apperrors"

	"github.com/jackc/pgx/v5/pgconn"
)

func mapRepoError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return apperrors.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return apperrors.ErrInvalidInput
		case "23505":
			return apperrors.ErrConflict
		}
	}

	return err
}
