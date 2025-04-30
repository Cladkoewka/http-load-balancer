package balancer

import (
	"fmt"
	"math/rand"
	"net/http"
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
	alive       map[string]bool
}

func NewLoadBalancer(backendURLs []string, strategy Strategy) *LoadBalancer {
	rand.New(rand.NewSource(time.Now().UnixMicro())) // init rand generator

	aliveMap := make(map[string]bool, len(backendURLs))
	for _, url := range backendURLs {
		aliveMap[url] = true
	}

	return &LoadBalancer{
		BackendURLs: backendURLs,
		Strategy:    strategy,
		alive:       aliveMap,
	}
}

func (lb *LoadBalancer) GetNextBackendURL() (*url.URL, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	liveBackends := make([]string, 0)
	for _, url := range lb.BackendURLs {
		if lb.alive[url] {
			liveBackends = append(liveBackends, url)
		}
	}

	if len(liveBackends) == 0 {
		return nil, fmt.Errorf("no backend URLs available")
	}

	var selectedURL string
	var err error

	// choose algorithm
	switch lb.Strategy {
	case RoundRobin:
		selectedURL, err = lb.selectRoundRobin(liveBackends)
	case Random:
		selectedURL, err = lb.selectRandom(liveBackends)
	default:
		return nil, fmt.Errorf("unknown load balancing strategy: %s", lb.Strategy)
	}

	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(selectedURL)
	if err != nil {
		logger.Log.Error("failed to parse URL", "backend URL", selectedURL)
		return nil, fmt.Errorf("failed to parse backend URL: %w", err)
	}

	logger.Log.Info("Selected backend URL", "backend", parsedURL)
	return parsedURL, nil
}

func (lb *LoadBalancer) selectRoundRobin(liveBackends []string) (string, error) {
	selected := liveBackends[lb.Index%len(liveBackends)]
	lb.Index = (lb.Index + 1) % len(liveBackends)
	return selected, nil
}

func (lb *LoadBalancer) selectRandom(liveBackends []string) (string, error) {
	return liveBackends[rand.Intn(len(liveBackends))], nil
}

func (lb *LoadBalancer) StartHealthCheck(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			lb.healthCheck()
		}
	}()
}

func (lb *LoadBalancer) healthCheck() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	logger.Log.Info("Health check started")

	for _, backendURL := range lb.BackendURLs {
		go func(url string) {
			client := http.Client{
				Timeout: 2 * time.Second,
			}
			resp, err := client.Get(url + "/health")
			lb.mu.Lock()
			defer lb.mu.Unlock()
			if err != nil || resp.StatusCode != http.StatusOK {
				lb.alive[url] = false
				logger.Log.Warn("Backend marked as DOWN", "backend", url)
			} else {
				lb.alive[url] = true
				logger.Log.Warn("Backend marked as UP", "backend", url)
			}
			if resp != nil {
				resp.Body.Close()
			}
		}(backendURL)
	}
}
