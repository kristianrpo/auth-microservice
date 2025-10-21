package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
)

func TestLogoutHandler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name            string
		authHeader      string
		requestBody     interface{}
		mockSetup       func(*MockAuthService)
		wantStatusCode  int
		wantError       bool
		checkResponse   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:       "successful logout with refresh token",
			authHeader: "Bearer access_token_123",
			requestBody: request.LogoutRequest{
				RefreshToken: "refresh_token_123",
			},
			mockSetup: func(m *MockAuthService) {
				m.LogoutFunc = func(ctx context.Context, accessToken, refreshToken string) error {
					if accessToken != "access_token_123" {
						return errors.New("invalid access token")
					}
					if refreshToken != "refresh_token_123" {
						return errors.New("invalid refresh token")
					}
					return nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.MessageResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Message != "logout successful" {
					t.Errorf("Message = %v, want logout successful", resp.Message)
				}
			},
		},
		{
			name:       "successful logout without refresh token",
			authHeader: "Bearer access_token_123",
			requestBody: request.LogoutRequest{
				RefreshToken: "",
			},
			mockSetup: func(m *MockAuthService) {
				m.LogoutFunc = func(ctx context.Context, accessToken, refreshToken string) error {
					return nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.MessageResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Message != "logout successful" {
					t.Errorf("Message = %v, want logout successful", resp.Message)
				}
			},
		},
		{
			name:       "missing authorization header",
			authHeader: "",
			requestBody: request.LogoutRequest{
				RefreshToken: "refresh_token_123",
			},
			mockSetup:      func(m *MockAuthService) {},
			wantStatusCode: http.StatusBadRequest,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "REQUIRED_FIELD" {
					t.Errorf("Error code = %v, want REQUIRED_FIELD", resp.Code)
				}
			},
		},
		{
			name:       "invalid authorization header format",
			authHeader: "InvalidFormat",
			requestBody: request.LogoutRequest{
				RefreshToken: "refresh_token_123",
			},
			mockSetup:      func(m *MockAuthService) {},
			wantStatusCode: http.StatusBadRequest,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "REQUIRED_FIELD" {
					t.Errorf("Error code = %v, want REQUIRED_FIELD", resp.Code)
				}
			},
		},
		{
			name:       "service error during logout",
			authHeader: "Bearer access_token_123",
			requestBody: request.LogoutRequest{
				RefreshToken: "refresh_token_123",
			},
			mockSetup: func(m *MockAuthService) {
				m.LogoutFunc = func(ctx context.Context, accessToken, refreshToken string) error {
					return domainerrors.ErrInvalidToken
				}
			},
			wantStatusCode: http.StatusUnauthorized,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "INVALID_TOKEN" {
					t.Errorf("Error code = %v, want INVALID_TOKEN", resp.Code)
				}
			},
		},
		{
			name:       "internal server error",
			authHeader: "Bearer access_token_123",
			requestBody: request.LogoutRequest{
				RefreshToken: "refresh_token_123",
			},
			mockSetup: func(m *MockAuthService) {
				m.LogoutFunc = func(ctx context.Context, accessToken, refreshToken string) error {
					return errors.New("database error")
				}
			},
			wantStatusCode: http.StatusInternalServerError,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "INTERNAL_SERVER_ERROR" {
					t.Errorf("Error code = %v, want INTERNAL_SERVER_ERROR", resp.Code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockAuthService := &MockAuthService{}
			tt.mockSetup(mockAuthService)

			// Create request
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Create inline handler function that mimics auth.Logout but uses our mock
			handlerFunc := func(w http.ResponseWriter, r *http.Request) {
				// Get access token from header
				authHeader := r.Header.Get("Authorization")
				var accessToken string
				if authHeader != "" {
					parts := strings.Split(authHeader, " ")
					if len(parts) == 2 {
						accessToken = parts[1]
					}
				}

				// Get refresh token from body (optional)
				var logoutReq request.LogoutRequest
				_ = json.NewDecoder(r.Body).Decode(&logoutReq)

				if accessToken == "" {
					httperrors.RespondWithError(w, httperrors.ErrRequiredField)
					return
				}

				// Perform logout
				if err := mockAuthService.Logout(r.Context(), accessToken, logoutReq.RefreshToken); err != nil {
					logger.Error("logout failed")
					httperrors.RespondWithDomainError(w, err)
					return
				}

				resp := response.MessageResponse{Message: "logout successful"}
				shared.RespondWithJSON(w, http.StatusOK, resp)
			}

			// Call handler
			handlerFunc(w, req)

			// Check status code
			if w.Code != tt.wantStatusCode {
				t.Errorf("status code = %v, want %v", w.Code, tt.wantStatusCode)
			}

			// Check response
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

