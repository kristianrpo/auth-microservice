package shared

import (
	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/application/services"
)

// OAuth2Handler manages OAuth2-related requests
type OAuth2Handler struct {
	OAuth2Service services.OAuth2ServiceInterface
	Logger        *zap.Logger
}

// NewOAuth2Handler creates a new instance of OAuth2Handler
func NewOAuth2Handler(oauth2Service services.OAuth2ServiceInterface, logger *zap.Logger) *OAuth2Handler {
	return &OAuth2Handler{
		OAuth2Service: oauth2Service,
		Logger:        logger,
	}
}
