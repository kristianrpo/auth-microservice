package tests

import (
	"encoding/json"
	"testing"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
)

func TestClientCredentialsRequest_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    request.ClientCredentialsRequest
		wantErr bool
	}{
		{
			name:  "valid client credentials request",
			input: `{"client_id":"test_client","client_secret":"secret123","grant_type":"client_credentials"}`,
			want: request.ClientCredentialsRequest{
				ClientID:     "test_client",
				ClientSecret: "secret123",
				GrantType:    "client_credentials",
			},
			wantErr: false,
		},
		{
			name:  "valid with different client",
			input: `{"client_id":"another_client","client_secret":"anothersecret","grant_type":"client_credentials"}`,
			want: request.ClientCredentialsRequest{
				ClientID:     "another_client",
				ClientSecret: "anothersecret",
				GrantType:    "client_credentials",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"client_id":"test_client","client_secret":}`,
			want:    request.ClientCredentialsRequest{},
			wantErr: true,
		},
		{
			name:  "empty fields",
			input: `{"client_id":"","client_secret":"","grant_type":""}`,
			want: request.ClientCredentialsRequest{
				ClientID:     "",
				ClientSecret: "",
				GrantType:    "",
			},
			wantErr: false,
		},
		{
			name:  "missing client_id",
			input: `{"client_secret":"secret123","grant_type":"client_credentials"}`,
			want: request.ClientCredentialsRequest{
				ClientID:     "",
				ClientSecret: "secret123",
				GrantType:    "client_credentials",
			},
			wantErr: false,
		},
		{
			name:  "missing client_secret",
			input: `{"client_id":"test_client","grant_type":"client_credentials"}`,
			want: request.ClientCredentialsRequest{
				ClientID:     "test_client",
				ClientSecret: "",
				GrantType:    "client_credentials",
			},
			wantErr: false,
		},
		{
			name:  "missing grant_type",
			input: `{"client_id":"test_client","client_secret":"secret123"}`,
			want: request.ClientCredentialsRequest{
				ClientID:     "test_client",
				ClientSecret: "secret123",
				GrantType:    "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got request.ClientCredentialsRequest
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.ClientID != tt.want.ClientID {
					t.Errorf("ClientCredentialsRequest.ClientID = %v, want %v", got.ClientID, tt.want.ClientID)
				}
				if got.ClientSecret != tt.want.ClientSecret {
					t.Errorf("ClientCredentialsRequest.ClientSecret = %v, want %v", got.ClientSecret, tt.want.ClientSecret)
				}
				if got.GrantType != tt.want.GrantType {
					t.Errorf("ClientCredentialsRequest.GrantType = %v, want %v", got.GrantType, tt.want.GrantType)
				}
			}
		})
	}
}

func TestClientCredentialsRequest_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		request request.ClientCredentialsRequest
		want    string
	}{
		{
			name: "marshal valid request",
			request: request.ClientCredentialsRequest{
				ClientID:     "test_client",
				ClientSecret: "secret123",
				GrantType:    "client_credentials",
			},
			want: `{"client_id":"test_client","client_secret":"secret123","grant_type":"client_credentials"}`,
		},
		{
			name: "marshal with empty fields",
			request: request.ClientCredentialsRequest{
				ClientID:     "",
				ClientSecret: "",
				GrantType:    "",
			},
			want: `{"client_id":"","client_secret":"","grant_type":""}`,
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

func TestClientCredentialsRequest_Fields(t *testing.T) {
	req := request.ClientCredentialsRequest{
		ClientID:     "test_client",
		ClientSecret: "secret123",
		GrantType:    "client_credentials",
	}

	if req.ClientID != "test_client" {
		t.Errorf("ClientCredentialsRequest.ClientID = %v, want test_client", req.ClientID)
	}

	if req.ClientSecret != "secret123" {
		t.Errorf("ClientCredentialsRequest.ClientSecret = %v, want secret123", req.ClientSecret)
	}

	if req.GrantType != "client_credentials" {
		t.Errorf("ClientCredentialsRequest.GrantType = %v, want client_credentials", req.GrantType)
	}
}
