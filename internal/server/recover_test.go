package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoverHandlerCatchesPanic(t *testing.T) {
	panicHandler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})
	handler := recoverHandler(slog.Default(), panicHandler)

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d want 500", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	errObj, _ := body["error"].(map[string]any)
	if errObj == nil || errObj["type"] != "internal_error" {
		t.Fatalf("body=%v", body)
	}
}

func TestChatInvalidJSONReturns400(t *testing.T) {
	srv := testServer(t, true, "", "secret-key")
	req := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:32126/v1/chat/completions", strings.NewReader(`{not-json`))
	req.Header.Set("Authorization", "Bearer secret-key")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.HTTP.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestChatNoAccountsReturns502(t *testing.T) {
	srv := testServer(t, true, "", "secret-key")
	body := `{"model":"auto","messages":[{"role":"user","content":"hi"}]}`
	req := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:32126/v1/chat/completions", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer secret-key")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.HTTP.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
