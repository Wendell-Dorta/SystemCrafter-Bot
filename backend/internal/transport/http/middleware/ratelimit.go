package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type clientBucket struct {
	tokens     float64
	lastRefill time.Time
}

// RateLimiter implements a token-bucket rate limiter per IP.
type RateLimiter struct {
	mu            sync.Mutex
	clients       map[string]*clientBucket
	rate          float64
	burst         int
	blockedCount  int64
	totalRequests int64
	stopCleanup   chan struct{}
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(rate float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		clients:     make(map[string]*clientBucket),
		rate:        rate,
		burst:       burst,
		stopCleanup: make(chan struct{}),
	}

	go rl.startCleanupLoop(5 * time.Minute)
	return rl
}

// Middleware wraps an http.Handler with rate limiting enforcement.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&rl.totalRequests, 1)

		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		ip := getClientIP(r)
		allowed := rl.allow(ip)

		if !allowed {
			atomic.AddInt64(&rl.blockedCount, 1)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "Rate limit exceeded. Please slow down.",
				"status":  http.StatusTooManyRequests,
				"success": false,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.clients[ip]
	if !exists {
		rl.clients[ip] = &clientBucket{
			tokens:     float64(rl.burst) - 1.0,
			lastRefill: now,
		}
		return true
	}

	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * rl.rate
	if bucket.tokens > float64(rl.burst) {
		bucket.tokens = float64(rl.burst)
	}
	bucket.lastRefill = now

	if bucket.tokens >= 1.0 {
		bucket.tokens -= 1.0
		return true
	}

	return false
}

// Stats returns rate limiting metrics.
func (rl *RateLimiter) Stats() map[string]interface{} {
	rl.mu.Lock()
	activeClients := len(rl.clients)
	rl.mu.Unlock()

	return map[string]interface{}{
		"activeTrackedIPs": activeClients,
		"rateLimitRPS":     rl.rate,
		"rateLimitBurst":   rl.burst,
		"blockedRequests":  atomic.LoadInt64(&rl.blockedCount),
		"totalRequests":    atomic.LoadInt64(&rl.totalRequests),
	}
}

func (rl *RateLimiter) startCleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-rl.stopCleanup:
			return
		case <-ticker.C:
			rl.cleanupStale()
		}
	}
}

func (rl *RateLimiter) cleanupStale() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, b := range rl.clients {
		if now.Sub(b.lastRefill) > 10*time.Minute {
			delete(rl.clients, ip)
		}
	}
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}

	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
