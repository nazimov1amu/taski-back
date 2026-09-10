package service

import (
	"context"
	"taski_backend/internal/apperrors"
)	

func getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", apperrors.NewAppError(apperrors.ErrInternal, "internal_server_error")
	}
	return userID, nil
}