package tests

import (
	"testing"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
	"github.com/kristianrpo/auth-microservice/internal/application/services"
)

func TestNewOAuth2Handler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name          string
		oauth2Service *services.OAuth2Service
		logger        *zap.Logger
		wantNil       bool
	}{
		{
			name:          "create oauth2 handler with valid service",
			oauth2Service: nil,
			logger:        logger,
			wantNil:       false,
		},
		{
			name:          "create oauth2 handler with nil logger",
			oauth2Service: nil,
			logger:        nil,
			wantNil:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := shared.NewOAuth2Handler(tt.oauth2Service, tt.logger)

			if tt.wantNil {
				if handler != nil {
					t.Errorf("NewOAuth2Handler() = %v, want nil", handler)
				}
				return
			}

			if handler == nil {
				t.Error("NewOAuth2Handler() returned nil handler")
				return
			}

			if handler.OAuth2Service != tt.oauth2Service {
				t.Errorf("NewOAuth2Handler() OAuth2Service = %v, want %v", handler.OAuth2Service, tt.oauth2Service)
			}

			if handler.Logger != tt.logger {
				t.Errorf("NewOAuth2Handler() Logger = %v, want %v", handler.Logger, tt.logger)
			}
		})
	}
}

func TestOAuth2Handler_Fields(t *testing.T) {
	logger := zap.NewNop()
	handler := shared.NewOAuth2Handler(nil, logger)

	if handler.Logger != logger {
		t.Errorf("OAuth2Handler.Logger = %v, want %v", handler.Logger, logger)
	}

	if handler.OAuth2Service != nil {
		t.Errorf("OAuth2Handler.OAuth2Service = %v, want nil", handler.OAuth2Service)
	}
}

