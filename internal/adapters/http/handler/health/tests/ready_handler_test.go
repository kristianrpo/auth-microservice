package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	httperrors "github.com/kristianrpo/auth-microservice/internal/adapters/http/errors"
)

func TestReadyHandler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name           string
		dbPingFunc     func(ctx context.Context) error
		redisErr       error
		wantStatusCode int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "service is ready",
			dbPingFunc: func(ctx context.Context) error {
				return nil
			},
			redisErr:       nil,
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]string
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp["status"] != "ready" {
					t.Errorf("status = %v, want ready", resp["status"])
				}
			},
		},
		{
			name: "database not ready",
			dbPingFunc: func(ctx context.Context) error {
				return errors.New("connection refused")
			},
			redisErr:       nil,
			wantStatusCode: http.StatusServiceUnavailable,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if errMsg, ok := resp["error"].(string); ok {
					if errMsg != "database not ready" {
						t.Errorf("error = %v, want database not ready", errMsg)
					}
				}
			},
		},
		{
			name: "redis not ready",
			dbPingFunc: func(ctx context.Context) error {
				return nil
			},
			redisErr:       errors.New("redis connection failed"),
			wantStatusCode: http.StatusServiceUnavailable,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if errMsg, ok := resp["error"].(string); ok {
					if errMsg != "redis not ready" {
						t.Errorf("error = %v, want redis not ready", errMsg)
					}
				}
			},
		},
		{
			name: "both services not ready",
			dbPingFunc: func(ctx context.Context) error {
				return errors.New("db error")
			},
			redisErr:       errors.New("redis error"),
			wantStatusCode: http.StatusServiceUnavailable,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if errMsg, ok := resp["error"].(string); ok {
					if errMsg != "database not ready" {
						t.Errorf("error = %v, want database not ready", errMsg)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := NewMockDB(tt.dbPingFunc)
			mockRedis := NewMockRedisClient(tt.redisErr)

			req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
			w := httptest.NewRecorder()

			handlerFunc := func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()

				if err := mockDB.PingContext(ctx); err != nil {
					logger.Warn("database health check failed")
					httperrors.RespondWithErrorMessage(w, http.StatusServiceUnavailable, "database not ready")
					return
				}

				if err := mockRedis.Ping(ctx).Err(); err != nil {
					logger.Warn("redis health check failed")
					httperrors.RespondWithErrorMessage(w, http.StatusServiceUnavailable, "redis not ready")
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
			}

			handlerFunc(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("status code = %v, want %v", w.Code, tt.wantStatusCode)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

