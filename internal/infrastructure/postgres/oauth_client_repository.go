package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"

	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// OAuthClientRepository is the PostgreSQL implementation of the OAuth client repository
type OAuthClientRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewOAuthClientRepository creates a new instance of OAuthClientRepository
func NewOAuthClientRepository(db *sql.DB, logger *zap.Logger) *OAuthClientRepository {
	return &OAuthClientRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new OAuth client in the database
func (r *OAuthClientRepository) Create(ctx context.Context, client *domain.OAuthClient) error {
	client.ID = uuid.New().String()
	client.CreatedAt = time.Now()
	client.UpdatedAt = time.Now()

	query := `
		INSERT INTO oauth_clients (id, client_id, client_secret, name, description, scopes, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		client.ID,
		client.ClientID,
		client.ClientSecret,
		client.Name,
		client.Description,
		pq.Array(client.Scopes),
		client.Active,
		client.CreatedAt,
		client.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("failed to create oauth client", zap.Error(err), zap.String("client_id", client.ClientID))
		return fmt.Errorf("failed to create oauth client: %w", err)
	}

	r.logger.Info("oauth client created successfully", zap.String("client_id", client.ClientID))
	return nil
}

// GetByClientID retrieves an OAuth client by their client_id
//
//nolint:dupl // Similar to GetByID but queries by client_id instead of id
func (r *OAuthClientRepository) GetByClientID(ctx context.Context, clientID string) (*domain.OAuthClient, error) {
	query := `
		SELECT id, client_id, client_secret, name, description, scopes, active, created_at, updated_at
		FROM oauth_clients
		WHERE client_id = $1 AND active = true
	`

	client := &domain.OAuthClient{}
	var scopes pq.StringArray

	err := r.db.QueryRowContext(ctx, query, clientID).Scan(
		&client.ID,
		&client.ClientID,
		&client.ClientSecret,
		&client.Name,
		&client.Description,
		&scopes,
		&client.Active,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domainerrors.ErrInvalidCredentials
	}
	if err != nil {
		r.logger.Error("failed to get oauth client by client_id", zap.Error(err), zap.String("client_id", clientID))
		return nil, fmt.Errorf("failed to get oauth client: %w", err)
	}

	client.Scopes = scopes
	return client, nil
}

// GetByID retrieves an OAuth client by their ID
//
//nolint:dupl // Similar to GetByClientID but queries by id instead of client_id
func (r *OAuthClientRepository) GetByID(ctx context.Context, id string) (*domain.OAuthClient, error) {
	query := `
		SELECT id, client_id, client_secret, name, description, scopes, active, created_at, updated_at
		FROM oauth_clients
		WHERE id = $1
	`

	client := &domain.OAuthClient{}
	var scopes pq.StringArray

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&client.ID,
		&client.ClientID,
		&client.ClientSecret,
		&client.Name,
		&client.Description,
		&scopes,
		&client.Active,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domainerrors.ErrClientNotFound
	}
	if err != nil {
		r.logger.Error("failed to get oauth client by ID", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("failed to get oauth client: %w", err)
	}

	client.Scopes = scopes
	return client, nil
}

// Update updates an existing OAuth client
func (r *OAuthClientRepository) Update(ctx context.Context, client *domain.OAuthClient) error {
	client.UpdatedAt = time.Now()

	query := `
		UPDATE oauth_clients
		SET name = $1, description = $2, scopes = $3, active = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		client.Name,
		client.Description,
		pq.Array(client.Scopes),
		client.Active,
		client.UpdatedAt,
		client.ID,
	)

	if err != nil {
		r.logger.Error("failed to update oauth client", zap.Error(err), zap.String("id", client.ID))
		return fmt.Errorf("failed to update oauth client: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domainerrors.ErrClientNotFound
	}

	r.logger.Info("oauth client updated successfully", zap.String("id", client.ID))
	return nil
}

// Delete deactivates an OAuth client (soft delete)
func (r *OAuthClientRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE oauth_clients
		SET active = false, updated_at = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		r.logger.Error("failed to delete oauth client", zap.Error(err), zap.String("id", id))
		return fmt.Errorf("failed to delete oauth client: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domainerrors.ErrClientNotFound
	}

	r.logger.Info("oauth client deleted successfully", zap.String("id", id))
	return nil
}

// List retrieves all active OAuth clients
func (r *OAuthClientRepository) List(ctx context.Context) ([]*domain.OAuthClient, error) {
	query := `
		SELECT id, client_id, client_secret, name, description, scopes, active, created_at, updated_at
		FROM oauth_clients
		WHERE active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Error("failed to list oauth clients", zap.Error(err))
		return nil, fmt.Errorf("failed to list oauth clients: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			r.logger.Error("failed to close rows", zap.Error(closeErr))
		}
	}()

	var clients []*domain.OAuthClient
	for rows.Next() {
		client := &domain.OAuthClient{}
		var scopes pq.StringArray

		err := rows.Scan(
			&client.ID,
			&client.ClientID,
			&client.ClientSecret,
			&client.Name,
			&client.Description,
			&scopes,
			&client.Active,
			&client.CreatedAt,
			&client.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("failed to scan oauth client", zap.Error(err))
			return nil, fmt.Errorf("failed to scan oauth client: %w", err)
		}

		client.Scopes = scopes
		clients = append(clients, client)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("error iterating oauth clients", zap.Error(err))
		return nil, fmt.Errorf("error iterating oauth clients: %w", err)
	}

	return clients, nil
}
