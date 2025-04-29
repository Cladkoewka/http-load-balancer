package ratelimiter

import (
	"sync"

	"github.com/Cladkoewka/http-load-balancer/internal/config"
)

type RateLimiter struct {
	buckets map[string]*TokenBucket
	mu sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*TokenBucket),
	}
}

func (rl *RateLimiter) GetOrCreateBucket(clientIP string) *TokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if bucket, exists := rl.buckets[clientIP]; exists {
		return bucket
	}
	bucket := NewTokenBucket(config.BucketCapacity, config.BucketRefillRate)
	rl.buckets[clientIP] = bucket
	return bucket
}

func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	bucket, exists := rl.buckets[clientIP]

	if !exists {
		return false
	}
	return bucket.Allow()
}