package handlers

import (
	"testing"
)

// TestBurstForLimit verifies the token-bucket burst is exactly double the
// per-second rate limit. Kills the ARITHMETIC_BASE mutants on `limit * 2`
// inside burstForLimit (server.go), which would halve the burst and reject
// legitimate short request spikes.
func TestBurstForLimit(t *testing.T) {
	t.Parallel()

	cases := []struct {
		limit int
		want  int
	}{
		{0, 0},
		{1, 2},
		{2, 4},
		{5, 10},
		{100, 200},
	}
	for _, tc := range cases {
		if got := burstForLimit(tc.limit); got != tc.want {
			t.Errorf("burstForLimit(%d) = %d, want %d", tc.limit, got, tc.want)
		}
	}
}
