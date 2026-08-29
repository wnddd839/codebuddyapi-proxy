package httputil

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestSSEStreamWriteEvent(t *testing.T) {
	rec := httptest.NewRecorder()
	sse, ok := NewSSEStream(rec, 1024)
	if !ok {
		t.Fatal("expected flusher")
	}
	if err := sse.WriteEvent(map[string]any{"ok": true}); err != nil {
		t.Fatal(err)
	}
	if err := sse.WriteDone(); err != nil {
		t.Fatal(err)
	}
	body := rec.Body.String()
	if body == "" || body[:5] != "data:" {
		t.Fatalf("body=%q", body)
	}
	var payload map[string]any
	line := body
	if idx := len("data: "); len(line) > idx {
		raw := line[idx:]
		end := 0
		for end < len(raw) && raw[end] != '\n' {
			end++
		}
		if err := json.Unmarshal([]byte(raw[:end]), &payload); err != nil {
			t.Fatal(err)
		}
	}
	if payload["ok"] != true {
		t.Fatalf("payload=%v body=%q", payload, body)
	}
}

func BenchmarkWriteSSE(b *testing.B) {
	payload := map[string]any{
		"id":     "chatcmpl_x",
		"object": "chat.completion.chunk",
		"choices": []map[string]any{{
			"index": 0,
			"delta": map[string]any{"content": "字"},
		}},
	}
	b.ReportAllocs()
	for b.Loop() {
		rec := httptest.NewRecorder()
		_ = WriteSSE(rec, payload)
	}
}

func BenchmarkSSEStreamWriteEvent(b *testing.B) {
	payload := map[string]any{
		"id":     "chatcmpl_x",
		"object": "chat.completion.chunk",
		"choices": []map[string]any{{
			"index": 0,
			"delta": map[string]any{"content": "字"},
		}},
	}
	rec := httptest.NewRecorder()
	sse, _ := NewSSEStream(rec, 16<<10)
	b.ReportAllocs()
	for b.Loop() {
		_ = sse.WriteEvent(payload)
	}
}
