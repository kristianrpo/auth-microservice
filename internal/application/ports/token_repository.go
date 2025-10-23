package ports

import (
	"context"
	"time"

	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// TokenRepository defines the cache operations for tokens
type TokenRepository interface {
	// StoreRefreshToken stores a refresh token in cache
	StoreRefreshToken(ctx context.Context, token string, data *domain.RefreshTokenData, ttl time.Duration) error

	// GetRefreshToken retrieves the data of a refresh token
	GetRefreshToken(ctx context.Context, token string) (*domain.RefreshTokenData, error)

	// DeleteRefreshToken deletes a refresh token from cache
	DeleteRefreshToken(ctx context.Context, token string) error

	// BlacklistToken adds a token to the blacklist
	BlacklistToken(ctx context.Context, token string, ttl time.Duration) error

	// IsTokenBlacklisted verifies if a token is in the blacklist
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)

	// DeleteUserTokens deletes all refresh tokens of a user
	DeleteUserTokens(ctx context.Context, idCitizen int) error
}
