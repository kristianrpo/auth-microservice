package ports

import (
	"context"

	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// UserRepository defines the persistence operations for users
type UserRepository interface {
	// Create creates a new user in the database
	Create(ctx context.Context, user *domain.User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (*domain.User, error)

	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *domain.User) error

	// Delete deletes a user (soft 	delete)
	Delete(ctx context.Context, id string) error

	// Exists verifies if a user exists by email
	Exists(ctx context.Context, email string) (bool, error)
}
