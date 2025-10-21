package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
)

func TestHealthResponse_JSON(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	testTimeStr := testTime.Format(time.RFC3339)

	tests := []struct {
		name    string
		input   string
		want    response.HealthResponse
		wantErr bool
	}{
		{
			name:  "valid health response",
			input: `{"status":"healthy","timestamp":"` + testTimeStr + `","version":"1.0.0","services":{"database":"connected","cache":"connected"}}`,
			want: response.HealthResponse{
				Status:    "healthy",
				Timestamp: testTime,
				Version:   "1.0.0",
				Services: map[string]string{
					"database": "connected",
					"cache":    "connected",
				},
			},
			wantErr: false,
		},
		{
			name:  "health response with degraded status",
			input: `{"status":"degraded","timestamp":"` + testTimeStr + `","version":"1.0.0","services":{"database":"connected","cache":"disconnected"}}`,
			want: response.HealthResponse{
				Status:    "degraded",
				Timestamp: testTime,
				Version:   "1.0.0",
				Services: map[string]string{
					"database": "connected",
					"cache":    "disconnected",
				},
			},
			wantErr: false,
		},
		{
			name:  "health response with empty services",
			input: `{"status":"healthy","timestamp":"` + testTimeStr + `","version":"1.0.0","services":{}}`,
			want: response.HealthResponse{
				Status:    "healthy",
				Timestamp: testTime,
				Version:   "1.0.0",
				Services:  map[string]string{},
			},
			wantErr: false,
		},
		{
			name:  "health response with null services",
			input: `{"status":"healthy","timestamp":"` + testTimeStr + `","version":"1.0.0","services":null}`,
			want: response.HealthResponse{
				Status:    "healthy",
				Timestamp: testTime,
				Version:   "1.0.0",
				Services:  nil,
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"status":"healthy","timestamp":}`,
			want:    response.HealthResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got response.HealthResponse
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Status != tt.want.Status {
					t.Errorf("HealthResponse.Status = %v, want %v", got.Status, tt.want.Status)
				}
				if !got.Timestamp.Equal(tt.want.Timestamp) {
					t.Errorf("HealthResponse.Timestamp = %v, want %v", got.Timestamp, tt.want.Timestamp)
				}
				if got.Version != tt.want.Version {
					t.Errorf("HealthResponse.Version = %v, want %v", got.Version, tt.want.Version)
				}
				if len(got.Services) != len(tt.want.Services) {
					t.Errorf("HealthResponse.Services length = %v, want %v", len(got.Services), len(tt.want.Services))
				}
				for key, value := range tt.want.Services {
					if got.Services[key] != value {
						t.Errorf("HealthResponse.Services[%s] = %v, want %v", key, got.Services[key], value)
					}
				}
			}
		})
	}
}

func TestHealthResponse_Marshal(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	testTimeStr := testTime.Format(time.RFC3339)

	tests := []struct {
		name     string
		response response.HealthResponse
		want     string
	}{
		{
			name: "marshal valid health response",
			response: response.HealthResponse{
				Status:    "healthy",
				Timestamp: testTime,
				Version:   "1.0.0",
				Services: map[string]string{
					"database": "connected",
				},
			},
			want: `{"status":"healthy","timestamp":"` + testTimeStr + `","version":"1.0.0","services":{"database":"connected"}}`,
		},
		{
			name: "marshal with null services",
			response: response.HealthResponse{
				Status:    "healthy",
				Timestamp: testTime,
				Version:   "1.0.0",
				Services:  nil,
			},
			want: `{"status":"healthy","timestamp":"` + testTimeStr + `","version":"1.0.0","services":null}`,
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

func TestHealthResponse_Fields(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	resp := response.HealthResponse{
		Status:    "healthy",
		Timestamp: testTime,
		Version:   "1.0.0",
		Services: map[string]string{
			"database": "connected",
			"cache":    "connected",
		},
	}

	if resp.Status != "healthy" {
		t.Errorf("HealthResponse.Status = %v, want healthy", resp.Status)
	}

	if !resp.Timestamp.Equal(testTime) {
		t.Errorf("HealthResponse.Timestamp = %v, want %v", resp.Timestamp, testTime)
	}

	if resp.Version != "1.0.0" {
		t.Errorf("HealthResponse.Version = %v, want 1.0.0", resp.Version)
	}

	if resp.Services["database"] != "connected" {
		t.Errorf("HealthResponse.Services[database] = %v, want connected", resp.Services["database"])
	}

	if resp.Services["cache"] != "connected" {
		t.Errorf("HealthResponse.Services[cache] = %v, want connected", resp.Services["cache"])
	}
}
