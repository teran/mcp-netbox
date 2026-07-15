package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestExtractClientIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		headers      map[string]string
		remote       string
		trustedProxy string
		expected     string
	}{
		{"RemoteAddr parsed", nil, "10.0.0.2:56789", "", "10.0.0.2"},
		{"RemoteAddr no port", nil, "10.0.0.3", "", "10.0.0.3"},
		{"X-Forwarded-For overridden by RemoteAddr", map[string]string{"X-Forwarded-For": "192.168.1.1"}, "10.0.0.4:12345", "", "10.0.0.4"},
		{"X-Client-IP overridden by RemoteAddr", map[string]string{"X-Client-IP": "10.0.0.1"}, "10.0.0.5:12345", "", "10.0.0.5"},
		{"trusted proxy X-Forwarded-For", map[string]string{"X-Forwarded-For": "192.168.1.1"}, "10.0.0.1:12345", "10.0.0.0/8", "192.168.1.1"},
		{"trusted proxy invalid IP in X-Forwarded-For", map[string]string{"X-Forwarded-For": "invalid-ip"}, "10.0.0.1:12345", "10.0.0.0/8", "10.0.0.1"},
		{"trusted proxy no X-Forwarded-For", nil, "10.0.0.1:12345", "10.0.0.0/8", "10.0.0.1"},
		{"trusted proxy non-matching remote", map[string]string{"X-Forwarded-For": "192.168.1.1"}, "192.168.1.1:12345", "10.0.0.0/8", "192.168.1.1"},
		{"trusted proxy empty string", map[string]string{"X-Forwarded-For": "10.0.0.2"}, "10.0.0.1:12345", "", "10.0.0.1"},
		{"trusted proxy invalid CIDR", map[string]string{"X-Forwarded-For": "10.0.0.2"}, "10.0.0.1:12345", "not-a-cidr", "10.0.0.1"},
		{"IPv6 remote address", nil, "[2001:db8::1]:443", "", "2001:db8::1"},
		{"IPv6 loopback", nil, "[::1]:8080", "", "::1"},
		{"IPv6 with trusted proxy", map[string]string{"X-Forwarded-For": "192.168.1.1"}, "[2001:db8::1]:443", "2001:db8::0/32", "192.168.1.1"},
		{"IPv6 with trusted proxy no XFF", nil, "[2001:db8::1]:443", "2001:db8::0/32", "2001:db8::1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			req.RemoteAddr = tc.remote

			got := extractClientIP(req, tc.trustedProxy)
			if got != tc.expected {
				t.Errorf("extractClientIP() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestNewRateLimiter(t *testing.T) {
	t.Parallel()

	rl := NewRateLimiter(RateLimiterConfig{
		GlobalLimit:    rate.Limit(100),
		GlobalBurst:    200,
		PerClientLimit: rate.Limit(10),
		PerClientBurst: 20,
	})
	defer rl.Stop()

	if !rl.Allow("127.0.0.1") {
		t.Error("Allow() = false, want true")
	}
}

func TestRateLimiter_GlobalLimit(t *testing.T) {
	t.Parallel()

	rl := NewRateLimiter(RateLimiterConfig{
		GlobalLimit:    rate.Limit(1),
		GlobalBurst:    1,
		PerClientLimit: rate.Limit(100),
		PerClientBurst: 200,
	})
	defer rl.Stop()

	if !rl.Allow("client-1") {
		t.Error("Allow() = false, want true")
	}
	// Second request should be rate limited
	if rl.Allow("client-2") {
		t.Error("Allow() = true, want false (global limit)")
	}
}

func TestRateLimiter_PerClientLimit(t *testing.T) {
	t.Parallel()

	rl := NewRateLimiter(RateLimiterConfig{
		GlobalLimit:    rate.Limit(100),
		GlobalBurst:    200,
		PerClientLimit: rate.Limit(1),
		PerClientBurst: 1,
	})
	defer rl.Stop()

	if !rl.Allow("client-1") {
		t.Error("Allow() = false, want true")
	}
	if rl.Allow("client-1") {
		t.Error("Allow() = true, want false (per-client limit)")
	}

	// Different client should be allowed
	if !rl.Allow("client-2") {
		t.Error("Allow() = false, want true (different client)")
	}
}

func TestRateLimiter_EvictExpired(t *testing.T) {
	t.Parallel()

	rl := &rateLimiter{
		config: RateLimiterConfig{
			GlobalLimit:    rate.Limit(100),
			GlobalBurst:    200,
			PerClientLimit: rate.Limit(10),
			PerClientBurst: 20,
		},
		global:  rate.NewLimiter(rate.Limit(100), 200),
		clients: make(map[string]*clientLimiter),
		stopCh:  make(chan struct{}),
	}

	rl.clients["stale"] = &clientLimiter{
		limiter:  rate.NewLimiter(rate.Limit(10), 20),
		lastSeen: time.Now().Add(-clientTTL - time.Minute),
	}

	rl.evictStaleClients()
	if _, ok := rl.clients["stale"]; ok {
		t.Error("stale client was not evicted")
	}
}

func TestRateLimiter_DefaultBurst(t *testing.T) {
	t.Parallel()

	rl := NewRateLimiter(RateLimiterConfig{
		GlobalLimit:    rate.Limit(100),
		GlobalBurst:    200,
		PerClientLimit: rate.Limit(10),
		PerClientBurst: 20,
	})
	defer rl.Stop()

	var wg sync.WaitGroup
	mu := sync.Mutex{}
	allowed := 0

	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if rl.Allow("test-client") {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if allowed == 0 {
		t.Error("no requests were allowed")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("200 within limit", func(t *testing.T) {
		mw, stop := RateLimitMiddleware(RateLimiterConfig{
			GlobalLimit:    rate.Limit(1000),
			GlobalBurst:    2000,
			PerClientLimit: rate.Limit(1000),
			PerClientBurst: 2000,
		}, "")
		defer stop()

		handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		req.RemoteAddr = "10.0.0.1:12345"
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("429 when rate limited", func(t *testing.T) {
		mw, stop := RateLimitMiddleware(RateLimiterConfig{
			GlobalLimit:    rate.Limit(1),
			GlobalBurst:    1,
			PerClientLimit: rate.Limit(1000),
			PerClientBurst: 2000,
		}, "")
		defer stop()

		handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// First request should succeed
		req1 := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		req1.RemoteAddr = "10.0.0.1:12345"
		rec1 := httptest.NewRecorder()
		handler.ServeHTTP(rec1, req1)

		// Second request should be rate-limited
		req2 := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		req2.RemoteAddr = "10.0.0.1:12345"
		rec2 := httptest.NewRecorder()
		handler.ServeHTTP(rec2, req2)

		if rec2.Code != http.StatusTooManyRequests {
			t.Errorf("status = %d, want %d", rec2.Code, http.StatusTooManyRequests)
		}
	})
}
