package tests

import (
	"errors"
	"testing"

	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
)

func TestDomainErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrUserNotFound",
			err:      domainerrors.ErrUserNotFound,
			expected: "user not found",
		},
		{
			name:     "ErrUserAlreadyExists",
			err:      domainerrors.ErrUserAlreadyExists,
			expected: "user already exists",
		},
		{
			name:     "ErrInvalidCredentials",
			err:      domainerrors.ErrInvalidCredentials,
			expected: "invalid credentials",
		},
		{
			name:     "ErrInvalidToken",
			err:      domainerrors.ErrInvalidToken,
			expected: "invalid token",
		},
		{
			name:     "ErrExpiredToken",
			err:      domainerrors.ErrExpiredToken,
			expected: "token has expired",
		},
		{
			name:     "ErrTokenRevoked",
			err:      domainerrors.ErrTokenRevoked,
			expected: "token has been revoked",
		},
		{
			name:     "ErrForbidden",
			err:      domainerrors.ErrForbidden,
			expected: "forbidden",
		},
		{
			name:     "ErrInternal",
			err:      domainerrors.ErrInternal,
			expected: "internal server error",
		},
		{
			name:     "ErrClientNotFound",
			err:      domainerrors.ErrClientNotFound,
			expected: "oauth client not found",
		},
		{
			name:     "ErrInvalidClient",
			err:      domainerrors.ErrInvalidClient,
			expected: "invalid oauth client",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Error message = %v, want %v", tt.err.Error(), tt.expected)
			}
		})
	}
}

func TestErrorComparison(t *testing.T) {
	tests := []struct {
		name     string
		err1     error
		err2     error
		expected bool
	}{
		{
			name:     "same error",
			err1:     domainerrors.ErrUserNotFound,
			err2:     domainerrors.ErrUserNotFound,
			expected: true,
		},
		{
			name:     "different errors",
			err1:     domainerrors.ErrUserNotFound,
			err2:     domainerrors.ErrUserAlreadyExists,
			expected: false,
		},
		{
			name:     "error vs nil",
			err1:     domainerrors.ErrInvalidToken,
			err2:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.Is(tt.err1, tt.err2)
			if result != tt.expected {
				t.Errorf("errors.Is() = %v, want %v", result, tt.expected)
			}
		})
	}
}

