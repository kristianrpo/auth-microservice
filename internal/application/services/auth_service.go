package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/application/ports"
	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// AuthService handles the business logic of authentication
type AuthService struct {
	userRepo   ports.UserRepository
	tokenRepo  ports.TokenRepository
	jwtService *JWTService
	logger     *zap.Logger
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(
	userRepo ports.UserRepository,
	tokenRepo ports.TokenRepository,
	jwtService *JWTService,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtService: jwtService,
		logger:     logger,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, email, password, name string) (*domain.UserPublic, error) {
	s.logger.Info("attempting to register user", zap.String("email", email))

	// Verificar si el usuario ya existe
	exists, err := s.userRepo.Exists(ctx, email)
	if err != nil {
		s.logger.Error("failed to check user existence", zap.Error(err))
		return nil, domainerrors.ErrInternal
	}

	if exists {
		s.logger.Warn("user already exists", zap.String("email", email))
		return nil, domainerrors.ErrUserAlreadyExists
	}

	// Create new user
	user, err := domain.NewUser(email, password, name)
	if err != nil {
		s.logger.Error("failed to create user entity", zap.Error(err))
		return nil, err
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("failed to save user", zap.Error(err))
		return nil, domainerrors.ErrInternal
	}

	s.logger.Info("user registered successfully", zap.String("user_id", user.ID), zap.String("email", email))
	return user.ToPublic(), nil
}

// Login authenticates a user and generates tokens
func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.TokenPair, error) {
	s.logger.Info("attempting login", zap.String("email", email))

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == domainerrors.ErrUserNotFound {
			s.logger.Warn("login failed: user not found", zap.String("email", email))
			return nil, domainerrors.ErrInvalidCredentials
		}
		s.logger.Error("failed to get user", zap.Error(err))
		return nil, domainerrors.ErrInternal
	}

	// Verify password
	if err := user.ComparePassword(password); err != nil {
		s.logger.Warn("login failed: invalid password", zap.String("email", email))
		return nil, domainerrors.ErrInvalidCredentials
	}

	// Generate token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		s.logger.Error("failed to generate token pair", zap.Error(err))
		return nil, domainerrors.ErrInternal
	}

	// Store refresh token in Redis
	refreshTokenData := &domain.RefreshTokenData{
		UserID:    user.ID,
		Email:     user.Email,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(s.jwtService.refreshTokenDuration),
	}

	err = s.tokenRepo.StoreRefreshToken(
		ctx,
		tokenPair.RefreshToken,
		refreshTokenData,
		s.jwtService.refreshTokenDuration,
	)
	if err != nil {
		s.logger.Error("failed to store refresh token", zap.Error(err))
		// No retornamos error aquí, el login fue exitoso
	}

	s.logger.Info("login successful", zap.String("user_id", user.ID))
	return tokenPair, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	s.logger.Debug("attempting to refresh token")

	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		s.logger.Warn("invalid refresh token", zap.Error(err))
		return nil, err
	}

	// Verify if the token is in the blacklist
	blacklisted, err := s.tokenRepo.IsTokenBlacklisted(ctx, refreshToken)
	if err != nil {
		s.logger.Error("failed to check token blacklist", zap.Error(err))
		return nil, domainerrors.ErrInternal
	}

	if blacklisted {
		s.logger.Warn("refresh token is blacklisted", zap.String("user_id", claims.UserID))
		return nil, domainerrors.ErrTokenRevoked
	}

	// Verify if the refresh token exists in Redis
	_, err = s.tokenRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.Warn("refresh token not found in cache", zap.Error(err))
		return nil, domainerrors.ErrInvalidToken
	}

	// Generate new token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(claims.UserID, claims.Email)
	if err != nil {
		s.logger.Error("failed to generate new token pair", zap.Error(err))
		return nil, domainerrors.ErrInternal
	}

	// Delete old refresh token
	if err := s.tokenRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		s.logger.Error("failed to delete old refresh token", zap.Error(err))
	}

	// Store new refresh token
	refreshTokenData := &domain.RefreshTokenData{
		UserID:    claims.UserID,
		Email:     claims.Email,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(s.jwtService.refreshTokenDuration),
	}

	err = s.tokenRepo.StoreRefreshToken(
		ctx,
		tokenPair.RefreshToken,
		refreshTokenData,
		s.jwtService.refreshTokenDuration,
	)
	if err != nil {
		s.logger.Error("failed to store new refresh token", zap.Error(err))
	}

	s.logger.Info("token refreshed successfully", zap.String("user_id", claims.UserID))
	return tokenPair, nil
}

// Logout invalidates the tokens of a user
func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	s.logger.Debug("attempting logout")

	// Validate access token
	claims, err := s.jwtService.ValidateAccessToken(accessToken)
	if err != nil {
		s.logger.Warn("invalid access token on logout", zap.Error(err))
		return err
	}

	// Add access token to blacklist
	expiresAt, err := s.jwtService.GetTokenExpiration(accessToken)
	if err == nil {
		ttl := time.Until(expiresAt)
		if ttl > 0 {
			if err := s.tokenRepo.BlacklistToken(ctx, accessToken, ttl); err != nil {
				s.logger.Error("failed to blacklist access token", zap.Error(err))
			}
		}
	}

	// Delete refresh token if provided
	if refreshToken != "" {
		if err := s.tokenRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
			s.logger.Error("failed to delete refresh token", zap.Error(err))
		}
	}

	s.logger.Info("logout successful", zap.String("user_id", claims.UserID))
	return nil
}

// GetUserByID retrieves a user by their ID
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*domain.UserPublic, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}

	return user.ToPublic(), nil
}

// ValidateAccessToken validates an access token and verifies that it is not revoked
func (s *AuthService) ValidateAccessToken(ctx context.Context, token string) (*domain.TokenClaims, error) {
	// Validate token
	claims, err := s.jwtService.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}

	// Verify if the token is in the blacklist
	blacklisted, err := s.tokenRepo.IsTokenBlacklisted(ctx, token)
	if err != nil {
		s.logger.Error("failed to check token blacklist", zap.Error(err))
		return nil, domainerrors.ErrInternal
	}

	if blacklisted {
		return nil, domainerrors.ErrTokenRevoked
	}

	return claims, nil
}

// RevokeAllUserTokens revokes all tokens of a user
func (s *AuthService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	s.logger.Info("revoking all user tokens", zap.String("user_id", userID))

	if err := s.tokenRepo.DeleteUserTokens(ctx, userID); err != nil {
		s.logger.Error("failed to revoke user tokens", zap.Error(err), zap.String("user_id", userID))
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}

	s.logger.Info("all user tokens revoked successfully", zap.String("user_id", userID))
	return nil
}
