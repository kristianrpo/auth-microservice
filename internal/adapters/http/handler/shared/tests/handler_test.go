package tests

import (
	"testing"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
	"github.com/kristianrpo/auth-microservice/internal/application/services"
)

func TestNewAuthHandler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name        string
		authService *services.AuthService
		logger      *zap.Logger
		wantNil     bool
	}{
		{
			name:        "create auth handler with valid service",
			authService: nil,
			logger:      logger,
			wantNil:     false,
		},
		{
			name:        "create auth handler with nil logger",
			authService: nil,
			logger:      nil,
			wantNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := shared.NewAuthHandler(tt.authService, tt.logger)

			if tt.wantNil {
				if handler != nil {
					t.Errorf("NewAuthHandler() = %v, want nil", handler)
				}
				return
			}

			if handler == nil {
				t.Error("NewAuthHandler() returned nil handler")
				return
			}

			if handler.AuthService != tt.authService {
				t.Errorf("NewAuthHandler() AuthService = %v, want %v", handler.AuthService, tt.authService)
			}

			if handler.Logger != tt.logger {
				t.Errorf("NewAuthHandler() Logger = %v, want %v", handler.Logger, tt.logger)
			}
		})
	}
}

func TestAuthHandler_Fields(t *testing.T) {
	logger := zap.NewNop()
	handler := shared.NewAuthHandler(nil, logger)

	if handler.Logger != logger {
		t.Errorf("AuthHandler.Logger = %v, want %v", handler.Logger, logger)
	}

	if handler.AuthService != nil {
		t.Errorf("AuthHandler.AuthService = %v, want nil", handler.AuthService)
	}
}

