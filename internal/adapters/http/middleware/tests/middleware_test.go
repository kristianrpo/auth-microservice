package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/middleware"
	domain "github.com/kristianrpo/auth-microservice/internal/domain/models"
)

// mock auth service to satisfy ValidateAccessToken
type mockAuthSvc struct {
	validate func(ctx context.Context, token string) (*domain.TokenClaims, error)
}

func (m *mockAuthSvc) ValidateAccessToken(ctx context.Context, token string) (*domain.TokenClaims, error) {
	return m.validate(ctx, token)
}

func TestCORSMiddleware_Options(t *testing.T) {
	handler := middleware.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for OPTIONS, got %d", w.Code)
	}
}

func TestCORSMiddleware_PassThrough(t *testing.T) {
	handler := middleware.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTeapot {
		t.Fatalf("expected 418, got %d", w.Code)
	}
}

func TestLoggingMiddleware_InvokesNext(t *testing.T) {
	logger := zap.NewNop()
	lm := middleware.LoggingMiddleware(logger)
	called := false
	handler := lm(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if !called {
		t.Fatal("next handler was not called")
	}
}

func TestRecoveryMiddleware_Recovers(t *testing.T) {
	logger := zap.NewNop()
	rm := middleware.RecoveryMiddleware(logger)

	handler := rm(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 after panic recovery, got %d", w.Code)
	}
}

func TestAuthenticate_MissingHeader(t *testing.T) {
	// authService not needed for this negative case - construct via NewAuthMiddleware
	logger := zap.NewNop()
	m := middleware.NewAuthMiddleware(nil, logger)

	handler := m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized && w.Code != http.StatusBadRequest {
		t.Fatalf("expected 401/400 for missing auth header, got %d", w.Code)
	}
}

func TestGetUserFromContext(t *testing.T) {
	claims := &domain.TokenClaims{IDCitizen: 1, Email: "e@e.com"}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)
	got, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		t.Fatal("expected claims in context")
	}
	if got.IDCitizen != 1 {
		t.Fatalf("unexpected id_citizen %d", got.IDCitizen)
	}
}
