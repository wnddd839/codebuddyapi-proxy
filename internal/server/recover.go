package server

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/wnddd839/codebuddy-proxy/internal/httputil"
	"github.com/wnddd839/codebuddy-proxy/internal/openai"
)

// recoverHandler 捕获 handler panic，记录日志并返回 500，避免整个进程退出。
func recoverHandler(log *slog.Logger, next http.Handler) http.Handler {
	if log == nil {
		log = slog.Default()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("http handler panic",
					"panic", rec,
					"method", r.Method,
					"path", r.URL.Path,
					"stack", string(debug.Stack()),
				)
				httputil.WriteJSON(w, http.StatusInternalServerError, openai.NewError("internal server error", "internal_error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// safeCall 在回调/goroutine 内调用，防止 panic 冒泡到进程级。
func safeCall(log *slog.Logger, name string, fn func()) {
	if fn == nil {
		return
	}
	defer func() {
		if rec := recover(); rec != nil {
			if log == nil {
				log = slog.Default()
			}
			log.Error("goroutine panic recovered",
				"name", name,
				"panic", rec,
				"stack", string(debug.Stack()),
			)
		}
	}()
	fn()
}
