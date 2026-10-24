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

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/redact"
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

// ErrTokenRequired indicates that no authorization token was provided.
var ErrTokenRequired = errors.New("authorization token is required")

const (
	// authHeaderParts is the expected number of fields in an Authorization header.
	authHeaderParts = 2
	// serverErrorStatus is the status-code threshold above which a request is
	// considered an error for observability purposes.
	serverErrorStatus = 400
)

// requestIDHeader is the standard header used to propagate a correlation ID.
const requestIDHeader = "X-Request-ID"

// RequestIDMiddleware injects a per-request correlation ID into the context.
// If the inbound request already carries an X-Request-ID (e.g. from a reverse
// proxy), that value is reused; otherwise a fresh ID is generated (L9/G11).
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id := request.Header.Get(requestIDHeader)
		if id == "" {
			id = newRequestID()
		}
		ctx := domain.WithRequestID(request.Context(), id)
		next.ServeHTTP(writer, request.WithContext(ctx))
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
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(writer, fmt.Sprintf(`{"error":"%s"}`, ErrTokenRequired.Error()), http.StatusUnauthorized)

			return
		}

		var token string
		parts := strings.SplitN(authHeader, " ", authHeaderParts)
		if len(parts) != authHeaderParts {
			http.Error(writer, `{"error":"malformed authorization header"}`, http.StatusUnauthorized)

			return
		}

		scheme := strings.ToLower(parts[0])
		if scheme != "bearer" && scheme != "token" {
			http.Error(writer, `{"error":"unsupported authorization scheme"}`, http.StatusUnauthorized)

			return
		}

		token = strings.TrimSpace(parts[1])
		if token == "" {
			http.Error(writer, `{"error":"empty token"}`, http.StatusUnauthorized)

			return
		}

		if len(token) > MaxTokenLength {
			http.Error(writer, `{"error":"token exceeds maximum length"}`, http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(request.Context(), tokenContextKey, application.NewToken(token))
		next.ServeHTTP(writer, request.WithContext(ctx))
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
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			request.Body = http.MaxBytesReader(writer, request.Body, maxBytes)
			next.ServeHTTP(writer, request)
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
func SanitizeLog(value string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\t' {
			return -1
		}

		return r
	}, value)
}

// SecurityHeadersMiddleware sets standard security HTTP headers on every response.
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("X-Frame-Options", "DENY")
		writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// HSTS is intentionally omitted — TLS termination is handled by the reverse proxy.
		// CSP is intentionally omitted — MCP uses JSON-RPC and does not serve HTML.
		next.ServeHTTP(writer, request)
	})
}

// HostValidationMiddleware rejects requests with empty or malformed Host headers.
func HostValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Host == "" {
			http.Error(writer, `{"error":"host header is required"}`, http.StatusBadRequest)

			return
		}
		next.ServeHTTP(writer, request)
	})
}

// RecoveryMiddleware catches panics in downstream handlers and returns 500.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				getLogger().WithError(fmt.Errorf("%v", rec)).WithField("stack", string(debug.Stack())).Error("panic recovered")
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusInternalServerError)
				_, _ = writer.Write([]byte(`{"error":"internal server error"}`))
			}
		}()
		next.ServeHTTP(writer, request)
	})
}

// LoggingMiddleware logs MCP requests with tool, source, duration, byte sizes,
// and outcome. It never logs the Authorization header content. When logging is
// enabled (LOG_LEVEL set), it emits a per-tool-call access-log line at info
// level (L8) with structured fields: tool, args (non-sensitive), source,
// duration, in_bytes, out_bytes, outcome — plus the request_id added by
// WithSession. The info level (not debug) keeps the access log visible whenever
// logging is enabled (L02), matching the L8 requirement of no debug gate.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			if MaxBytesError(err) {
				http.Error(writer, `{"error":"request body too large"}`, http.StatusRequestEntityTooLarge)

				return
			}
			http.Error(writer, `{"error":"failed to read request body"}`, http.StatusInternalServerError)

			return
		}
		request.Body = io.NopCloser(bytes.NewReader(body))

		if err := checkBatchSize(body); err != nil {
			http.Error(writer, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)

			return
		}

		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: writer, statusCode: http.StatusOK}

		method := mcpRequestMethod(body)
		// Extract tool name from tools/call/<toolname> or method field
		toolName := method
		if after, ok := strings.CutPrefix(method, "tools/call/"); ok {
			toolName = after
		}

		next.ServeHTTP(lrw, request)

		duration := time.Since(start)
		outcome := "success"
		if lrw.statusCode >= serverErrorStatus {
			outcome = "error"
		}

		WithSession(request.Context(), getLogger()).WithFields(logrus.Fields{
			"tool":        toolName,
			"args":        logArgs(body),
			"source":      requestSource(request),
			"duration":    duration,
			"in_bytes":    len(body),
			"out_bytes":   lrw.bodySize.Load(),
			"outcome":     outcome,
			"http_method": request.Method,
			"path":        request.URL.Path,
			"status":      lrw.statusCode,
		}).Info("mcp_request")
	})
}

// toolAccessLogEnabled gates the per-tool access log emitted from
// WrapToolHandler. It is enabled only for the stdio transport (L08), which has
// no HTTP LoggingMiddleware; the HTTP transport already logs every request via
// LoggingMiddleware and must not emit a second per-tool line.
var toolAccessLogEnabled atomic.Bool

// EnableToolAccessLog turns on the per-tool info-level access log emitted by
// WrapToolHandler. It is set once by the stdio transport before serving.
func EnableToolAccessLog() {
	toolAccessLogEnabled.Store(true)
}

// logToolAccess records a per-tool access-log line at info level (no debug
// gate, L08) with the tool name, redacted arguments, duration and outcome,
// correlated by request_id via WithSession. Arguments are the tool's typed
// input passed through redact.MarshalJSON so any secret:"true" field is masked
// and no Authorization header exists in this path to leak. outcome is "error"
// when the handler returned an error or an IsError result, otherwise "success".
func logToolAccess(
	ctx context.Context,
	toolName string,
	input any,
	duration time.Duration,
	err error,
	result *mcp.CallToolResult,
) {
	outcome := "success"
	if err != nil || (result != nil && result.IsError) {
		outcome = "error"
	}

	args, _ := redact.MarshalJSON(input)

	WithSession(ctx, getLogger()).WithFields(logrus.Fields{
		"tool":     toolName,
		"args":     string(args),
		"source":   "stdio",
		"duration": duration,
		"outcome":  outcome,
	}).Info("mcp_request")
}

// requestArgs extracts the JSON-RPC params from a request body for logging.
// The value is only the params object (if present) and never contains the
// Authorization header, which carries the token. Empty bodies yield nil.
func requestArgs(body []byte) any {
	if len(body) == 0 {
		return nil
	}
	var req struct {
		Params json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(body, &req); err != nil || len(req.Params) == 0 {
		return nil
	}

	return req.Params
}

// logArgs derives the value logged in the HTTP access-log "args" field
// (L08/S02). It extracts the JSON-RPC params from the body, applies the same
// redaction mechanism as the stdio path (redact.MarshalJSON) and then
// sanitizes the result of ANSI escape and C0 control sequences, so neither
// secret-tagged fields nor terminal escape sequences reach the log. Empty or
// non-JSON bodies yield nil; otherwise the result is a JSON string matching
// the stdio path's string(args).
func logArgs(body []byte) any {
	params := requestArgs(body)
	if params == nil {
		return nil
	}

	raw, ok := params.(json.RawMessage)
	if !ok {
		return nil
	}

	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil
	}

	return redactArgs(v)
}

// redactArgs redacts v via redact.MarshalJSON (masking any secret:"true"
// field) and returns the JSON serialization sanitized of ANSI/C0 sequences.
// A serialization failure yields an empty string.
func redactArgs(v any) string {
	serialized, err := redact.MarshalJSON(v)
	if err != nil {
		return ""
	}

	return sanitizeOutput(string(serialized))
}

// requestSource derives the client source from the request headers, in order:
// X-Real-IP, then X-Forwarded-For, then the peer address. Values are
// comma-joined so all proxy hops are visible (L8).
func requestSource(request *http.Request) string {
	var parts []string
	if v := request.Header.Get("X-Real-IP"); v != "" {
		parts = append(parts, v)
	}
	if v := request.Header.Get("X-Forwarded-For"); v != "" {
		parts = append(parts, v)
	}
	parts = append(parts, request.RemoteAddr)

	return strings.Join(parts, ", ")
}
