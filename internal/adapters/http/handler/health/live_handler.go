package health

import (
	"encoding/json"
	nethttp "net/http"
)

// Live checks if the service is alive (liveness probe)
// @Summary Liveness check
// @Description Check if the service is alive (used by Kubernetes)
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Service is alive"
// @Router /health/live [get]
func (h *HealthHandler) Live(w nethttp.ResponseWriter, r *nethttp.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(nethttp.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
}
