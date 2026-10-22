package netbox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAuthHeader verifies the keyword selection for the outbound Authorization
// header: v2 tokens (nbt_ prefix) use Bearer, everything else (v1 plaintext)
// uses Token.
func TestAuthHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "v2 token uses Bearer",
			input: "nbt_abcdef.123456",
			want:  "Bearer nbt_abcdef.123456",
		},
		{
			name:  "v2 token with only key prefix",
			input: "nbt_single",
			want:  "Bearer nbt_single",
		},
		{
			name:  "v1 plaintext token uses Token",
			input: "0123456789abcdef",
			want:  "Token 0123456789abcdef",
		},
		{
			name:  "empty token is treated as v1",
			input: "",
			want:  "Token ",
		},
		{
			name:  "nbt_ not at the start is v1",
			input: "xxx_nbt_yyy",
			want:  "Token xxx_nbt_yyy",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := authHeader(testCase.input); got != testCase.want {
				t.Errorf("authHeader(%q) = %q, want %q", testCase.input, got, testCase.want)
			}
		})
	}
}

// TestDoRequest_AuthHeaderVersions asserts the outbound Authorization header is
// set correctly on the resty request for both v1 and v2 tokens.
func TestDoRequest_AuthHeaderVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		wantHeader string
	}{
		{
			name:       "v2 token",
			input:      "nbt_key.data",
			wantHeader: "Bearer nbt_key.data",
		},
		{
			name:       "v1 token",
			input:      "legacy-plaintext-token",
			wantHeader: "Token legacy-plaintext-token",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var sawAuth string
			srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				sawAuth = req.Header.Get("Authorization")
				writer.Header().Set("Content-Type", "application/json")
				_, _ = writer.Write([]byte(`{"name":"test"}`))
			}))
			defer srv.Close()

			client := NewClient(srv.URL, http.DefaultClient)
			_, err := client.doRequest(
				context.Background(), testCase.input, http.MethodGet, "/api/dcim/sites/", nil, nil,
			)
			if err != nil {
				t.Fatalf("doRequest() returned error: %v", err)
			}

			if sawAuth != testCase.wantHeader {
				t.Errorf("Authorization = %q, want %q", sawAuth, testCase.wantHeader)
			}
		})
	}
}
