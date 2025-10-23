package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// JWTService handles the generation and validation of JWT tokens
type JWTService struct {
	secret               []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	logger               *zap.Logger
}

// CustomClaims extends the standard JWT claims
type CustomClaims struct {
	IDCitizen int         `json:"id_citizen"`
	Email     string      `json:"email"`
	Role      domain.Role `json:"role"`
	Type      string      `json:"type"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new instance of JWTService
func NewJWTService(secret string, accessDuration, refreshDuration time.Duration, logger *zap.Logger) *JWTService {
	return &JWTService{
		secret:               []byte(secret),
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
		logger:               logger,
	}
}

// GenerateAccessToken generates a new access token
func (s *JWTService) GenerateAccessToken(idCitizen int, email string, role domain.Role) (string, error) {
	return s.generateToken(idCitizen, email, role, domain.TokenTypeAccess, s.accessTokenDuration)
}

// GenerateRefreshToken generates a new refresh token
func (s *JWTService) GenerateRefreshToken(idCitizen int, email string, role domain.Role) (string, error) {
	return s.generateToken(idCitizen, email, role, domain.TokenTypeRefresh, s.refreshTokenDuration)
}

// GenerateTokenPair generates a token pair (access and refresh)
func (s *JWTService) GenerateTokenPair(idCitizen int, email string, role domain.Role) (*domain.TokenPair, error) {
	accessToken, err := s.GenerateAccessToken(idCitizen, email, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.GenerateRefreshToken(idCitizen, email, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    domain.TokenTypeBearer,
		ExpiresIn:    int64(s.accessTokenDuration.Seconds()),
	}, nil
}

// ValidateToken validates a token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		s.logger.Debug("token validation failed", zap.Error(err))
		return nil, domainerrors.ErrInvalidToken
	}

	if !token.Valid {
		return nil, domainerrors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, domainerrors.ErrInvalidToken
	}

	// Verify expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, domainerrors.ErrExpiredToken
	}

	return &domain.TokenClaims{
		IDCitizen: claims.IDCitizen,
		Email:  claims.Email,
		Role:   claims.Role,
		Type:   claims.Type,
	}, nil
}

// ValidateAccessToken validates a specific access token
func (s *JWTService) ValidateAccessToken(tokenString string) (*domain.TokenClaims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != domain.TokenTypeAccess {
		return nil, domainerrors.ErrInvalidTokenType
	}

	return claims, nil
}

// ValidateRefreshToken validates a specific refresh token
func (s *JWTService) ValidateRefreshToken(tokenString string) (*domain.TokenClaims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != domain.TokenTypeRefresh {
		return nil, domainerrors.ErrInvalidTokenType
	}

	return claims, nil
}

// GetTokenExpiration returns the expiration time of a token
func (s *JWTService) GetTokenExpiration(tokenString string) (time.Time, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})

	if err != nil {
		return time.Time{}, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || claims.ExpiresAt == nil {
		return time.Time{}, domainerrors.ErrInvalidToken
	}

	return claims.ExpiresAt.Time, nil
}

// generateToken is a helper method to generate tokens
func (s *JWTService) generateToken(idCitizen int, email string, role domain.Role, tokenType string, duration time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(duration)

	claims := CustomClaims{
		IDCitizen: idCitizen,
		Email:     email,
		Role:      role,
		Type:      tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "auth-microservice",
			Subject:   fmt.Sprintf("%d", idCitizen),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		s.logger.Error("failed to sign token", zap.Error(err), zap.String("type", tokenType))
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	s.logger.Debug("token generated successfully",
		zap.String("type", tokenType),
		zap.Int("id_citizen", idCitizen),
		zap.Time("expires_at", expiresAt))

	return tokenString, nil
}
