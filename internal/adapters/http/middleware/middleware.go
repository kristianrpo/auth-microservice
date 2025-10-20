package middleware

import (
	"context"
	nethttp "net/http"
	"strings"

	"go.uber.org/zap"

	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	"github.com/kristianrpo/auth-microservice/internal/application/services"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

type contextKey string

const (
	// UserContextKey is the key to get the user from context
	UserContextKey contextKey = "user"
)

// AuthMiddleware is the authentication middleware
type AuthMiddleware struct {
	authService *services.AuthService
	logger      *zap.Logger
}

// NewAuthMiddleware creates a new instance of the authentication middleware
func NewAuthMiddleware(authService *services.AuthService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// Authenticate verifies the JWT token in the Authorization header
func (m *AuthMiddleware) Authenticate(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.Debug("missing authorization header")
			httperrors.RespondWithError(w, httperrors.ErrMissingAuthHeader)
			return
		}

		// Verify format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Debug("invalid authorization header format",
				zap.String("header", authHeader),
				zap.Int("parts_count", len(parts)),
				zap.Strings("parts", parts))
			httperrors.RespondWithError(w, httperrors.ErrInvalidAuthHeader)
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.authService.ValidateAccessToken(r.Context(), token)
		if err != nil {
			m.logger.Debug("invalid token", zap.Error(err))
			httperrors.RespondWithDomainError(w, err)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORSMiddleware maneja CORS
func CORSMiddleware(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(nethttp.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(logger *zap.Logger) func(nethttp.Handler) nethttp.Handler {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			logger.Info("http request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)
			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(logger *zap.Logger) func(nethttp.Handler) nethttp.Handler {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered",
						zap.Any("error", err),
						zap.String("path", r.URL.Path),
					)
					httperrors.RespondWithError(w, httperrors.ErrInternalServer)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext retrieves the user claims from the context
func GetUserFromContext(ctx context.Context) (*domain.TokenClaims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*domain.TokenClaims)
	return claims, ok
}
