package tests

import (
	"context"

	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// AuthServiceInterface defines the interface for auth operations used by handlers
type AuthServiceInterface interface {
	Login(ctx context.Context, email, password string) (*domain.TokenPair, error)
	Register(ctx context.Context, email, password, name string, idCitizen int) (*domain.UserPublic, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
	GetUserByIDCitizen(ctx context.Context, idCitizen int) (*domain.UserPublic, error)
}

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	LoginFunc        func(ctx context.Context, email, password string) (*domain.TokenPair, error)
	RegisterFunc     func(ctx context.Context, email, password, name string, idCitizen int) (*domain.UserPublic, error)
	RefreshTokenFunc func(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
	LogoutFunc       func(ctx context.Context, accessToken, refreshToken string) error
	GetUserByIDCitizenFunc  func(ctx context.Context, idCitizen int) (*domain.UserPublic, error)
	// Additional mocked functions to satisfy services.AuthServiceInterface
	ValidateAccessTokenFunc func(ctx context.Context, token string) (*domain.TokenClaims, error)
	RevokeAllUserTokensFunc func(ctx context.Context, idCitizen int) error
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (*domain.TokenPair, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, email, password)
	}
	return nil, nil
}

func (m *MockAuthService) Register(ctx context.Context, email, password, name string, idCitizen int) (*domain.UserPublic, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, email, password, name, idCitizen)
	}
	return nil, nil
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(ctx, refreshToken)
	}
	return nil, nil
}

func (m *MockAuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	if m.LogoutFunc != nil {
		return m.LogoutFunc(ctx, accessToken, refreshToken)
	}
	return nil
}

func (m *MockAuthService) GetUserByIDCitizen(ctx context.Context, idCitizen int) (*domain.UserPublic, error) {
	if m.GetUserByIDCitizenFunc != nil {
		return m.GetUserByIDCitizenFunc(ctx, idCitizen)
	}
	return nil, nil
}

func (m *MockAuthService) ValidateAccessToken(ctx context.Context, token string) (*domain.TokenClaims, error) {
	if m.ValidateAccessTokenFunc != nil {
		return m.ValidateAccessTokenFunc(ctx, token)
	}
	return nil, nil
}

func (m *MockAuthService) RevokeAllUserTokens(ctx context.Context, idCitizen int) error {
	if m.RevokeAllUserTokensFunc != nil {
		return m.RevokeAllUserTokensFunc(ctx, idCitizen)
	}
	return nil
}
