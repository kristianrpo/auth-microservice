package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/handler/shared"
)

func TestRespondWithJSON(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		payload        interface{}
		wantStatusCode int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "respond with valid map payload",
			status:         http.StatusOK,
			payload:        map[string]string{"message": "success"},
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]string
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp["message"] != "success" {
					t.Errorf("message = %v, want success", resp["message"])
				}
			},
		},
		{
			name:           "respond with struct payload",
			status:         http.StatusCreated,
			payload:        struct{ ID string }{ ID: "123" },
			wantStatusCode: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp struct{ ID string }
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.ID != "123" {
					t.Errorf("ID = %v, want 123", resp.ID)
				}
			},
		},
		{
			name:           "respond with array payload",
			status:         http.StatusOK,
			payload:        []string{"item1", "item2"},
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp []string
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(resp) != 2 {
					t.Errorf("length = %v, want 2", len(resp))
				}
				if resp[0] != "item1" {
					t.Errorf("resp[0] = %v, want item1", resp[0])
				}
			},
		},
		{
			name:           "respond with empty object",
			status:         http.StatusOK,
			payload:        struct{}{},
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				if w.Body.String() != "{}\n" {
					t.Errorf("body = %v, want {}", w.Body.String())
				}
			},
		},
		{
			name:           "respond with 404 status",
			status:         http.StatusNotFound,
			payload:        map[string]string{"error": "not found"},
			wantStatusCode: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]string
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp["error"] != "not found" {
					t.Errorf("error = %v, want not found", resp["error"])
				}
			},
		},
		{
			name:           "respond with 500 status",
			status:         http.StatusInternalServerError,
			payload:        map[string]string{"error": "internal error"},
			wantStatusCode: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]string
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp["error"] != "internal error" {
					t.Errorf("error = %v, want internal error", resp["error"])
				}
			},
		},
		{
			name:           "respond with null payload",
			status:         http.StatusOK,
			payload:        nil,
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				body := w.Body.String()
				if body != "null\n" {
					t.Errorf("body = %v, want null", body)
				}
			},
		},
		{
			name:           "respond with string payload",
			status:         http.StatusOK,
			payload:        "simple string",
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp string
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp != "simple string" {
					t.Errorf("resp = %v, want simple string", resp)
				}
			},
		},
		{
			name:           "respond with number payload",
			status:         http.StatusOK,
			payload:        42,
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp int
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp != 42 {
					t.Errorf("resp = %v, want 42", resp)
				}
			},
		},
		{
			name:           "respond with boolean payload",
			status:         http.StatusOK,
			payload:        true,
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp bool
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp != true {
					t.Errorf("resp = %v, want true", resp)
				}
			},
		},
		{
			name:           "respond with nested structure",
			status:         http.StatusOK,
			payload: map[string]interface{}{
				"user": map[string]string{
					"name":  "John",
					"email": "john@example.com",
				},
				"active": true,
			},
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if user, ok := resp["user"].(map[string]interface{}); ok {
					if user["name"] != "John" {
						t.Errorf("user.name = %v, want John", user["name"])
					}
				} else {
					t.Error("user field is not a map")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			shared.RespondWithJSON(w, tt.status, tt.payload)

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

func TestRespondWithJSON_ContentType(t *testing.T) {
	w := httptest.NewRecorder()
	payload := map[string]string{"key": "value"}

	shared.RespondWithJSON(w, http.StatusOK, payload)

	if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type = %v, want application/json", contentType)
	}
}

func TestRespondWithJSON_StatusCode(t *testing.T) {
	statusCodes := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNoContent,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusConflict,
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
	}

	for _, status := range statusCodes {
		t.Run(http.StatusText(status), func(t *testing.T) {
			w := httptest.NewRecorder()
			payload := map[string]string{"status": http.StatusText(status)}

			shared.RespondWithJSON(w, status, payload)

			if w.Code != status {
				t.Errorf("status code = %v, want %v", w.Code, status)
			}
		})
	}
}

