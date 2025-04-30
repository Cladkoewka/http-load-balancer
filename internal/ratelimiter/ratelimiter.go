package ratelimiter

import (
	"sync"
)

type RateLimiter struct {
	buckets map[string]*TokenBucket
	mu      sync.Mutex
	defaultCapacity   int
	defaultRefillRate int
}



func NewRateLimiter(defaultCapacity, defaultRefill int) *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*TokenBucket),
		defaultCapacity:   defaultCapacity,
		defaultRefillRate: defaultRefill,
	}
}

func (rl *RateLimiter) GetOrCreateBucket(clientIP string) *TokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if bucket, exists := rl.buckets[clientIP]; exists {
		return bucket
	}

	bucket := NewTokenBucket(rl.defaultCapacity, rl.defaultRefillRate)
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
