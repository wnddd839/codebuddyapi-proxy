package gateway

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/wnddd839/codebuddy-proxy/internal/accounts"
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

func TestShouldRetryNextAccount(t *testing.T) {
	svc := New(config.Config{Site: "domestic"}, slog.Default())
	sel := accounts.Selection{Account: accounts.Account{ID: "acc-1", RefreshToken: "rt"}}
	base := CompleteOptions{Model: "auto"}

	tests := []struct {
		name string
		err  string
		opts CompleteOptions
		want bool
	}{
		{name: "429", err: "CodeBuddy chat completion failed with 429: too many requests", want: true},
		{name: "503", err: "CodeBuddy chat completion failed with 503: service unavailable", want: true},
		{name: "rate limit text", err: "upstream rate limit exceeded", want: true},
		{name: "11140", err: "request illegal 11140", want: true},
		{name: "11128 no retry", err: "unapproved channel 11128", want: false},
		{name: "11101 no retry", err: "tool_choice unmarshal 11101", want: false},
		{name: "401 refresh not switch", err: "failed with 401: unauthorized", want: false},
		{name: "pinned account", err: "failed with 429", opts: CompleteOptions{AccountID: "acc-1"}, want: false},
		{name: "already excluded", err: "failed with 429", opts: CompleteOptions{ExcludeIDs: []string{"acc-1"}}, want: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts := base
			if tc.opts.AccountID != "" || len(tc.opts.ExcludeIDs) > 0 {
				opts = tc.opts
			}
			got := svc.shouldRetryNextAccount(errors.New(tc.err), sel, opts)
			if got != tc.want {
				t.Fatalf("shouldRetryNextAccount(%q)=%v want %v", tc.err, got, tc.want)
			}
		})
	}
}

func TestPoolSelectRoundRobin(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/accounts.json"
	pool := accounts.NewPool(path)
	t.Cleanup(func() { _ = pool.Close() })

	mk := func(label string) accounts.Account {
		return accounts.CreateAccount(accounts.Account{
			Label:       label,
			Site:        "domestic",
			BearerToken: "token-" + label,
			Enabled:     true,
		})
	}
	a1, _, err := pool.Upsert(mk("a1"))
	if err != nil {
		t.Fatal(err)
	}
	a2, _, err := pool.Upsert(mk("a2"))
	if err != nil {
		t.Fatal(err)
	}

	s1, err := pool.Select(accounts.SelectOptions{Site: "domestic"})
	if err != nil {
		t.Fatal(err)
	}
	s2, err := pool.Select(accounts.SelectOptions{Site: "domestic"})
	if err != nil {
		t.Fatal(err)
	}
	if s1.Account.ID == s2.Account.ID {
		t.Fatalf("expected round-robin to pick different accounts, got %s twice", s1.Account.ID)
	}

	s3, err := pool.Select(accounts.SelectOptions{Site: "domestic", ExcludeIDs: []string{s1.Account.ID}})
	if err != nil {
		t.Fatal(err)
	}
	if s3.Account.ID != s2.Account.ID && s3.Account.ID != a1.ID && s3.Account.ID != a2.ID {
		t.Fatalf("unexpected account %s", s3.Account.ID)
	}
	if s3.Account.ID == s1.Account.ID {
		t.Fatalf("exclude should skip %s", s1.Account.ID)
	}
}

func BenchmarkResolveProviderModel(b *testing.B) {
	models := []string{"auto", "codebuddy:deepseek-v4", "codebuddy/deepseek-v4-flash", "deepseek-v4"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ResolveProviderModel(models[i%len(models)])
	}
}
