package admin

import (
	"encoding/json"
	nethttp "net/http"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
)

// CreateOAuthClient creates a new OAuth2 client (ADMIN only)
// @Summary Create OAuth2 Client
// @Description Creates a new OAuth2 client for service-to-service authentication. Only administrators can create clients.
// @Tags Admin - OAuth Clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateOAuthClientRequest true "OAuth Client data"
// @Success 201 {object} response.OAuthClientResponse "OAuth client created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - Admin role required"
// @Failure 409 {object} response.ErrorResponse "Client already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /admin/oauth-clients [post]
func CreateOAuthClient(h *shared.AdminOAuthClientsHandler) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		var req request.CreateOAuthClientRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Logger.Debug("invalid request body", zap.Error(err))
			httperrors.RespondWithError(w, httperrors.ErrInvalidRequestBody)
			return
		}

		// Validate required fields
		if req.ClientID == "" || req.ClientSecret == "" || req.Name == "" {
			httperrors.RespondWithError(w, httperrors.ErrRequiredField)
			return
		}

		// Create OAuth client
		client, err := h.OAuth2Service.CreateClient(
			r.Context(),
			req.ClientID,
			req.ClientSecret,
			req.Name,
			req.Description,
			req.Scopes,
		)
		if err != nil {
			h.Logger.Error("failed to create oauth client", zap.Error(err))
			httperrors.RespondWithDomainError(w, err)
			return
		}

		// Convert to DTO
		resp := response.OAuthClientResponse{
			ID:          client.ID,
			ClientID:    client.ClientID,
			Name:        client.Name,
			Description: client.Description,
			Scopes:      client.Scopes,
			Active:      client.Active,
			CreatedAt:   client.CreatedAt,
			UpdatedAt:   client.UpdatedAt,
		}

		shared.RespondWithJSON(w, nethttp.StatusCreated, resp)
	}
}
