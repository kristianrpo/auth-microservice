package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/application/services"
	domainerrors "github.com/kristianrpo/auth-microservice/internal/domain/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

func TestAuthService_Register(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name        string
		email       string
		password    string
		userName    string
		idCitizen   int
		existsFunc  func(ctx context.Context, email string) (bool, error)
		createFunc  func(ctx context.Context, user *domain.User) error
		wantErr     bool
		expectedErr error
	}{
		{
			name:      "successful registration",
			email:     "test@example.com",
			password:  "password123",
			userName:  "Test User",
			idCitizen: 12345,
			existsFunc: func(ctx context.Context, email string) (bool, error) {
				return false, nil
			},
			createFunc: func(ctx context.Context, user *domain.User) error {
				user.ID = "user-123"
				return nil
			},
			wantErr: false,
		},
		{
			name:      "user already exists",
			email:     "existing@example.com",
			password:  "password123",
			userName:  "Existing User",
			idCitizen: 12345,
			existsFunc: func(ctx context.Context, email string) (bool, error) {
				return true, nil
			},
			wantErr:     true,
			expectedErr: domainerrors.ErrUserAlreadyExists,
		},
		{
			name:      "repository error on exists check",
			email:     "test@example.com",
			password:  "password123",
			userName:  "Test User",
			idCitizen: 12345,
			existsFunc: func(ctx context.Context, email string) (bool, error) {
				return false, domainerrors.ErrInternal
			},
			wantErr:     true,
			expectedErr: domainerrors.ErrInternal,
		},
		{
			name:      "repository error on create",
			email:     "test@example.com",
			password:  "password123",
			userName:  "Test User",
			idCitizen: 12345,
			existsFunc: func(ctx context.Context, email string) (bool, error) {
				return false, nil
			},
			createFunc: func(ctx context.Context, user *domain.User) error {
				return domainerrors.ErrInternal
			},
			wantErr:     true,
			expectedErr: domainerrors.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockUserRepository{
				ExistsFunc: tt.existsFunc,
				CreateFunc: tt.createFunc,
			}
			mockTokenRepo := &MockTokenRepository{}
			jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)
			authService := services.NewAuthService(mockUserRepo, mockTokenRepo, jwtService, logger)

			user, err := authService.Register(context.Background(), tt.email, tt.password, tt.userName, tt.idCitizen)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Register() expected error but got none")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("Register() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Register() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("Register() returned nil user")
				return
			}

			if user.Email != tt.email {
				t.Errorf("Register() Email = %v, want %v", user.Email, tt.email)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Create a test user with known password
	testUser, _ := domain.NewUser("test@example.com", "password123", "Test User", 12345)
	testUser.ID = "user-123"

	tests := []struct {
		name           string
		email          string
		password       string
		getByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
		wantErr        bool
		expectedErr    error
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "password123",
			getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return testUser, nil
			},
			wantErr: false,
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: "password123",
			getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return nil, domainerrors.ErrUserNotFound
			},
			wantErr:     true,
			expectedErr: domainerrors.ErrInvalidCredentials,
		},
		{
			name:     "invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return testUser, nil
			},
			wantErr:     true,
			expectedErr: domainerrors.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockUserRepository{
				GetByEmailFunc: tt.getByEmailFunc,
			}
			mockTokenRepo := &MockTokenRepository{}
			jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)
			authService := services.NewAuthService(mockUserRepo, mockTokenRepo, jwtService, logger)

			tokenPair, err := authService.Login(context.Background(), tt.email, tt.password)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Login() expected error but got none")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("Login() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Login() unexpected error: %v", err)
				return
			}

			if tokenPair == nil {
				t.Errorf("Login() returned nil token pair")
				return
			}

			if tokenPair.AccessToken == "" {
				t.Errorf("Login() AccessToken is empty")
			}

			if tokenPair.RefreshToken == "" {
				t.Errorf("Login() RefreshToken is empty")
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)

	// Generate a valid refresh token
	validRefreshToken, _ := jwtService.GenerateRefreshToken("user-123", "test@example.com", domain.RoleUser)

	tests := []struct {
		name                   string
		refreshToken           string
		isTokenBlacklistedFunc func(ctx context.Context, token string) (bool, error)
		getRefreshTokenFunc    func(ctx context.Context, token string) (*domain.RefreshTokenData, error)
		wantErr                bool
		expectedErr            error
	}{
		{
			name:         "successful refresh",
			refreshToken: validRefreshToken,
			isTokenBlacklistedFunc: func(ctx context.Context, token string) (bool, error) {
				return false, nil
			},
			getRefreshTokenFunc: func(ctx context.Context, token string) (*domain.RefreshTokenData, error) {
				return &domain.RefreshTokenData{
					UserID: "user-123",
					Email:  "test@example.com",
				}, nil
			},
			wantErr: false,
		},
		{
			name:         "blacklisted token",
			refreshToken: validRefreshToken,
			isTokenBlacklistedFunc: func(ctx context.Context, token string) (bool, error) {
				return true, nil
			},
			wantErr:     true,
			expectedErr: domainerrors.ErrTokenRevoked,
		},
		{
			name:         "token not in cache",
			refreshToken: validRefreshToken,
			isTokenBlacklistedFunc: func(ctx context.Context, token string) (bool, error) {
				return false, nil
			},
			getRefreshTokenFunc: func(ctx context.Context, token string) (*domain.RefreshTokenData, error) {
				return nil, domainerrors.ErrInvalidToken
			},
			wantErr:     true,
			expectedErr: domainerrors.ErrInvalidToken,
		},
		{
			name:         "invalid token",
			refreshToken: "invalid.token.here",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockUserRepository{}
			mockTokenRepo := &MockTokenRepository{
				IsTokenBlacklistedFunc: tt.isTokenBlacklistedFunc,
				GetRefreshTokenFunc:    tt.getRefreshTokenFunc,
			}
			authService := services.NewAuthService(mockUserRepo, mockTokenRepo, jwtService, logger)

			tokenPair, err := authService.RefreshToken(context.Background(), tt.refreshToken)

			if tt.wantErr {
				if err == nil {
					t.Errorf("RefreshToken() expected error but got none")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("RefreshToken() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("RefreshToken() unexpected error: %v", err)
				return
			}

			if tokenPair == nil {
				t.Errorf("RefreshToken() returned nil token pair")
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)

	// Generate valid tokens
	validAccessToken, _ := jwtService.GenerateAccessToken("user-123", "test@example.com", domain.RoleUser)
	validRefreshToken, _ := jwtService.GenerateRefreshToken("user-123", "test@example.com", domain.RoleUser)

	tests := []struct {
		name         string
		accessToken  string
		refreshToken string
		wantErr      bool
	}{
		{
			name:         "successful logout",
			accessToken:  validAccessToken,
			refreshToken: validRefreshToken,
			wantErr:      false,
		},
		{
			name:         "invalid access token",
			accessToken:  "invalid.token",
			refreshToken: validRefreshToken,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockUserRepository{}
			mockTokenRepo := &MockTokenRepository{}
			authService := services.NewAuthService(mockUserRepo, mockTokenRepo, jwtService, logger)

			err := authService.Logout(context.Background(), tt.accessToken, tt.refreshToken)

			if tt.wantErr && err == nil {
				t.Errorf("Logout() expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Logout() unexpected error: %v", err)
			}
		})
	}
}

func TestAuthService_GetUserByID(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	testUser, _ := domain.NewUser("test@example.com", "password123", "Test User", 12345)
	testUser.ID = "user-123"

	tests := []struct {
		name        string
		userID      string
		getByIDFunc func(ctx context.Context, id string) (*domain.User, error)
		wantErr     bool
	}{
		{
			name:   "user found",
			userID: "user-123",
			getByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
				return testUser, nil
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			getByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
				return nil, domainerrors.ErrUserNotFound
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockUserRepository{
				GetByIDFunc: tt.getByIDFunc,
			}
			mockTokenRepo := &MockTokenRepository{}
			jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)
			authService := services.NewAuthService(mockUserRepo, mockTokenRepo, jwtService, logger)

			user, err := authService.GetUserByID(context.Background(), tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUserByID() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GetUserByID() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("GetUserByID() returned nil user")
			}
		})
	}
}

func TestAuthService_Register_DuplicateIDCitizen(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	mockUserRepo := &MockUserRepository{
		ExistsFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
		GetByIDCitizenFunc: func(ctx context.Context, idCitizen int) (*domain.User, error) {
			// simulate existing user with same id_citizen
			u, _ := domain.NewUser("other@example.com", "password123", "Other", idCitizen)
			return u, nil
		},
	}

	mockTokenRepo := &MockTokenRepository{}
	jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)
	authService := services.NewAuthService(mockUserRepo, mockTokenRepo, jwtService, logger)

	_, err := authService.Register(context.Background(), "test@example.com", "password123", "Name", 123)
	if err == nil {
		t.Fatalf("expected error for duplicate id_citizen, got nil")
	}
	if !errors.Is(err, domainerrors.ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestAuthService_Register_CreateReturnsAlreadyExists(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	mockUserRepo := &MockUserRepository{
		ExistsFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
		GetByIDCitizenFunc: func(ctx context.Context, idCitizen int) (*domain.User, error) {
			return nil, domainerrors.ErrUserNotFound
		},
		CreateFunc: func(ctx context.Context, user *domain.User) error {
			return domainerrors.ErrUserAlreadyExists
		},
	}

	mockTokenRepo := &MockTokenRepository{}
	jwtService := services.NewJWTService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 7*24*time.Hour, logger)
	authService := services.NewAuthService(mockUserRepo, mockTokenRepo, jwtService, logger)

	_, err := authService.Register(context.Background(), "test@example.com", "password123", "Name", 123)
	if err == nil {
		t.Fatalf("expected error when repo Create returns ErrUserAlreadyExists, got nil")
	}
	if !errors.Is(err, domainerrors.ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}
}
