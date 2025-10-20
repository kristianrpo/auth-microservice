package response

// ClientCredentialsResponse represents the OAuth2 token response
type ClientCredentialsResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}
