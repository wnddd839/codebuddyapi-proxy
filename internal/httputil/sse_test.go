package httputil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wnddd839/codebuddy-proxy/internal/openai"
)

type writeCounter struct {
	http.ResponseWriter
	writes int
}

func (w *writeCounter) Write(p []byte) (int, error) {
	w.writes++
	return w.ResponseWriter.Write(p)
}

func (w *writeCounter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func TestSSEStreamWriteStreamChunk(t *testing.T) {
	rec := httptest.NewRecorder()
	sse, ok := NewSSEStream(rec, 1024)
	if !ok {
		t.Fatal("expected flusher")
	}
	chunk := openai.StreamChunkOf("id1", "auto", openai.Delta{Content: "hello"}, nil)
	want, err := json.Marshal(chunk)
	if err != nil {
		t.Fatal(err)
	}
	if err := sse.WriteEvent(chunk); err != nil {
		t.Fatal(err)
	}
	if err := sse.Flush(); err != nil {
		t.Fatal(err)
	}
	body := rec.Body.String()
	if len(body) < 6 {
		t.Fatalf("body=%q", body)
	}
	line := body[len("data: "):]
	end := 0
	for end < len(line) && line[end] != '\n' {
		end++
	}
	if string(line[:end]) != string(want) {
		t.Fatalf("payload mismatch\nwant: %s\ngot:  %s", want, line[:end])
	}
}

func TestSSEStreamWriteEvent(t *testing.T) {
	rec := httptest.NewRecorder()
	sse, ok := NewSSEStream(rec, 1024)
	if !ok {
		t.Fatal("expected flusher")
	}
	if err := sse.WriteEvent(map[string]any{"ok": true}); err != nil {
		t.Fatal(err)
	}
	if err := sse.Flush(); err != nil {
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

func TestSSEStreamAggregatesBurstWrites(t *testing.T) {
	rec := httptest.NewRecorder()
	counter := &writeCounter{ResponseWriter: rec}
	sse, ok := NewSSEStream(counter, 16<<10)
	if !ok {
		t.Fatal("expected flusher")
	}

	const frames = 3000
	payload := map[string]any{
		"id":     "chatcmpl_x",
		"object": "chat.completion.chunk",
		"choices": []map[string]any{{
			"index": 0,
			"delta": map[string]any{"content": "字"},
		}},
	}
	for range frames {
		if err := sse.WriteEvent(payload); err != nil {
			t.Fatal(err)
		}
	}
	if err := sse.WriteDone(); err != nil {
		t.Fatal(err)
	}

	ratio := float64(counter.writes) / float64(frames)
	t.Logf("frames=%d underlying_writes=%d ratio=%.3f", frames, counter.writes, ratio)
	if counter.writes >= frames/10 {
		t.Fatalf("buffering ineffective: writes=%d frames=%d ratio=%.3f want << 1", counter.writes, frames, ratio)
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

func BenchmarkSSEStreamWriteEventBurst(b *testing.B) {
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
		b.StopTimer()
		rec := httptest.NewRecorder()
		sse, _ := NewSSEStream(rec, 16<<10)
		b.StartTimer()
		for range 100 {
			_ = sse.WriteEvent(payload)
		}
		_ = sse.WriteDone()
	}
}
