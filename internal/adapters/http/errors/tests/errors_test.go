package tests

import (
	"errors"
	"testing"

	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
)

func TestNewHTTPError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		code       string
		want       *httperrors.HTTPError
	}{
		{
			name:       "create bad request error",
			statusCode: 400,
			message:    "Bad request",
			code:       "BAD_REQUEST",
			want: &httperrors.HTTPError{
				StatusCode: 400,
				Message:    "Bad request",
				Code:       "BAD_REQUEST",
			},
		},
		{
			name:       "create unauthorized error",
			statusCode: 401,
			message:    "Unauthorized",
			code:       "UNAUTHORIZED",
			want: &httperrors.HTTPError{
				StatusCode: 401,
				Message:    "Unauthorized",
				Code:       "UNAUTHORIZED",
			},
		},
		{
			name:       "create internal server error",
			statusCode: 500,
			message:    "Internal server error",
			code:       "INTERNAL_SERVER_ERROR",
			want: &httperrors.HTTPError{
				StatusCode: 500,
				Message:    "Internal server error",
				Code:       "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name:       "create custom error",
			statusCode: 418,
			message:    "I'm a teapot",
			code:       "TEAPOT",
			want: &httperrors.HTTPError{
				StatusCode: 418,
				Message:    "I'm a teapot",
				Code:       "TEAPOT",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := httperrors.NewHTTPError(tt.statusCode, tt.message, tt.code)

			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("NewHTTPError() StatusCode = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
			if got.Message != tt.want.Message {
				t.Errorf("NewHTTPError() Message = %v, want %v", got.Message, tt.want.Message)
			}
			if got.Code != tt.want.Code {
				t.Errorf("NewHTTPError() Code = %v, want %v", got.Code, tt.want.Code)
			}
		})
	}
}

func TestHTTPError_Error(t *testing.T) {
	tests := []struct {
		name     string
		httpErr  *httperrors.HTTPError
		wantMsg  string
	}{
		{
			name: "error method returns message",
			httpErr: &httperrors.HTTPError{
				StatusCode: 400,
				Message:    "Bad request",
				Code:       "BAD_REQUEST",
			},
			wantMsg: "Bad request",
		},
		{
			name: "error method returns custom message",
			httpErr: &httperrors.HTTPError{
				StatusCode: 404,
				Message:    "User not found",
				Code:       "USER_NOT_FOUND",
			},
			wantMsg: "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.httpErr.Error(); got != tt.wantMsg {
				t.Errorf("HTTPError.Error() = %v, want %v", got, tt.wantMsg)
			}
		})
	}
}

func TestPredefinedHTTPErrors(t *testing.T) {
	tests := []struct {
		name       string
		httpErr    *httperrors.HTTPError
		statusCode int
		message    string
		code       string
	}{
		{
			name:       "ErrBadRequest",
			httpErr:    httperrors.ErrBadRequest,
			statusCode: 400,
			message:    "Bad request",
			code:       "BAD_REQUEST",
		},
		{
			name:       "ErrUnauthorized",
			httpErr:    httperrors.ErrUnauthorized,
			statusCode: 401,
			message:    "Unauthorized",
			code:       "UNAUTHORIZED",
		},
		{
			name:       "ErrForbidden",
			httpErr:    httperrors.ErrForbidden,
			statusCode: 403,
			message:    "Forbidden",
			code:       "FORBIDDEN",
		},
		{
			name:       "ErrNotFound",
			httpErr:    httperrors.ErrNotFound,
			statusCode: 404,
			message:    "Resource not found",
			code:       "NOT_FOUND",
		},
		{
			name:       "ErrConflict",
			httpErr:    httperrors.ErrConflict,
			statusCode: 409,
			message:    "Resource conflict",
			code:       "CONFLICT",
		},
		{
			name:       "ErrInternalServer",
			httpErr:    httperrors.ErrInternalServer,
			statusCode: 500,
			message:    "Internal server error",
			code:       "INTERNAL_SERVER_ERROR",
		},
		{
			name:       "ErrInvalidCredentials",
			httpErr:    httperrors.ErrInvalidCredentials,
			statusCode: 401,
			message:    "Invalid credentials",
			code:       "INVALID_CREDENTIALS",
		},
		{
			name:       "ErrInvalidToken",
			httpErr:    httperrors.ErrInvalidToken,
			statusCode: 401,
			message:    "Invalid or expired token",
			code:       "INVALID_TOKEN",
		},
		{
			name:       "ErrTokenRevoked",
			httpErr:    httperrors.ErrTokenRevoked,
			statusCode: 401,
			message:    "Token has been revoked",
			code:       "TOKEN_REVOKED",
		},
		{
			name:       "ErrUserAlreadyExists",
			httpErr:    httperrors.ErrUserAlreadyExists,
			statusCode: 409,
			message:    "User already exists",
			code:       "USER_ALREADY_EXISTS",
		},
		{
			name:       "ErrUserNotFound",
			httpErr:    httperrors.ErrUserNotFound,
			statusCode: 404,
			message:    "User not found",
			code:       "USER_NOT_FOUND",
		},
		{
			name:       "ErrMissingAuthHeader",
			httpErr:    httperrors.ErrMissingAuthHeader,
			statusCode: 401,
			message:    "Missing authorization header",
			code:       "MISSING_AUTH_HEADER",
		},
		{
			name:       "ErrInvalidAuthHeader",
			httpErr:    httperrors.ErrInvalidAuthHeader,
			statusCode: 401,
			message:    "Invalid authorization header format",
			code:       "INVALID_AUTH_HEADER",
		},
		{
			name:       "ErrRequiredField",
			httpErr:    httperrors.ErrRequiredField,
			statusCode: 400,
			message:    "Required field is missing",
			code:       "REQUIRED_FIELD",
		},
		{
			name:       "ErrInvalidRequestBody",
			httpErr:    httperrors.ErrInvalidRequestBody,
			statusCode: 400,
			message:    "Invalid request body",
			code:       "INVALID_REQUEST_BODY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.httpErr.StatusCode != tt.statusCode {
				t.Errorf("%s StatusCode = %v, want %v", tt.name, tt.httpErr.StatusCode, tt.statusCode)
			}
			if tt.httpErr.Message != tt.message {
				t.Errorf("%s Message = %v, want %v", tt.name, tt.httpErr.Message, tt.message)
			}
			if tt.httpErr.Code != tt.code {
				t.Errorf("%s Code = %v, want %v", tt.name, tt.httpErr.Code, tt.code)
			}
		})
	}
}

func TestMapDomainError(t *testing.T) {
	tests := []struct {
		name        string
		domainErr   error
		wantHTTPErr *httperrors.HTTPError
	}{
		{
			name:        "nil error returns nil",
			domainErr:   nil,
			wantHTTPErr: nil,
		},
		{
			name:        "ErrUserNotFound maps to ErrUserNotFound",
			domainErr:   domainerrors.ErrUserNotFound,
			wantHTTPErr: httperrors.ErrUserNotFound,
		},
		{
			name:        "ErrUserAlreadyExists maps to ErrUserAlreadyExists",
			domainErr:   domainerrors.ErrUserAlreadyExists,
			wantHTTPErr: httperrors.ErrUserAlreadyExists,
		},
		{
			name:        "ErrInvalidCredentials maps to ErrInvalidCredentials",
			domainErr:   domainerrors.ErrInvalidCredentials,
			wantHTTPErr: httperrors.ErrInvalidCredentials,
		},
		{
			name:        "ErrInvalidToken maps to ErrInvalidToken",
			domainErr:   domainerrors.ErrInvalidToken,
			wantHTTPErr: httperrors.ErrInvalidToken,
		},
		{
			name:        "ErrExpiredToken maps to ErrInvalidToken",
			domainErr:   domainerrors.ErrExpiredToken,
			wantHTTPErr: httperrors.ErrInvalidToken,
		},
		{
			name:        "ErrTokenRevoked maps to ErrTokenRevoked",
			domainErr:   domainerrors.ErrTokenRevoked,
			wantHTTPErr: httperrors.ErrTokenRevoked,
		},
		{
			name:        "ErrInvalidTokenType maps to ErrInvalidToken",
			domainErr:   domainerrors.ErrInvalidTokenType,
			wantHTTPErr: httperrors.ErrInvalidToken,
		},
		{
			name:        "ErrUnauthorized maps to ErrUnauthorized",
			domainErr:   domainerrors.ErrUnauthorized,
			wantHTTPErr: httperrors.ErrUnauthorized,
		},
		{
			name:        "ErrForbidden maps to ErrForbidden",
			domainErr:   domainerrors.ErrForbidden,
			wantHTTPErr: httperrors.ErrForbidden,
		},
		{
			name:        "ErrBadRequest maps to ErrBadRequest",
			domainErr:   domainerrors.ErrBadRequest,
			wantHTTPErr: httperrors.ErrBadRequest,
		},
		{
			name:        "unknown error maps to ErrInternalServer",
			domainErr:   errors.New("some unknown error"),
			wantHTTPErr: httperrors.ErrInternalServer,
		},
		{
			name:        "wrapped ErrUserNotFound maps correctly",
			domainErr:   errors.Join(domainerrors.ErrUserNotFound, errors.New("additional context")),
			wantHTTPErr: httperrors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := httperrors.MapDomainError(tt.domainErr)

			if tt.wantHTTPErr == nil {
				if got != nil {
					t.Errorf("MapDomainError() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Errorf("MapDomainError() = nil, want %v", tt.wantHTTPErr)
				return
			}

			if got.StatusCode != tt.wantHTTPErr.StatusCode {
				t.Errorf("MapDomainError() StatusCode = %v, want %v", got.StatusCode, tt.wantHTTPErr.StatusCode)
			}
			if got.Message != tt.wantHTTPErr.Message {
				t.Errorf("MapDomainError() Message = %v, want %v", got.Message, tt.wantHTTPErr.Message)
			}
			if got.Code != tt.wantHTTPErr.Code {
				t.Errorf("MapDomainError() Code = %v, want %v", got.Code, tt.wantHTTPErr.Code)
			}
		})
	}
}

