package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// TokenRepository is the Redis implementation of the token repository
type TokenRepository struct {
	client *redis.Client
	logger *zap.Logger
}

// NewTokenRepository creates a new instance of TokenRepository
func NewTokenRepository(client *redis.Client, logger *zap.Logger) *TokenRepository {
	return &TokenRepository{
		client: client,
		logger: logger,
	}
}

// StoreRefreshToken stores a refresh token in cache
func (r *TokenRepository) StoreRefreshToken(ctx context.Context, token string, data *domain.RefreshTokenData, ttl time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s", token)

	jsonData, err := json.Marshal(data)
	if err != nil {
		r.logger.Error("failed to marshal refresh token data", zap.Error(err))
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	err = r.client.Set(ctx, key, jsonData, ttl).Err()
	if err != nil {
		r.logger.Error("failed to store refresh token", zap.Error(err), zap.String("key", key))
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	r.logger.Debug("refresh token stored successfully", zap.String("user_id", data.UserID))
	return nil
}

// GetRefreshToken retrieves the data of a refresh token
func (r *TokenRepository) GetRefreshToken(ctx context.Context, token string) (*domain.RefreshTokenData, error) {
	key := fmt.Sprintf("refresh_token:%s", token)

	jsonData, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, domainerrors.ErrInvalidToken
	}
	if err != nil {
		r.logger.Error("failed to get refresh token", zap.Error(err), zap.String("key", key))
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	var data domain.RefreshTokenData
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		r.logger.Error("failed to unmarshal refresh token data", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	return &data, nil
}

// DeleteRefreshToken deletes a refresh token from cache
func (r *TokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("refresh_token:%s", token)

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.logger.Error("failed to delete refresh token", zap.Error(err), zap.String("key", key))
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	r.logger.Debug("refresh token deleted successfully")
	return nil
}

// BlacklistToken adds a token to the blacklist
func (r *TokenRepository) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", token)

	err := r.client.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		r.logger.Error("failed to blacklist token", zap.Error(err), zap.String("key", key))
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	r.logger.Debug("token blacklisted successfully")
	return nil
}

// IsTokenBlacklisted verifies if a token is in the blacklist
func (r *TokenRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		r.logger.Error("failed to check if token is blacklisted", zap.Error(err), zap.String("key", key))
		return false, fmt.Errorf("failed to check blacklist: %w", err)
	}

	return exists > 0, nil
}

// DeleteUserTokens deletes all refresh tokens of a user
func (r *TokenRepository) DeleteUserTokens(ctx context.Context, userID string) error {
	pattern := "refresh_token:*"

	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		jsonData, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var data domain.RefreshTokenData
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			continue
		}

		if data.UserID == userID {
			if err := r.client.Del(ctx, key).Err(); err != nil {
				r.logger.Error("failed to delete user token", zap.Error(err), zap.String("key", key))
			}
		}
	}

	if err := iter.Err(); err != nil {
		r.logger.Error("failed to iterate user tokens", zap.Error(err), zap.String("user_id", userID))
		return fmt.Errorf("failed to delete user tokens: %w", err)
	}

	r.logger.Info("user tokens deleted successfully", zap.String("user_id", userID))
	return nil
}

// NewRedisClient creates a new connection to Redis
func NewRedisClient(address, password string, db int, logger *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         address,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	logger.Info("redis connection established successfully")
	return client, nil
}
