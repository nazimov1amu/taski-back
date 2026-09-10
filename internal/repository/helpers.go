package repository

import (
	"context"
	"database/sql"
	"taski_backend/internal/apperrors"
	"taski_backend/internal/constants"
)

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nullBool(b bool) any {
	if !b {
		return nil
	}
	return b
}

func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(constants.LocalDateTimeLayout)
}

func getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return "", apperrors.NewAppError(apperrors.ErrForbidden, "unauthorized")
	}
	return userID, nil
}