package gateway

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/wnddd839/codebuddy-proxy/internal/config"
)

func TestResolveProviderModel(t *testing.T) {
	tests := []struct {
		in     string
		model  string
		public string
	}{
		{"", "auto", "auto"},
		{"codebuddy", "auto", "auto"},
		{"codebuddy:deepseek-v4", "deepseek-v4", "deepseek-v4"},
		{"codebuddy/deepseek-v4", "deepseek-v4", "deepseek-v4"},
		{"deepseek-v4-flash", "deepseek-v4-flash", "deepseek-v4-flash"},
	}
	for _, tc := range tests {
		got := ResolveProviderModel(tc.in)
		if got.Model != tc.model || got.PublicModel != tc.public {
			t.Fatalf("ResolveProviderModel(%q) = %+v, want model=%q public=%q", tc.in, got, tc.model, tc.public)
		}
	}
}

func TestCompleteFromPoolRetryDepthExceeded(t *testing.T) {
	cfg := config.Config{
		Host:         "127.0.0.1",
		Port:         32126,
		AccountsPath: t.TempDir() + "/accounts.json",
		Site:         "domestic",
	}
	svc := New(cfg, slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})))

	_, err := svc.CompleteFromPool(context.Background(), CompleteOptions{
		Model:      "auto",
		Messages:   []map[string]any{{"role": "user", "content": "hi"}},
		RetryDepth: maxCompleteRetryDepth,
	})
	if err == nil {
		t.Fatal("expected retry depth error")
	}
	if !strings.Contains(err.Error(), "重试深度") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func BenchmarkResolveProviderModel(b *testing.B) {
	models := []string{"auto", "codebuddy:deepseek-v4", "codebuddy/deepseek-v4-flash", "deepseek-v4"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ResolveProviderModel(models[i%len(models)])
	}
}
