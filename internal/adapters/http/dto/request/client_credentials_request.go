package request

// ClientCredentialsRequest represents the OAuth2 client credentials request
type ClientCredentialsRequest struct {
	ClientID     string `json:"client_id" form:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret" form:"client_secret" validate:"required"`
	GrantType    string `json:"grant_type" form:"grant_type" validate:"required,eq=client_credentials"`
}
