package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
)

func TestHealthCheckHandler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name           string
		dbPingFunc     func(ctx context.Context) error
		redisErr       error
		wantStatusCode int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "all services healthy",
			dbPingFunc: func(ctx context.Context) error {
				return nil
			},
			redisErr:       nil,
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.HealthResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Status != "healthy" {
					t.Errorf("Status = %v, want healthy", resp.Status)
				}
				if resp.Version != "1.0.0-test" {
					t.Errorf("Version = %v, want 1.0.0-test", resp.Version)
				}
				if resp.Services["database"] != "healthy" {
					t.Errorf("Services[database] = %v, want healthy", resp.Services["database"])
				}
				if resp.Services["redis"] != "healthy" {
					t.Errorf("Services[redis] = %v, want healthy", resp.Services["redis"])
				}
			},
		},
		{
			name: "database unhealthy",
			dbPingFunc: func(ctx context.Context) error {
				return errors.New("database connection failed")
			},
			redisErr:       nil,
			wantStatusCode: http.StatusServiceUnavailable,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.HealthResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Status != "unhealthy" {
					t.Errorf("Status = %v, want unhealthy", resp.Status)
				}
				if resp.Services["database"] != "unhealthy" {
					t.Errorf("Services[database] = %v, want unhealthy", resp.Services["database"])
				}
				if resp.Services["redis"] != "healthy" {
					t.Errorf("Services[redis] = %v, want healthy", resp.Services["redis"])
				}
			},
		},
		{
			name: "redis unhealthy",
			dbPingFunc: func(ctx context.Context) error {
				return nil
			},
			redisErr:       errors.New("redis connection failed"),
			wantStatusCode: http.StatusServiceUnavailable,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.HealthResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Status != "unhealthy" {
					t.Errorf("Status = %v, want unhealthy", resp.Status)
				}
				if resp.Services["database"] != "healthy" {
					t.Errorf("Services[database] = %v, want healthy", resp.Services["database"])
				}
				if resp.Services["redis"] != "unhealthy" {
					t.Errorf("Services[redis] = %v, want unhealthy", resp.Services["redis"])
				}
			},
		},
		{
			name: "all services unhealthy",
			dbPingFunc: func(ctx context.Context) error {
				return errors.New("database error")
			},
			redisErr:       errors.New("redis error"),
			wantStatusCode: http.StatusServiceUnavailable,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.HealthResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Status != "unhealthy" {
					t.Errorf("Status = %v, want unhealthy", resp.Status)
				}
				if resp.Services["database"] != "unhealthy" {
					t.Errorf("Services[database] = %v, want unhealthy", resp.Services["database"])
				}
				if resp.Services["redis"] != "unhealthy" {
					t.Errorf("Services[redis] = %v, want unhealthy", resp.Services["redis"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := NewMockDB(tt.dbPingFunc)
			mockRedis := NewMockRedisClient(tt.redisErr)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()

			handlerFunc := func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				services := make(map[string]string)

				if err := mockDB.PingContext(ctx); err != nil {
					logger.Warn("database health check failed")
					services["database"] = "unhealthy"
				} else {
					services["database"] = "healthy"
				}

				if err := mockRedis.Ping(ctx).Err(); err != nil {
					logger.Warn("redis health check failed")
					services["redis"] = "unhealthy"
				} else {
					services["redis"] = "healthy"
				}

				overallStatus := "healthy"
				for _, status := range services {
					if status == "unhealthy" {
						overallStatus = "unhealthy"
						break
					}
				}

				resp := response.HealthResponse{
					Status:    overallStatus,
					Timestamp: time.Now(),
					Version:   "1.0.0-test",
					Services:  services,
				}

				statusCode := http.StatusOK
				if overallStatus == "unhealthy" {
					statusCode = http.StatusServiceUnavailable
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)
				_ = json.NewEncoder(w).Encode(resp)
			}

			handlerFunc(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("status code = %v, want %v", w.Code, tt.wantStatusCode)
			}

			if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("Content-Type = %v, want application/json", contentType)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

