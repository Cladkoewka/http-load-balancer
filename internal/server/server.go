package server

import (
	"fmt"
	"net/http"

	"github.com/Cladkoewka/http-load-balancer/internal/balancer"
	"github.com/Cladkoewka/http-load-balancer/internal/logger"
	"github.com/Cladkoewka/http-load-balancer/internal/proxy"
	"github.com/Cladkoewka/http-load-balancer/internal/ratelimiter"
)

func StartServer(lb *balancer.LoadBalancer, port int) error {
	rl := ratelimiter.NewRateLimiter()

	// Wrap the main handler with the logging middleware and rate limit middleware
	var handler http.Handler
	handler = http.HandlerFunc(proxy.ProxyHandler(lb))
	handler = RateLimitMiddleware(rl, handler)
	handler = LoggingMiddleware(handler)
	

	logger.Log.Info("Starting HTTP server", "port", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
} 