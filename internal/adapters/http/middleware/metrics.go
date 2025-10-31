package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/kristianrpo/auth-microservice/internal/observability/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// MetricsMiddleware records basic HTTP metrics for each request.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)

		endpoint := r.URL.Path
		if route := mux.CurrentRoute(r); route != nil {
			if template, err := route.GetPathTemplate(); err == nil {
				endpoint = template
			}
		}

		metrics.ObserveHTTPRequest(r.Method, endpoint, strconv.Itoa(recorder.status), duration)
	})
}
