package tests

import (
	"encoding/json"
	"testing"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
)

func TestLoginRequest_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    request.LoginRequest
		wantErr bool
	}{
		{
			name:  "valid login request",
			input: `{"email":"test@example.com","password":"password123"}`,
			want: request.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name:  "valid login request with extra fields",
			input: `{"email":"user@test.com","password":"secure123","extra":"ignored"}`,
			want: request.LoginRequest{
				Email:    "user@test.com",
				Password: "secure123",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"email":"test@example.com","password":}`,
			want:    request.LoginRequest{},
			wantErr: true,
		},
		{
			name:  "empty fields",
			input: `{"email":"","password":""}`,
			want: request.LoginRequest{
				Email:    "",
				Password: "",
			},
			wantErr: false,
		},
		{
			name:  "missing email field",
			input: `{"password":"password123"}`,
			want: request.LoginRequest{
				Email:    "",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name:  "missing password field",
			input: `{"email":"test@example.com"}`,
			want: request.LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got request.LoginRequest
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Email != tt.want.Email {
					t.Errorf("LoginRequest.Email = %v, want %v", got.Email, tt.want.Email)
				}
				if got.Password != tt.want.Password {
					t.Errorf("LoginRequest.Password = %v, want %v", got.Password, tt.want.Password)
				}
			}
		})
	}
}

func TestLoginRequest_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		request request.LoginRequest
		want    string
	}{
		{
			name: "marshal valid login request",
			request: request.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			want: `{"email":"test@example.com","password":"password123"}`,
		},
		{
			name: "marshal with empty fields",
			request: request.LoginRequest{
				Email:    "",
				Password: "",
			},
			want: `{"email":"","password":""}`,
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

func TestLoginRequest_Fields(t *testing.T) {
	req := request.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	if req.Email != "test@example.com" {
		t.Errorf("LoginRequest.Email = %v, want test@example.com", req.Email)
	}

	if req.Password != "password123" {
		t.Errorf("LoginRequest.Password = %v, want password123", req.Password)
	}
}
