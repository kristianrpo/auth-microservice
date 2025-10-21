package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	middleware "github.com/kristianrpo/auth-microservice/internal/adapters/http/middleware"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
	"go.uber.org/zap"
)

func TestRequireRole_AdminAllowed(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	am := middleware.NewRoleMiddleware(logger)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// Put token claims in context with admin role
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, &domain.TokenClaims{UserID: "user-123", Role: domain.RoleAdmin})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler := am.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 got %d", w.Result().StatusCode)
	}
}

func TestRequireRole_UserAllowed(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	am := middleware.NewRoleMiddleware(logger)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// Put token claims in context with user role
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, &domain.TokenClaims{UserID: "user-456", Role: domain.RoleUser})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler := am.RequireUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 got %d", w.Result().StatusCode)
	}
}
