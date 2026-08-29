package httputil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdminMutationAllowed(t *testing.T) {
	get := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:32126/direct-admin/api/status", nil)
	if !AdminMutationAllowed(get) {
		t.Fatal("GET should always be allowed")
	}

	postCLI := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:32126/direct-admin/api/pool-site", nil)
	if !AdminMutationAllowed(postCLI) {
		t.Fatal("CLI POST without Origin/Referer should be allowed")
	}

	postOK := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:32126/direct-admin/api/pool-site", nil)
	postOK.Host = "127.0.0.1:32126"
	postOK.Header.Set("Origin", "http://127.0.0.1:32126")
	if !AdminMutationAllowed(postOK) {
		t.Fatal("same-origin Origin should be allowed")
	}

	postBad := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:32126/direct-admin/api/pool-site", nil)
	postBad.Host = "127.0.0.1:32126"
	postBad.Header.Set("Origin", "https://evil.example")
	if AdminMutationAllowed(postBad) {
		t.Fatal("cross-origin Origin must be blocked")
	}

	postRefOK := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:32126/direct-admin/api/pool-site", nil)
	postRefOK.Host = "127.0.0.1:32126"
	postRefOK.Header.Set("Referer", "http://127.0.0.1:32126/direct-admin/")
	if !AdminMutationAllowed(postRefOK) {
		t.Fatal("same-origin Referer should be allowed")
	}
}
