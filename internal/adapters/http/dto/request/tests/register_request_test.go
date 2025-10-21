package tests

import (
	"encoding/json"
	"testing"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/request"
)

func TestRegisterRequest_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    request.RegisterRequest
		wantErr bool
	}{
		{
			name:  "valid register request",
			input: `{"id_citizen":12345,"email":"test@example.com","password":"password123","name":"Test User"}`,
			want: request.RegisterRequest{
				IDCitizen: 12345,
				Email:     "test@example.com",
				Password:  "password123",
				Name:      "Test User",
			},
			wantErr: false,
		},
		{
			name:  "valid register with minimum password",
			input: `{"id_citizen":54321,"email":"user@test.com","password":"pass1234","name":"User"}`,
			want: request.RegisterRequest{
				IDCitizen: 54321,
				Email:     "user@test.com",
				Password:  "pass1234",
				Name:      "User",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"id_citizen":12345,"email":"test@example.com","password":}`,
			want:    request.RegisterRequest{},
			wantErr: true,
		},
		{
			name:  "empty fields",
			input: `{"id_citizen":0,"email":"","password":"","name":""}`,
			want: request.RegisterRequest{
				IDCitizen: 0,
				Email:     "",
				Password:  "",
				Name:      "",
			},
			wantErr: false,
		},
		{
			name:  "missing optional fields",
			input: `{"id_citizen":12345,"email":"test@example.com","password":"password123","name":"Test"}`,
			want: request.RegisterRequest{
				IDCitizen: 12345,
				Email:     "test@example.com",
				Password:  "password123",
				Name:      "Test",
			},
			wantErr: false,
		},
		{
			name:  "negative id_citizen",
			input: `{"id_citizen":-1,"email":"test@example.com","password":"password123","name":"Test User"}`,
			want: request.RegisterRequest{
				IDCitizen: -1,
				Email:     "test@example.com",
				Password:  "password123",
				Name:      "Test User",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got request.RegisterRequest
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.IDCitizen != tt.want.IDCitizen {
					t.Errorf("RegisterRequest.IDCitizen = %v, want %v", got.IDCitizen, tt.want.IDCitizen)
				}
				if got.Email != tt.want.Email {
					t.Errorf("RegisterRequest.Email = %v, want %v", got.Email, tt.want.Email)
				}
				if got.Password != tt.want.Password {
					t.Errorf("RegisterRequest.Password = %v, want %v", got.Password, tt.want.Password)
				}
				if got.Name != tt.want.Name {
					t.Errorf("RegisterRequest.Name = %v, want %v", got.Name, tt.want.Name)
				}
			}
		})
	}
}

func TestRegisterRequest_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		request request.RegisterRequest
		want    string
	}{
		{
			name: "marshal valid register request",
			request: request.RegisterRequest{
				IDCitizen: 12345,
				Email:     "test@example.com",
				Password:  "password123",
				Name:      "Test User",
			},
			want: `{"id_citizen":12345,"email":"test@example.com","password":"password123","name":"Test User"}`,
		},
		{
			name: "marshal with empty strings",
			request: request.RegisterRequest{
				IDCitizen: 0,
				Email:     "",
				Password:  "",
				Name:      "",
			},
			want: `{"id_citizen":0,"email":"","password":"","name":""}`,
		},
		{
			name: "marshal with special characters in name",
			request: request.RegisterRequest{
				IDCitizen: 99999,
				Email:     "special@test.com",
				Password:  "secure123",
				Name:      "José María O'Brien",
			},
			want: `{"id_citizen":99999,"email":"special@test.com","password":"secure123","name":"José María O'Brien"}`,
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

func TestRegisterRequest_Fields(t *testing.T) {
	req := request.RegisterRequest{
		IDCitizen: 12345,
		Email:     "test@example.com",
		Password:  "password123",
		Name:      "Test User",
	}

	if req.IDCitizen != 12345 {
		t.Errorf("RegisterRequest.IDCitizen = %v, want 12345", req.IDCitizen)
	}

	if req.Email != "test@example.com" {
		t.Errorf("RegisterRequest.Email = %v, want test@example.com", req.Email)
	}

	if req.Password != "password123" {
		t.Errorf("RegisterRequest.Password = %v, want password123", req.Password)
	}

	if req.Name != "Test User" {
		t.Errorf("RegisterRequest.Name = %v, want Test User", req.Name)
	}
}
