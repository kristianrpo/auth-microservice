package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	admin "github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/admin"
	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
)

func TestTokenHandler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name           string
		contentType    string
		requestBody    interface{}
		formData       url.Values
		mockSetup      func(*MockOAuth2Service)
		wantStatusCode int
		wantError      bool
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "successful token generation with JSON",
			contentType: "application/json",
			requestBody: request.ClientCredentialsRequest{
				ClientID:     "test_client",
				ClientSecret: "secret123",
				GrantType:    "client_credentials",
			},
			mockSetup: func(m *MockOAuth2Service) {
				m.ClientCredentialsFunc = func(ctx context.Context, clientID, clientSecret string) (string, int64, error) {
					return "access_token_123", 3600, nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ClientCredentialsResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.AccessToken != "access_token_123" {
					t.Errorf("AccessToken = %v, want access_token_123", resp.AccessToken)
				}
				if resp.TokenType != "Bearer" {
					t.Errorf("TokenType = %v, want Bearer", resp.TokenType)
				}
				if resp.ExpiresIn != 3600 {
					t.Errorf("ExpiresIn = %v, want 3600", resp.ExpiresIn)
				}
			},
		},
		{
			name:        "successful token generation with form data",
			contentType: "application/x-www-form-urlencoded",
			formData: url.Values{
				"client_id":     []string{"test_client"},
				"client_secret": []string{"secret123"},
				"grant_type":    []string{"client_credentials"},
			},
			mockSetup: func(m *MockOAuth2Service) {
				m.ClientCredentialsFunc = func(ctx context.Context, clientID, clientSecret string) (string, int64, error) {
					return "access_token_456", 7200, nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ClientCredentialsResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.AccessToken != "access_token_456" {
					t.Errorf("AccessToken = %v, want access_token_456", resp.AccessToken)
				}
			},
		},
		{
			name:           "invalid json body",
			contentType:    "application/json",
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
			name:        "missing client_id",
			contentType: "application/json",
			requestBody: request.ClientCredentialsRequest{
				ClientID:     "",
				ClientSecret: "secret123",
				GrantType:    "client_credentials",
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
			name:        "missing client_secret",
			contentType: "application/json",
			requestBody: request.ClientCredentialsRequest{
				ClientID:     "test_client",
				ClientSecret: "",
				GrantType:    "client_credentials",
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
			name:        "missing grant_type",
			contentType: "application/json",
			requestBody: request.ClientCredentialsRequest{
				ClientID:     "test_client",
				ClientSecret: "secret123",
				GrantType:    "",
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
			name:        "invalid grant_type",
			contentType: "application/json",
			requestBody: request.ClientCredentialsRequest{
				ClientID:     "test_client",
				ClientSecret: "secret123",
				GrantType:    "authorization_code",
			},
			mockSetup:      func(m *MockOAuth2Service) {},
			wantStatusCode: http.StatusBadRequest,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if !strings.Contains(resp.Error, "grant_type") {
					t.Errorf("Error message should mention grant_type")
				}
			},
		},
		{
			name:        "invalid client credentials",
			contentType: "application/json",
			requestBody: request.ClientCredentialsRequest{
				ClientID:     "wrong_client",
				ClientSecret: "wrong_secret",
				GrantType:    "client_credentials",
			},
			mockSetup: func(m *MockOAuth2Service) {
				m.ClientCredentialsFunc = func(ctx context.Context, clientID, clientSecret string) (string, int64, error) {
					return "", 0, domainerrors.ErrInvalidCredentials
				}
			},
			wantStatusCode: http.StatusUnauthorized,
			wantError:      true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Code != "INVALID_CREDENTIALS" {
					t.Errorf("Error code = %v, want INVALID_CREDENTIALS", resp.Code)
				}
			},
		},
		{
			name:        "internal server error",
			contentType: "application/json",
			requestBody: request.ClientCredentialsRequest{
				ClientID:     "test_client",
				ClientSecret: "secret123",
				GrantType:    "client_credentials",
			},
			mockSetup: func(m *MockOAuth2Service) {
				m.ClientCredentialsFunc = func(ctx context.Context, clientID, clientSecret string) (string, int64, error) {
					return "", 0, errors.New("database error")
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockOAuth2Service := &MockOAuth2Service{}
			tt.mockSetup(mockOAuth2Service)

			// Create request
			var req *http.Request
			if tt.formData != nil {
				req = httptest.NewRequest(http.MethodPost, "/auth/token", strings.NewReader(tt.formData.Encode()))
			} else {
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
				req = httptest.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))
			}

			req.Header.Set("Content-Type", tt.contentType)

			// Create response recorder
			w := httptest.NewRecorder()

			h := shared.NewOAuth2Handler(mockOAuth2Service, logger)
			handler := admin.Token(h)
			handler(w, req)

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
