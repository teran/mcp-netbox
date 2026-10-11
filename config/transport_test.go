package config

import "testing"

func TestValidateTransport(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		transport   string
		netboxToken string
		wantError   bool
	}{
		{"empty defaults to http", "", "", false},
		{"http without token", TransportHTTP, "", false},
		{"http with token", TransportHTTP, "secret", false},
		{"stdio with token", TransportStdio, "secret", false},
		{"stdio without token", TransportStdio, "", true},
		{"stdio with blank token", TransportStdio, "   ", true},
		{"invalid transport", "ws", "", true},
		{"non-string value", "1", "", true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := Config{
				NetBoxURL:          "http://192.168.1.1:8000",
				AllowPrivateNetBox: true,
				NetBoxToken:        testCase.netboxToken,
				Transport:          testCase.transport,
				RateLimitGlobal:    100,
				RateLimitPerClient: 10,
				WriteTimeout:       300e9,
			}
			err := cfg.validate()
			if testCase.wantError && err == nil {
				t.Errorf("validate() = nil, want error")
			}
			if !testCase.wantError && err != nil {
				t.Errorf("validate() = %v, want nil", err)
			}
		})
	}
}
