package shared

import (
	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/application/services"
)

// AdminOAuthClientsHandler manages OAuth clients administration (ADMIN only)
type AdminOAuthClientsHandler struct {
	OAuth2Service *services.OAuth2Service
	Logger        *zap.Logger
}

// NewAdminOAuthClientsHandler creates a new instance of AdminOAuthClientsHandler
func NewAdminOAuthClientsHandler(oauth2Service *services.OAuth2Service, logger *zap.Logger) *AdminOAuthClientsHandler {
	return &AdminOAuthClientsHandler{
		OAuth2Service: oauth2Service,
		Logger:        logger,
	}
}
