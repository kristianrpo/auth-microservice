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

// Login handles user authentication
// @Summary User login
// @Description Authenticates a user and returns access and refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login credentials"
// @Success 200 {object} response.TokenResponse "Login successful, tokens generated"
// @Failure 400 {object} response.ErrorResponse "Invalid request or missing data"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /login [post]
func Login(h *shared.AuthHandler) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		var req request.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Logger.Debug("invalid request body", zap.Error(err))
			httperrors.RespondWithError(w, httperrors.ErrInvalidRequestBody)
			return
		}

		// Basic validations
		if req.Email == "" || req.Password == "" {
			httperrors.RespondWithError(w, httperrors.ErrRequiredField)
			return
		}

		// Authenticate user
		tokenPair, err := h.AuthService.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			h.Logger.Warn("login failed", zap.Error(err), zap.String("email", req.Email))
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
