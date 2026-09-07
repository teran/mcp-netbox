package handlers

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
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

const (
	cleanupInterval = 10 * time.Minute
	clientTTL       = 30 * time.Minute
)

type rateLimiter struct {
	config    RateLimiterConfig
	global    *rate.Limiter
	clients   map[string]*clientLimiter
	mu        sync.Mutex
	stopCh    chan struct{}
	closeOnce sync.Once
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
	rl.closeOnce.Do(func() {
		close(rl.stopCh)
	})
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
func RateLimitMiddleware(cfg RateLimiterConfig, trustedProxy string) (func(http.Handler) http.Handler, func()) {
	rl := NewRateLimiter(cfg)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := extractClientIP(r, trustedProxy)
			if !rl.Allow(clientIP) {
				getLogger().WithField("client_ip", SanitizeLog(clientIP)).Warn("rate limit exceeded")
				w.Header().Set("Retry-After", "1")
				http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}, rl.Stop
}

func extractClientIP(r *http.Request, trustedProxy string) string {
	proxyPrefix := parseTrustedProxy(trustedProxy)
	if proxyPrefix == nil {
		// No trusted proxy configured — use RemoteAddr only (secure default).
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return r.RemoteAddr
		}
		return host
	}

	// Check if the connecting IP is from a trusted proxy
	remoteIPStr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	remoteIP, err := netip.ParseAddr(remoteIPStr)
	if err != nil {
		return remoteIPStr
	}

	if proxyPrefix.Contains(remoteIP) {
		// Trusted proxy — extract first IP from X-Forwarded-For
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			if firstIP := extractFirstIP(xff); firstIP != "" {
				return firstIP
			}
		}
	}

	return remoteIPStr
}

func extractFirstIP(xff string) string {
	if xff == "" {
		return ""
	}
	parts := strings.Split(xff, ",")
	ip := strings.TrimSpace(parts[0])
	if _, err := netip.ParseAddr(ip); err != nil {
		return ""
	}
	return ip
}

func parseTrustedProxy(cidr string) *netip.Prefix {
	if cidr == "" {
		return nil
	}
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return nil
	}
	return &prefix
}
