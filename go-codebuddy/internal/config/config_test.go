package config_test

import (
	"os"
	"testing"

	"github.com/wnddd839/codebuddy-proxy/internal/config"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("CODEBUDDY_PROXY_HOST", "")
	t.Setenv("CURSOR_DIRECT_HOST", "")
	t.Setenv("CODEBUDDY_PROXY_PORT", "")
	t.Setenv("CURSOR_DIRECT_PORT", "")
	t.Setenv("CODEBUDDY_PROXY_API_KEY", "")
	t.Setenv("CURSOR_DIRECT_API_KEY", "")
	t.Setenv("CURSOR_GATEWAY_API_KEY", "")
	_ = os.Unsetenv("CODEBUDDY_PROXY_HOST")
	_ = os.Unsetenv("CURSOR_DIRECT_HOST")

	cfg := config.Load()
	if cfg.Host != "127.0.0.1" {
		t.Fatalf("host=%q", cfg.Host)
	}
	if cfg.Port != 32126 {
		t.Fatalf("port=%d", cfg.Port)
	}
	if cfg.Transport != "protocol_direct" {
		t.Fatalf("transport=%q", cfg.Transport)
	}
	if cfg.ChatCompletionsPath != "/v2/chat/completions" {
		t.Fatalf("path=%q", cfg.ChatCompletionsPath)
	}
}

func TestNormalizeSite(t *testing.T) {
	cases := map[string]string{
		"":              "global",
		"global":        "global",
		"international": "global",
		"domestic":      "domestic",
		"CN":            "domestic",
		"china":         "domestic",
		"internal":      "domestic",
	}
	for in, want := range cases {
		if got := config.NormalizeSite(in); got != want {
			t.Fatalf("NormalizeSite(%q)=%q want %q", in, got, want)
		}
	}
}
