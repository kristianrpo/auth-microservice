package shared

import (
	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/application/services"
)

// AuthHandler manages the requests related to authentication
type AuthHandler struct {
	AuthService services.AuthServiceInterface
	Logger      *zap.Logger
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService services.AuthServiceInterface, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		Logger:      logger,
	}
}
