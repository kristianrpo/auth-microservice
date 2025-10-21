package tests

import (
	"encoding/json"
	"testing"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
)

func TestErrorResponse_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    response.ErrorResponse
		wantErr bool
	}{
		{
			name:  "valid error response with all fields",
			input: `{"error":"Bad request","code":"BAD_REQUEST","details":"Invalid email format"}`,
			want: response.ErrorResponse{
				Error:   "Bad request",
				Code:    "BAD_REQUEST",
				Details: "Invalid email format",
			},
			wantErr: false,
		},
		{
			name:  "valid error response without optional fields",
			input: `{"error":"Internal server error"}`,
			want: response.ErrorResponse{
				Error:   "Internal server error",
				Code:    "",
				Details: "",
			},
			wantErr: false,
		},
		{
			name:  "valid error response with code only",
			input: `{"error":"Unauthorized","code":"UNAUTHORIZED"}`,
			want: response.ErrorResponse{
				Error:   "Unauthorized",
				Code:    "UNAUTHORIZED",
				Details: "",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"error":"Bad request","code":}`,
			want:    response.ErrorResponse{},
			wantErr: true,
		},
		{
			name:  "empty fields",
			input: `{"error":"","code":"","details":""}`,
			want: response.ErrorResponse{
				Error:   "",
				Code:    "",
				Details: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got response.ErrorResponse
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Error != tt.want.Error {
					t.Errorf("ErrorResponse.Error = %v, want %v", got.Error, tt.want.Error)
				}
				if got.Code != tt.want.Code {
					t.Errorf("ErrorResponse.Code = %v, want %v", got.Code, tt.want.Code)
				}
				if got.Details != tt.want.Details {
					t.Errorf("ErrorResponse.Details = %v, want %v", got.Details, tt.want.Details)
				}
			}
		})
	}
}

func TestErrorResponse_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		response response.ErrorResponse
		want     string
	}{
		{
			name: "marshal with all fields",
			response: response.ErrorResponse{
				Error:   "Bad request",
				Code:    "BAD_REQUEST",
				Details: "Invalid email format",
			},
			want: `{"error":"Bad request","code":"BAD_REQUEST","details":"Invalid email format"}`,
		},
		{
			name: "marshal with error only",
			response: response.ErrorResponse{
				Error:   "Internal server error",
				Code:    "",
				Details: "",
			},
			want: `{"error":"Internal server error"}`,
		},
		{
			name: "marshal with error and code",
			response: response.ErrorResponse{
				Error:   "Unauthorized",
				Code:    "UNAUTHORIZED",
				Details: "",
			},
			want: `{"error":"Unauthorized","code":"UNAUTHORIZED"}`,
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

func TestErrorResponse_Fields(t *testing.T) {
	resp := response.ErrorResponse{
		Error:   "Bad request",
		Code:    "BAD_REQUEST",
		Details: "Invalid email format",
	}

	if resp.Error != "Bad request" {
		t.Errorf("ErrorResponse.Error = %v, want Bad request", resp.Error)
	}

	if resp.Code != "BAD_REQUEST" {
		t.Errorf("ErrorResponse.Code = %v, want BAD_REQUEST", resp.Code)
	}

	if resp.Details != "Invalid email format" {
		t.Errorf("ErrorResponse.Details = %v, want Invalid email format", resp.Details)
	}
}
