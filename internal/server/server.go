package server

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/wnddd839/codebuddy-proxy/internal/accounts"
	"github.com/wnddd839/codebuddy-proxy/internal/admin"
	"github.com/wnddd839/codebuddy-proxy/internal/billing"
	"github.com/wnddd839/codebuddy-proxy/internal/config"
	"github.com/wnddd839/codebuddy-proxy/internal/gateway"
	"github.com/wnddd839/codebuddy-proxy/internal/httputil"
	"github.com/wnddd839/codebuddy-proxy/internal/oauth"
	"github.com/wnddd839/codebuddy-proxy/internal/openai"
	"github.com/wnddd839/codebuddy-proxy/internal/provider"
	"github.com/wnddd839/codebuddy-proxy/internal/strutil"
)

type Server struct {
	Svc  *gateway.Service
	HTTP *http.Server
}

func New(cfg config.Config, svc *gateway.Service) *Server {
	s := &Server{Svc: svc}
	mux := http.NewServeMux()
	// Go 1.22+ method-aware patterns (http_servemux_patterns)
	// Avoid method-specific trailing-slash subtree patterns colliding with deeper exact routes.
	mux.HandleFunc("OPTIONS /{$}", s.handleOptions)
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("HEAD /health", s.handleHealth)
	mux.HandleFunc("GET /v1/models", s.handleModelsAuth)
	mux.HandleFunc("GET /models", s.handleModelsAuth)
	mux.HandleFunc("POST /v1/chat/completions", s.handleChatAuth)
	mux.HandleFunc("POST /chat/completions", s.handleChatAuth)
	mux.HandleFunc("GET /direct-admin", s.handleAdminPage)
	mux.HandleFunc("GET /direct-admin/{$}", s.handleAdminPage)
	mux.HandleFunc("HEAD /direct-admin", s.handleAdminPage)
	mux.HandleFunc("HEAD /direct-admin/{$}", s.handleAdminPage)
	mux.HandleFunc("GET /direct-admin/codebuddy/oauth/launch", s.handleOAuthLaunch)
	mux.HandleFunc("HEAD /direct-admin/codebuddy/oauth/launch", s.handleOAuthLaunch)
	mux.HandleFunc("GET /direct-admin/codebuddy/oauth/callback", s.handleOAuthCallback)
	mux.HandleFunc("HEAD /direct-admin/codebuddy/oauth/callback", s.handleOAuthCallback)
	mux.HandleFunc("/direct-admin/api/", s.handleAdminAPIEntry)
	mux.HandleFunc("/", s.handleFallback)
	s.HTTP = &http.Server{
		Addr:              cfg.Addr(),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       0,
		WriteTimeout:      0,
		IdleTimeout:       120 * time.Second,
	}
	return s
}

func (s *Server) ListenAndServe() error {
	return s.HTTP.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.HTTP.Shutdown(ctx)
	if s.Svc != nil {
		if cerr := s.Svc.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}
	return err
}

func (s *Server) handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-API-Key")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "provider": "codebuddy", "transport": s.Svc.Config().Transport})
}

func (s *Server) handleModelsAuth(w http.ResponseWriter, r *http.Request) {
	if !s.authorizeAPI(w, r) {
		return
	}
	s.handleModels(w, r)
}

func (s *Server) handleChatAuth(w http.ResponseWriter, r *http.Request) {
	if !s.authorizeAPI(w, r) {
		return
	}
	s.handleChatCompletions(w, r)
}

func (s *Server) handleAdminPage(w http.ResponseWriter, r *http.Request) {
	if !s.authorizeAdmin(w, r) {
		return
	}
	httputil.WriteHTML(w, http.StatusOK, admin.PageHTML())
}

func (s *Server) handleAdminAPIEntry(w http.ResponseWriter, r *http.Request) {
	path := httputil.NormalizePath(r.URL.Path)
	s.handleAdminAPI(w, r, path)
}

func (s *Server) handleFallback(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		s.handleOptions(w, r)
		return
	}
	path := httputil.NormalizePath(r.URL.Path)
	httputil.WriteJSON(w, http.StatusNotFound, openai.NewError(fmt.Sprintf("Unsupported route: %s", path), "not_found_error"))
}

func (s *Server) authorizeAPI(w http.ResponseWriter, r *http.Request) bool {
	cfg := s.Svc.Config()
	if !cfg.RequireAPIKey {
		return true
	}
	token := httputil.BearerToken(r)
	if token == "" {
		token = strings.TrimSpace(r.Header.Get("X-API-Key"))
	}
	if token == "" || !secretEqual(token, cfg.APIKey) {
		httputil.WriteJSON(w, http.StatusUnauthorized, openai.NewError("Missing or invalid API key", "authentication_error"))
		return false
	}
	return true
}

func (s *Server) authorizeAdmin(w http.ResponseWriter, r *http.Request) bool {
	cfg := s.Svc.Config()
	if cfg.AdminPassword == "" {
		return true
	}
	user, pass, ok := r.BasicAuth()
	if ok && secretEqual(pass, cfg.AdminPassword) && (user == "" || user == "admin") {
		return true
	}
	token := httputil.BearerToken(r)
	if token != "" && (secretEqual(token, cfg.AdminPassword) || secretEqual(token, cfg.APIKey)) {
		return true
	}
	// Intentionally no ?password= query auth: secrets in URLs leak into proxy logs,
	// browser history, and Referer. Use Basic Auth or Bearer instead.
	w.Header().Set("WWW-Authenticate", `Basic realm="CodeBuddy Admin"`)
	httputil.WriteJSON(w, http.StatusUnauthorized, map[string]any{"ok": false, "error": "admin auth required"})
	return false
}

func secretEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	fresh := queryTruthy(r.URL.Query().Get("fresh"))
	listed, err := s.Svc.ListModels(r.Context(), fresh)
	models := listed.Models
	if err != nil || len(models) == 0 {
		models = s.Svc.ConfiguredModels()
	}
	created := time.Now().Unix()
	data := make([]map[string]any, 0, len(models))
	for _, model := range models {
		item := map[string]any{
			"id":       model.ID,
			"object":   "model",
			"created":  created,
			"owned_by": strutil.First(model.OwnedBy, "codebuddy"),
		}
		if model.Credits != "" {
			item["credits"] = model.Credits
		}
		if model.CreditMultiplier != nil {
			item["credit_multiplier"] = *model.CreditMultiplier
		}
		if model.Free != nil {
			item["free"] = *model.Free
		}
		if model.Description != "" {
			item["description"] = model.Description
		}
		data = append(data, item)
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]any{"object": "list", "data": data})
}

func (s *Server) handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Model               string           `json:"model"`
		Messages            []map[string]any `json:"messages"`
		Stream              bool             `json:"stream"`
		Tools               any              `json:"tools"`
		ToolChoice          any              `json:"tool_choice"`
		Temperature         *float64         `json:"temperature"`
		TopP                *float64         `json:"top_p"`
		MaxTokens           *int             `json:"max_tokens"`
		MaxCompletionTokens *int             `json:"max_completion_tokens"`
	}
	if err := httputil.ReadJSON(r, &body); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, openai.NewError("Invalid JSON body", "invalid_request_error"))
		return
	}
	providerModel := gateway.ResolveProviderModel(body.Model)
	promptChars := estimatePromptChars(body.Messages)
	finish := s.Svc.BeginRequest(providerModel.PublicModel, promptChars, body.Stream)
	maxTokens := 0
	if body.MaxCompletionTokens != nil && *body.MaxCompletionTokens > 0 {
		maxTokens = *body.MaxCompletionTokens
	} else if body.MaxTokens != nil && *body.MaxTokens > 0 {
		maxTokens = *body.MaxTokens
	}
	completeOpts := gateway.CompleteOptions{
		Model:               providerModel.Model,
		Messages:            body.Messages,
		Tools:               body.Tools,
		ToolChoice:          body.ToolChoice,
		Temperature:         body.Temperature,
		TopP:                body.TopP,
		MaxCompletionTokens: maxTokens,
	}

	if body.Stream {
		s.streamChat(w, r, providerModel, completeOpts, finish)
		return
	}

	completeOpts.Stream = false
	result, err := s.Svc.CompleteFromPool(r.Context(), completeOpts)
	if err != nil {
		if openai.IsClientCanceled(err) {
			// Client aborted before/during non-stream aggregation — not a gateway fault.
			finish(true, 0, 0, 0, 0, "")
			return
		}
		finish(false, 0, 0, 0, 0, err.Error())
		typ, code := openai.ClassifyUpstream(err)
		httputil.WriteJSON(w, http.StatusBadGateway, openai.NewErrorWithCode(err.Error(), typ, code))
		return
	}
	finish(true, len(result.Turn.Text), result.Bytes, int64(result.EventCount), int64(result.DeltaCount), "")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	httputil.WriteJSON(w, http.StatusOK, openai.FromTurn(result.Turn, "", providerModel.PublicModel))
}

func (s *Server) streamChat(
	w http.ResponseWriter,
	r *http.Request,
	providerModel gateway.ProviderModel,
	opts gateway.CompleteOptions,
	finish func(bool, int, int64, int64, int64, string),
) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		httputil.WriteJSON(w, http.StatusInternalServerError, openai.NewError("streaming unsupported", "internal_error"))
		return
	}
	id := fmt.Sprintf("chatcmpl_codebuddy_%d", time.Now().UnixMilli())
	var (
		mu                sync.Mutex
		started           bool
		done              bool
		streamedToolCalls int
		streamedChars     int
	)
	writeLocked := func(fn func()) {
		mu.Lock()
		defer mu.Unlock()
		if done {
			return
		}
		fn()
	}
	startStream := func() {
		if started {
			return
		}
		started = true
		w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-transform")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		_ = httputil.WriteSSE(w, openai.StreamChunkOf(id, providerModel.PublicModel, openai.Delta{Role: "assistant"}, nil))
		flusher.Flush()
	}

	// Open SSE immediately so clients don't treat upstream connect latency as a hang,
	// and so keep-alive comments have a live response to write into.
	writeLocked(startStream)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	keepAlive := s.Svc.Config().StreamKeepAlive
	if keepAlive <= 0 {
		// Frequent enough that ZCode/NewAPI do not treat long upstream TTFB as a hang.
		keepAlive = 5 * time.Second
	}
	keepAliveStop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(keepAlive)
		defer ticker.Stop()
		for {
			select {
			case <-keepAliveStop:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				writeLocked(func() {
					startStream()
					httputil.WriteSSEComment(w, "keep-alive")
				})
			}
		}
	}()
	defer close(keepAliveStop)

	opts.Stream = true
	opts.OnEvent = func(event provider.Event) {
		if event.Type != "tool_use" {
			return
		}
		writeLocked(func() {
			startStream()
			streamedToolCalls++
			_ = httputil.WriteSSE(w, openai.StreamChunkOf(id, providerModel.PublicModel, openai.Delta{
				ToolCalls: []openai.ToolCall{{
					Index: streamedToolCalls - 1,
					ID:    strutil.First(event.ID, fmt.Sprintf("call_%d", time.Now().UnixNano())),
					Type:  "function",
					Function: openai.ToolFunction{
						Name:      strutil.First(event.Name, "tool"),
						Arguments: mustJSON(event.Input),
					},
				}},
			}, nil))
		})
	}
	opts.OnDelta = func(delta string) {
		if delta == "" {
			return
		}
		writeLocked(func() {
			startStream()
			streamedChars += len(delta)
			_ = httputil.WriteSSE(w, openai.StreamChunkOf(id, providerModel.PublicModel, openai.Delta{Content: delta}, nil))
		})
	}
	result, err := s.Svc.CompleteFromPool(ctx, opts)
	if err != nil {
		if openai.IsClientCanceled(err) {
			// ZCode/browser often abort slow models (hy4-preview) and retry —
			// that surfaces as context canceled, not an upstream outage.
			finish(true, streamedChars, 0, 0, 0, "")
			writeLocked(func() { done = true })
			return
		}
		finish(false, streamedChars, 0, 0, 0, err.Error())
		typ, code := openai.ClassifyUpstream(err)
		writeLocked(func() {
			if !started {
				httputil.WriteJSON(w, http.StatusBadGateway, openai.NewErrorWithCode(err.Error(), typ, code))
				done = true
				return
			}
			errObj := map[string]any{"message": err.Error(), "type": typ}
			if code != nil {
				errObj["code"] = code
			}
			_ = httputil.WriteSSE(w, map[string]any{"error": errObj})
			httputil.WriteSSEDone(w)
			done = true
		})
		return
	}
	finish(true, len(result.Turn.Text), result.Bytes, int64(result.EventCount), int64(result.DeltaCount), "")
	writeLocked(func() {
		startStream()
		if streamedChars == 0 && result.Turn.Text != "" {
			_ = httputil.WriteSSE(w, openai.StreamChunkOf(id, providerModel.PublicModel, openai.Delta{Content: result.Turn.Text}, nil))
		}
		if len(result.Turn.ToolUses) > 0 && streamedToolCalls == 0 {
			for i, tool := range result.Turn.ToolUses {
				_ = httputil.WriteSSE(w, openai.StreamChunkOf(id, providerModel.PublicModel, openai.Delta{
					ToolCalls: []openai.ToolCall{{
						Index: i,
						ID:    strutil.First(tool.ID, fmt.Sprintf("call_%d", time.Now().UnixNano())),
						Type:  "function",
						Function: openai.ToolFunction{
							Name:      strutil.First(tool.Name, "tool"),
							Arguments: mustJSON(tool.Input),
						},
					}},
				}, nil))
			}
		}
		finishReason := "stop"
		if len(result.Turn.ToolUses) > 0 {
			finishReason = "tool_calls"
		}
		_ = httputil.WriteSSE(w, openai.StreamChunkOf(id, providerModel.PublicModel, openai.Delta{}, &finishReason))
		httputil.WriteSSEDone(w)
		done = true
	})
}

func (s *Server) handleAdminAPI(w http.ResponseWriter, r *http.Request, path string) {
	if !s.authorizeAdmin(w, r) {
		return
	}
	if !httputil.AdminMutationAllowed(r) {
		httputil.WriteJSON(w, http.StatusForbidden, map[string]any{"ok": false, "error": "cross-origin admin mutation blocked"})
		return
	}
	publicOrigin := httputil.PublicOrigin(r, s.Svc.Config().PublicBaseURL)
	switch {
	case path == "/direct-admin/api/status" && r.Method == http.MethodGet:
		httputil.WriteJSON(w, http.StatusOK, s.Svc.Status())
		return
	case path == "/direct-admin/api/pool-site" && (r.Method == http.MethodPost || r.Method == http.MethodPut):
		var body struct {
			Site string `json:"site"`
		}
		_ = httputil.ReadJSON(r, &body)
		payload, err := s.Svc.SetPoolSite(body.Site)
		if err != nil {
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		httputil.WriteJSON(w, http.StatusOK, payload)
		return
	case path == "/direct-admin/api/client-config" && r.Method == http.MethodGet:
		httputil.WriteJSON(w, http.StatusOK, s.clientConfigPayload(publicOrigin))
		return
	case path == "/direct-admin/api/client-config/generate-key" && r.Method == http.MethodPost:
		key, err := GenerateProxyAPIKey()
		if err != nil {
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		envPath := config.ResolveEnvFilePath()
		values := map[string]string{
			"CODEBUDDY_PROXY_API_KEY":         key,
			"CODEBUDDY_PROXY_REQUIRE_API_KEY": "true",
		}
		if err := config.UpsertEnvFile(envPath, values); err != nil {
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"ok":    false,
				"error": "API Key 已生成但写入 .env 失败: " + err.Error(),
			})
			return
		}
		s.Svc.SetAPIKey(key)
		payload := s.clientConfigPayload(publicOrigin)
		payload["generated"] = true
		payload["envFile"] = envPath
		payload["note"] = "新 API Key 已写入 " + envPath + "，重启后仍然有效。请同步更新 ZCode / NewAPI 等客户端里的 Key。"
		httputil.WriteJSON(w, http.StatusOK, payload)
		return
	case path == "/direct-admin/api/codebuddy/status" && r.Method == http.MethodGet:
		store, err := s.Svc.Pool.Read()
		if err != nil {
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		httputil.WriteJSON(w, http.StatusOK, accounts.SummarizeStore(store, s.Svc.Pool.Path()))
		return
	case path == "/direct-admin/api/codebuddy/oauth/session" && r.Method == http.MethodGet:
		httputil.WriteJSON(w, http.StatusOK, s.Svc.OAuthSession(publicOrigin))
		return
	case path == "/direct-admin/api/codebuddy/oauth/start" && r.Method == http.MethodPost:
		var body struct {
			Site          string `json:"site"`
			Label         string `json:"label"`
			ReuseExisting bool   `json:"reuseExisting"`
		}
		_ = httputil.ReadJSON(r, &body)
		payload, err := s.Svc.StartOAuth(r.Context(), body.Site, body.Label, publicOrigin, body.ReuseExisting)
		if err != nil {
			httputil.WriteJSON(w, http.StatusBadGateway, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		httputil.WriteJSON(w, http.StatusOK, payload)
		return
	case path == "/direct-admin/api/codebuddy/oauth/poll" && r.Method == http.MethodPost:
		payload, err := s.Svc.PollOAuth(r.Context(), publicOrigin)
		if err != nil {
			httputil.WriteJSON(w, http.StatusBadGateway, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		httputil.WriteJSON(w, http.StatusOK, payload)
		return
	case path == "/direct-admin/api/codebuddy/oauth/callback" && r.Method == http.MethodPost:
		payload, err := s.Svc.PollOAuth(r.Context(), publicOrigin)
		if err != nil {
			httputil.WriteJSON(w, http.StatusBadGateway, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		httputil.WriteJSON(w, http.StatusOK, payload)
		return
	case path == "/direct-admin/api/codebuddy/accounts" && r.Method == http.MethodGet:
		store, err := s.Svc.Pool.Read()
		if err != nil {
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		httputil.WriteJSON(w, http.StatusOK, accounts.SummarizeStore(store, s.Svc.Pool.Path()))
		return
	case path == "/direct-admin/api/codebuddy/models" && r.Method == http.MethodGet:
		fresh := true
		if q := r.URL.Query().Get("fresh"); q != "" {
			fresh = queryTruthy(q)
		}
		listed, err := s.Svc.ListModels(r.Context(), fresh)
		if err != nil {
			httputil.WriteJSON(w, http.StatusBadGateway, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		httputil.WriteJSON(w, http.StatusOK, listed)
		return
	case path == "/direct-admin/api/codebuddy/probe" && r.Method == http.MethodPost:
		listed, err := s.Svc.ListModels(r.Context(), true)
		if err != nil {
			httputil.WriteJSON(w, http.StatusBadGateway, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		httputil.WriteJSON(w, http.StatusOK, listed)
		return
	}

	if strings.HasPrefix(path, "/direct-admin/api/codebuddy/accounts/") {
		s.handleAccountAction(w, r, path)
		return
	}
	httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"ok": false, "error": "unknown admin api"})
}

func (s *Server) handleAccountAction(w http.ResponseWriter, r *http.Request, path string) {
	rest := strings.TrimPrefix(path, "/direct-admin/api/codebuddy/accounts/")
	parts := strings.Split(rest, "/")
	if len(parts) == 0 || parts[0] == "" {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "missing account id"})
		return
	}
	id := parts[0]
	action := ""
	if len(parts) > 1 {
		action = parts[1]
	}
	switch {
	case r.Method == http.MethodDelete && action == "":
		store, err := s.Svc.Pool.Delete(id)
		if err != nil {
			httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		httputil.WriteJSON(w, http.StatusOK, accounts.SummarizeStore(store, s.Svc.Pool.Path()))
		return
	case r.Method == http.MethodPost && (action == "enable" || action == "disable"):
		account, store, err := s.Svc.Pool.SetEnabled(id, action == "enable")
		if err != nil {
			httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		httputil.WriteJSON(w, http.StatusOK, map[string]any{
			"ok":       true,
			"account":  accounts.SummarizeAccount(account),
			"accounts": accounts.SummarizeStore(store, s.Svc.Pool.Path()),
		})
		return
	case r.Method == http.MethodGet && action == "usage":
		store, err := s.Svc.Pool.Read()
		if err != nil {
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		var selected accounts.Account
		found := false
		for _, item := range store.Accounts {
			if item.ID == id {
				selected = item
				found = true
				break
			}
		}
		if !found {
			httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"ok": false, "error": "account not found"})
			return
		}
		if !accounts.HasCredentials(selected) {
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "account has no credentials"})
			return
		}
		usage, err := billing.FetchAccountUsage(r.Context(), s.Svc.Provider, selected, s.Svc.Config())
		if err != nil {
			httputil.WriteJSON(w, http.StatusBadGateway, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		payload := map[string]any{
			"ok":               usage.OK,
			"provider":         usage.Provider,
			"accountId":        usage.AccountID,
			"site":             usage.Site,
			"endpoint":         usage.Endpoint,
			"officialUsageUrl": usage.OfficialUsageURL,
			"note":             usage.Note,
			"credits":          usage.Credits,
			"notify":           usage.Notify,
			"account":          accounts.SummarizeAccount(selected),
		}
		httputil.WriteJSON(w, http.StatusOK, payload)
		return
	case r.Method == http.MethodPost && action == "refresh-token":
		store, err := s.Svc.Pool.Read()
		if err != nil {
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		var selected accounts.Account
		found := false
		for _, item := range store.Accounts {
			if item.ID == id {
				selected = item
				found = true
				break
			}
		}
		if !found {
			httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"ok": false, "error": "account not found"})
			return
		}
		token, err := s.Svc.OAuth.Refresh(r.Context(), oauth.RefreshOptions{
			Site:         strutil.First(selected.Site, s.Svc.Config().Site),
			BaseURL:      strutil.First(selected.BaseURL, s.Svc.Config().BaseURL),
			AccessToken:  selected.BearerToken,
			RefreshToken: selected.RefreshToken,
		})
		if err != nil {
			httputil.WriteJSON(w, http.StatusBadGateway, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		updated := oauth.AccountFromTokenData(token, strutil.First(selected.Site, s.Svc.Config().Site), selected.Label)
		updated.ID = selected.ID
		updated.Enabled = selected.Enabled
		updated.BaseURL = selected.BaseURL
		updated.InternetEnvironment = selected.InternetEnvironment
		updated.APIEndpoint = selected.APIEndpoint
		updated.ChatCompletionsPath = selected.ChatCompletionsPath
		updated.Transport = selected.Transport
		updated.CreatedAt = selected.CreatedAt
		updated.SuccessRequests = selected.SuccessRequests
		updated.FailedRequests = selected.FailedRequests
		updated.LastUsedAt = selected.LastUsedAt
		updated.LastSelectedAt = selected.LastSelectedAt
		saved, store, err := s.Svc.Pool.ReplaceAccount(updated)
		if err != nil {
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		httputil.WriteJSON(w, http.StatusOK, map[string]any{
			"ok":       true,
			"account":  accounts.SummarizeAccount(saved),
			"accounts": accounts.SummarizeStore(store, s.Svc.Pool.Path()),
		})
		return
	default:
		httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"ok": false, "error": "unknown account action"})
	}
}

func (s *Server) handleOAuthLaunch(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	setCookie := fmt.Sprintf("cursor_codebuddy_oauth=%s; Path=/; Max-Age=900; HttpOnly; SameSite=Lax", token)
	// Authorize under the service lock; never race on a shared *OAuthSession.
	if !s.Svc.OAuthLaunchAuthorized(id, token) {
		w.Header().Set("Set-Cookie", setCookie)
		httputil.WriteHTML(w, http.StatusForbidden, admin.LaunchPage("登录入口已失效或参数不正确，请回到管理台重新生成 CodeBuddy 登录入口。", false))
		return
	}
	session := s.Svc.CurrentOAuth()
	if session.URL != "" && session.Status == "waiting" {
		w.Header().Set("Set-Cookie", setCookie)
		w.Header().Set("Cache-Control", "no-store")
		http.Redirect(w, r, session.URL, http.StatusFound)
		return
	}
	w.Header().Set("Set-Cookie", setCookie)
	httputil.WriteHTML(w, http.StatusConflict, admin.LaunchPage(strutil.First(session.Error, "请回到管理台重新发起 CodeBuddy OAuth 登录。"), false))
}

func (s *Server) handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	publicOrigin := httputil.PublicOrigin(r, s.Svc.Config().PublicBaseURL)
	setCookie := fmt.Sprintf("cursor_codebuddy_oauth=%s; Path=/; Max-Age=900; HttpOnly; SameSite=Lax", token)
	w.Header().Set("Set-Cookie", setCookie)
	if !s.Svc.OAuthLaunchAuthorized(id, token) {
		httputil.WriteHTML(w, http.StatusForbidden, admin.LaunchPage("登录回调无效，请回到管理台重试。", false))
		return
	}
	payload, err := s.Svc.PollOAuth(r.Context(), publicOrigin)
	if err != nil {
		httputil.WriteHTML(w, http.StatusBadGateway, admin.LaunchPage("CodeBuddy 登录回调失败："+err.Error(), false))
		return
	}
	ok, _ := payload["ok"].(bool)
	// Re-read a value snapshot after poll.
	liveAfter := s.Svc.CurrentOAuth()
	msg := "CodeBuddy 登录已确认，账号已导入账号池。"
	if !ok {
		msg = strutil.First(liveAfter.Error, "CodeBuddy 登录尚未完成。")
	}
	status := http.StatusOK
	if !ok {
		status = http.StatusConflict
	}
	httputil.WriteHTML(w, status, admin.LaunchPage(msg, ok))
}

func (s *Server) clientConfigPayload(publicOrigin string) map[string]any {
	cfg := s.Svc.Config()
	apiBase := strings.TrimRight(publicOrigin, "/") + "/v1"
	key := strings.TrimSpace(cfg.APIKey)
	return map[string]any{
		"ok":                 true,
		"baseUrl":            apiBase,
		"apiBase":            apiBase,
		"apiBasePath":        "/v1",
		"chatCompletionsUrl": apiBase + "/chat/completions",
		"recommendedModel":   "auto",
		"requireApiKey":      cfg.RequireAPIKey,
		"apiKeyConfigured":   key != "",
		"apiKeyPreview":      strutil.MaskSecret(key, 6),
		"apiKey":             key,
		"transport":          cfg.Transport,
		"site":               s.Svc.ActivePoolSite(),
		"poolSite":           s.Svc.ActivePoolSite(),
	}
}

// GenerateProxyAPIKey creates a gateway API key (cbp_...).
func GenerateProxyAPIKey() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return "cbp_" + hex.EncodeToString(buf), nil
}

func estimatePromptChars(messages []map[string]any) int {
	total := 0
	for _, message := range messages {
		switch v := message["content"].(type) {
		case string:
			total += len(v)
		default:
			raw, _ := json.Marshal(v)
			total += len(raw)
		}
	}
	return total
}

func mustJSON(value map[string]any) string {
	if value == nil {
		return "{}"
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func queryTruthy(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
