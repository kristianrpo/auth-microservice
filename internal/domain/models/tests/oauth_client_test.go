package tests

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

func TestNewOAuthClient(t *testing.T) {
	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		clientName   string
		description  string
		scopes       []string
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid client",
			clientID:     "client123",
			clientSecret: "secret123",
			clientName:   "Test Client",
			description:  "Test Description",
			scopes:       []string{"read", "write"},
			wantErr:      false,
		},
		{
			name:         "empty client_id",
			clientID:     "",
			clientSecret: "secret123",
			clientName:   "Test Client",
			description:  "Test Description",
			scopes:       []string{"read"},
			wantErr:      true,
			errMsg:       "validation error",
		},
		{
			name:         "empty client_secret",
			clientID:     "client123",
			clientSecret: "",
			clientName:   "Test Client",
			description:  "Test Description",
			scopes:       []string{"read"},
			wantErr:      true,
			errMsg:       "validation error",
		},
		{
			name:         "empty name",
			clientID:     "client123",
			clientSecret: "secret123",
			clientName:   "",
			description:  "Test Description",
			scopes:       []string{"read"},
			wantErr:      true,
			errMsg:       "validation error",
		},
		{
			name:         "empty scopes",
			clientID:     "client123",
			clientSecret: "secret123",
			clientName:   "Test Client",
			description:  "Test Description",
			scopes:       []string{},
			wantErr:      false,
		},
		{
			name:         "nil scopes",
			clientID:     "client123",
			clientSecret: "secret123",
			clientName:   "Test Client",
			description:  "Test Description",
			scopes:       nil,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := domain.NewOAuthClient(tt.clientID, tt.clientSecret, tt.clientName, tt.description, tt.scopes)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewOAuthClient() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("NewOAuthClient() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("NewOAuthClient() unexpected error: %v", err)
				return
			}

			if client.ClientID != tt.clientID {
				t.Errorf("NewOAuthClient() ClientID = %v, want %v", client.ClientID, tt.clientID)
			}

			if client.Name != tt.clientName {
				t.Errorf("NewOAuthClient() Name = %v, want %v", client.Name, tt.clientName)
			}

			if client.Description != tt.description {
				t.Errorf("NewOAuthClient() Description = %v, want %v", client.Description, tt.description)
			}

			if !client.Active {
				t.Errorf("NewOAuthClient() Active = %v, want true", client.Active)
			}

			// Verify secret is hashed
			if err := bcrypt.CompareHashAndPassword([]byte(client.ClientSecret), []byte(tt.clientSecret)); err != nil {
				t.Errorf("NewOAuthClient() secret not properly hashed")
			}

			// Verify scopes
			if tt.scopes == nil && len(client.Scopes) != 0 {
				t.Errorf("NewOAuthClient() Scopes = %v, want empty", client.Scopes)
			}
		})
	}
}

func TestOAuthClient_ValidateSecret(t *testing.T) {
	client, _ := domain.NewOAuthClient("client123", "secret123", "Test Client", "Description", []string{"read"})

	tests := []struct {
		name   string
		secret string
		want   bool
	}{
		{
			name:   "correct secret",
			secret: "secret123",
			want:   true,
		},
		{
			name:   "incorrect secret",
			secret: "wrongsecret",
			want:   false,
		},
		{
			name:   "empty secret",
			secret: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := client.ValidateSecret(tt.secret)

			if valid != tt.want {
				t.Errorf("ValidateSecret() = %v, want %v", valid, tt.want)
			}
		})
	}
}
