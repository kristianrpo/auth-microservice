package handler

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"time"

	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
)

// Ready checks if the service is ready to receive traffic
// @Summary Readiness check
// @Description Check if the service is ready to receive traffic (used by Kubernetes)
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Service is ready"
// @Failure 503 {object} response.ErrorResponse "Service is not ready"
// @Router /health/ready [get]
func (h *HealthHandler) Ready(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Check critical connections
	if err := h.db.PingContext(ctx); err != nil {
		httperrors.RespondWithErrorMessage(w, nethttp.StatusServiceUnavailable, "database not ready")
		return
	}

	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		httperrors.RespondWithErrorMessage(w, nethttp.StatusServiceUnavailable, "redis not ready")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(nethttp.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}
