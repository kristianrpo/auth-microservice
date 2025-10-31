package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/application/ports"
	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
	"github.com/kristianrpo/auth-microservice/internal/observability/metrics"
)

// OAuth2Service handles OAuth2 Client Credentials flow
type OAuth2Service struct {
	clientRepo        ports.OAuthClientRepository
	jwtSecret         string
	accessTokenExpiry time.Duration
	logger            *zap.Logger
}

// OAuth2ServiceInterface defines the subset of methods used by handlers so tests can inject mocks.
type OAuth2ServiceInterface interface {
	CreateClient(ctx context.Context, clientID, clientSecret, name, description string, scopes []string) (*domain.OAuthClient, error)
	ListClients(ctx context.Context) ([]*domain.OAuthClient, error)
	ClientCredentials(ctx context.Context, clientID, clientSecret string) (string, int64, error)
	ValidateAccessToken(ctx context.Context, tokenString string) (*domain.OAuthTokenClaims, error)
	GetClient(ctx context.Context, id string) (*domain.OAuthClient, error)
	DeleteClient(ctx context.Context, id string) error
}

// NewOAuth2Service creates a new instance of OAuth2Service
func NewOAuth2Service(
	clientRepo ports.OAuthClientRepository,
	jwtSecret string,
	accessTokenExpiry time.Duration,
	logger *zap.Logger,
) *OAuth2Service {
	return &OAuth2Service{
		clientRepo:        clientRepo,
		jwtSecret:         jwtSecret,
		accessTokenExpiry: accessTokenExpiry,
		logger:            logger,
	}
}

// ClientCredentials authenticates a client and generates an access token
func (s *OAuth2Service) ClientCredentials(ctx context.Context, clientID, clientSecret string) (string, int64, error) {
	// Retrieve client from database
	client, err := s.clientRepo.GetByClientID(ctx, clientID)
	if err != nil {
		s.logger.Warn("client not found", zap.String("client_id", clientID), zap.Error(err))
		return "", 0, domainerrors.ErrInvalidCredentials
	}

	// Validate client is active
	if !client.Active {
		s.logger.Warn("inactive client attempted authentication", zap.String("client_id", clientID))
		return "", 0, domainerrors.ErrInvalidClient
	}

	// Validate client secret
	if !client.ValidateSecret(clientSecret) {
		s.logger.Warn("invalid client secret", zap.String("client_id", clientID))
		return "", 0, domainerrors.ErrInvalidCredentials
	}

	// Generate access token
	accessToken, expiresIn, err := s.generateAccessToken(client)
	if err != nil {
		s.logger.Error("failed to generate access token", zap.Error(err), zap.String("client_id", clientID))
		return "", 0, fmt.Errorf("failed to generate access token: %w", err)
	}

	metrics.AddJWTTokensGenerated(1)

	s.logger.Info("client credentials token generated successfully",
		zap.String("client_id", clientID),
		zap.Int64("expires_in", expiresIn),
	)

	return accessToken, expiresIn, nil
}

// generateAccessToken creates a JWT access token for the OAuth client
func (s *OAuth2Service) generateAccessToken(client *domain.OAuthClient) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(s.accessTokenExpiry)

	claims := jwt.MapClaims{
		"client_id": client.ClientID,
		"scopes":    client.Scopes,
		"jti":       uuid.New().String(),
		"iat":       now.Unix(),
		"exp":       expiresAt.Unix(),
		"type":      "client_credentials",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", 0, err
	}

	expiresIn := int64(s.accessTokenExpiry.Seconds())
	return tokenString, expiresIn, nil
}

// ValidateAccessToken validates an OAuth2 access token
func (s *OAuth2Service) ValidateAccessToken(ctx context.Context, tokenString string) (*domain.OAuthTokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		s.logger.Debug("token parsing failed", zap.Error(err))
		return nil, domainerrors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, domainerrors.ErrInvalidToken
	}

	// Validate token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "client_credentials" {
		return nil, domainerrors.ErrInvalidTokenType
	}

	// Validate expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, domainerrors.ErrInvalidToken
	}
	if time.Unix(int64(exp), 0).Before(time.Now()) {
		return nil, domainerrors.ErrExpiredToken
	}

	// Extract client_id
	clientID, ok := claims["client_id"].(string)
	if !ok {
		return nil, domainerrors.ErrInvalidToken
	}

	// Extract scopes
	scopesInterface, ok := claims["scopes"].([]interface{})
	if !ok {
		return nil, domainerrors.ErrInvalidToken
	}

	scopes := make([]string, len(scopesInterface))
	for i, scope := range scopesInterface {
		scopes[i], ok = scope.(string)
		if !ok {
			return nil, domainerrors.ErrInvalidToken
		}
	}

	// Extract other claims
	jti, _ := claims["jti"].(string)
	iat, _ := claims["iat"].(float64)

	tokenClaims := &domain.OAuthTokenClaims{
		ClientID: clientID,
		Scopes:   scopes,
		TokenID:  jti,
		IssuedAt: int64(iat),
		ExpireAt: int64(exp),
		Type:     tokenType,
	}

	return tokenClaims, nil
}

// CreateClient creates a new OAuth2 client
func (s *OAuth2Service) CreateClient(ctx context.Context, clientID, clientSecret, name, description string, scopes []string) (*domain.OAuthClient, error) {
	// Check if client already exists
	existing, err := s.clientRepo.GetByClientID(ctx, clientID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("client with id %s already exists", clientID)
	}

	// Create new client
	client, err := domain.NewOAuthClient(clientID, clientSecret, name, description, scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth client: %w", err)
	}

	// Save to database
	if err := s.clientRepo.Create(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to save oauth client: %w", err)
	}

	s.logger.Info("oauth client created", zap.String("client_id", clientID))
	return client, nil
}

// ListClients retrieves all OAuth2 clients
func (s *OAuth2Service) ListClients(ctx context.Context) ([]*domain.OAuthClient, error) {
	return s.clientRepo.List(ctx)
}

// GetClient retrieves an OAuth2 client by ID
func (s *OAuth2Service) GetClient(ctx context.Context, id string) (*domain.OAuthClient, error) {
	return s.clientRepo.GetByID(ctx, id)
}

// DeleteClient deletes an OAuth2 client
func (s *OAuth2Service) DeleteClient(ctx context.Context, id string) error {
	return s.clientRepo.Delete(ctx, id)
}
