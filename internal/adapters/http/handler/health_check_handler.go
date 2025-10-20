package handler

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"time"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
)

// Health checks service health status
// @Summary Complete health check
// @Description Check the health status of the service and its dependencies (database, Redis)
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} response.HealthResponse "Service is healthy"
// @Failure 503 {object} response.HealthResponse "Service is unhealthy"
// @Router /health [get]
func (h *HealthHandler) Health(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	services := make(map[string]string)

	// Check database
	if err := h.db.PingContext(ctx); err != nil {
		h.logger.Warn("database health check failed", zap.Error(err))
		services["database"] = "unhealthy"
	} else {
		services["database"] = "healthy"
	}

	// Check Redis
	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		h.logger.Warn("redis health check failed", zap.Error(err))
		services["redis"] = "unhealthy"
	} else {
		services["redis"] = "healthy"
	}

	// Determine overall status
	overallStatus := "healthy"
	for _, status := range services {
		if status == "unhealthy" {
			overallStatus = "unhealthy"
			break
		}
	}

	resp := response.HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   h.version,
		Services:  services,
	}

	statusCode := nethttp.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = nethttp.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}
