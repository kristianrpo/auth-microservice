package tests

import (
	"context"

	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// OAuth2ServiceInterface defines the interface for OAuth2 operations used by handlers
type OAuth2ServiceInterface interface {
	CreateClient(ctx context.Context, clientID, clientSecret, name, description string, scopes []string) (*domain.OAuthClient, error)
	ListClients(ctx context.Context) ([]*domain.OAuthClient, error)
	ClientCredentials(ctx context.Context, clientID, clientSecret string) (string, int64, error)
}

// MockOAuth2Service is a mock implementation of OAuth2Service
type MockOAuth2Service struct {
	CreateClientFunc      func(ctx context.Context, clientID, clientSecret, name, description string, scopes []string) (*domain.OAuthClient, error)
	ListClientsFunc       func(ctx context.Context) ([]*domain.OAuthClient, error)
	ClientCredentialsFunc func(ctx context.Context, clientID, clientSecret string) (string, int64, error)
}

func (m *MockOAuth2Service) CreateClient(ctx context.Context, clientID, clientSecret, name, description string, scopes []string) (*domain.OAuthClient, error) {
	if m.CreateClientFunc != nil {
		return m.CreateClientFunc(ctx, clientID, clientSecret, name, description, scopes)
	}
	return nil, nil
}

func (m *MockOAuth2Service) ListClients(ctx context.Context) ([]*domain.OAuthClient, error) {
	if m.ListClientsFunc != nil {
		return m.ListClientsFunc(ctx)
	}
	return nil, nil
}

func (m *MockOAuth2Service) ClientCredentials(ctx context.Context, clientID, clientSecret string) (string, int64, error) {
	if m.ClientCredentialsFunc != nil {
		return m.ClientCredentialsFunc(ctx, clientID, clientSecret)
	}
	return "", 0, nil
}

// Additional stub methods to satisfy the OAuth2ServiceInterface used by handlers
func (m *MockOAuth2Service) ValidateAccessToken(ctx context.Context, tokenString string) (*domain.OAuthTokenClaims, error) {
	return &domain.OAuthTokenClaims{ClientID: "client-123", Scopes: []string{"read"}, TokenID: "jti", IssuedAt: 0, ExpireAt: 0, Type: "client_credentials"}, nil
}

func (m *MockOAuth2Service) GetClient(ctx context.Context, id string) (*domain.OAuthClient, error) {
	return &domain.OAuthClient{ID: id, ClientID: "client-123", Name: "Test", Active: true}, nil
}

func (m *MockOAuth2Service) DeleteClient(ctx context.Context, id string) error {
	return nil
}
