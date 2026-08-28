package gateway

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/wnddd839/codebuddy-proxy/internal/config"
)

func TestOAuthSessionResetKeepsPointerIdentity(t *testing.T) {
	cfg := config.Config{OAuthSessionTTL: config.DefaultOAuthSessionTTL, Site: "domestic"}
	svc := New(cfg, slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})))

	liveBefore := svc.LiveOAuthSession()
	if liveBefore == nil {
		t.Fatal("expected live session")
	}
	liveBefore.Status = "idle"

	// Simulate start reset without calling upstream plugin.
	svc.mu.Lock()
	svc.resetOAuthSessionLocked("domestic", "test", "http://127.0.0.1:32126")
	svc.oauth.Status = "waiting"
	svc.oauth.URL = "https://www.codebuddy.cn/login?platform=CLI&state=abc"
	svc.oauth.AuthState = "abc"
	id, token := svc.oauth.ID, svc.oauth.Token
	svc.mu.Unlock()

	liveAfter := svc.LiveOAuthSession()
	if liveBefore != liveAfter {
		t.Fatalf("session pointer identity changed; stale launch refs would break")
	}
	if !svc.OAuthLaunchAuthorized(id, token) {
		t.Fatalf("expected live session to authorize freshly generated launch credentials")
	}
	if liveAfter.Status != "waiting" || liveAfter.URL == "" {
		t.Fatalf("unexpected live session: %+v", liveAfter)
	}

	// Ensure CurrentOAuth snapshot matches live fields.
	snap := svc.CurrentOAuth()
	if snap.ID != id || snap.Token != token || snap.Status != "waiting" {
		t.Fatalf("snapshot mismatch: %+v", snap)
	}

	_ = context.Background()
}
