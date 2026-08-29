package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadDotEnv loads KEY=VALUE pairs from nearby .env files into the process
// environment. Existing process env vars always win (are not overwritten).
//
// Search order:
//  1. CODEBUDDY_PROXY_ENV_FILE (if set)
//  2. ./.env (current working directory)
//  3. <executable-dir>/.env
//  4. ../.env (repo root when running from go-codebuddy/)
func LoadDotEnv() (loaded []string, err error) {
	candidates := make([]string, 0, 4)
	if custom := strings.TrimSpace(os.Getenv("CODEBUDDY_PROXY_ENV_FILE")); custom != "" {
		candidates = append(candidates, expandHome(custom))
	}
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		candidates = append(candidates, filepath.Join(cwd, ".env"))
		candidates = append(candidates, filepath.Join(cwd, "..", ".env"))
	}
	if exe, exeErr := os.Executable(); exeErr == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), ".env"))
	}

	seen := map[string]struct{}{}
	for _, path := range candidates {
		abs, absErr := filepath.Abs(path)
		if absErr != nil {
			continue
		}
		if _, ok := seen[abs]; ok {
			continue
		}
		seen[abs] = struct{}{}
		info, statErr := os.Stat(abs)
		if statErr != nil || info.IsDir() {
			continue
		}
		if loadErr := loadDotEnvFile(abs); loadErr != nil {
			return loaded, fmt.Errorf("load %s: %w", abs, loadErr)
		}
		loaded = append(loaded, abs)
	}
	return loaded, nil
}

func loadDotEnvFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// Allow long values (JWT-ish) without failing.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		value = strings.TrimSpace(value)
		value = unquoteEnvValue(value)
		// Never override real process env.
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("line %d: set %s: %w", lineNo, key, err)
		}
	}
	return scanner.Err()
}

func unquoteEnvValue(value string) string {
	if len(value) < 2 {
		return value
	}
	if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
		return value[1 : len(value)-1]
	}
	// Strip inline comments for unquoted values: KEY=val # comment
	if i := strings.Index(value, " #"); i >= 0 {
		return strings.TrimSpace(value[:i])
	}
	return value
}

// ResolveEnvFilePath picks where generated settings should be persisted.
func ResolveEnvFilePath() string {
	if custom := strings.TrimSpace(os.Getenv("CODEBUDDY_PROXY_ENV_FILE")); custom != "" {
		return expandHome(custom)
	}
	if cwd, err := os.Getwd(); err == nil {
		cwdEnv := filepath.Join(cwd, ".env")
		if _, err := os.Stat(cwdEnv); err == nil {
			return cwdEnv
		}
		parentEnv := filepath.Join(cwd, "..", ".env")
		if abs, absErr := filepath.Abs(parentEnv); absErr == nil {
			if _, err := os.Stat(abs); err == nil {
				return abs
			}
		}
		if exe, exeErr := os.Executable(); exeErr == nil {
			exeEnv := filepath.Join(filepath.Dir(exe), ".env")
			if _, err := os.Stat(exeEnv); err == nil {
				return exeEnv
			}
		}
		return cwdEnv
	}
	if exe, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exe), ".env")
	}
	return ".env"
}

// UpsertEnvFile sets keys in a .env file (create if missing). Preserves other
// lines/comments. Values are written unquoted when safe.
func UpsertEnvFile(path string, values map[string]string) error {
	if path == "" {
		return fmt.Errorf("empty env file path")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && !os.IsExist(err) {
		// Dir may be "." — ignore.
		if filepath.Dir(path) != "." {
			return err
		}
	}

	existing := []string{}
	if raw, err := os.ReadFile(path); err == nil {
		text := strings.ReplaceAll(string(raw), "\r\n", "\n")
		text = strings.ReplaceAll(text, "\r", "\n")
		existing = strings.Split(text, "\n")
		// Drop trailing empty split artifact.
		if len(existing) > 0 && existing[len(existing)-1] == "" {
			existing = existing[:len(existing)-1]
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	pending := map[string]string{}
	for k, v := range values {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		pending[k] = strings.TrimSpace(v)
	}

	out := make([]string, 0, len(existing)+len(pending)+2)
	seen := map[string]struct{}{}
	for _, line := range existing {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || !strings.Contains(trimmed, "=") {
			out = append(out, line)
			continue
		}
		work := trimmed
		if strings.HasPrefix(work, "export ") {
			work = strings.TrimSpace(strings.TrimPrefix(work, "export "))
		}
		key, _, ok := strings.Cut(work, "=")
		if !ok {
			out = append(out, line)
			continue
		}
		key = strings.TrimSpace(key)
		if val, hit := pending[key]; hit {
			out = append(out, key+"="+escapeEnvValue(val))
			seen[key] = struct{}{}
			continue
		}
		out = append(out, line)
	}

	// Append missing keys.
	ordered := []string{
		"CODEBUDDY_PROXY_API_KEY",
		"CODEBUDDY_PROXY_ADMIN_PASSWORD",
		"CODEBUDDY_PROXY_REQUIRE_API_KEY",
		"CODEBUDDY_SITE",
		"CODEBUDDY_BASE_URL",
		"CODEBUDDY_INTERNET_ENVIRONMENT",
	}
	for _, key := range ordered {
		if _, ok := seen[key]; ok {
			continue
		}
		val, ok := pending[key]
		if !ok {
			continue
		}
		out = append(out, key+"="+escapeEnvValue(val))
		seen[key] = struct{}{}
	}
	for key, val := range pending {
		if _, ok := seen[key]; ok {
			continue
		}
		out = append(out, key+"="+escapeEnvValue(val))
	}

	content := strings.Join(out, "\n") + "\n"
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(content), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func escapeEnvValue(value string) string {
	if value == "" {
		return ""
	}
	if strings.ContainsAny(value, " \t#\"'`") || strings.Contains(value, "\n") {
		escaped := strings.ReplaceAll(value, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `"`, `\"`)
		return `"` + escaped + `"`
	}
	return value
}
