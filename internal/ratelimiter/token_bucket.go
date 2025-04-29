package ratelimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity     int
	tokens       int
	refillRate   int // count of adding tokens per time
	lastRefill   time.Time
	mu           sync.Mutex
	refillTicker *time.Ticker
}

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	tb := &TokenBucket{
		capacity:     capacity,
		tokens:       capacity,
		refillRate:   refillRate,
		lastRefill:   time.Now(),
		refillTicker: time.NewTicker(time.Second),
	}
	go tb.refillTokens()
	return tb
}

func (tb *TokenBucket) refillTokens() {
	for range tb.refillTicker.C {
		tb.mu.Lock()
		elapsed := time.Since(tb.lastRefill)
		tb.tokens += int(elapsed.Seconds()) * tb.refillRate
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = time.Now()
		tb.mu.Unlock()
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func (tb *TokenBucket) TokensLeft() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.tokens
}
