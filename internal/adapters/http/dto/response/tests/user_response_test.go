package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

func TestUserResponse_JSON(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	testTimeStr := testTime.Format(time.RFC3339)

	tests := []struct {
		name    string
		input   string
		want    response.UserResponse
		wantErr bool
	}{
		{
			name:  "valid user response with user role",
			input: `{"id":"123e4567-e89b-12d3-a456-426614174000","id_citizen":12345,"email":"test@example.com","name":"Test User","role":"USER","created_at":"` + testTimeStr + `","updated_at":"` + testTimeStr + `"}`,
			want: response.UserResponse{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				IDCitizen: 12345,
				Email:     "test@example.com",
				Name:      "Test User",
				Role:      domain.RoleUser,
				CreatedAt: testTime,
				UpdatedAt: testTime,
			},
			wantErr: false,
		},
		{
			name:  "valid user response with admin role",
			input: `{"id":"123e4567-e89b-12d3-a456-426614174001","id_citizen":54321,"email":"admin@example.com","name":"Admin User","role":"ADMIN","created_at":"` + testTimeStr + `","updated_at":"` + testTimeStr + `"}`,
			want: response.UserResponse{
				ID:        "123e4567-e89b-12d3-a456-426614174001",
				IDCitizen: 54321,
				Email:     "admin@example.com",
				Name:      "Admin User",
				Role:      domain.RoleAdmin,
				CreatedAt: testTime,
				UpdatedAt: testTime,
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"id":"123","email":"test@example.com","name":}`,
			want:    response.UserResponse{},
			wantErr: true,
		},
		{
			name:  "empty fields",
			input: `{"id":"","id_citizen":0,"email":"","name":"","role":"","created_at":"` + testTimeStr + `","updated_at":"` + testTimeStr + `"}`,
			want: response.UserResponse{
				ID:        "",
				IDCitizen: 0,
				Email:     "",
				Name:      "",
				Role:      "",
				CreatedAt: testTime,
				UpdatedAt: testTime,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got response.UserResponse
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.ID != tt.want.ID {
					t.Errorf("UserResponse.ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.IDCitizen != tt.want.IDCitizen {
					t.Errorf("UserResponse.IDCitizen = %v, want %v", got.IDCitizen, tt.want.IDCitizen)
				}
				if got.Email != tt.want.Email {
					t.Errorf("UserResponse.Email = %v, want %v", got.Email, tt.want.Email)
				}
				if got.Name != tt.want.Name {
					t.Errorf("UserResponse.Name = %v, want %v", got.Name, tt.want.Name)
				}
				if got.Role != tt.want.Role {
					t.Errorf("UserResponse.Role = %v, want %v", got.Role, tt.want.Role)
				}
				if !got.CreatedAt.Equal(tt.want.CreatedAt) {
					t.Errorf("UserResponse.CreatedAt = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
				}
				if !got.UpdatedAt.Equal(tt.want.UpdatedAt) {
					t.Errorf("UserResponse.UpdatedAt = %v, want %v", got.UpdatedAt, tt.want.UpdatedAt)
				}
			}
		})
	}
}

func TestUserResponse_Marshal(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	testTimeStr := testTime.Format(time.RFC3339)

	tests := []struct {
		name     string
		response response.UserResponse
		want     string
	}{
		{
			name: "marshal user with user role",
			response: response.UserResponse{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				IDCitizen: 12345,
				Email:     "test@example.com",
				Name:      "Test User",
				Role:      domain.RoleUser,
				CreatedAt: testTime,
				UpdatedAt: testTime,
			},
			want: `{"id":"123e4567-e89b-12d3-a456-426614174000","id_citizen":12345,"email":"test@example.com","name":"Test User","role":"USER","created_at":"` + testTimeStr + `","updated_at":"` + testTimeStr + `"}`,
		},
		{
			name: "marshal user with admin role",
			response: response.UserResponse{
				ID:        "123e4567-e89b-12d3-a456-426614174001",
				IDCitizen: 54321,
				Email:     "admin@example.com",
				Name:      "Admin User",
				Role:      domain.RoleAdmin,
				CreatedAt: testTime,
				UpdatedAt: testTime,
			},
			want: `{"id":"123e4567-e89b-12d3-a456-426614174001","id_citizen":54321,"email":"admin@example.com","name":"Admin User","role":"ADMIN","created_at":"` + testTimeStr + `","updated_at":"` + testTimeStr + `"}`,
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

func TestUserResponse_Fields(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	resp := response.UserResponse{
		ID:        "123e4567-e89b-12d3-a456-426614174000",
		IDCitizen: 12345,
		Email:     "test@example.com",
		Name:      "Test User",
		Role:      domain.RoleUser,
		CreatedAt: testTime,
		UpdatedAt: testTime,
	}

	if resp.ID != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("UserResponse.ID = %v, want 123e4567-e89b-12d3-a456-426614174000", resp.ID)
	}

	if resp.IDCitizen != 12345 {
		t.Errorf("UserResponse.IDCitizen = %v, want 12345", resp.IDCitizen)
	}

	if resp.Email != "test@example.com" {
		t.Errorf("UserResponse.Email = %v, want test@example.com", resp.Email)
	}

	if resp.Name != "Test User" {
		t.Errorf("UserResponse.Name = %v, want Test User", resp.Name)
	}

	if resp.Role != domain.RoleUser {
		t.Errorf("UserResponse.Role = %v, want USER", resp.Role)
	}

	if !resp.CreatedAt.Equal(testTime) {
		t.Errorf("UserResponse.CreatedAt = %v, want %v", resp.CreatedAt, testTime)
	}

	if !resp.UpdatedAt.Equal(testTime) {
		t.Errorf("UserResponse.UpdatedAt = %v, want %v", resp.UpdatedAt, testTime)
	}
}
