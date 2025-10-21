package request

// CreateOAuthClientRequest represents the request to create an OAuth client
type CreateOAuthClientRequest struct {
	ClientID     string   `json:"client_id" validate:"required,min=3"`
	ClientSecret string   `json:"client_secret" validate:"required,min=8"`
	Name         string   `json:"name" validate:"required,min=3"`
	Description  string   `json:"description"`
	Scopes       []string `json:"scopes"`
}
