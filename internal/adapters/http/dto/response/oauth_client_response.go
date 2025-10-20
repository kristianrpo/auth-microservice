package response

import "time"

// OAuthClientResponse represents the response with OAuth client data
type OAuthClientResponse struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Scopes      []string  `json:"scopes"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

