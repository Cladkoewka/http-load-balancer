package balancer

import (
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/Cladkoewka/http-load-balancer/internal/logger"
)

type Strategy string

const (
	RoundRobin Strategy = "round-robin"
	Random     Strategy = "random"
)

type LoadBalancer struct {
	BackendURLs []string
	Index       int
	Strategy    Strategy
	mu          sync.Mutex
}

func NewLoadBalancer(backendURLs []string, strategy Strategy) *LoadBalancer {
	rand.Seed(time.Now().UnixNano()) // init rand generator
	return &LoadBalancer{
		BackendURLs: backendURLs,
		Strategy:    strategy,
	}
}

func (lb *LoadBalancer) GetNextBackendURL() (*url.URL, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.BackendURLs) == 0 {
		return nil, fmt.Errorf("no backend URLs available")
	}

	var selectedIndex int
	var err error

	// choose algorithm
	switch lb.Strategy {
	case RoundRobin:
		selectedIndex, err = lb.selectRoundRobin()
	case Random:
		selectedIndex, err = lb.selectRandom()
	default:
		return nil, fmt.Errorf("unknown load balancing strategy: %s", lb.Strategy)
	}

	if err != nil {
		return nil, err
	}

	rawURL := lb.BackendURLs[selectedIndex]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		logger.Log.Error("failed to parse URL", "backend URL", rawURL)
		return nil, fmt.Errorf("failed to parse backend URL: %w", err)
	}

	logger.Log.Info("Selected backend URL", "backend", parsedURL)
	return parsedURL, nil
}

func (lb *LoadBalancer) selectRoundRobin() (int, error) {
	index := lb.Index
	lb.Index = (lb.Index + 1) % len(lb.BackendURLs)
	return index, nil
}

func (lb *LoadBalancer) selectRandom() (int, error) {
	return rand.Intn(len(lb.BackendURLs)), nil
}
