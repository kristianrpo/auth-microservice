package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// UserRepository is the PostgreSQL implementation of the user repository
type UserRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, email, password, name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.Name,
		user.Role.String(),
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("failed to create user", zap.Error(err), zap.String("email", user.Email))
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.Info("user created successfully", zap.String("user_id", user.ID), zap.String("email", user.Email))
	return nil
}

// GetByID retrieves a user by their ID
//
//nolint:dupl // Similar to GetByEmail but queries by ID instead of email
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := &domain.User{}
	var roleStr string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&roleStr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domainerrors.ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("failed to get user by ID", zap.Error(err), zap.String("user_id", id))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	role, _ := domain.ParseRole(roleStr)
	user.Role = role
	return user, nil
}

// GetByEmail retrieves a user by their email
//
//nolint:dupl // Similar to GetByID but queries by email instead of ID
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	user := &domain.User{}
	var roleStr string
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&roleStr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domainerrors.ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("failed to get user by email", zap.Error(err), zap.String("email", email))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	role, _ := domain.ParseRole(roleStr)
	user.Role = role
	return user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET email = $2, password = $3, name = $4, role = $5, updated_at = $6
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.Name,
		user.Role.String(),
		user.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("failed to update user", zap.Error(err), zap.String("user_id", user.ID))
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domainerrors.ErrUserNotFound
	}

	r.logger.Info("user updated successfully", zap.String("user_id", user.ID))
	return nil
}

// Delete performs a soft delete of a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET deleted_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, time.Now())
	if err != nil {
		r.logger.Error("failed to delete user", zap.Error(err), zap.String("user_id", id))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domainerrors.ErrUserNotFound
	}

	r.logger.Info("user deleted successfully", zap.String("user_id", id))
	return nil
}

// Exists verifies if a user exists by email
func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM users
			WHERE email = $1 AND deleted_at IS NULL
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		r.logger.Error("failed to check user existence", zap.Error(err), zap.String("email", email))
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// NewDB creates a new connection to PostgreSQL
func NewDB(connectionString string, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connection established successfully")
	return db, nil
}

// InitSchema initializes the database schema
func InitSchema(db *sql.DB) error {
	// First, create tables
	createTables := `
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'USER',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS oauth_clients (
			id VARCHAR(36) PRIMARY KEY,
			client_id VARCHAR(255) UNIQUE NOT NULL,
			client_secret VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			scopes TEXT[] NOT NULL DEFAULT '{}',
			active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`

	if _, err := db.Exec(createTables); err != nil {
		return err
	}

	// Then, create indexes
	createIndexes := `
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;
		CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
		CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
		CREATE INDEX IF NOT EXISTS idx_oauth_clients_client_id ON oauth_clients(client_id);
		CREATE INDEX IF NOT EXISTS idx_oauth_clients_active ON oauth_clients(active);
	`

	_, err := db.Exec(createIndexes)
	return err
}
