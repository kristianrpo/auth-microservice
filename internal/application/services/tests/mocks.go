package tests

import (
	"context"
	"time"

	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// MockUserRepository is a mock implementation of ports.UserRepository
type MockUserRepository struct {
	CreateFunc         func(ctx context.Context, user *domain.User) error
	GetByIDFunc        func(ctx context.Context, id string) (*domain.User, error)
	GetByEmailFunc     func(ctx context.Context, email string) (*domain.User, error)
	GetByIDCitizenFunc func(ctx context.Context, idCitizen int) (*domain.User, error)
	UpdateFunc         func(ctx context.Context, user *domain.User) error
	DeleteFunc         func(ctx context.Context, id string) error
	ExistsFunc         func(ctx context.Context, email string) (bool, error)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) GetByIDCitizen(ctx context.Context, idCitizen int) (*domain.User, error) {
	if m.GetByIDCitizenFunc != nil {
		return m.GetByIDCitizenFunc(ctx, idCitizen)
	}
	// Default behavior: not found
	return nil, domainerrors.ErrUserNotFound
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, email)
	}
	return false, nil
}

// MockTokenRepository is a mock implementation of ports.TokenRepository
type MockTokenRepository struct {
	StoreRefreshTokenFunc  func(ctx context.Context, token string, data *domain.RefreshTokenData, ttl time.Duration) error
	GetRefreshTokenFunc    func(ctx context.Context, token string) (*domain.RefreshTokenData, error)
	DeleteRefreshTokenFunc func(ctx context.Context, token string) error
	BlacklistTokenFunc     func(ctx context.Context, token string, ttl time.Duration) error
	IsTokenBlacklistedFunc func(ctx context.Context, token string) (bool, error)
	DeleteUserTokensFunc   func(ctx context.Context, userID string) error
}

func (m *MockTokenRepository) StoreRefreshToken(ctx context.Context, token string, data *domain.RefreshTokenData, ttl time.Duration) error {
	if m.StoreRefreshTokenFunc != nil {
		return m.StoreRefreshTokenFunc(ctx, token, data, ttl)
	}
	return nil
}

func (m *MockTokenRepository) GetRefreshToken(ctx context.Context, token string) (*domain.RefreshTokenData, error) {
	if m.GetRefreshTokenFunc != nil {
		return m.GetRefreshTokenFunc(ctx, token)
	}
	return nil, nil
}

func (m *MockTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	if m.DeleteRefreshTokenFunc != nil {
		return m.DeleteRefreshTokenFunc(ctx, token)
	}
	return nil
}

func (m *MockTokenRepository) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	if m.BlacklistTokenFunc != nil {
		return m.BlacklistTokenFunc(ctx, token, ttl)
	}
	return nil
}

func (m *MockTokenRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	if m.IsTokenBlacklistedFunc != nil {
		return m.IsTokenBlacklistedFunc(ctx, token)
	}
	return false, nil
}

func (m *MockTokenRepository) DeleteUserTokens(ctx context.Context, userID string) error {
	if m.DeleteUserTokensFunc != nil {
		return m.DeleteUserTokensFunc(ctx, userID)
	}
	return nil
}

// MockOAuthClientRepository is a mock implementation of ports.OAuthClientRepository
type MockOAuthClientRepository struct {
	CreateFunc        func(ctx context.Context, client *domain.OAuthClient) error
	GetByIDFunc       func(ctx context.Context, id string) (*domain.OAuthClient, error)
	GetByClientIDFunc func(ctx context.Context, clientID string) (*domain.OAuthClient, error)
	UpdateFunc        func(ctx context.Context, client *domain.OAuthClient) error
	DeleteFunc        func(ctx context.Context, id string) error
	ListFunc          func(ctx context.Context) ([]*domain.OAuthClient, error)
}

func (m *MockOAuthClientRepository) Create(ctx context.Context, client *domain.OAuthClient) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, client)
	}
	return nil
}

func (m *MockOAuthClientRepository) GetByID(ctx context.Context, id string) (*domain.OAuthClient, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockOAuthClientRepository) GetByClientID(ctx context.Context, clientID string) (*domain.OAuthClient, error) {
	if m.GetByClientIDFunc != nil {
		return m.GetByClientIDFunc(ctx, clientID)
	}
	return nil, nil
}

func (m *MockOAuthClientRepository) Update(ctx context.Context, client *domain.OAuthClient) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, client)
	}
	return nil
}

func (m *MockOAuthClientRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockOAuthClientRepository) List(ctx context.Context) ([]*domain.OAuthClient, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}
