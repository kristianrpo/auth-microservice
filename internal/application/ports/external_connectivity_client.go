package ports

import "context"

// ExternalConnectivityClient defines the interface for communicating with the external-connectivity microservice
type ExternalConnectivityClient interface {
	// CheckCitizenExists verifies if a citizen exists in the centralizer
	// Returns true if citizen exists (HTTP 200), false if not exists (HTTP 204)
	CheckCitizenExists(ctx context.Context, idCitizen int) (bool, error)
}
