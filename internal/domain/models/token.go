package domain

import (
	"time"
)

// TokenPair representa un par de tokens (access y refresh)
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // Segundos hasta la expiración
}

// TokenClaims representa los claims personalizados del JWT
type TokenClaims struct {
	IDCitizen int    `json:"id_citizen"`
	Email     string `json:"email"`
	Role      Role   `json:"role"`
	Type      string `json:"type"` // "access" o "refresh"
}

// RefreshTokenData represents the data stored in Redis for a refresh token
type RefreshTokenData struct {
	IDCitizen int       `json:"id_citizen"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// BlacklistedToken represents a revoked/blacklisted token
type BlacklistedToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

const (
	// TokenTypeAccess represents an access token
	TokenTypeAccess = "access"
	// TokenTypeRefresh represents a refresh token
	TokenTypeRefresh = "refresh"
	// TokenTypeBearer represents the token type in the Authorization header
	TokenTypeBearer = "Bearer"
)
