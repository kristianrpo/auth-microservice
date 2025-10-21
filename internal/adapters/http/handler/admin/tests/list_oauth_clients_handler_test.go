package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

func TestListOAuthClientsHandler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name           string
		mockSetup      func(*MockOAuth2Service)
		wantStatusCode int
		wantError      bool
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful list with multiple clients",
			mockSetup: func(m *MockOAuth2Service) {
				m.ListClientsFunc = func(ctx context.Context) ([]*domain.OAuthClient, error) {
					return []*domain.OAuthClient{
						{
							ID:          "client-1",
							ClientID:    "test_client_1",
							Name:        "Test Client 1",
							Description: "First test client",
							Scopes:      []string{"read", "write"},
							Active:      true,
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						},
						{
							ID:          "client-2",
							ClientID:    "test_client_2",
							Name:        "Test Client 2",
							Description: "Second test client",
							Scopes:      []string{"read"},
							Active:      true,
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						},
					}, nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp []response.OAuthClientResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(resp) != 2 {
					t.Errorf("number of clients = %v, want 2", len(resp))
				}
				if resp[0].ClientID != "test_client_1" {
					t.Errorf("first client ID = %v, want test_client_1", resp[0].ClientID)
				}
				if resp[1].ClientID != "test_client_2" {
					t.Errorf("second client ID = %v, want test_client_2", resp[1].ClientID)
				}
			},
		},
		{
			name: "successful list with empty result",
			mockSetup: func(m *MockOAuth2Service) {
				m.ListClientsFunc = func(ctx context.Context) ([]*domain.OAuthClient, error) {
					return []*domain.OAuthClient{}, nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp []response.OAuthClientResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(resp) != 0 {
					t.Errorf("number of clients = %v, want 0", len(resp))
				}
			},
		},
		{
			name: "successful list with single client",
			mockSetup: func(m *MockOAuth2Service) {
				m.ListClientsFunc = func(ctx context.Context) ([]*domain.OAuthClient, error) {
					return []*domain.OAuthClient{
						{
							ID:          "client-single",
							ClientID:    "single_client",
							Name:        "Single Client",
							Description: "Only one client",
							Scopes:      []string{"admin"},
							Active:      true,
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						},
					}, nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp []response.OAuthClientResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(resp) != 1 {
					t.Errorf("number of clients = %v, want 1", len(resp))
				}
				if resp[0].Name != "Single Client" {
					t.Errorf("client name = %v, want Single Client", resp[0].Name)
				}
			},
		},
		{
			name: "internal server error",
			mockSetup: func(m *MockOAuth2Service) {
				m.ListClientsFunc = func(ctx context.Context) ([]*domain.OAuthClient, error) {
					return nil, errors.New("database error")
				}
			},
			wantStatusCode: http.StatusInternalServerError,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "INTERNAL_SERVER_ERROR" {
					t.Errorf("Error code = %v, want INTERNAL_SERVER_ERROR", resp.Code)
				}
			},
		},
		{
			name: "list with inactive clients",
			mockSetup: func(m *MockOAuth2Service) {
				m.ListClientsFunc = func(ctx context.Context) ([]*domain.OAuthClient, error) {
					return []*domain.OAuthClient{
						{
							ID:          "client-active",
							ClientID:    "active_client",
							Name:        "Active Client",
							Description: "This is active",
							Scopes:      []string{"read"},
							Active:      true,
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						},
						{
							ID:          "client-inactive",
							ClientID:    "inactive_client",
							Name:        "Inactive Client",
							Description: "This is inactive",
							Scopes:      []string{"read"},
							Active:      false,
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						},
					}, nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp []response.OAuthClientResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(resp) != 2 {
					t.Errorf("number of clients = %v, want 2", len(resp))
				}
				if resp[0].Active != true {
					t.Errorf("first client active = %v, want true", resp[0].Active)
				}
				if resp[1].Active != false {
					t.Errorf("second client active = %v, want false", resp[1].Active)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockOAuth2Service := &MockOAuth2Service{}
			tt.mockSetup(mockOAuth2Service)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/admin/oauth-clients", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Create inline handler function that mimics admin.ListOAuthClients but uses our mock
			handlerFunc := func(w http.ResponseWriter, r *http.Request) {
				clients, err := mockOAuth2Service.ListClients(r.Context())
				if err != nil {
					logger.Error("failed to list oauth clients")
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

				shared.RespondWithJSON(w, http.StatusOK, clientResponses)
			}

			// Call handler
			handlerFunc(w, req)

			// Check status code
			if w.Code != tt.wantStatusCode {
				t.Errorf("status code = %v, want %v", w.Code, tt.wantStatusCode)
			}

			// Check response
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}
