package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
)

// contextKey is an unexported type for context value keys.
// Using an unexported type prevents context-key collisions from other packages.
type contextKey string

const (
	// tokenContextKey is the context key storing the NetBox API token.
	tokenContextKey contextKey = "netbox_token"

	// MaxTokenLength is the maximum allowed length for the Authorization token.
	// Increased to 4096 to support NetBox v2 tokens (nbt_<key>.<token>).
	MaxTokenLength = 4096
)

var ErrTokenRequired = errors.New("authorization token is required")

// requestIDHeader is the standard header used to propagate a correlation ID.
const requestIDHeader = "X-Request-ID"

// RequestIDMiddleware injects a per-request correlation ID into the context.
// If the inbound request already carries an X-Request-ID (e.g. from a reverse
// proxy), that value is reused; otherwise a fresh ID is generated (L9/G11).
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestIDHeader)
		if id == "" {
			id = newRequestID()
		}
		ctx := domain.WithRequestID(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithSession returns a logrus entry derived from logger and decorated with the
// request-scoped correlation fields (request_id, and session_id when known)
// read from ctx. It is the base entry for all request-scoped log lines (L9).
func WithSession(ctx context.Context, logger *logrus.Logger) *logrus.Entry {
	entry := logger.WithFields(logrus.Fields{})
	if id := domain.RequestIDFromContext(ctx); id != "" {
		entry = entry.WithField("request_id", id)
	}
	if sid := sessionIDFromContext(ctx); sid != "" {
		entry = entry.WithField("session_id", sid)
	}
	return entry
}

// sessionIDFromContext returns an optional session id from the context. MCP
// sessions are not tracked by this server, so it always returns "".
func sessionIDFromContext(context.Context) string {
	return ""
}

// newRequestID returns a fresh correlation identifier.
func newRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// TokenMiddleware extracts the NetBox API token from the Authorization header.
// Supports both "Bearer <token>" and "Token <token>" schemes (case-insensitive).
func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, ErrTokenRequired.Error()), http.StatusUnauthorized)
			return
		}

		var token string
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			http.Error(w, `{"error":"malformed authorization header"}`, http.StatusUnauthorized)
			return
		}

		scheme := strings.ToLower(parts[0])
		if scheme != "bearer" && scheme != "token" {
			http.Error(w, `{"error":"unsupported authorization scheme"}`, http.StatusUnauthorized)
			return
		}

		token = strings.TrimSpace(parts[1])
		if token == "" {
			http.Error(w, `{"error":"empty token"}`, http.StatusUnauthorized)
			return
		}

		if len(token) > MaxTokenLength {
			http.Error(w, `{"error":"token exceeds maximum length"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), tokenContextKey, application.NewToken(token))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkBatchSize(body []byte) error {
	if len(body) == 0 {
		return nil
	}
	// Check if it's a batch request (starts with '[')
	trimmed := bytes.TrimLeft(body, " \t\r\n")
	if len(trimmed) == 0 || trimmed[0] != '[' {
		return nil
	}

	var batch []json.RawMessage
	if err := json.Unmarshal(body, &batch); err != nil {
		return fmt.Errorf("invalid batch request: %w", err)
	}
	if len(batch) > MaxBatchSize {
		return fmt.Errorf("batch request exceeds maximum size of %d", MaxBatchSize)
	}
	return nil
}

// DefaultMaxRequestBodySize is 1 MB.
const DefaultMaxRequestBodySize = 1 << 20

// MaxBatchSize is the maximum number of items in a JSON-RPC batch request.
const MaxBatchSize = 100

// BodyLimitMiddleware limits the request body size using http.MaxBytesReader.
func BodyLimitMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}

// MaxBytesError checks if an error is an http.MaxBytesError.
func MaxBytesError(err error) bool {
	var maxBytesErr *http.MaxBytesError
	return errors.As(err, &maxBytesErr)
}

type loggingResponseWriter struct {
	http.ResponseWriter

	statusCode int
	bodySize   atomic.Int64
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bodySize.Add(int64(n))
	return n, err
}

func mcpRequestMethod(body []byte) string {
	if len(body) == 0 {
		return "empty_body"
	}

	var req struct {
		Method string `json:"method"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return "parse_error"
	}
	if req.Method == "" {
		return "no_method"
	}

	return req.Method
}

// SanitizeLog removes control characters from a string for safe logging.
func SanitizeLog(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\t' {
			return -1
		}
		return r
	}, s)
}

// SecurityHeadersMiddleware sets standard security HTTP headers on every response.
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// HSTS is intentionally omitted — TLS termination is handled by the reverse proxy.
		// CSP is intentionally omitted — MCP uses JSON-RPC and does not serve HTML.
		next.ServeHTTP(w, r)
	})
}

// HostValidationMiddleware rejects requests with empty or malformed Host headers.
func HostValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host == "" {
			http.Error(w, `{"error":"host header is required"}`, http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware catches panics in downstream handlers and returns 500.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				getLogger().WithError(fmt.Errorf("%v", rec)).WithField("stack", string(debug.Stack())).Error("panic recovered")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs MCP requests with method, duration, size, and status.
// It never logs the Authorization header content.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			if MaxBytesError(err) {
				http.Error(w, `{"error":"request body too large"}`, http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, `{"error":"failed to read request body"}`, http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		if err := checkBatchSize(body); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		method := mcpRequestMethod(body)
		// Extract tool name from tools/call/<toolname> or method field
		toolName := method
		if strings.HasPrefix(method, "tools/call/") {
			toolName = strings.TrimPrefix(method, "tools/call/")
		}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		WithSession(r.Context(), getLogger()).WithFields(logrus.Fields{
			"http_method": r.Method,
			"path":        r.URL.Path,
			"method":      toolName,
			"duration":    duration,
			"req_size":    len(body),
			"resp_size":   lrw.bodySize.Load(),
			"status":      lrw.statusCode,
		}).Info("mcp_request")
	})
}
