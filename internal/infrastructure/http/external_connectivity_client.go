package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ExternalConnectivityClient implements the client for external-connectivity microservice
type ExternalConnectivityClient struct {
	baseURL      string
	authURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	logger       *zap.Logger
	
	// Token caching
	tokenMutex   sync.RWMutex
	accessToken  string
	tokenExpiry  time.Time
}

// TokenResponse represents the OAuth2 token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// NewExternalConnectivityClient creates a new external connectivity client
func NewExternalConnectivityClient(baseURL, authURL, clientID, clientSecret string, logger *zap.Logger) *ExternalConnectivityClient {
	return &ExternalConnectivityClient{
		baseURL:      baseURL,
		authURL:      authURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// getAccessToken obtains an OAuth2 access token using client credentials
func (c *ExternalConnectivityClient) getAccessToken(ctx context.Context) (string, error) {
	// Check if we have a valid cached token
	c.tokenMutex.RLock()
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		token := c.accessToken
		c.tokenMutex.RUnlock()
		return token, nil
	}
	c.tokenMutex.RUnlock()

	// Need to get a new token
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	// Double-check after acquiring write lock
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	c.logger.Debug("obtaining new access token from auth service")

	// Prepare token request
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.authURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("failed to obtain access token", zap.Error(err))
		return "", fmt.Errorf("failed to obtain access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("token request failed",
			zap.Int("status_code", resp.StatusCode))
		return "", fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		c.logger.Error("failed to decode token response", zap.Error(err))
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	// Cache the token (with 30 seconds buffer before expiry)
	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-30) * time.Second)

	c.logger.Debug("access token obtained successfully",
		zap.Int("expires_in", tokenResp.ExpiresIn))

	return c.accessToken, nil
}

// CheckCitizenExists verifies if a citizen exists in the centralizer
// Returns true if citizen exists (HTTP 200), false if not exists (HTTP 204)
func (c *ExternalConnectivityClient) CheckCitizenExists(ctx context.Context, idCitizen int) (bool, error) {
	// Get access token
	token, err := c.getAccessToken(ctx)
	if err != nil {
		c.logger.Error("failed to get access token", zap.Error(err))
		return false, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/api/connectivity/external/citizens/%d/exists", c.baseURL, idCitizen)

	c.logger.Debug("checking citizen existence in external-connectivity",
		zap.String("url", url),
		zap.Int("id_citizen", idCitizen))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error("failed to create request", zap.Error(err))
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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

	case http.StatusUnauthorized: // 401 - Token expired or invalid
		c.logger.Warn("unauthorized response, clearing token cache")
		c.tokenMutex.Lock()
		c.accessToken = ""
		c.tokenExpiry = time.Time{}
		c.tokenMutex.Unlock()
		return false, fmt.Errorf("unauthorized: token may be invalid")

	default:
		c.logger.Warn("unexpected status code from external-connectivity",
			zap.Int("status_code", resp.StatusCode),
			zap.Int("id_citizen", idCitizen))
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
