package server

import (
	"net"
	"net/http"
	"github.com/Cladkoewka/http-load-balancer/internal/ratelimiter"
)

func RateLimitMiddleware(rl *ratelimiter.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(rw, "Unable to parse clientID", http.StatusInternalServerError)
			return
		}

		bucket := rl.GetOrCreateBucket(clientIP)

		if !bucket.Allow() {
			http.Error(rw, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(rw, r)
	})
}