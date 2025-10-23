package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	authhandler "github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/auth"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/middleware"
	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

func TestGetMeHandler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name           string
		setupContext   func(context.Context) context.Context
		mockSetup      func(*MockAuthService)
		wantStatusCode int
		wantError      bool
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful get user info",
			setupContext: func(ctx context.Context) context.Context {
				claims := &domain.TokenClaims{
					IDCitizen: 12345,
					Email:     "test@example.com",
					Role:      domain.RoleUser,
				}
				return context.WithValue(ctx, middleware.UserContextKey, claims)
			},
			mockSetup: func(m *MockAuthService) {
				m.GetUserByIDCitizenFunc = func(ctx context.Context, idCitizen int) (*domain.UserPublic, error) {
					return &domain.UserPublic{
						ID:        "user-123",
						IDCitizen: 12345,
						Email:     "test@example.com",
						Name:      "Test User",
						Role:      domain.RoleUser,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.UserResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.ID != "user-123" {
					t.Errorf("ID = %v, want user-123", resp.ID)
				}
				if resp.Email != "test@example.com" {
					t.Errorf("Email = %v, want test@example.com", resp.Email)
				}
				if resp.Name != "Test User" {
					t.Errorf("Name = %v, want Test User", resp.Name)
				}
				if resp.Role != domain.RoleUser {
					t.Errorf("Role = %v, want USER", resp.Role)
				}
			},
		},
		{
			name: "missing user context",
			setupContext: func(ctx context.Context) context.Context {
				// Return context without user claims
				return ctx
			},
			mockSetup:      func(m *MockAuthService) {},
			wantStatusCode: http.StatusUnauthorized,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "UNAUTHORIZED" {
					t.Errorf("Error code = %v, want UNAUTHORIZED", resp.Code)
				}
			},
		},
		{
			name: "user not found",
			setupContext: func(ctx context.Context) context.Context {
				claims := &domain.TokenClaims{
					IDCitizen: 54321,
					Email:     "nonexistent@example.com",
					Role:      domain.RoleUser,
				}
				return context.WithValue(ctx, middleware.UserContextKey, claims)
			},
			mockSetup: func(m *MockAuthService) {
				m.GetUserByIDCitizenFunc = func(ctx context.Context, idCitizen int) (*domain.UserPublic, error) {
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
			setupContext: func(ctx context.Context) context.Context {
				claims := &domain.TokenClaims{
					IDCitizen: 12345,
					Email:     "test@example.com",
					Role:      domain.RoleUser,
				}
				return context.WithValue(ctx, middleware.UserContextKey, claims)
			},
			mockSetup: func(m *MockAuthService) {
				m.GetUserByIDCitizenFunc = func(ctx context.Context, idCitizen int) (*domain.UserPublic, error) {
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
		{
			name: "admin user",
			setupContext: func(ctx context.Context) context.Context {
				claims := &domain.TokenClaims{
					IDCitizen: 99999,
					Email:     "admin@example.com",
					Role:      domain.RoleAdmin,
				}
				return context.WithValue(ctx, middleware.UserContextKey, claims)
			},
			mockSetup: func(m *MockAuthService) {
				m.GetUserByIDCitizenFunc = func(ctx context.Context, idCitizen int) (*domain.UserPublic, error) {
					return &domain.UserPublic{
						ID:        "admin-123",
						IDCitizen: 99999,
						Email:     "admin@example.com",
						Name:      "Admin User",
						Role:      domain.RoleAdmin,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.UserResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Role != domain.RoleAdmin {
					t.Errorf("Role = %v, want ADMIN", resp.Role)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := &MockAuthService{}
			tt.mockSetup(mockAuthService)

			req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)

			// Setup context with user claims
			ctx := tt.setupContext(req.Context())
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			// Use real handler with mock service injected
			h := shared.NewAuthHandler(mockAuthService, logger)
			handler := authhandler.GetMe(h)
			handler(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("status code = %v, want %v", w.Code, tt.wantStatusCode)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}
