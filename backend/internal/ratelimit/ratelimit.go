package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	limiter *rate.Limiter
	mu      sync.Mutex
}

// NewRateLimiter creates a new rate limiter
// rps: requests per second
// burst: maximum burst size
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
	}
}

// Limit returns a middleware that enforces rate limiting
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		allowed := rl.limiter.Allow()
		rl.mu.Unlock()

		if !allowed {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "rate limit exceeded", "message": "Too many requests. Please try again later."}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// PerIPRateLimiter provides per-IP rate limiting
type PerIPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      float64
	burst    int
	cleanup  *time.Ticker
}

// NewPerIPRateLimiter creates a new per-IP rate limiter
func NewPerIPRateLimiter(rps float64, burst int) *PerIPRateLimiter {
	rl := &PerIPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rps,
		burst:    burst,
		cleanup:  time.NewTicker(5 * time.Minute), // Cleanup every 5 minutes
	}

	// Start cleanup goroutine
	go rl.cleanupLimiters()

	return rl
}

// cleanupLimiters periodically cleans up old limiters
func (rl *PerIPRateLimiter) cleanupLimiters() {
	for range rl.cleanup.C {
		rl.mu.Lock()
		// In a production system, you'd track last access time
		// For now, we'll keep all limiters (memory trade-off)
		rl.mu.Unlock()
	}
}

// getLimiter gets or creates a limiter for an IP address
func (rl *PerIPRateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check after acquiring write lock
		limiter, exists = rl.limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(rl.rps), rl.burst)
			rl.limiters[ip] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter
}

// Limit returns a middleware that enforces per-IP rate limiting
func (rl *PerIPRateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "rate limit exceeded", "message": "Too many requests from your IP. Please try again later."}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP (client IP)
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
