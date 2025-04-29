package balancer

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/Cladkoewka/http-load-balancer/internal/logger"
)

type LoadBalancer struct {
	BackendURLs []string
	Index int
	mu sync.Mutex
}

func NewLoadBalancer(backendURLs []string) *LoadBalancer {
	return &LoadBalancer{
		BackendURLs: backendURLs,
		Index: 0,
	}
}

func (lb *LoadBalancer) GetNextBackendURL() (*url.URL, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.BackendURLs) == 0 {
		return nil, fmt.Errorf("no backend URLs available")
	}

	backendURL := lb.BackendURLs[lb.Index]
	lb.Index = (lb.Index + 1 ) % len(lb.BackendURLs) // Round-Robin algorithm 

	parsedURL, err := url.Parse(backendURL)
	if err != nil {
		logger.Log.Error("failed to parse URL %s", backendURL)
		return nil, fmt.Errorf("failed to parse backend URL: %w", err)
	}

	
	logger.Log.Info("Selected backend URL", "backend", parsedURL)
	return parsedURL, nil
}