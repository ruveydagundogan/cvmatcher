package errors

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrValidation    = errors.New("validation failed")
	ErrRateLimited   = errors.New("rate limited")
	ErrNotImplemented = errors.New("not implemented")
	ErrInternal      = errors.New("internal error")
)

type DomainError struct {
	Kind    error
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if os.Getenv("GO_ENV") == "production" && e.Kind == ErrInternal {
		return "an internal error occurred"
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func New(kind error, message string, err error) *DomainError {
	return &DomainError{
		Kind:    kind,
		Message: message,
		Err:     err,
	}
}

func NotFound(message string) *DomainError {
	return &DomainError{Kind: ErrNotFound, Message: message}
}

func AlreadyExists(message string) *DomainError {
	return &DomainError{Kind: ErrAlreadyExists, Message: message}
}

func Unauthorized(message string) *DomainError {
	return &DomainError{Kind: ErrUnauthorized, Message: message}
}

func Forbidden(message string) *DomainError {
	return &DomainError{Kind: ErrForbidden, Message: message}
}

func Validation(message string) *DomainError {
	return &DomainError{Kind: ErrValidation, Message: message}
}

func Internal(message string, err error) *DomainError {
	return &DomainError{Kind: ErrInternal, Message: message, Err: err}
}

func HTTPStatusCode(err error) int {
	var domainErr *DomainError
	if !errors.As(err, &domainErr) {
		return http.StatusInternalServerError
	}

	switch domainErr.Kind {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists:
		return http.StatusConflict
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrValidation:
		return http.StatusUnprocessableEntity
	case ErrRateLimited:
		return http.StatusTooManyRequests
	case ErrNotImplemented:
		return http.StatusNotImplemented
	default:
		return http.StatusInternalServerError
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
