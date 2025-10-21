package admin

import (
	nethttp "net/http"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
)

// ListOAuthClients retrieves all OAuth2 clients (ADMIN only)
// @Summary List OAuth2 Clients
// @Description Retrieves all OAuth2 clients. Only administrators can list clients.
// @Tags Admin - OAuth Clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} response.OAuthClientResponse "List of OAuth clients"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - Admin role required"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /admin/oauth-clients [get]
func ListOAuthClients(h *shared.AdminOAuthClientsHandler) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		clients, err := h.OAuth2Service.ListClients(r.Context())
		if err != nil {
			h.Logger.Error("failed to list oauth clients", zap.Error(err))
			httperrors.RespondWithError(w, httperrors.ErrInternalServer)
			return
		}

		// Convert to DTOs
		var clientResponses []response.OAuthClientResponse
		for _, client := range clients {
			clientResponses = append(clientResponses, response.OAuthClientResponse{
				ID:          client.ID,
				ClientID:    client.ClientID,
				Name:        client.Name,
				Description: client.Description,
				Scopes:      client.Scopes,
				Active:      client.Active,
				CreatedAt:   client.CreatedAt,
				UpdatedAt:   client.UpdatedAt,
			})
		}

		shared.RespondWithJSON(w, nethttp.StatusOK, clientResponses)
	}
}
