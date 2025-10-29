package admin

import (
	"encoding/json"
	nethttp "net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
)

// Token handles OAuth2 Client Credentials flow
// @Summary OAuth2 Client Credentials
// @Description Authenticates a client application and returns an access token for service-to-service communication.
// @Description
// @Description **Test Credentials (use in Swagger):**
// @Description ```json
// @Description {
// @Description   "client_id": "123",
// @Description   "client_secret": "123",
// @Description   "grant_type": "client_credentials"
// @Description }
// @Description ```
// @Tags OAuth2
// @Accept json
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param request body request.ClientCredentialsRequest true "Client Credentials"
// @Success 200 {object} response.ClientCredentialsResponse "Access token generated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request or missing parameters"
// @Failure 401 {object} response.ErrorResponse "Invalid client credentials"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /token [post]
func Token(h *shared.OAuth2Handler) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		var req request.ClientCredentialsRequest

		// Check Content-Type to handle both JSON and form-urlencoded
		contentType := r.Header.Get("Content-Type")

		// OAuth2 spec prefers application/x-www-form-urlencoded, so we check for JSON explicitly
		if strings.Contains(contentType, "application/json") {
			// Parse JSON
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				h.Logger.Debug("invalid request body (JSON)", zap.Error(err))
				httperrors.RespondWithError(w, httperrors.ErrInvalidRequestBody)
				return
			}
		} else {
			// Parse form data (default for OAuth2)
			if err := r.ParseForm(); err != nil {
				h.Logger.Debug("failed to parse form", zap.Error(err))
				httperrors.RespondWithError(w, httperrors.ErrInvalidRequestBody)
				return
			}

			req.ClientID = r.FormValue("client_id")
			req.ClientSecret = r.FormValue("client_secret")
			req.GrantType = r.FormValue("grant_type")
		}

		// Validate required fields
		if req.ClientID == "" || req.ClientSecret == "" || req.GrantType == "" {
			httperrors.RespondWithError(w, httperrors.ErrRequiredField)
			return
		}

		// Validate grant_type
		if req.GrantType != "client_credentials" {
			httperrors.RespondWithErrorMessage(w, nethttp.StatusBadRequest, "unsupported grant_type, must be 'client_credentials'")
			return
		}

		// Authenticate client and generate token
		accessToken, expiresIn, err := h.OAuth2Service.ClientCredentials(r.Context(), req.ClientID, req.ClientSecret)
		if err != nil {
			h.Logger.Warn("client credentials authentication failed", zap.Error(err), zap.String("client_id", req.ClientID))
			httperrors.RespondWithDomainError(w, err)
			return
		}

		// Return token response
		resp := response.ClientCredentialsResponse{
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresIn:   expiresIn,
		}

		shared.RespondWithJSON(w, nethttp.StatusOK, resp)
	}
}
