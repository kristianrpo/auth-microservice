package errors

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/kristianrpo/auth-microservice/internal/adapters/http/dto/response"
)

// RespondWithError sends an HTTP error response
func RespondWithError(w nethttp.ResponseWriter, err *HTTPError) {
	resp := response.ErrorResponse{
		Error:   err.Message,
		Code:    err.Code,
		Details: "",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

// RespondWithErrorMessage sends an HTTP error response with a custom message
func RespondWithErrorMessage(w nethttp.ResponseWriter, statusCode int, message string) {
	resp := response.ErrorResponse{
		Error:   message,
		Code:    "",
		Details: "",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

// RespondWithDomainError maps a domain error and sends the HTTP response
func RespondWithDomainError(w nethttp.ResponseWriter, err error) {
	httpErr := MapDomainError(err)
	RespondWithError(w, httpErr)
}
