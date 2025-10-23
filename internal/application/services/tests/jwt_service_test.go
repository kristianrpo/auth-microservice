package tests

import (
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/application/services"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

func TestJWTService_GenerateAccessToken(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)

	tests := []struct {
		name    string
		idCitizen  int
		email   string
		role    domain.Role
		wantErr bool
	}{
		{
			name:    "valid token generation",
		idCitizen: 123,
			email:   "test@example.com",
			role:    domain.RoleUser,
			wantErr: false,
		},
		{
			name:    "empty user id",
		idCitizen: 0,
			email:   "test@example.com",
			role:    domain.RoleUser,
			wantErr: false, // Should still generate token
		},
		{
			name:    "admin role",
		idCitizen: 999,
			email:   "admin@example.com",
			role:    domain.RoleAdmin,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		token, err := jwtService.GenerateAccessToken(tt.idCitizen, tt.email, tt.role)

			if tt.wantErr && err == nil {
				t.Errorf("GenerateAccessToken() expected error but got none")
				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("GenerateAccessToken() unexpected error: %v", err)
				return
			}

			if !tt.wantErr && token == "" {
				t.Errorf("GenerateAccessToken() returned empty token")
			}
		})
	}
}

func TestJWTService_GenerateRefreshToken(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)

	token, err := jwtService.GenerateRefreshToken(123, "test@example.com", domain.RoleUser)

	if err != nil {
		t.Errorf("GenerateRefreshToken() unexpected error: %v", err)
	}

	if token == "" {
		t.Errorf("GenerateRefreshToken() returned empty token")
	}
}

func TestJWTService_GenerateTokenPair(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)

	tokenPair, err := jwtService.GenerateTokenPair(123, "test@example.com", domain.RoleUser)

	if err != nil {
		t.Errorf("GenerateTokenPair() unexpected error: %v", err)
	}

	if tokenPair == nil {
		t.Errorf("GenerateTokenPair() returned nil")
		return
	}

	if tokenPair.AccessToken == "" {
		t.Errorf("GenerateTokenPair() AccessToken is empty")
	}

	if tokenPair.RefreshToken == "" {
		t.Errorf("GenerateTokenPair() RefreshToken is empty")
	}

	if tokenPair.TokenType != domain.TokenTypeBearer {
		t.Errorf("GenerateTokenPair() TokenType = %v, want %v", tokenPair.TokenType, domain.TokenTypeBearer)
	}

	if tokenPair.ExpiresIn != int64((15 * time.Minute).Seconds()) {
		t.Errorf("GenerateTokenPair() ExpiresIn = %v, want %v", tokenPair.ExpiresIn, int64((15 * time.Minute).Seconds()))
	}
}

func TestJWTService_ValidateToken(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)

	// Generate a valid token
	validToken, _ := jwtService.GenerateAccessToken(123, "test@example.com", domain.RoleUser)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "malformed token",
			token:   "not-a-jwt-token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := jwtService.ValidateToken(tt.token)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateToken() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateToken() unexpected error: %v", err)
				return
			}

			if claims == nil {
				t.Errorf("ValidateToken() returned nil claims")
				return
			}

			if claims.IDCitizen != 123 {
				t.Errorf("ValidateToken() IDCitizen = %v, want %v", claims.IDCitizen, 123)
			}

			if claims.Email != "test@example.com" {
				t.Errorf("ValidateToken() Email = %v, want %v", claims.Email, "test@example.com")
			}

			if claims.Role != domain.RoleUser {
				t.Errorf("ValidateToken() Role = %v, want %v", claims.Role, domain.RoleUser)
			}
		})
	}
}

func TestJWTService_ValidateAccessToken(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)

	accessToken, _ := jwtService.GenerateAccessToken(123, "test@example.com", domain.RoleUser)
	refreshToken, _ := jwtService.GenerateRefreshToken(123, "test@example.com", domain.RoleUser)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid access token",
			token:   accessToken,
			wantErr: false,
		},
		{
			name:    "refresh token as access token",
			token:   refreshToken,
			wantErr: true,
		},
		{
			name:    "invalid token",
			token:   "invalid.token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := jwtService.ValidateAccessToken(tt.token)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateAccessToken() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateAccessToken() unexpected error: %v", err)
				return
			}

			if claims.Type != domain.TokenTypeAccess {
				t.Errorf("ValidateAccessToken() Type = %v, want %v", claims.Type, domain.TokenTypeAccess)
			}
		})
	}
}

func TestJWTService_ValidateRefreshToken(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)

	accessToken, _ := jwtService.GenerateAccessToken(123, "test@example.com", domain.RoleUser)
	refreshToken, _ := jwtService.GenerateRefreshToken(123, "test@example.com", domain.RoleUser)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid refresh token",
			token:   refreshToken,
			wantErr: false,
		},
		{
			name:    "access token as refresh token",
			token:   accessToken,
			wantErr: true,
		},
		{
			name:    "invalid token",
			token:   "invalid.token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := jwtService.ValidateRefreshToken(tt.token)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateRefreshToken() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateRefreshToken() unexpected error: %v", err)
				return
			}

			if claims.Type != domain.TokenTypeRefresh {
				t.Errorf("ValidateRefreshToken() Type = %v, want %v", claims.Type, domain.TokenTypeRefresh)
			}
		})
	}
}

func TestJWTService_GetTokenExpiration(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)

	token, _ := jwtService.GenerateAccessToken(123, "test@example.com", domain.RoleUser)

	expiresAt, err := jwtService.GetTokenExpiration(token)

	if err != nil {
		t.Errorf("GetTokenExpiration() unexpected error: %v", err)
	}

	if expiresAt.IsZero() {
		t.Errorf("GetTokenExpiration() returned zero time")
	}

	// Token should expire in approximately 15 minutes
	expectedExpiration := time.Now().Add(15 * time.Minute)
	diff := expiresAt.Sub(expectedExpiration).Abs()

	if diff > 5*time.Second {
		t.Errorf("GetTokenExpiration() expiration time difference too large: %v", diff)
	}
}

func TestJWTService_ExpiredToken(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	// Create service with very short expiration
	jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 1*time.Millisecond, 1*time.Millisecond, logger)

	token, _ := jwtService.GenerateAccessToken(123, "test@example.com", domain.RoleUser)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	_, err := jwtService.ValidateToken(token)

	if err == nil {
		t.Errorf("ValidateToken() expected error for expired token but got none")
	}
}
