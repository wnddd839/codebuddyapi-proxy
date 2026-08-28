package server

import (
	"context"
	"crypto/rand"
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
	Cfg  config.Config
	Svc  *gateway.Service
	HTTP *http.Server
}

func New(cfg config.Config, svc *gateway.Service) *Server {
	s := &Server{Cfg: cfg, Svc: svc}
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
	return s.HTTP.Shutdown(ctx)
}

func (s *Server) handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-API-Key")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "provider": "codebuddy", "transport": s.Cfg.Transport})
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

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	s.handleFallback(w, r)
}

func (s *Server) authorizeAPI(w http.ResponseWriter, r *http.Request) bool {
	if !s.Cfg.RequireAPIKey {
		return true
	}
	token := httputil.BearerToken(r)
	if token == "" {
		token = strings.TrimSpace(r.Header.Get("X-API-Key"))
	}
	if token == "" || token != s.Cfg.APIKey {
		httputil.WriteJSON(w, http.StatusUnauthorized, openai.NewError("Missing or invalid API key", "authentication_error"))
		return false
	}
	return true
}

func (s *Server) authorizeAdmin(w http.ResponseWriter, r *http.Request) bool {
	if s.Cfg.AdminPassword == "" {
		return true
	}
	user, pass, ok := r.BasicAuth()
	if ok && pass == s.Cfg.AdminPassword && (user == "" || user == "admin") {
		return true
	}
	token := httputil.BearerToken(r)
	if token != "" && (token == s.Cfg.AdminPassword || token == s.Cfg.APIKey) {
		return true
	}
	if q := strings.TrimSpace(r.URL.Query().Get("password")); q != "" && q == s.Cfg.AdminPassword {
		return true
	}
	w.Header().Set("WWW-Authenticate", `Basic realm="CodeBuddy Admin"`)
	httputil.WriteJSON(w, http.StatusUnauthorized, map[string]any{"ok": false, "error": "admin auth required"})
	return false
}

func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	listed, err := s.Svc.ListModels(r.Context(), true)
	models := listed.Models
	if err != nil || len(models) == 0 {
		models = s.Svc.ConfiguredModels()
	}
	created := time.Now().Unix()
	data := make([]map[string]any, 0, len(models))
	for _, model := range models {
		data = append(data, map[string]any{
			"id":       model.ID,
			"object":   "model",
			"created":  created,
			"owned_by": strutil.First(model.OwnedBy, "codebuddy"),
		})
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]any{"object": "list", "data": data})
}

func (s *Server) handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Model      string           `json:"model"`
		Messages   []map[string]any `json:"messages"`
		Stream     bool             `json:"stream"`
		Tools      any              `json:"tools"`
		ToolChoice any              `json:"tool_choice"`
	}
	if err := httputil.ReadJSON(r, &body); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, openai.NewError("Invalid JSON body", "invalid_request_error"))
		return
	}
	providerModel := gateway.ResolveProviderModel(body.Model)
	promptChars := estimatePromptChars(body.Messages)
	finish := s.Svc.BeginRequest(providerModel.PublicModel, promptChars, body.Stream)

	if body.Stream {
		s.streamChat(w, r, body.Messages, providerModel, body.Tools, body.ToolChoice, finish)
		return
	}

	result, err := s.Svc.CompleteFromPool(r.Context(), gateway.CompleteOptions{
		Model:      providerModel.Model,
		Messages:   body.Messages,
		Stream:     false,
		Tools:      body.Tools,
		ToolChoice: body.ToolChoice,
	})
	if err != nil {
		finish(false, 0, 0, 0, 0, err.Error())
		httputil.WriteJSON(w, http.StatusBadGateway, openai.NewError(err.Error(), "upstream_error"))
		return
	}
	finish(true, len(result.Turn.Text), result.Bytes, int64(result.EventCount), int64(result.DeltaCount), "")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	httputil.WriteJSON(w, http.StatusOK, openai.FromTurn(result.Turn, "", providerModel.PublicModel))
}

func (s *Server) streamChat(
	w http.ResponseWriter,
	r *http.Request,
	messages []map[string]any,
	providerModel gateway.ProviderModel,
	tools any,
	toolChoice any,
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

	keepAlive := s.Cfg.StreamKeepAlive
	if keepAlive <= 0 {
		keepAlive = 15 * time.Second
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

	result, err := s.Svc.CompleteFromPool(ctx, gateway.CompleteOptions{
		Model:      providerModel.Model,
		Messages:   messages,
		Stream:     true,
		Tools:      tools,
		ToolChoice: toolChoice,
		OnEvent: func(event provider.Event) {
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
		},
		OnDelta: func(delta string) {
			if delta == "" {
				return
			}
			writeLocked(func() {
				startStream()
				streamedChars += len(delta)
				_ = httputil.WriteSSE(w, openai.StreamChunkOf(id, providerModel.PublicModel, openai.Delta{Content: delta}, nil))
			})
		},
	})
	if err != nil {
		finish(false, streamedChars, 0, 0, 0, err.Error())
		writeLocked(func() {
			if !started {
				httputil.WriteJSON(w, http.StatusBadGateway, openai.NewError(err.Error(), "upstream_error"))
				done = true
				return
			}
			_ = httputil.WriteSSE(w, map[string]any{"error": map[string]any{"message": err.Error(), "type": "upstream_error"}})
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
	publicOrigin := httputil.PublicOrigin(r, s.Cfg.PublicBaseURL)
	switch {
	case path == "/direct-admin/api/status" && r.Method == http.MethodGet:
		httputil.WriteJSON(w, http.StatusOK, s.Svc.Status())
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
		if strings.TrimSpace(s.Cfg.AdminPassword) == "" {
			values["CODEBUDDY_PROXY_ADMIN_PASSWORD"] = key
			s.Cfg.AdminPassword = key
		}
		if err := config.UpsertEnvFile(envPath, values); err != nil {
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"ok":    false,
				"error": "API Key 已生成但写入 .env 失败: " + err.Error(),
			})
			return
		}
		s.Cfg.APIKey = key
		s.Cfg.RequireAPIKey = true
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
		listed, err := s.Svc.ListModels(r.Context(), true)
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
		usage, err := billing.FetchAccountUsage(r.Context(), s.Svc.Provider, selected, s.Cfg)
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
			Site:         strutil.First(selected.Site, s.Cfg.Site),
			BaseURL:      strutil.First(selected.BaseURL, s.Cfg.BaseURL),
			AccessToken:  selected.BearerToken,
			RefreshToken: selected.RefreshToken,
		})
		if err != nil {
			httputil.WriteJSON(w, http.StatusBadGateway, map[string]any{"ok": false, "error": err.Error()})
			return
		}
		updated := oauth.AccountFromTokenData(token, strutil.First(selected.Site, s.Cfg.Site), selected.Label)
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
	publicOrigin := httputil.PublicOrigin(r, s.Cfg.PublicBaseURL)
	setCookie := fmt.Sprintf("cursor_codebuddy_oauth=%s; Path=/; Max-Age=900; HttpOnly; SameSite=Lax", token)
	// Always read the live session object; never cache a start-time snapshot.
	live := s.Svc.LiveOAuthSession()
	if id == "" || token == "" || id != live.ID || token != live.Token {
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
	_ = publicOrigin
}

func (s *Server) handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	publicOrigin := httputil.PublicOrigin(r, s.Cfg.PublicBaseURL)
	setCookie := fmt.Sprintf("cursor_codebuddy_oauth=%s; Path=/; Max-Age=900; HttpOnly; SameSite=Lax", token)
	w.Header().Set("Set-Cookie", setCookie)
	live := s.Svc.LiveOAuthSession()
	if id == "" || token == "" || id != live.ID || token != live.Token {
		httputil.WriteHTML(w, http.StatusForbidden, admin.LaunchPage("登录回调无效，请回到管理台重试。", false))
		return
	}
	payload, err := s.Svc.PollOAuth(r.Context(), publicOrigin)
	if err != nil {
		httputil.WriteHTML(w, http.StatusBadGateway, admin.LaunchPage("CodeBuddy 登录回调失败："+err.Error(), false))
		return
	}
	ok, _ := payload["ok"].(bool)
	// Re-read live session after poll; do not trust a stale snapshot.
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
	apiBase := strings.TrimRight(publicOrigin, "/") + "/v1"
	key := strings.TrimSpace(s.Cfg.APIKey)
	return map[string]any{
		"ok":                 true,
		"baseUrl":            apiBase,
		"apiBase":            apiBase,
		"apiBasePath":        "/v1",
		"chatCompletionsUrl": apiBase + "/chat/completions",
		"recommendedModel":   "auto",
		"requireApiKey":      s.Cfg.RequireAPIKey,
		"apiKeyConfigured":   key != "",
		"apiKeyPreview":      maskAPIKey(key, 6),
		"apiKey":             key,
		"transport":          s.Cfg.Transport,
		"site":               s.Cfg.Site,
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

func maskAPIKey(value string, visible int) string {
	text := strings.TrimSpace(value)
	if text == "" {
		return ""
	}
	if visible < 1 {
		visible = 1
	}
	if len(text) <= visible*2 {
		if visible > len(text) {
			visible = len(text)
		}
		return text[:visible] + "..."
	}
	return text[:visible] + "..." + text[len(text)-visible:]
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
