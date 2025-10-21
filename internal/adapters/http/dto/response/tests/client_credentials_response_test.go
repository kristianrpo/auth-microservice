package tests

import (
	"encoding/json"
	"testing"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
)

func TestClientCredentialsResponse_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    response.ClientCredentialsResponse
		wantErr bool
	}{
		{
			name:  "valid client credentials response",
			input: `{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9","token_type":"Bearer","expires_in":3600}`,
			want: response.ClientCredentialsResponse{
				AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			},
			wantErr: false,
		},
		{
			name:  "valid with different expiry",
			input: `{"access_token":"access_token_string","token_type":"Bearer","expires_in":7200}`,
			want: response.ClientCredentialsResponse{
				AccessToken: "access_token_string",
				TokenType:   "Bearer",
				ExpiresIn:   7200,
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"access_token":"token","token_type":}`,
			want:    response.ClientCredentialsResponse{},
			wantErr: true,
		},
		{
			name:  "empty fields",
			input: `{"access_token":"","token_type":"","expires_in":0}`,
			want: response.ClientCredentialsResponse{
				AccessToken: "",
				TokenType:   "",
				ExpiresIn:   0,
			},
			wantErr: false,
		},
		{
			name:  "zero expires_in",
			input: `{"access_token":"token","token_type":"Bearer","expires_in":0}`,
			want: response.ClientCredentialsResponse{
				AccessToken: "token",
				TokenType:   "Bearer",
				ExpiresIn:   0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got response.ClientCredentialsResponse
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.AccessToken != tt.want.AccessToken {
					t.Errorf("ClientCredentialsResponse.AccessToken = %v, want %v", got.AccessToken, tt.want.AccessToken)
				}
				if got.TokenType != tt.want.TokenType {
					t.Errorf("ClientCredentialsResponse.TokenType = %v, want %v", got.TokenType, tt.want.TokenType)
				}
				if got.ExpiresIn != tt.want.ExpiresIn {
					t.Errorf("ClientCredentialsResponse.ExpiresIn = %v, want %v", got.ExpiresIn, tt.want.ExpiresIn)
				}
			}
		})
	}
}

func TestClientCredentialsResponse_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		response response.ClientCredentialsResponse
		want     string
	}{
		{
			name: "marshal valid response",
			response: response.ClientCredentialsResponse{
				AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			},
			want: `{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9","token_type":"Bearer","expires_in":3600}`,
		},
		{
			name: "marshal with empty fields",
			response: response.ClientCredentialsResponse{
				AccessToken: "",
				TokenType:   "",
				ExpiresIn:   0,
			},
			want: `{"access_token":"","token_type":"","expires_in":0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.response)
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

func TestClientCredentialsResponse_Fields(t *testing.T) {
	resp := response.ClientCredentialsResponse{
		AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	if resp.AccessToken != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" {
		t.Errorf("ClientCredentialsResponse.AccessToken = %v, want eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", resp.AccessToken)
	}

	if resp.TokenType != "Bearer" {
		t.Errorf("ClientCredentialsResponse.TokenType = %v, want Bearer", resp.TokenType)
	}

	if resp.ExpiresIn != 3600 {
		t.Errorf("ClientCredentialsResponse.ExpiresIn = %v, want 3600", resp.ExpiresIn)
	}
}
