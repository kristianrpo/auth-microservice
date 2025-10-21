package tests

import (
	"encoding/json"
	"testing"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
)

func TestLogoutRequest_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    request.LogoutRequest
		wantErr bool
	}{
		{
			name:  "valid logout request",
			input: `{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"}`,
			want: request.LogoutRequest{
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
			wantErr: false,
		},
		{
			name:  "logout request with empty token",
			input: `{"refresh_token":""}`,
			want: request.LogoutRequest{
				RefreshToken: "",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"refresh_token":}`,
			want:    request.LogoutRequest{},
			wantErr: true,
		},
		{
			name:  "missing refresh_token field",
			input: `{}`,
			want: request.LogoutRequest{
				RefreshToken: "",
			},
			wantErr: false,
		},
		{
			name:  "logout request with null token",
			input: `{"refresh_token":null}`,
			want: request.LogoutRequest{
				RefreshToken: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got request.LogoutRequest
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.RefreshToken != tt.want.RefreshToken {
					t.Errorf("LogoutRequest.RefreshToken = %v, want %v", got.RefreshToken, tt.want.RefreshToken)
				}
			}
		})
	}
}

func TestLogoutRequest_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		request request.LogoutRequest
		want    string
	}{
		{
			name: "marshal valid logout request",
			request: request.LogoutRequest{
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
			want: `{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"}`,
		},
		{
			name: "marshal with empty token",
			request: request.LogoutRequest{
				RefreshToken: "",
			},
			want: `{"refresh_token":""}`,
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

func TestLogoutRequest_Fields(t *testing.T) {
	req := request.LogoutRequest{
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
	}

	if req.RefreshToken != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" {
		t.Errorf("LogoutRequest.RefreshToken = %v, want eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", req.RefreshToken)
	}
}
