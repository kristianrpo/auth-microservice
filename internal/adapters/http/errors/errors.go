package errors

import (
	"errors"
	nethttp "net/http"

	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
)

// HTTPError represents an HTTP error with status code and message
type HTTPError struct {
	StatusCode int
	Message    string
	Code       string
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	return e.Message
}

// NewHTTPError creates a new HTTPError
func NewHTTPError(statusCode int, message, code string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
		Code:       code,
	}
}

// Predefined HTTP errors
var (
	ErrBadRequest         = NewHTTPError(nethttp.StatusBadRequest, "Bad request", "BAD_REQUEST")
	ErrUnauthorized       = NewHTTPError(nethttp.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
	ErrForbidden          = NewHTTPError(nethttp.StatusForbidden, "Forbidden", "FORBIDDEN")
	ErrNotFound           = NewHTTPError(nethttp.StatusNotFound, "Resource not found", "NOT_FOUND")
	ErrConflict           = NewHTTPError(nethttp.StatusConflict, "Resource conflict", "CONFLICT")
	ErrInternalServer     = NewHTTPError(nethttp.StatusInternalServerError, "Internal server error", "INTERNAL_SERVER_ERROR")
	ErrInvalidCredentials = NewHTTPError(nethttp.StatusUnauthorized, "Invalid credentials", "INVALID_CREDENTIALS")
	ErrInvalidToken       = NewHTTPError(nethttp.StatusUnauthorized, "Invalid or expired token", "INVALID_TOKEN")
	ErrTokenRevoked       = NewHTTPError(nethttp.StatusUnauthorized, "Token has been revoked", "TOKEN_REVOKED")
	ErrUserAlreadyExists  = NewHTTPError(nethttp.StatusConflict, "User already exists", "USER_ALREADY_EXISTS")
	ErrUserNotFound       = NewHTTPError(nethttp.StatusNotFound, "User not found", "USER_NOT_FOUND")
	ErrMissingAuthHeader  = NewHTTPError(nethttp.StatusUnauthorized, "Missing authorization header", "MISSING_AUTH_HEADER")
	ErrInvalidAuthHeader  = NewHTTPError(nethttp.StatusUnauthorized, "Invalid authorization header format", "INVALID_AUTH_HEADER")
	ErrRequiredField      = NewHTTPError(nethttp.StatusBadRequest, "Required field is missing", "REQUIRED_FIELD")
	ErrInvalidRequestBody = NewHTTPError(nethttp.StatusBadRequest, "Invalid request body", "INVALID_REQUEST_BODY")
)

// MapDomainError maps domain errors to HTTP errors
func MapDomainError(err error) *HTTPError {
	if err == nil {
		return nil
	}

	// Mapping of domain errors to HTTP errors
	switch {
	case errors.Is(err, domainerrors.ErrUserNotFound):
		return ErrUserNotFound
	case errors.Is(err, domainerrors.ErrUserAlreadyExists):
		return ErrUserAlreadyExists
	case errors.Is(err, domainerrors.ErrInvalidCredentials):
		return ErrInvalidCredentials
	case errors.Is(err, domainerrors.ErrInvalidToken):
		return ErrInvalidToken
	case errors.Is(err, domainerrors.ErrExpiredToken):
		return ErrInvalidToken
	case errors.Is(err, domainerrors.ErrTokenRevoked):
		return ErrTokenRevoked
	case errors.Is(err, domainerrors.ErrInvalidTokenType):
		return ErrInvalidToken
	case errors.Is(err, domainerrors.ErrUnauthorized):
		return ErrUnauthorized
	case errors.Is(err, domainerrors.ErrForbidden):
		return ErrForbidden
	case errors.Is(err, domainerrors.ErrBadRequest):
		return ErrBadRequest
	default:
		// Error genérico
		return ErrInternalServer
	}
}
