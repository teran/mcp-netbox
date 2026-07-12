package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
	"unicode"
)

// ContextKey is an unexported type for context value keys.
type ContextKey string

const (
	// TokenContextKey is the context key storing the NetBox API token.
	TokenContextKey ContextKey = "netbox_token"

	// MaxTokenLength is the maximum allowed length for the Authorization token.
	MaxTokenLength = 512
)

var ErrTokenRequired = errors.New("authorization token is required")

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

		ctx := context.WithValue(r.Context(), TokenContextKey, token)
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
		return nil
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
	bodySize   int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bodySize += n
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
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, s)
}

// RecoveryMiddleware catches panics in downstream handlers and returns 500.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC: %v\n%s", rec, debug.Stack())
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
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

		log.Printf("mcp_log http_method=%s path=%s method=%s duration=%s req_size=%d resp_size=%d status=%d", //nolint:gosec
			SanitizeLog(r.Method),
			SanitizeLog(r.URL.Path),
			SanitizeLog(toolName),
			duration,
			len(body),
			lrw.bodySize,
			lrw.statusCode,
		)
	})
}
