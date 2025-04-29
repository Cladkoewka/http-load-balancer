package server

import (
	"net"
	"net/http"

	"github.com/Cladkoewka/http-load-balancer/internal/logger"
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
		logger.Log.Info("Rate limit check", "ClientIP", clientIP, "RemainingTokens", bucket.TokensLeft())

		if !bucket.Allow() {
			http.Error(rw, "Too Many Requests", http.StatusTooManyRequests)
			logger.Log.Warn("Rate limit exceeded", "ClientIP", clientIP)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
