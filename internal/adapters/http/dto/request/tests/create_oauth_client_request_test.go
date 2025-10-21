package tests

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
)

func TestCreateOAuthClientRequest_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    request.CreateOAuthClientRequest
		wantErr bool
	}{
		{
			name:  "valid create oauth client request",
			input: `{"client_id":"test_client","client_secret":"secret123","name":"Test Client","description":"A test client","scopes":["read","write"]}`,
			want: request.CreateOAuthClientRequest{
				ClientID:     "test_client",
				ClientSecret: "secret123",
				Name:         "Test Client",
				Description:  "A test client",
				Scopes:       []string{"read", "write"},
			},
			wantErr: false,
		},
		{
			name:  "valid without optional fields",
			input: `{"client_id":"minimal_client","client_secret":"secret456","name":"Minimal"}`,
			want: request.CreateOAuthClientRequest{
				ClientID:     "minimal_client",
				ClientSecret: "secret456",
				Name:         "Minimal",
				Description:  "",
				Scopes:       nil,
			},
			wantErr: false,
		},
		{
			name:  "valid with empty scopes array",
			input: `{"client_id":"empty_scopes","client_secret":"secret789","name":"Empty Scopes","description":"No scopes","scopes":[]}`,
			want: request.CreateOAuthClientRequest{
				ClientID:     "empty_scopes",
				ClientSecret: "secret789",
				Name:         "Empty Scopes",
				Description:  "No scopes",
				Scopes:       []string{},
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"client_id":"test_client","client_secret":}`,
			want:    request.CreateOAuthClientRequest{},
			wantErr: true,
		},
		{
			name:  "empty required fields",
			input: `{"client_id":"","client_secret":"","name":""}`,
			want: request.CreateOAuthClientRequest{
				ClientID:     "",
				ClientSecret: "",
				Name:         "",
				Description:  "",
				Scopes:       nil,
			},
			wantErr: false,
		},
		{
			name:  "valid with single scope",
			input: `{"client_id":"single_scope","client_secret":"secret999","name":"Single","scopes":["admin"]}`,
			want: request.CreateOAuthClientRequest{
				ClientID:     "single_scope",
				ClientSecret: "secret999",
				Name:         "Single",
				Description:  "",
				Scopes:       []string{"admin"},
			},
			wantErr: false,
		},
		{
			name:  "valid with multiple scopes",
			input: `{"client_id":"multi_scope","client_secret":"secret111","name":"Multi","scopes":["read","write","admin","delete"]}`,
			want: request.CreateOAuthClientRequest{
				ClientID:     "multi_scope",
				ClientSecret: "secret111",
				Name:         "Multi",
				Description:  "",
				Scopes:       []string{"read", "write", "admin", "delete"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got request.CreateOAuthClientRequest
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.ClientID != tt.want.ClientID {
					t.Errorf("CreateOAuthClientRequest.ClientID = %v, want %v", got.ClientID, tt.want.ClientID)
				}
				if got.ClientSecret != tt.want.ClientSecret {
					t.Errorf("CreateOAuthClientRequest.ClientSecret = %v, want %v", got.ClientSecret, tt.want.ClientSecret)
				}
				if got.Name != tt.want.Name {
					t.Errorf("CreateOAuthClientRequest.Name = %v, want %v", got.Name, tt.want.Name)
				}
				if got.Description != tt.want.Description {
					t.Errorf("CreateOAuthClientRequest.Description = %v, want %v", got.Description, tt.want.Description)
				}
				if !reflect.DeepEqual(got.Scopes, tt.want.Scopes) {
					t.Errorf("CreateOAuthClientRequest.Scopes = %v, want %v", got.Scopes, tt.want.Scopes)
				}
			}
		})
	}
}

func TestCreateOAuthClientRequest_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		request request.CreateOAuthClientRequest
		want    string
	}{
		{
			name: "marshal complete request",
			request: request.CreateOAuthClientRequest{
				ClientID:     "test_client",
				ClientSecret: "secret123",
				Name:         "Test Client",
				Description:  "A test client",
				Scopes:       []string{"read", "write"},
			},
			want: `{"client_id":"test_client","client_secret":"secret123","name":"Test Client","description":"A test client","scopes":["read","write"]}`,
		},
		{
			name: "marshal minimal request",
			request: request.CreateOAuthClientRequest{
				ClientID:     "minimal",
				ClientSecret: "secret",
				Name:         "Min",
				Description:  "",
				Scopes:       nil,
			},
			want: `{"client_id":"minimal","client_secret":"secret","name":"Min","description":"","scopes":null}`,
		},
		{
			name: "marshal with empty scopes array",
			request: request.CreateOAuthClientRequest{
				ClientID:     "empty",
				ClientSecret: "secret",
				Name:         "Empty",
				Description:  "",
				Scopes:       []string{},
			},
			want: `{"client_id":"empty","client_secret":"secret","name":"Empty","description":"","scopes":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.request)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			if string(got) != tt.want {
				t.Errorf("json.Marshal() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestCreateOAuthClientRequest_Fields(t *testing.T) {
	req := request.CreateOAuthClientRequest{
		ClientID:     "test_client",
		ClientSecret: "secret123",
		Name:         "Test Client",
		Description:  "A test client",
		Scopes:       []string{"read", "write"},
	}

	if req.ClientID != "test_client" {
		t.Errorf("CreateOAuthClientRequest.ClientID = %v, want test_client", req.ClientID)
	}

	if req.ClientSecret != "secret123" {
		t.Errorf("CreateOAuthClientRequest.ClientSecret = %v, want secret123", req.ClientSecret)
	}

	if req.Name != "Test Client" {
		t.Errorf("CreateOAuthClientRequest.Name = %v, want Test Client", req.Name)
	}

	if req.Description != "A test client" {
		t.Errorf("CreateOAuthClientRequest.Description = %v, want A test client", req.Description)
	}

	expectedScopes := []string{"read", "write"}
	if !reflect.DeepEqual(req.Scopes, expectedScopes) {
		t.Errorf("CreateOAuthClientRequest.Scopes = %v, want %v", req.Scopes, expectedScopes)
	}
}

