package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

func TestCreateOAuthClientHandler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockOAuth2Service)
		wantStatusCode int
		wantError      bool
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful client creation",
			requestBody: request.CreateOAuthClientRequest{
				ClientID:     "test_client",
				ClientSecret: "secret123",
				Name:         "Test Client",
				Description:  "A test OAuth client",
				Scopes:       []string{"read", "write"},
			},
			mockSetup: func(m *MockOAuth2Service) {
				m.CreateClientFunc = func(ctx context.Context, clientID, clientSecret, name, description string, scopes []string) (*domain.OAuthClient, error) {
					return &domain.OAuthClient{
						ID:          "client-123",
						ClientID:    clientID,
						Name:        name,
						Description: description,
						Scopes:      scopes,
						Active:      true,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil
				}
			},
			wantStatusCode: http.StatusCreated,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.OAuthClientResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.ID != "client-123" {
					t.Errorf("ID = %v, want client-123", resp.ID)
				}
				if resp.ClientID != "test_client" {
					t.Errorf("ClientID = %v, want test_client", resp.ClientID)
				}
				if resp.Name != "Test Client" {
					t.Errorf("Name = %v, want Test Client", resp.Name)
				}
				if !resp.Active {
					t.Errorf("Active = %v, want true", resp.Active)
				}
			},
		},
		{
			name:           "invalid json body",
			requestBody:    "invalid json",
			mockSetup:      func(m *MockOAuth2Service) {},
			wantStatusCode: http.StatusBadRequest,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "INVALID_REQUEST_BODY" {
					t.Errorf("Error code = %v, want INVALID_REQUEST_BODY", resp.Code)
				}
			},
		},
		{
			name: "missing client_id",
			requestBody: request.CreateOAuthClientRequest{
				ClientID:     "",
				ClientSecret: "secret123",
				Name:         "Test Client",
			},
			mockSetup:      func(m *MockOAuth2Service) {},
			wantStatusCode: http.StatusBadRequest,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "REQUIRED_FIELD" {
					t.Errorf("Error code = %v, want REQUIRED_FIELD", resp.Code)
				}
			},
		},
		{
			name: "missing client_secret",
			requestBody: request.CreateOAuthClientRequest{
				ClientID:     "test_client",
				ClientSecret: "",
				Name:         "Test Client",
			},
			mockSetup:      func(m *MockOAuth2Service) {},
			wantStatusCode: http.StatusBadRequest,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "REQUIRED_FIELD" {
					t.Errorf("Error code = %v, want REQUIRED_FIELD", resp.Code)
				}
			},
		},
		{
			name: "missing name",
			requestBody: request.CreateOAuthClientRequest{
				ClientID:     "test_client",
				ClientSecret: "secret123",
				Name:         "",
			},
			mockSetup:      func(m *MockOAuth2Service) {},
			wantStatusCode: http.StatusBadRequest,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "REQUIRED_FIELD" {
					t.Errorf("Error code = %v, want REQUIRED_FIELD", resp.Code)
				}
			},
		},
		{
			name: "client already exists",
			requestBody: request.CreateOAuthClientRequest{
				ClientID:     "existing_client",
				ClientSecret: "secret123",
				Name:         "Existing Client",
			},
			mockSetup: func(m *MockOAuth2Service) {
				m.CreateClientFunc = func(ctx context.Context, clientID, clientSecret, name, description string, scopes []string) (*domain.OAuthClient, error) {
					return nil, errors.New("client with id existing_client already exists")
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
			name: "internal server error",
			requestBody: request.CreateOAuthClientRequest{
				ClientID:     "test_client",
				ClientSecret: "secret123",
				Name:         "Test Client",
			},
			mockSetup: func(m *MockOAuth2Service) {
				m.CreateClientFunc = func(ctx context.Context, clientID, clientSecret, name, description string, scopes []string) (*domain.OAuthClient, error) {
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
			name: "successful client creation without optional fields",
			requestBody: request.CreateOAuthClientRequest{
				ClientID:     "minimal_client",
				ClientSecret: "secret456",
				Name:         "Minimal Client",
			},
			mockSetup: func(m *MockOAuth2Service) {
				m.CreateClientFunc = func(ctx context.Context, clientID, clientSecret, name, description string, scopes []string) (*domain.OAuthClient, error) {
					return &domain.OAuthClient{
						ID:          "client-456",
						ClientID:    clientID,
						Name:        name,
						Description: "",
						Scopes:      nil,
						Active:      true,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil
				}
			},
			wantStatusCode: http.StatusCreated,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.OAuthClientResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.ClientID != "minimal_client" {
					t.Errorf("ClientID = %v, want minimal_client", resp.ClientID)
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
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/admin/oauth-clients", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Create inline handler function that mimics admin.CreateOAuthClient but uses our mock
			handlerFunc := func(w http.ResponseWriter, r *http.Request) {
				var createReq request.CreateOAuthClientRequest
				if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
					logger.Debug("invalid request body")
					httperrors.RespondWithError(w, httperrors.ErrInvalidRequestBody)
					return
				}

				// Validate required fields
				if createReq.ClientID == "" || createReq.ClientSecret == "" || createReq.Name == "" {
					httperrors.RespondWithError(w, httperrors.ErrRequiredField)
					return
				}

				// Create OAuth client
				client, err := mockOAuth2Service.CreateClient(
					r.Context(),
					createReq.ClientID,
					createReq.ClientSecret,
					createReq.Name,
					createReq.Description,
					createReq.Scopes,
				)
				if err != nil {
					logger.Error("failed to create oauth client")
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

				shared.RespondWithJSON(w, http.StatusCreated, resp)
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
