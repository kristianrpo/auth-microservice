package handler

import (
	"encoding/json"
	nethttp "net/http"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
)

// Register handles new user registration
// @Summary Register a new user
// @Description Create a new user account in the system
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "User registration data"
// @Success 201 {object} response.UserResponse "User created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request or missing data"
// @Failure 409 {object} response.ErrorResponse "User already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req request.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Debug("invalid request body", zap.Error(err))
		httperrors.RespondWithError(w, httperrors.ErrInvalidRequestBody)
		return
	}

	// Basic validations
	if req.Email == "" || req.Password == "" || req.Name == "" {
		httperrors.RespondWithError(w, httperrors.ErrRequiredField)
		return
	}

	// Register user
	user, err := h.authService.Register(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		h.logger.Error("failed to register user", zap.Error(err))
		httperrors.RespondWithDomainError(w, err)
		return
	}

	// Convert to DTO
	resp := response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	respondWithJSON(w, nethttp.StatusCreated, resp)
}
