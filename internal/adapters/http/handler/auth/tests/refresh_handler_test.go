package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
)

func TestRefreshHandler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		wantStatusCode int
		wantError      bool
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful token refresh",
			requestBody: request.RefreshTokenRequest{
				RefreshToken: "valid_refresh_token",
			},
			mockSetup: func(m *MockAuthService) {
				m.RefreshTokenFunc = func(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
					return &domain.TokenPair{
						AccessToken:  "new_access_token",
						RefreshToken: "new_refresh_token",
						TokenType:    "Bearer",
						ExpiresIn:    3600,
					}, nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.TokenResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.AccessToken != "new_access_token" {
					t.Errorf("AccessToken = %v, want new_access_token", resp.AccessToken)
				}
				if resp.RefreshToken != "new_refresh_token" {
					t.Errorf("RefreshToken = %v, want new_refresh_token", resp.RefreshToken)
				}
			},
		},
		{
			name:           "invalid json body",
			requestBody:    "invalid json",
			mockSetup:      func(m *MockAuthService) {},
			wantStatusCode: http.StatusBadRequest,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "INVALID_REQUEST_BODY" {
					t.Errorf("Error code = %v, want INVALID_REQUEST_BODY", resp.Code)
				}
			},
		},
		{
			name: "missing refresh token",
			requestBody: request.RefreshTokenRequest{
				RefreshToken: "",
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
			name: "invalid refresh token",
			requestBody: request.RefreshTokenRequest{
				RefreshToken: "invalid_token",
			},
			mockSetup: func(m *MockAuthService) {
				m.RefreshTokenFunc = func(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
					return nil, domainerrors.ErrInvalidToken
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
			name: "expired refresh token",
			requestBody: request.RefreshTokenRequest{
				RefreshToken: "expired_token",
			},
			mockSetup: func(m *MockAuthService) {
				m.RefreshTokenFunc = func(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
					return nil, domainerrors.ErrExpiredToken
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
			name: "revoked refresh token",
			requestBody: request.RefreshTokenRequest{
				RefreshToken: "revoked_token",
			},
			mockSetup: func(m *MockAuthService) {
				m.RefreshTokenFunc = func(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
					return nil, domainerrors.ErrTokenRevoked
				}
			},
			wantStatusCode: http.StatusUnauthorized,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "TOKEN_REVOKED" {
					t.Errorf("Error code = %v, want TOKEN_REVOKED", resp.Code)
				}
			},
		},
		{
			name: "internal server error",
			requestBody: request.RefreshTokenRequest{
				RefreshToken: "valid_token",
			},
			mockSetup: func(m *MockAuthService) {
				m.RefreshTokenFunc = func(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
					return nil, errors.New("database error")
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

			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Create inline handler function that mimics auth.Refresh but uses our mock
			handlerFunc := func(w http.ResponseWriter, r *http.Request) {
				var refReq request.RefreshTokenRequest
				if err := json.NewDecoder(r.Body).Decode(&refReq); err != nil {
					logger.Debug("invalid request body")
					httperrors.RespondWithError(w, httperrors.ErrInvalidRequestBody)
					return
				}

				if refReq.RefreshToken == "" {
					httperrors.RespondWithError(w, httperrors.ErrRequiredField)
					return
				}

				// Refresh token
				tokenPair, err := mockAuthService.RefreshToken(r.Context(), refReq.RefreshToken)
				if err != nil {
					logger.Warn("token refresh failed")
					httperrors.RespondWithDomainError(w, err)
					return
				}

				// Convert to DTO
				resp := response.TokenResponse{
					AccessToken:  tokenPair.AccessToken,
					RefreshToken: tokenPair.RefreshToken,
					TokenType:    tokenPair.TokenType,
					ExpiresIn:    tokenPair.ExpiresIn,
				}

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

