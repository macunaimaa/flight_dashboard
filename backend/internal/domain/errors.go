package domain

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrForbidden     = errors.New("access forbidden")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrConflict      = errors.New("resource conflict")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInternal      = errors.New("internal error")
)
