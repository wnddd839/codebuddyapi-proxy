package httputil

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	raw, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, `{"error":{"message":"encode failed","type":"internal_error"}}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write(raw)
}

func WriteHTML(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = io.WriteString(w, body)
}

// SSEStream 在 http.ResponseWriter 上套 16KB bufio，把一帧 SSE 的多次 Write 聚成一次 syscall，
// 仍在每帧后 Flush，保证流式低延迟。
type SSEStream struct {
	buf     *bufio.Writer
	flusher http.Flusher
}

// NewSSEStream 创建带缓冲的 SSE 写出器。size<=0 时默认 16KiB。
func NewSSEStream(w http.ResponseWriter, size int) (*SSEStream, bool) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, false
	}
	if size <= 0 {
		size = 16 << 10
	}
	return &SSEStream{buf: bufio.NewWriterSize(w, size), flusher: flusher}, true
}

func (s *SSEStream) WriteEvent(payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	frame := fmt.Appendf(nil, "data: %s\n\n", raw)
	if _, err := s.buf.Write(frame); err != nil {
		return err
	}
	if err := s.buf.Flush(); err != nil {
		return err
	}
	s.flusher.Flush()
	return nil
}

func (s *SSEStream) WriteDone() error {
	if _, err := s.buf.WriteString("data: [DONE]\n\n"); err != nil {
		return err
	}
	if err := s.buf.Flush(); err != nil {
		return err
	}
	s.flusher.Flush()
	return nil
}

func (s *SSEStream) WriteComment(comment string) error {
	frame := fmt.Appendf(nil, ": %s\n\n", comment)
	if _, err := s.buf.Write(frame); err != nil {
		return err
	}
	if err := s.buf.Flush(); err != nil {
		return err
	}
	s.flusher.Flush()
	return nil
}

// WriteSSE 无缓冲兼容路径（测试/简单 handler）；生产流式请用 SSEStream。
func WriteSSE(w http.ResponseWriter, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	frame := fmt.Appendf(nil, "data: %s\n\n", raw)
	_, err = w.Write(frame)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return err
}

func WriteSSEDone(w http.ResponseWriter) {
	_, _ = io.WriteString(w, "data: [DONE]\n\n")
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func WriteSSEComment(w http.ResponseWriter, comment string) {
	frame := fmt.Appendf(nil, ": %s\n\n", comment)
	_, _ = w.Write(frame)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func ReadJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(r.Body, 8<<20))
	decoder.UseNumber()
	return decoder.Decode(dst)
}

func BearerToken(r *http.Request) string {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if auth == "" {
		return ""
	}
	if token, ok := strings.CutPrefix(auth, "Bearer "); ok {
		return strings.TrimSpace(token)
	}
	if token, ok := strings.CutPrefix(auth, "bearer "); ok {
		return strings.TrimSpace(token)
	}
	return ""
}

func PublicOrigin(r *http.Request, configured string) string {
	if configured != "" {
		return strings.TrimRight(configured, "/")
	}
	scheme := "http"
	if r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}
	return scheme + "://" + host
}

func NormalizePath(path string) string {
	if path == "" {
		return "/"
	}
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		return strings.TrimRight(path, "/")
	}
	return path
}
