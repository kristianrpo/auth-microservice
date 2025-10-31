package auth

import (
	"encoding/json"
	nethttp "net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
	"github.com/kristianrpo/auth-microservice/internal/observability/metrics"
)

// Logout handles user logout
// @Summary User logout
// @Description Invalidates user tokens (access and refresh)
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.LogoutRequest false "Refresh token (optional)"
// @Success 200 {object} response.MessageResponse "Logout successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /logout [post]
func Logout(h *shared.AuthHandler) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		metrics.IncLogoutRequests()

		// Get access token from header
		authHeader := r.Header.Get("Authorization")
		var accessToken string
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 {
				accessToken = parts[1]
			}
		}

		// Get refresh token from body (optional)
		var req request.LogoutRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		if accessToken == "" {
			httperrors.RespondWithError(w, httperrors.ErrRequiredField)
			return
		}

		// Perform logout
		if err := h.AuthService.Logout(r.Context(), accessToken, req.RefreshToken); err != nil {
			h.Logger.Error("logout failed", zap.Error(err))
			httperrors.RespondWithDomainError(w, err)
			return
		}

		resp := response.MessageResponse{Message: "logout successful"}
		shared.RespondWithJSON(w, nethttp.StatusOK, resp)
	}
}
