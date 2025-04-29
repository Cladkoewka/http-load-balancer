package server

import (
	"fmt"
	"net/http"

	"github.com/Cladkoewka/http-load-balancer/internal/balancer"
	"github.com/Cladkoewka/http-load-balancer/internal/logger"
	"github.com/Cladkoewka/http-load-balancer/internal/proxy"
)

func StartServer(lb *balancer.LoadBalancer, port int) error {
	// Wrap the main handler with the logging middleware
	handler := LoggingMiddleware(http.HandlerFunc(proxy.ProxyHandler(lb)))

	logger.Log.Info("Starting HTTP server", "port", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
} 