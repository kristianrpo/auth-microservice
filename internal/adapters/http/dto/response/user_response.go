package response

import (
	"time"

	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// UserResponse represents the response with user data
type UserResponse struct {
	ID        string      `json:"id"`
	IDCitizen int         `json:"id_citizen"`
	Email     string      `json:"email"`
	Name      string      `json:"name"`
	Role      domain.Role `json:"role"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
