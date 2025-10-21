package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	authhandler "github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/auth"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

func TestRegisterHandler(t *testing.T) {
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
			name: "successful registration",
			requestBody: request.RegisterRequest{
				IDCitizen: 12345,
				Email:     "newuser@example.com",
				Password:  "password123",
				Name:      "New User",
			},
			mockSetup: func(m *MockAuthService) {
				m.RegisterFunc = func(ctx context.Context, email, password, name string, idCitizen int) (*domain.UserPublic, error) {
					return &domain.UserPublic{
						ID:        "user-123",
						IDCitizen: 12345,
						Email:     "newuser@example.com",
						Name:      "New User",
						Role:      domain.RoleUser,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil
				}
			},
			wantStatusCode: http.StatusCreated,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.UserResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.ID != "user-123" {
					t.Errorf("ID = %v, want user-123", resp.ID)
				}
				if resp.Email != "newuser@example.com" {
					t.Errorf("Email = %v, want newuser@example.com", resp.Email)
				}
				if resp.Name != "New User" {
					t.Errorf("Name = %v, want New User", resp.Name)
				}
				if resp.IDCitizen != 12345 {
					t.Errorf("IDCitizen = %v, want 12345", resp.IDCitizen)
				}
				if resp.Role != domain.RoleUser {
					t.Errorf("Role = %v, want USER", resp.Role)
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
			requestBody: request.RegisterRequest{
				IDCitizen: 12345,
				Email:     "",
				Password:  "password123",
				Name:      "Test User",
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
			requestBody: request.RegisterRequest{
				IDCitizen: 12345,
				Email:     "test@example.com",
				Password:  "",
				Name:      "Test User",
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
			name: "missing name",
			requestBody: request.RegisterRequest{
				IDCitizen: 12345,
				Email:     "test@example.com",
				Password:  "password123",
				Name:      "",
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
			name: "invalid id_citizen",
			requestBody: request.RegisterRequest{
				IDCitizen: 0,
				Email:     "test@example.com",
				Password:  "password123",
				Name:      "Test User",
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
			name: "user already exists",
			requestBody: request.RegisterRequest{
				IDCitizen: 12345,
				Email:     "existing@example.com",
				Password:  "password123",
				Name:      "Existing User",
			},
			mockSetup: func(m *MockAuthService) {
				m.RegisterFunc = func(ctx context.Context, email, password, name string, idCitizen int) (*domain.UserPublic, error) {
					return nil, domainerrors.ErrUserAlreadyExists
				}
			},
			wantStatusCode: http.StatusConflict,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "USER_ALREADY_EXISTS" {
					t.Errorf("Error code = %v, want USER_ALREADY_EXISTS", resp.Code)
				}
			},
		},
		{
			name: "internal server error",
			requestBody: request.RegisterRequest{
				IDCitizen: 12345,
				Email:     "test@example.com",
				Password:  "password123",
				Name:      "Test User",
			},
			mockSetup: func(m *MockAuthService) {
				m.RegisterFunc = func(ctx context.Context, email, password, name string, idCitizen int) (*domain.UserPublic, error) {
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

			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Use real handler with mock service injected
			h := shared.NewAuthHandler(mockAuthService, logger)
			handler := authhandler.Register(h)
			handler(w, req)

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
