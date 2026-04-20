package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/stride-pro/backend/pkg/response"
)

// bucket tracks the token state for a single IP.
type bucket struct {
	tokens    float64
	lastCheck time.Time
}

// RateLimiter implements a per-IP token bucket rate limiter.
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     float64 // tokens added per second
	capacity float64 // max tokens
}

// NewRateLimiter creates a rate limiter with the given rate (requests/sec) and burst capacity.
func NewRateLimiter(rate, capacity float64) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*bucket),
		rate:     rate,
		capacity: capacity,
	}
	go rl.cleanup()
	return rl
}

// NewAuthRateLimiter creates a stricter rate limiter suited for auth endpoints.
// Allows 5 requests per minute per IP with a burst of 5.
func NewAuthRateLimiter() *RateLimiter {
	return NewRateLimiter(5.0/60.0, 5)
}

// Middleware returns an HTTP middleware that enforces rate limits per client IP.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)

		if !rl.allow(ip) {
			w.Header().Set("Retry-After", "60")
			response.Error(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, ok := rl.buckets[ip]
	if !ok {
		rl.buckets[ip] = &bucket{
			tokens:    rl.capacity - 1,
			lastCheck: time.Now(),
		}
		return true
	}

	now := time.Now()
	elapsed := now.Sub(b.lastCheck).Seconds()
	b.lastCheck = now

	// Add tokens based on elapsed time
	b.tokens += elapsed * rl.rate
	if b.tokens > rl.capacity {
		b.tokens = rl.capacity
	}

	if b.tokens < 1 {
		return false
	}

	b.tokens--
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-10 * time.Minute)
		for ip, b := range rl.buckets {
			if b.lastCheck.Before(cutoff) {
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// extractIP returns the best-effort client IP from the request.
// It reads only the first address from X-Forwarded-For to avoid using
// the full proxy chain as the key (which could be spoofed to bypass limits).
func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can be a comma-separated list: "client, proxy1, proxy2"
		// Take only the first entry (reported client IP) and trim whitespace.
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	return r.RemoteAddr
}
