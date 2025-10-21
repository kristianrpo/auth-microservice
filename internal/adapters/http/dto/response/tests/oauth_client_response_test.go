package tests

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
)

func TestOAuthClientResponse_JSON(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	testTimeStr := testTime.Format(time.RFC3339)

	tests := []struct {
		name    string
		input   string
		want    response.OAuthClientResponse
		wantErr bool
	}{
		{
			name:  "valid oauth client response",
			input: `{"id":"123e4567-e89b-12d3-a456-426614174000","client_id":"test_client","name":"Test Client","description":"A test client","scopes":["read","write"],"active":true,"created_at":"` + testTimeStr + `","updated_at":"` + testTimeStr + `"}`,
			want: response.OAuthClientResponse{
				ID:          "123e4567-e89b-12d3-a456-426614174000",
				ClientID:    "test_client",
				Name:        "Test Client",
				Description: "A test client",
				Scopes:      []string{"read", "write"},
				Active:      true,
				CreatedAt:   testTime,
				UpdatedAt:   testTime,
			},
			wantErr: false,
		},
		{
			name:  "valid oauth client response without scopes",
			input: `{"id":"123e4567-e89b-12d3-a456-426614174001","client_id":"minimal_client","name":"Minimal","description":"","scopes":null,"active":false,"created_at":"` + testTimeStr + `","updated_at":"` + testTimeStr + `"}`,
			want: response.OAuthClientResponse{
				ID:          "123e4567-e89b-12d3-a456-426614174001",
				ClientID:    "minimal_client",
				Name:        "Minimal",
				Description: "",
				Scopes:      nil,
				Active:      false,
				CreatedAt:   testTime,
				UpdatedAt:   testTime,
			},
			wantErr: false,
		},
		{
			name:  "valid oauth client response with empty scopes",
			input: `{"id":"123e4567-e89b-12d3-a456-426614174002","client_id":"empty_scopes","name":"Empty","description":"","scopes":[],"active":true,"created_at":"` + testTimeStr + `","updated_at":"` + testTimeStr + `"}`,
			want: response.OAuthClientResponse{
				ID:          "123e4567-e89b-12d3-a456-426614174002",
				ClientID:    "empty_scopes",
				Name:        "Empty",
				Description: "",
				Scopes:      []string{},
				Active:      true,
				CreatedAt:   testTime,
				UpdatedAt:   testTime,
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"id":"123","client_id":"test","name":}`,
			want:    response.OAuthClientResponse{},
			wantErr: true,
		},
		{
			name:  "valid with multiple scopes",
			input: `{"id":"123e4567-e89b-12d3-a456-426614174003","client_id":"multi_scope","name":"Multi","description":"Multiple scopes","scopes":["read","write","admin","delete"],"active":true,"created_at":"` + testTimeStr + `","updated_at":"` + testTimeStr + `"}`,
			want: response.OAuthClientResponse{
				ID:          "123e4567-e89b-12d3-a456-426614174003",
				ClientID:    "multi_scope",
				Name:        "Multi",
				Description: "Multiple scopes",
				Scopes:      []string{"read", "write", "admin", "delete"},
				Active:      true,
				CreatedAt:   testTime,
				UpdatedAt:   testTime,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got response.OAuthClientResponse
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.ID != tt.want.ID {
					t.Errorf("OAuthClientResponse.ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.ClientID != tt.want.ClientID {
					t.Errorf("OAuthClientResponse.ClientID = %v, want %v", got.ClientID, tt.want.ClientID)
				}
				if got.Name != tt.want.Name {
					t.Errorf("OAuthClientResponse.Name = %v, want %v", got.Name, tt.want.Name)
				}
				if got.Description != tt.want.Description {
					t.Errorf("OAuthClientResponse.Description = %v, want %v", got.Description, tt.want.Description)
				}
				if !reflect.DeepEqual(got.Scopes, tt.want.Scopes) {
					t.Errorf("OAuthClientResponse.Scopes = %v, want %v", got.Scopes, tt.want.Scopes)
				}
				if got.Active != tt.want.Active {
					t.Errorf("OAuthClientResponse.Active = %v, want %v", got.Active, tt.want.Active)
				}
				if !got.CreatedAt.Equal(tt.want.CreatedAt) {
					t.Errorf("OAuthClientResponse.CreatedAt = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
				}
				if !got.UpdatedAt.Equal(tt.want.UpdatedAt) {
					t.Errorf("OAuthClientResponse.UpdatedAt = %v, want %v", got.UpdatedAt, tt.want.UpdatedAt)
				}
			}
		})
	}
}

func TestOAuthClientResponse_Marshal(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	testTimeStr := testTime.Format(time.RFC3339)

	tests := []struct {
		name     string
		response response.OAuthClientResponse
		want     string
	}{
		{
			name: "marshal complete client",
			response: response.OAuthClientResponse{
				ID:          "123e4567-e89b-12d3-a456-426614174000",
				ClientID:    "test_client",
				Name:        "Test Client",
				Description: "A test client",
				Scopes:      []string{"read", "write"},
				Active:      true,
				CreatedAt:   testTime,
				UpdatedAt:   testTime,
			},
			want: `{"id":"123e4567-e89b-12d3-a456-426614174000","client_id":"test_client","name":"Test Client","description":"A test client","scopes":["read","write"],"active":true,"created_at":"` + testTimeStr + `","updated_at":"` + testTimeStr + `"}`,
		},
		{
			name: "marshal inactive client with null scopes",
			response: response.OAuthClientResponse{
				ID:          "123e4567-e89b-12d3-a456-426614174001",
				ClientID:    "inactive",
				Name:        "Inactive",
				Description: "",
				Scopes:      nil,
				Active:      false,
				CreatedAt:   testTime,
				UpdatedAt:   testTime,
			},
			want: `{"id":"123e4567-e89b-12d3-a456-426614174001","client_id":"inactive","name":"Inactive","description":"","scopes":null,"active":false,"created_at":"` + testTimeStr + `","updated_at":"` + testTimeStr + `"}`,
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

func TestOAuthClientResponse_Fields(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	resp := response.OAuthClientResponse{
		ID:          "123e4567-e89b-12d3-a456-426614174000",
		ClientID:    "test_client",
		Name:        "Test Client",
		Description: "A test client",
		Scopes:      []string{"read", "write"},
		Active:      true,
		CreatedAt:   testTime,
		UpdatedAt:   testTime,
	}

	if resp.ID != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("OAuthClientResponse.ID = %v, want 123e4567-e89b-12d3-a456-426614174000", resp.ID)
	}

	if resp.ClientID != "test_client" {
		t.Errorf("OAuthClientResponse.ClientID = %v, want test_client", resp.ClientID)
	}

	if resp.Name != "Test Client" {
		t.Errorf("OAuthClientResponse.Name = %v, want Test Client", resp.Name)
	}

	if resp.Description != "A test client" {
		t.Errorf("OAuthClientResponse.Description = %v, want A test client", resp.Description)
	}

	expectedScopes := []string{"read", "write"}
	if !reflect.DeepEqual(resp.Scopes, expectedScopes) {
		t.Errorf("OAuthClientResponse.Scopes = %v, want %v", resp.Scopes, expectedScopes)
	}

	if !resp.Active {
		t.Errorf("OAuthClientResponse.Active = %v, want true", resp.Active)
	}

	if !resp.CreatedAt.Equal(testTime) {
		t.Errorf("OAuthClientResponse.CreatedAt = %v, want %v", resp.CreatedAt, testTime)
	}

	if !resp.UpdatedAt.Equal(testTime) {
		t.Errorf("OAuthClientResponse.UpdatedAt = %v, want %v", resp.UpdatedAt, testTime)
	}
}
