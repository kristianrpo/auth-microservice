package auth

import (
	nethttp "net/http"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/middleware"
)

// GetMe retrieves authenticated user information
// @Summary Get current user
// @Description Get the authenticated user's information using the JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.UserResponse "User information"
// @Failure 401 {object} response.ErrorResponse "Unauthorized or invalid token"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/me [get]
func GetMe(h *shared.AuthHandler) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		// Get claims from context
		claims, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			httperrors.RespondWithError(w, httperrors.ErrUnauthorized)
			return
		}

		// Get complete user
		user, err := h.AuthService.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			h.Logger.Error("failed to get user", zap.Error(err), zap.String("user_id", claims.UserID))
			httperrors.RespondWithDomainError(w, err)
			return
		}

		// Convert to DTO
		resp := response.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		shared.RespondWithJSON(w, nethttp.StatusOK, resp)
	}
}
