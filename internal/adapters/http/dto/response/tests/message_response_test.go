package tests

import (
	"encoding/json"
	"testing"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
)

func TestMessageResponse_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    response.MessageResponse
		wantErr bool
	}{
		{
			name:  "valid message response",
			input: `{"message":"Operation successful"}`,
			want: response.MessageResponse{
				Message: "Operation successful",
			},
			wantErr: false,
		},
		{
			name:  "valid empty message",
			input: `{"message":""}`,
			want: response.MessageResponse{
				Message: "",
			},
			wantErr: false,
		},
		{
			name:  "valid long message",
			input: `{"message":"This is a very long success message with lots of details about the operation"}`,
			want: response.MessageResponse{
				Message: "This is a very long success message with lots of details about the operation",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{"message":}`,
			want:    response.MessageResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got response.MessageResponse
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Message != tt.want.Message {
					t.Errorf("MessageResponse.Message = %v, want %v", got.Message, tt.want.Message)
				}
			}
		})
	}
}

func TestMessageResponse_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		response response.MessageResponse
		want     string
	}{
		{
			name: "marshal valid message",
			response: response.MessageResponse{
				Message: "Operation successful",
			},
			want: `{"message":"Operation successful"}`,
		},
		{
			name: "marshal empty message",
			response: response.MessageResponse{
				Message: "",
			},
			want: `{"message":""}`,
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

func TestMessageResponse_Fields(t *testing.T) {
	resp := response.MessageResponse{
		Message: "Operation successful",
	}

	if resp.Message != "Operation successful" {
		t.Errorf("MessageResponse.Message = %v, want Operation successful", resp.Message)
	}
}
