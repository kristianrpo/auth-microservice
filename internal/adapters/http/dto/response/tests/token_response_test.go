package tests

import (
	"encoding/json"
	"testing"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
)

func TestTokenResponse_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    response.TokenResponse
		wantErr bool
	}{
		{
			name:  "valid token response",
			input: `{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9","refresh_token":"refresh_token_string","token_type":"Bearer","expires_in":3600}`,
			want: response.TokenResponse{
				AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				RefreshToken: "refresh_token_string",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
			},
			wantErr: false,
		},
		{
			name:  "valid token response with different expiry",
			input: `{"access_token":"access_token","refresh_token":"refresh_token","token_type":"Bearer","expires_in":7200}`,
			want: response.TokenResponse{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
				TokenType:    "Bearer",
				ExpiresIn:    7200,
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"access_token":"token","refresh_token":}`,
			want:    response.TokenResponse{},
			wantErr: true,
		},
		{
			name:  "empty fields",
			input: `{"access_token":"","refresh_token":"","token_type":"","expires_in":0}`,
			want: response.TokenResponse{
				AccessToken:  "",
				RefreshToken: "",
				TokenType:    "",
				ExpiresIn:    0,
			},
			wantErr: false,
		},
		{
			name:  "negative expires_in",
			input: `{"access_token":"token","refresh_token":"refresh","token_type":"Bearer","expires_in":-1}`,
			want: response.TokenResponse{
				AccessToken:  "token",
				RefreshToken: "refresh",
				TokenType:    "Bearer",
				ExpiresIn:    -1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got response.TokenResponse
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.AccessToken != tt.want.AccessToken {
					t.Errorf("TokenResponse.AccessToken = %v, want %v", got.AccessToken, tt.want.AccessToken)
				}
				if got.RefreshToken != tt.want.RefreshToken {
					t.Errorf("TokenResponse.RefreshToken = %v, want %v", got.RefreshToken, tt.want.RefreshToken)
				}
				if got.TokenType != tt.want.TokenType {
					t.Errorf("TokenResponse.TokenType = %v, want %v", got.TokenType, tt.want.TokenType)
				}
				if got.ExpiresIn != tt.want.ExpiresIn {
					t.Errorf("TokenResponse.ExpiresIn = %v, want %v", got.ExpiresIn, tt.want.ExpiresIn)
				}
			}
		})
	}
}

func TestTokenResponse_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		response response.TokenResponse
		want     string
	}{
		{
			name: "marshal valid token response",
			response: response.TokenResponse{
				AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				RefreshToken: "refresh_token_string",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
			},
			want: `{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9","refresh_token":"refresh_token_string","token_type":"Bearer","expires_in":3600}`,
		},
		{
			name: "marshal with empty tokens",
			response: response.TokenResponse{
				AccessToken:  "",
				RefreshToken: "",
				TokenType:    "",
				ExpiresIn:    0,
			},
			want: `{"access_token":"","refresh_token":"","token_type":"","expires_in":0}`,
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

func TestTokenResponse_Fields(t *testing.T) {
	resp := response.TokenResponse{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		RefreshToken: "refresh_token_string",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	if resp.AccessToken != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" {
		t.Errorf("TokenResponse.AccessToken = %v, want eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", resp.AccessToken)
	}

	if resp.RefreshToken != "refresh_token_string" {
		t.Errorf("TokenResponse.RefreshToken = %v, want refresh_token_string", resp.RefreshToken)
	}

	if resp.TokenType != "Bearer" {
		t.Errorf("TokenResponse.TokenType = %v, want Bearer", resp.TokenType)
	}

	if resp.ExpiresIn != 3600 {
		t.Errorf("TokenResponse.ExpiresIn = %v, want 3600", resp.ExpiresIn)
	}
}

