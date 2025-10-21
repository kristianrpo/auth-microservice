package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
)

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name           string
		httpErr        *httperrors.HTTPError
		wantStatusCode int
		wantError      string
		wantCode       string
	}{
		{
			name:           "respond with bad request",
			httpErr:        httperrors.ErrBadRequest,
			wantStatusCode: http.StatusBadRequest,
			wantError:      "Bad request",
			wantCode:       "BAD_REQUEST",
		},
		{
			name:           "respond with unauthorized",
			httpErr:        httperrors.ErrUnauthorized,
			wantStatusCode: http.StatusUnauthorized,
			wantError:      "Unauthorized",
			wantCode:       "UNAUTHORIZED",
		},
		{
			name:           "respond with not found",
			httpErr:        httperrors.ErrNotFound,
			wantStatusCode: http.StatusNotFound,
			wantError:      "Resource not found",
			wantCode:       "NOT_FOUND",
		},
		{
			name:           "respond with internal server error",
			httpErr:        httperrors.ErrInternalServer,
			wantStatusCode: http.StatusInternalServerError,
			wantError:      "Internal server error",
			wantCode:       "INTERNAL_SERVER_ERROR",
		},
		{
			name:           "respond with custom error",
			httpErr:        httperrors.NewHTTPError(418, "I'm a teapot", "TEAPOT"),
			wantStatusCode: 418,
			wantError:      "I'm a teapot",
			wantCode:       "TEAPOT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a response recorder
			w := httptest.NewRecorder()

			// Call the function
			httperrors.RespondWithError(w, tt.httpErr)

			// Check status code
			if w.Code != tt.wantStatusCode {
				t.Errorf("RespondWithError() status code = %v, want %v", w.Code, tt.wantStatusCode)
			}

			// Check content type
			if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("RespondWithError() Content-Type = %v, want application/json", contentType)
			}

			// Decode response body
			var errResp response.ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			// Check error message
			if errResp.Error != tt.wantError {
				t.Errorf("RespondWithError() error = %v, want %v", errResp.Error, tt.wantError)
			}

			// Check error code
			if errResp.Code != tt.wantCode {
				t.Errorf("RespondWithError() code = %v, want %v", errResp.Code, tt.wantCode)
			}
		})
	}
}

func TestRespondWithErrorMessage(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		message        string
		wantStatusCode int
		wantError      string
	}{
		{
			name:           "respond with custom bad request message",
			statusCode:     http.StatusBadRequest,
			message:        "Invalid input data",
			wantStatusCode: http.StatusBadRequest,
			wantError:      "Invalid input data",
		},
		{
			name:           "respond with custom unauthorized message",
			statusCode:     http.StatusUnauthorized,
			message:        "Session expired",
			wantStatusCode: http.StatusUnauthorized,
			wantError:      "Session expired",
		},
		{
			name:           "respond with custom internal error message",
			statusCode:     http.StatusInternalServerError,
			message:        "Database connection failed",
			wantStatusCode: http.StatusInternalServerError,
			wantError:      "Database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a response recorder
			w := httptest.NewRecorder()

			// Call the function
			httperrors.RespondWithErrorMessage(w, tt.statusCode, tt.message)

			// Check status code
			if w.Code != tt.wantStatusCode {
				t.Errorf("RespondWithErrorMessage() status code = %v, want %v", w.Code, tt.wantStatusCode)
			}

			// Check content type
			if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("RespondWithErrorMessage() Content-Type = %v, want application/json", contentType)
			}

			// Decode response body
			var errResp response.ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			// Check error message
			if errResp.Error != tt.wantError {
				t.Errorf("RespondWithErrorMessage() error = %v, want %v", errResp.Error, tt.wantError)
			}

			// Code should be empty for custom messages
			if errResp.Code != "" {
				t.Errorf("RespondWithErrorMessage() code = %v, want empty string", errResp.Code)
			}
		})
	}
}

func TestRespondWithDomainError(t *testing.T) {
	tests := []struct {
		name           string
		domainErr      error
		wantStatusCode int
		wantError      string
		wantCode       string
	}{
		{
			name:           "respond with user not found error",
			domainErr:      domainerrors.ErrUserNotFound,
			wantStatusCode: http.StatusNotFound,
			wantError:      "User not found",
			wantCode:       "USER_NOT_FOUND",
		},
		{
			name:           "respond with user already exists error",
			domainErr:      domainerrors.ErrUserAlreadyExists,
			wantStatusCode: http.StatusConflict,
			wantError:      "User already exists",
			wantCode:       "USER_ALREADY_EXISTS",
		},
		{
			name:           "respond with invalid credentials error",
			domainErr:      domainerrors.ErrInvalidCredentials,
			wantStatusCode: http.StatusUnauthorized,
			wantError:      "Invalid credentials",
			wantCode:       "INVALID_CREDENTIALS",
		},
		{
			name:           "respond with invalid token error",
			domainErr:      domainerrors.ErrInvalidToken,
			wantStatusCode: http.StatusUnauthorized,
			wantError:      "Invalid or expired token",
			wantCode:       "INVALID_TOKEN",
		},
		{
			name:           "respond with expired token error",
			domainErr:      domainerrors.ErrExpiredToken,
			wantStatusCode: http.StatusUnauthorized,
			wantError:      "Invalid or expired token",
			wantCode:       "INVALID_TOKEN",
		},
		{
			name:           "respond with token revoked error",
			domainErr:      domainerrors.ErrTokenRevoked,
			wantStatusCode: http.StatusUnauthorized,
			wantError:      "Token has been revoked",
			wantCode:       "TOKEN_REVOKED",
		},
		{
			name:           "respond with unauthorized error",
			domainErr:      domainerrors.ErrUnauthorized,
			wantStatusCode: http.StatusUnauthorized,
			wantError:      "Unauthorized",
			wantCode:       "UNAUTHORIZED",
		},
		{
			name:           "respond with forbidden error",
			domainErr:      domainerrors.ErrForbidden,
			wantStatusCode: http.StatusForbidden,
			wantError:      "Forbidden",
			wantCode:       "FORBIDDEN",
		},
		{
			name:           "respond with bad request error",
			domainErr:      domainerrors.ErrBadRequest,
			wantStatusCode: http.StatusBadRequest,
			wantError:      "Bad request",
			wantCode:       "BAD_REQUEST",
		},
		{
			name:           "respond with unknown error as internal server error",
			domainErr:      errors.New("unknown database error"),
			wantStatusCode: http.StatusInternalServerError,
			wantError:      "Internal server error",
			wantCode:       "INTERNAL_SERVER_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a response recorder
			w := httptest.NewRecorder()

			// Call the function
			httperrors.RespondWithDomainError(w, tt.domainErr)

			// Check status code
			if w.Code != tt.wantStatusCode {
				t.Errorf("RespondWithDomainError() status code = %v, want %v", w.Code, tt.wantStatusCode)
			}

			// Check content type
			if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("RespondWithDomainError() Content-Type = %v, want application/json", contentType)
			}

			// Decode response body
			var errResp response.ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			// Check error message
			if errResp.Error != tt.wantError {
				t.Errorf("RespondWithDomainError() error = %v, want %v", errResp.Error, tt.wantError)
			}

			// Check error code
			if errResp.Code != tt.wantCode {
				t.Errorf("RespondWithDomainError() code = %v, want %v", errResp.Code, tt.wantCode)
			}
		})
	}
}
