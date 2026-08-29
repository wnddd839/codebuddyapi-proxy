package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wnddd839/codebuddy-proxy/internal/config"
	"github.com/wnddd839/codebuddy-proxy/internal/gateway"
)

func testServer(t *testing.T, requireAPIKey bool, adminPassword, apiKey string) *Server {
	t.Helper()
	cfg := config.Config{
		Host:          "127.0.0.1",
		Port:          32126,
		RequireAPIKey: requireAPIKey,
		APIKey:        apiKey,
		AdminPassword: adminPassword,
		AccountsPath:  filepath.Join(t.TempDir(), "accounts.json"),
		Site:          "domestic",
		Transport:     config.DefaultTransport,
	}
	svc := gateway.New(cfg, slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})))
	return New(cfg, svc)
}

func TestAuthorizeAPIKeyRequired(t *testing.T) {
	srv := testServer(t, true, "", "secret-key")
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:32126/v1/models", nil)
	rec := httptest.NewRecorder()
	srv.HTTP.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d want 401", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:32126/v1/models", nil)
	req.Header.Set("Authorization", "Bearer secret-key")
	rec = httptest.NewRecorder()
	srv.HTTP.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("with key status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestAdminOpenWithoutPassword(t *testing.T) {
	srv := testServer(t, true, "", "secret-key")
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:32126/direct-admin/", nil)
	rec := httptest.NewRecorder()
	srv.HTTP.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("admin status=%d", rec.Code)
	}
	body, _ := io.ReadAll(rec.Body)
	if !strings.Contains(string(body), "CodeBuddy") {
		t.Fatalf("unexpected admin body")
	}
}

func TestAdminCSRFBlocksCrossOriginMutation(t *testing.T) {
	srv := testServer(t, false, "", "")
	req := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:32126/direct-admin/api/pool-site", strings.NewReader(`{"site":"global"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://evil.example")
	req.Host = "127.0.0.1:32126"
	rec := httptest.NewRecorder()
	srv.HTTP.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("csrf status=%d body=%s", rec.Code, rec.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["ok"] != false {
		t.Fatalf("payload=%v", payload)
	}
}

func TestAdminCSRFAllowsSameOriginMutation(t *testing.T) {
	srv := testServer(t, false, "", "")
	t.Setenv("CODEBUDDY_PROXY_ENV_FILE", filepath.Join(t.TempDir(), ".env"))
	req := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:32126/direct-admin/api/pool-site", strings.NewReader(`{"site":"domestic"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://127.0.0.1:32126")
	req.Host = "127.0.0.1:32126"
	rec := httptest.NewRecorder()
	srv.HTTP.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("same-origin status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestAdminPasswordRequiredWhenConfigured(t *testing.T) {
	srv := testServer(t, false, "admin-pass", "")
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:32126/direct-admin/api/status", nil)
	rec := httptest.NewRecorder()
	srv.HTTP.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d want 401", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:32126/direct-admin/api/status", nil)
	req.SetBasicAuth("admin", "admin-pass")
	rec = httptest.NewRecorder()
	srv.HTTP.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("authed status=%d body=%s", rec.Code, rec.Body.String())
	}

	// Bearer admin password also works.
	req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:32126/direct-admin/api/status", nil)
	req.Header.Set("Authorization", "Bearer admin-pass")
	rec = httptest.NewRecorder()
	srv.HTTP.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("bearer status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestAdminRejectsPasswordQueryParam(t *testing.T) {
	srv := testServer(t, false, "admin-pass", "")
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:32126/direct-admin/api/status?password=admin-pass", nil)
	rec := httptest.NewRecorder()
	srv.HTTP.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("query password must be rejected, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestHealth(t *testing.T) {
	srv := testServer(t, false, "", "")
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:32126/health", nil)
	rec := httptest.NewRecorder()
	srv.HTTP.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("health=%d", rec.Code)
	}
}
