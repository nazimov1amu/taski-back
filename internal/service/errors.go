package service

import (
	"database/sql"
	"errors"
	"fmt"

	"taski_backend/internal/apperrors"

	"github.com/jackc/pgx/v5/pgconn"
)

func mapRepoError(err error, entity string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return apperrors.NewAppError(apperrors.ErrNotFound, fmt.Sprintf("%s_not_found", entity))
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return apperrors.NewAppError(apperrors.ErrInvalidInput, fmt.Sprintf("%s_invalid_input", entity))
		case "23505":
			return apperrors.NewAppError(apperrors.ErrConflict, fmt.Sprintf("%s_conflict", entity))
		}
	}

	return apperrors.NewAppError(err, fmt.Sprintf("%s_error", entity))
}
