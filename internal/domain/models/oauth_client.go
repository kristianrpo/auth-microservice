package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// OAuthClient represents an OAuth2 client application for service-to-service communication
type OAuthClient struct {
	ID           string    `json:"id"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"-"` // Never expose in JSON
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Scopes       []string  `json:"scopes"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewOAuthClient creates a new OAuth client with hashed secret
func NewOAuthClient(clientID, clientSecret, name, description string, scopes []string) (*OAuthClient, error) {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &OAuthClient{
		ClientID:     clientID,
		ClientSecret: string(hashedSecret),
		Name:         name,
		Description:  description,
		Scopes:       scopes,
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// ValidateSecret checks if the provided secret matches the stored hash
func (c *OAuthClient) ValidateSecret(secret string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.ClientSecret), []byte(secret))
	return err == nil
}

// HasScope checks if the client has a specific scope
func (c *OAuthClient) HasScope(scope string) bool {
	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// OAuthTokenClaims represents the claims for an OAuth access token
type OAuthTokenClaims struct {
	ClientID string   `json:"client_id"`
	Scopes   []string `json:"scopes"`
	TokenID  string   `json:"jti"`
	IssuedAt int64    `json:"iat"`
	ExpireAt int64    `json:"exp"`
	Type     string   `json:"type"` // "client_credentials"
}
