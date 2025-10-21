package tests

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
	
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		userName  string
		idCitizen int
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid user",
			email:     "test@example.com",
			password:  "password123",
			userName:  "Test User",
			idCitizen: 12345,
			wantErr:   false,
		},
		{
			name:      "empty email",
			email:     "",
			password:  "password123",
			userName:  "Test User",
			idCitizen: 12345,
			wantErr:   true,
			errMsg:    "email is required",
		},
		{
			name:      "empty password",
			email:     "test@example.com",
			password:  "",
			userName:  "Test User",
			idCitizen: 12345,
			wantErr:   true,
			errMsg:    "password is required",
		},
		{
			name:      "short password",
			email:     "test@example.com",
			password:  "short",
			userName:  "Test User",
			idCitizen: 12345,
			wantErr:   true,
			errMsg:    "password must be at least 8 characters",
		},
		{
			name:      "empty name",
			email:     "test@example.com",
			password:  "password123",
			userName:  "",
			idCitizen: 12345,
			wantErr:   true,
			errMsg:    "name is required",
		},
		{
			name:      "zero id_citizen",
			email:     "test@example.com",
			password:  "password123",
			userName:  "Test User",
			idCitizen: 0,
			wantErr:   true,
			errMsg:    "id_citizen is required and must be positive",
		},
		{
			name:      "negative id_citizen",
			email:     "test@example.com",
			password:  "password123",
			userName:  "Test User",
			idCitizen: -1,
			wantErr:   true,
			errMsg:    "id_citizen is required and must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := domain.NewUser(tt.email, tt.password, tt.userName, tt.idCitizen)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUser() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("NewUser() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUser() unexpected error: %v", err)
				return
			}

			if user.Email != tt.email {
				t.Errorf("NewUser() email = %v, want %v", user.Email, tt.email)
			}

			if user.Name != tt.userName {
				t.Errorf("NewUser() name = %v, want %v", user.Name, tt.userName)
			}

			if user.IDCitizen != tt.idCitizen {
				t.Errorf("NewUser() idCitizen = %v, want %v", user.IDCitizen, tt.idCitizen)
			}

			if user.Role != domain.RoleUser {
				t.Errorf("NewUser() role = %v, want %v", user.Role, domain.RoleUser)
			}

			// Verify password is hashed
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tt.password)); err != nil {
				t.Errorf("NewUser() password not properly hashed")
			}
		})
	}
}

func TestUser_ComparePassword(t *testing.T) {
	user, _ := domain.NewUser("test@example.com", "password123", "Test User", 12345)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "correct password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "incorrect password",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.ComparePassword(tt.password)

			if tt.wantErr && err == nil {
				t.Errorf("ComparePassword() expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("ComparePassword() unexpected error: %v", err)
			}
		})
	}
}

func TestUser_ToPublic(t *testing.T) {
	user, _ := domain.NewUser("test@example.com", "password123", "Test User", 12345)
	user.ID = "test-id"

	publicUser := user.ToPublic()

	if publicUser.ID != user.ID {
		t.Errorf("ToPublic() ID = %v, want %v", publicUser.ID, user.ID)
	}

	if publicUser.IDCitizen != user.IDCitizen {
		t.Errorf("ToPublic() IDCitizen = %v, want %v", publicUser.IDCitizen, user.IDCitizen)
	}

	if publicUser.Email != user.Email {
		t.Errorf("ToPublic() Email = %v, want %v", publicUser.Email, user.Email)
	}

	if publicUser.Name != user.Name {
		t.Errorf("ToPublic() Name = %v, want %v", publicUser.Name, user.Name)
	}

	if publicUser.Role != user.Role {
		t.Errorf("ToPublic() Role = %v, want %v", publicUser.Role, user.Role)
	}

	if publicUser.CreatedAt != user.CreatedAt {
		t.Errorf("ToPublic() CreatedAt = %v, want %v", publicUser.CreatedAt, user.CreatedAt)
	}

	if publicUser.UpdatedAt != user.UpdatedAt {
		t.Errorf("ToPublic() UpdatedAt = %v, want %v", publicUser.UpdatedAt, user.UpdatedAt)
	}
}

