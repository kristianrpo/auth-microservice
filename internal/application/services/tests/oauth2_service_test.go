package tests

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"

	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
	"github.com/kristianrpo/auth-microservice/internal/application/services"
)

func TestOAuth2Service_ClientCredentials(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Create a test OAuth client with known secret
	testClient, _ := domain.NewOAuthClient("client-123", "secret123", "Test Client", "Test Description", []string{"read"})

	tests := []struct {
		name              string
		clientID          string
		clientSecret      string
		getByClientIDFunc func(ctx context.Context, clientID string) (*domain.OAuthClient, error)
		wantErr           bool
		expectedErr       error
	}{
		{
			name:         "successful authentication",
			clientID:     "client-123",
			clientSecret: "secret123",
			getByClientIDFunc: func(ctx context.Context, clientID string) (*domain.OAuthClient, error) {
				return testClient, nil
			},
			wantErr: false,
		},
		{
			name:         "client not found",
			clientID:     "nonexistent",
			clientSecret: "secret123",
			getByClientIDFunc: func(ctx context.Context, clientID string) (*domain.OAuthClient, error) {
				return nil, domainerrors.ErrClientNotFound
			},
			wantErr:     true,
			expectedErr: domainerrors.ErrInvalidCredentials,
		},
		{
			name:         "invalid secret",
			clientID:     "client-123",
			clientSecret: "wrongsecret",
			getByClientIDFunc: func(ctx context.Context, clientID string) (*domain.OAuthClient, error) {
				return testClient, nil
			},
			wantErr:     true,
			expectedErr: domainerrors.ErrInvalidCredentials,
		},
		{
			name:         "inactive client",
			clientID:     "client-123",
			clientSecret: "secret123",
			getByClientIDFunc: func(ctx context.Context, clientID string) (*domain.OAuthClient, error) {
				inactiveClient, _ := domain.NewOAuthClient("client-123", "secret123", "Test Client", "Test Description", []string{"read"})
				inactiveClient.Active = false
				return inactiveClient, nil
			},
			wantErr:     true,
			expectedErr: domainerrors.ErrInvalidClient,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClientRepo := &MockOAuthClientRepository{
				GetByClientIDFunc: tt.getByClientIDFunc,
			}
			oauth2Service := services.NewOAuth2Service(mockClientRepo, "test-secret-key-at-least-32-chars-long", 15*time.Minute, logger)

			token, expiresIn, err := oauth2Service.ClientCredentials(context.Background(), tt.clientID, tt.clientSecret)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ClientCredentials() expected error but got none")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("ClientCredentials() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("ClientCredentials() unexpected error: %v", err)
				return
			}

			if token == "" {
				t.Errorf("ClientCredentials() returned empty token")
			}

			if expiresIn <= 0 {
				t.Errorf("ClientCredentials() expiresIn = %v, want > 0", expiresIn)
			}
		})
	}
}

func TestOAuth2Service_CreateClient(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name              string
		clientID          string
		clientSecret      string
		clientName        string
		description       string
		scopes            []string
		getByClientIDFunc func(ctx context.Context, clientID string) (*domain.OAuthClient, error)
		createFunc        func(ctx context.Context, client *domain.OAuthClient) error
		wantErr           bool
	}{
		{
			name:         "successful creation",
			clientID:     "new-client",
			clientSecret: "newsecret123",
			clientName:   "New Client",
			description:  "New Description",
			scopes:       []string{"read", "write"},
			getByClientIDFunc: func(ctx context.Context, clientID string) (*domain.OAuthClient, error) {
				return nil, domainerrors.ErrClientNotFound
			},
			createFunc: func(ctx context.Context, client *domain.OAuthClient) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:         "client already exists",
			clientID:     "existing-client",
			clientSecret: "secret123",
			clientName:   "Existing Client",
			description:  "Description",
			scopes:       []string{"read"},
			getByClientIDFunc: func(ctx context.Context, clientID string) (*domain.OAuthClient, error) {
				existing, _ := domain.NewOAuthClient("existing-client", "secret123", "Existing", "Desc", []string{"read"})
				return existing, nil
			},
			wantErr: true,
		},
		{
			name:         "repository error",
			clientID:     "new-client",
			clientSecret: "newsecret123",
			clientName:   "New Client",
			description:  "Description",
			scopes:       []string{"read"},
			getByClientIDFunc: func(ctx context.Context, clientID string) (*domain.OAuthClient, error) {
				return nil, domainerrors.ErrClientNotFound
			},
			createFunc: func(ctx context.Context, client *domain.OAuthClient) error {
				return domainerrors.ErrInternal
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClientRepo := &MockOAuthClientRepository{
				GetByClientIDFunc: tt.getByClientIDFunc,
				CreateFunc:        tt.createFunc,
			}
			oauth2Service := services.NewOAuth2Service(mockClientRepo, "test-secret-key-at-least-32-chars-long", 15*time.Minute, logger)

			client, err := oauth2Service.CreateClient(context.Background(), tt.clientID, tt.clientSecret, tt.clientName, tt.description, tt.scopes)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateClient() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateClient() unexpected error: %v", err)
				return
			}

			if client == nil {
				t.Errorf("CreateClient() returned nil client")
				return
			}

			if client.ClientID != tt.clientID {
				t.Errorf("CreateClient() ClientID = %v, want %v", client.ClientID, tt.clientID)
			}
		})
	}
}

func TestOAuth2Service_ListClients(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	client1, _ := domain.NewOAuthClient("client-1", "secret1", "Client 1", "Desc 1", []string{"read"})
	client2, _ := domain.NewOAuthClient("client-2", "secret2", "Client 2", "Desc 2", []string{"write"})

	tests := []struct {
		name     string
		listFunc func(ctx context.Context) ([]*domain.OAuthClient, error)
		wantErr  bool
		wantLen  int
	}{
		{
			name: "list all clients",
			listFunc: func(ctx context.Context) ([]*domain.OAuthClient, error) {
				return []*domain.OAuthClient{client1, client2}, nil
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "empty list",
			listFunc: func(ctx context.Context) ([]*domain.OAuthClient, error) {
				return []*domain.OAuthClient{}, nil
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "repository error",
			listFunc: func(ctx context.Context) ([]*domain.OAuthClient, error) {
				return nil, domainerrors.ErrInternal
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClientRepo := &MockOAuthClientRepository{
				ListFunc: tt.listFunc,
			}
			oauth2Service := services.NewOAuth2Service(mockClientRepo, "test-secret-key-at-least-32-chars-long", 15*time.Minute, logger)

			clients, err := oauth2Service.ListClients(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Errorf("ListClients() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ListClients() unexpected error: %v", err)
				return
			}

			if len(clients) != tt.wantLen {
				t.Errorf("ListClients() len = %v, want %v", len(clients), tt.wantLen)
			}
		})
	}
}

func TestOAuth2Service_GetClient(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	testClient, _ := domain.NewOAuthClient("client-123", "secret123", "Test Client", "Test Description", []string{"read"})

	tests := []struct {
		name        string
		clientID    string
		getByIDFunc func(ctx context.Context, id string) (*domain.OAuthClient, error)
		wantErr     bool
	}{
		{
			name:     "client found",
			clientID: "client-123",
			getByIDFunc: func(ctx context.Context, id string) (*domain.OAuthClient, error) {
				return testClient, nil
			},
			wantErr: false,
		},
		{
			name:     "client not found",
			clientID: "nonexistent",
			getByIDFunc: func(ctx context.Context, id string) (*domain.OAuthClient, error) {
				return nil, domainerrors.ErrClientNotFound
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClientRepo := &MockOAuthClientRepository{
				GetByIDFunc: tt.getByIDFunc,
			}
			oauth2Service := services.NewOAuth2Service(mockClientRepo, "test-secret-key-at-least-32-chars-long", 15*time.Minute, logger)

			client, err := oauth2Service.GetClient(context.Background(), tt.clientID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetClient() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GetClient() unexpected error: %v", err)
				return
			}

			if client == nil {
				t.Errorf("GetClient() returned nil client")
			}
		})
	}
}

func TestOAuth2Service_DeleteClient(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name       string
		clientID   string
		deleteFunc func(ctx context.Context, id string) error
		wantErr    bool
	}{
		{
			name:     "successful deletion",
			clientID: "client-123",
			deleteFunc: func(ctx context.Context, id string) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:     "deletion error",
			clientID: "client-123",
			deleteFunc: func(ctx context.Context, id string) error {
				return domainerrors.ErrInternal
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClientRepo := &MockOAuthClientRepository{
				DeleteFunc: tt.deleteFunc,
			}
			oauth2Service := services.NewOAuth2Service(mockClientRepo, "test-secret-key-at-least-32-chars-long", 15*time.Minute, logger)

			err := oauth2Service.DeleteClient(context.Background(), tt.clientID)

			if tt.wantErr && err == nil {
				t.Errorf("DeleteClient() expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("DeleteClient() unexpected error: %v", err)
			}
		})
	}
}

