package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Cladkoewka/http-load-balancer/internal/balancer"
	"github.com/Cladkoewka/http-load-balancer/internal/logger"
	"github.com/Cladkoewka/http-load-balancer/internal/proxy"
	"github.com/Cladkoewka/http-load-balancer/internal/ratelimiter"
	"github.com/Cladkoewka/http-load-balancer/internal/config"
)

func StartServer(cfg *config.Config) error {
	lb := balancer.NewLoadBalancer(cfg.BackendURLs, balancer.RoundRobin) // balancer algorithm
	lb.StartHealthCheck(time.Second * time.Duration(cfg.HealthCheckInterval))

	rl := ratelimiter.NewRateLimiter(cfg.RateLimitDefaults.BucketCapacity, cfg.RateLimitDefaults.BucketRefillRate)

	mux := http.NewServeMux()

	// Wrap the proxy handler with middlewares
	var handler http.Handler
	handler = http.HandlerFunc(proxy.ProxyHandler(lb))
	handler = RateLimitMiddleware(rl, handler)
	handler = LoggingMiddleware(handler)

	mux.Handle("/test", handler)

	addr := fmt.Sprintf(":%d", cfg.Port)

	// Create aт HTTP server to support graceful shutdown
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start server in a separate goroutine to allow signal listening
	go func() {
		logger.Log.Info("Starting HTTP server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("Server error", "error", err)
			os.Exit(1) // Exit immediately on critical server error
		}
	}()

	// Create channel to receive OS shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-stop
	logger.Log.Info("Shutdown signal received")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Graceful shutdown failed", "error", err)
		return err
	}

	logger.Log.Info("Server stopped gracefully")
	return nil
}
