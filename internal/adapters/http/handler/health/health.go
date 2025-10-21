package health

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// DBPinger is the minimal interface used by health checks for DB
type DBPinger interface {
	PingContext(ctx context.Context) error
}

// RedisPinger is the minimal interface used by health checks for Redis
type RedisPinger interface {
	Ping(ctx context.Context) *redis.StatusCmd
}

// HealthHandler manages the health check
type HealthHandler struct {
	db          DBPinger
	redisClient RedisPinger
	logger      *zap.Logger
	version     string
}

// NewHealthHandler creates a new instance of HealthHandler
func NewHealthHandler(db DBPinger, redisClient RedisPinger, logger *zap.Logger, version string) *HealthHandler {
	return &HealthHandler{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
		version:     version,
	}
}
