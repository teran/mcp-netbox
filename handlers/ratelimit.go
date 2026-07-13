package handlers

import (
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiterConfig struct {
	GlobalLimit    rate.Limit
	GlobalBurst    int
	PerClientLimit rate.Limit
	PerClientBurst int
}

const cleanupInterval = 10 * time.Minute
const clientTTL = 30 * time.Minute

type rateLimiter struct {
	config  RateLimiterConfig
	global  *rate.Limiter
	clients map[string]*clientLimiter
	mu      sync.Mutex
	stopCh  chan struct{}
}

func NewRateLimiter(config RateLimiterConfig) *rateLimiter {
	rl := &rateLimiter{
		config:  config,
		global:  rate.NewLimiter(config.GlobalLimit, config.GlobalBurst),
		clients: make(map[string]*clientLimiter),
		stopCh:  make(chan struct{}),
	}
	go rl.evictExpired()
	return rl
}

func (rl *rateLimiter) Stop() {
	close(rl.stopCh)
}

func (rl *rateLimiter) Allow(clientIP string) bool {
	// Global limit first
	if !rl.global.Allow() {
		return false
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	cl, ok := rl.clients[clientIP]
	if !ok {
		cl = &clientLimiter{
			limiter:  rate.NewLimiter(rl.config.PerClientLimit, rl.config.PerClientBurst),
			lastSeen: time.Now(),
		}
		rl.clients[clientIP] = cl
	}
	cl.lastSeen = time.Now()

	return cl.limiter.Allow()
}

func (rl *rateLimiter) evictExpired() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.evictStaleClients()
		case <-rl.stopCh:
			return
		}
	}
}

func (rl *rateLimiter) evictStaleClients() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, cl := range rl.clients {
		if now.Sub(cl.lastSeen) > clientTTL {
			delete(rl.clients, ip)
		}
	}
}

// RateLimitMiddleware returns a middleware and a stop function for the rate limiter.
func RateLimitMiddleware(cfg RateLimiterConfig) (func(http.Handler) http.Handler, func()) {
	rl := NewRateLimiter(cfg)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := extractClientIP(r)
			if !rl.Allow(clientIP) {
				slog.Warn("rate limit exceeded", "client_ip", SanitizeLog(clientIP))
				w.Header().Set("Retry-After", "1")
				http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}, rl.Stop
}

func extractClientIP(r *http.Request) string {
	// Only use RemoteAddr to prevent IP spoofing via X-Forwarded-For / X-Client-IP headers.
	// If running behind a reverse proxy, ensure the proxy sets X-Envoy-External-Address
	// or use a trusted header via configuration.
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
