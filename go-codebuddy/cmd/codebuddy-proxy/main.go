package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/wnddd839/codebuddy-proxy/internal/config"
	"github.com/wnddd839/codebuddy-proxy/internal/gateway"
	"github.com/wnddd839/codebuddy-proxy/internal/server"
	"github.com/wnddd839/codebuddy-proxy/internal/strutil"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()

	// First-run / binary UX: if no gateway key is configured, create one and
	// persist it to .env so clients don't break after restart.
	if strings.TrimSpace(cfg.APIKey) == "" {
		key, err := server.GenerateProxyAPIKey()
		if err != nil {
			logger.Error("failed to generate default api key", "error", err)
			os.Exit(1)
		}
		envPath := config.ResolveEnvFilePath()
		values := map[string]string{
			"CODEBUDDY_PROXY_API_KEY":         key,
			"CODEBUDDY_PROXY_REQUIRE_API_KEY": "true",
		}
		if err := config.UpsertEnvFile(envPath, values); err != nil {
			logger.Error("failed to persist default api key", "path", envPath, "error", err)
			os.Exit(1)
		}
		cfg.APIKey = key
		cfg.RequireAPIKey = true
		_ = os.Setenv("CODEBUDDY_PROXY_API_KEY", key)
		_ = os.Setenv("CODEBUDDY_PROXY_REQUIRE_API_KEY", "true")
		logger.Info("generated and saved gateway api key", "path", envPath, "apiKey", key)
	}

	svc := gateway.New(cfg, logger)
	srv := server.New(cfg, svc)

	logger.Info("codebuddy proxy starting",
		"addr", "http://"+cfg.Addr()+"/v1",
		"admin", "http://"+cfg.Addr()+"/direct-admin/",
		"transport", cfg.Transport,
		"accountsPath", cfg.AccountsPath,
		"requireApiKey", cfg.RequireAPIKey,
		"apiKeyPreview", strutil.MaskSecret(cfg.APIKey, 6),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil {
			logger.Error("server exited", "error", err)
			os.Exit(1)
		}
	case sig := <-sigCh:
		logger.Info("shutting down", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("shutdown error", "error", err)
			os.Exit(1)
		}
	}
}
