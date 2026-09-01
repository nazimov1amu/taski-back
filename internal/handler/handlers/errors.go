package handlers

import (
	"errors"
	"net/http"

	"taski_backend/internal/apperrors"
)

func writeServiceError(w http.ResponseWriter, err error, entity string) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		http.Error(w, entity+" not found", http.StatusNotFound)
	case errors.Is(err, apperrors.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrConflict):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
