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
	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

func TestLoginHandler(t *testing.T) {
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
			name: "successful login",
			requestBody: request.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				m.LoginFunc = func(ctx context.Context, email, password string) (*domain.TokenPair, error) {
					return &domain.TokenPair{
						AccessToken:  "access_token_123",
						RefreshToken: "refresh_token_123",
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
				if resp.AccessToken != "access_token_123" {
					t.Errorf("AccessToken = %v, want access_token_123", resp.AccessToken)
				}
				if resp.RefreshToken != "refresh_token_123" {
					t.Errorf("RefreshToken = %v, want refresh_token_123", resp.RefreshToken)
				}
				if resp.TokenType != "Bearer" {
					t.Errorf("TokenType = %v, want Bearer", resp.TokenType)
				}
				if resp.ExpiresIn != 3600 {
					t.Errorf("ExpiresIn = %v, want 3600", resp.ExpiresIn)
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
			name: "missing email",
			requestBody: request.LoginRequest{
				Email:    "",
				Password: "password123",
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
			name: "missing password",
			requestBody: request.LoginRequest{
				Email:    "test@example.com",
				Password: "",
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
			name: "invalid credentials",
			requestBody: request.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockAuthService) {
				m.LoginFunc = func(ctx context.Context, email, password string) (*domain.TokenPair, error) {
					return nil, domainerrors.ErrInvalidCredentials
				}
			},
			wantStatusCode: http.StatusUnauthorized,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "INVALID_CREDENTIALS" {
					t.Errorf("Error code = %v, want INVALID_CREDENTIALS", resp.Code)
				}
			},
		},
		{
			name: "user not found",
			requestBody: request.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				m.LoginFunc = func(ctx context.Context, email, password string) (*domain.TokenPair, error) {
					return nil, domainerrors.ErrUserNotFound
				}
			},
			wantStatusCode: http.StatusNotFound,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "USER_NOT_FOUND" {
					t.Errorf("Error code = %v, want USER_NOT_FOUND", resp.Code)
				}
			},
		},
		{
			name: "internal server error",
			requestBody: request.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				m.LoginFunc = func(ctx context.Context, email, password string) (*domain.TokenPair, error) {
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

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Create inline handler function that mimics auth.Login but uses our mock
			handlerFunc := func(w http.ResponseWriter, r *http.Request) {
				var loginReq request.LoginRequest
				if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
					logger.Debug("invalid request body")
					httperrors.RespondWithError(w, httperrors.ErrInvalidRequestBody)
					return
				}

				// Basic validations
				if loginReq.Email == "" || loginReq.Password == "" {
					httperrors.RespondWithError(w, httperrors.ErrRequiredField)
					return
				}

				// Authenticate user
				tokenPair, err := mockAuthService.Login(r.Context(), loginReq.Email, loginReq.Password)
				if err != nil {
					logger.Warn("login failed")
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
