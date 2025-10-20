package health

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// HealthHandler manages the health check
type HealthHandler struct {
	db          *sql.DB
	redisClient *redis.Client
	logger      *zap.Logger
	version     string
}

// NewHealthHandler creates a new instance of HealthHandler
func NewHealthHandler(db *sql.DB, redisClient *redis.Client, logger *zap.Logger, version string) *HealthHandler {
	return &HealthHandler{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
		version:     version,
	}
}
