package ports

import (
	"context"

	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// OAuthClientRepository defines the interface for OAuth client persistence operations
type OAuthClientRepository interface {
	// Create creates a new OAuth client
	Create(ctx context.Context, client *domain.OAuthClient) error

	// GetByClientID retrieves an OAuth client by client_id
	GetByClientID(ctx context.Context, clientID string) (*domain.OAuthClient, error)

	// GetByID retrieves an OAuth client by ID
	GetByID(ctx context.Context, id string) (*domain.OAuthClient, error)

	// Update updates an existing OAuth client
	Update(ctx context.Context, client *domain.OAuthClient) error

	// Delete soft deletes an OAuth client
	Delete(ctx context.Context, id string) error

	// List retrieves all active OAuth clients
	List(ctx context.Context) ([]*domain.OAuthClient, error)
}
