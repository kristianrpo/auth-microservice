package middleware

import (
	nethttp "net/http"

	"go.uber.org/zap"

	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// RoleMiddleware checks if the user has the required role
type RoleMiddleware struct {
	logger *zap.Logger
}

// NewRoleMiddleware creates a new instance of RoleMiddleware
func NewRoleMiddleware(logger *zap.Logger) *RoleMiddleware {
	return &RoleMiddleware{
		logger: logger,
	}
}

// RequireRole creates a middleware that checks if the user has one of the required roles
func (m *RoleMiddleware) RequireRole(allowedRoles ...domain.Role) func(nethttp.Handler) nethttp.Handler {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			// Get user claims from context (set by AuthMiddleware)
			claims, ok := GetUserFromContext(r.Context())
			if !ok {
				m.logger.Debug("no user claims in context")
				httperrors.RespondWithError(w, httperrors.ErrUnauthorized)
				return
			}

			// Check if user has one of the allowed roles
			hasRole := false
			for _, allowedRole := range allowedRoles {
				if claims.Role == allowedRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				m.logger.Warn("user does not have required role",
					zap.Int("id_citizen", claims.IDCitizen),
					zap.String("user_role", claims.Role.String()),
					zap.Any("required_roles", allowedRoles))
				httperrors.RespondWithError(w, httperrors.ErrForbidden)
				return
			}

			m.logger.Debug("user has required role",
				zap.Int("id_citizen", claims.IDCitizen),
				zap.String("role", claims.Role.String()))

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin is a shorthand for requiring ADMIN role
func (m *RoleMiddleware) RequireAdmin(next nethttp.Handler) nethttp.Handler {
	return m.RequireRole(domain.RoleAdmin)(next)
}

// RequireUser is a shorthand for requiring USER role (or higher)
func (m *RoleMiddleware) RequireUser(next nethttp.Handler) nethttp.Handler {
	return m.RequireRole(domain.RoleUser, domain.RoleAdmin)(next)
}
