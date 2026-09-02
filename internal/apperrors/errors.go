package apperrors

import (
	"errors"
)

var (
    ErrNotFound     = errors.New("not_found")
    ErrConflict     = errors.New("conflict")
    ErrInvalidInput = errors.New("invalid_input")
)

type AppError struct {
    Kind error
    Message string
}

func (e *AppError) Error() string {
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Kind
}

func NewAppError(kind error, message string) *AppError {
    return &AppError{Kind: kind, Message: message}
}

