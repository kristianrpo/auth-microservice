package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ExternalConnectivityClient implements the client for external-connectivity microservice
type ExternalConnectivityClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
}

// NewExternalConnectivityClient creates a new external connectivity client
func NewExternalConnectivityClient(baseURL string, logger *zap.Logger) *ExternalConnectivityClient {
	return &ExternalConnectivityClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// CheckCitizenExists verifies if a citizen exists in the centralizer
// Returns true if citizen exists (HTTP 200), false if not exists (HTTP 204)
func (c *ExternalConnectivityClient) CheckCitizenExists(ctx context.Context, idCitizen int) (bool, error) {
	url := fmt.Sprintf("%s/api/external/citizen/%d", c.baseURL, idCitizen)

	c.logger.Debug("checking citizen existence in external-connectivity",
		zap.String("url", url),
		zap.Int("id_citizen", idCitizen))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error("failed to create request", zap.Error(err))
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("failed to call external-connectivity service",
			zap.Error(err),
			zap.String("url", url))
		return false, fmt.Errorf("failed to call external-connectivity service: %w", err)
	}
	defer resp.Body.Close()

	c.logger.Debug("external-connectivity response",
		zap.Int("status_code", resp.StatusCode),
		zap.Int("id_citizen", idCitizen))

	switch resp.StatusCode {
	case http.StatusOK: // 200 - Citizen exists
		c.logger.Info("citizen exists in centralizer",
			zap.Int("id_citizen", idCitizen))
		return true, nil

	case http.StatusNoContent: // 204 - Citizen does not exist
		c.logger.Info("citizen does not exist in centralizer",
			zap.Int("id_citizen", idCitizen))
		return false, nil

	default:
		c.logger.Warn("unexpected status code from external-connectivity",
			zap.Int("status_code", resp.StatusCode),
			zap.Int("id_citizen", idCitizen))
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
