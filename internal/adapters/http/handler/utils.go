package handler

import (
	"encoding/json"
	nethttp "net/http"
)

// respondWithJSON sends a JSON response
func respondWithJSON(w nethttp.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		nethttp.Error(w, "Failed to encode response", nethttp.StatusInternalServerError)
	}
}
