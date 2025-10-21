package auth

import (
	"encoding/json"
	nethttp "net/http"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
)

// Refresh handles token renewal
// @Summary Refresh tokens
// @Description Generate a new token pair using a valid refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.TokenResponse "Tokens refreshed successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request or missing data"
// @Failure 401 {object} response.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/refresh [post]
func Refresh(h *shared.AuthHandler) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		var req request.RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Logger.Debug("invalid request body", zap.Error(err))
			httperrors.RespondWithError(w, httperrors.ErrInvalidRequestBody)
			return
		}

		if req.RefreshToken == "" {
			httperrors.RespondWithError(w, httperrors.ErrRequiredField)
			return
		}

		// Refresh token
		tokenPair, err := h.AuthService.RefreshToken(r.Context(), req.RefreshToken)
		if err != nil {
			h.Logger.Warn("token refresh failed", zap.Error(err))
			httperrors.RespondWithDomainError(w, err)
			return
		}

		// Convert to DTO
		resp := response.TokenResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			TokenType:    tokenPair.TokenType,
			ExpiresIn:    tokenPair.ExpiresIn,
		}

		shared.RespondWithJSON(w, nethttp.StatusOK, resp)
	}
}
