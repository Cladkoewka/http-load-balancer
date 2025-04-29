package server

import (
	"net/http"

	"github.com/Cladkoewka/http-load-balancer/internal/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Wrap ResponseWriter to intercept status code
		rr := &responseRecorder{ResponseWriter: rw, statusCode: http.StatusOK}
		next.ServeHTTP(rr, r)

		logger.Log.Info("HTTP request", 
			"method", r.Method, 
			"path", r.URL.Path, 
			"client", r.RemoteAddr,
			"status", rr.statusCode,
		)
	})
}

// Wrapper to capture status codes from HTTP responses
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}