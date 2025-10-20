package handler

import (
	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/application/services"
)

// AuthHandler manages the requests related to authentication
type AuthHandler struct {
	authService *services.AuthService
	logger      *zap.Logger
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService *services.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}
