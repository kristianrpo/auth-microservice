package tests

import (
	"encoding/json"
	"testing"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
)

func TestRefreshTokenRequest_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    request.RefreshTokenRequest
		wantErr bool
	}{
		{
			name:  "valid refresh token request",
			input: `{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"}`,
			want: request.RefreshTokenRequest{
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
			wantErr: false,
		},
		{
			name:  "valid refresh token with long token",
			input: `{"refresh_token":"verylongtokenstringverylongtokenstringverylongtokenstringverylongtokenstring"}`,
			want: request.RefreshTokenRequest{
				RefreshToken: "verylongtokenstringverylongtokenstringverylongtokenstringverylongtokenstring",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"refresh_token":}`,
			want:    request.RefreshTokenRequest{},
			wantErr: true,
		},
		{
			name:  "empty refresh token",
			input: `{"refresh_token":""}`,
			want: request.RefreshTokenRequest{
				RefreshToken: "",
			},
			wantErr: false,
		},
		{
			name:  "missing refresh_token field",
			input: `{}`,
			want: request.RefreshTokenRequest{
				RefreshToken: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got request.RefreshTokenRequest
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.RefreshToken != tt.want.RefreshToken {
					t.Errorf("RefreshTokenRequest.RefreshToken = %v, want %v", got.RefreshToken, tt.want.RefreshToken)
				}
			}
		})
	}
}

func TestRefreshTokenRequest_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		request request.RefreshTokenRequest
		want    string
	}{
		{
			name: "marshal valid refresh token request",
			request: request.RefreshTokenRequest{
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
			want: `{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"}`,
		},
		{
			name: "marshal with empty token",
			request: request.RefreshTokenRequest{
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

func TestRefreshTokenRequest_Fields(t *testing.T) {
	req := request.RefreshTokenRequest{
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
	}

	if req.RefreshToken != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" {
		t.Errorf("RefreshTokenRequest.RefreshToken = %v, want eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", req.RefreshToken)
	}
}

