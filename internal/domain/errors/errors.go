package domain

import "errors"

// Authentication errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password is too weak")
	ErrClientNotFound     = errors.New("oauth client not found")
	ErrInvalidClient      = errors.New("invalid oauth client")
)

// Token errors
var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrTokenRevoked     = errors.New("token has been revoked")
	ErrInvalidTokenType = errors.New("invalid token type")
)

// Generic errors
var (
	ErrInternal       = errors.New("internal server error")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrBadRequest     = errors.New("bad request")
	ErrValidation     = errors.New("validation error")
	ErrNotImplemented = errors.New("not implemented")
)
