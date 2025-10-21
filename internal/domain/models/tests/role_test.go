package tests

import (
	"testing"

	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

func TestRole_String(t *testing.T) {
	tests := []struct {
		name string
		role domain.Role
		want string
	}{
		{
			name: "user role",
			role: domain.RoleUser,
			want: "USER",
		},
		{
			name: "admin role",
			role: domain.RoleAdmin,
			want: "ADMIN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.String(); got != tt.want {
				t.Errorf("Role.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		role domain.Role
		want bool
	}{
		{
			name: "valid user role",
			role: domain.RoleUser,
			want: true,
		},
		{
			name: "valid admin role",
			role: domain.RoleAdmin,
			want: true,
		},
		{
			name: "invalid role",
			role: domain.Role("INVALID"),
			want: false,
		},
		{
			name: "empty role",
			role: domain.Role(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsValid(); got != tt.want {
				t.Errorf("Role.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRole(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    domain.Role
		wantErr bool
	}{
		{
			name:    "parse USER",
			input:   "USER",
			want:    domain.RoleUser,
			wantErr: false,
		},
		{
			name:    "parse ADMIN",
			input:   "ADMIN",
			want:    domain.RoleAdmin,
			wantErr: false,
		},
		{
			name:    "parse invalid role",
			input:   "INVALID",
			want:    "",
			wantErr: true,
		},
		{
			name:    "parse empty string",
			input:   "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ParseRole(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseRole() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseRole() unexpected error: %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("ParseRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
