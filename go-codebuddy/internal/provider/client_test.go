package provider_test

import (
	"fmt"
	"testing"

	"github.com/wnddd839/codebuddy-proxy/internal/provider"
)

func TestResolveProtocolDirectDomestic(t *testing.T) {
	opts := provider.ChatOptions{
		Site:    "domestic",
		BaseURL: "https://www.codebuddy.cn",
		Domain:  "www.codebuddy.cn", // portal host must not win over chat endpoint
	}
	endpoint := provider.ResolveProtocolDirectEndpoint(opts)
	if endpoint != "https://copilot.tencent.com/v2/chat/completions" {
		t.Fatalf("unexpected endpoint: %s", endpoint)
	}
	domain := provider.ResolveProtocolDirectDomain(opts)
	if domain != "copilot.tencent.com" {
		t.Fatalf("unexpected domain: %s", domain)
	}
}

func TestAccountSiteBeatsProxyBaseURL(t *testing.T) {
	// Domestic account must stay on CN chat host even if process BaseURL is global
	// (overseas VPS / VPN / mis-set CODEBUDDY_BASE_URL).
	domestic := provider.ChatOptions{
		Site:        "domestic",
		BaseURL:     "https://www.codebuddy.ai",
		APIEndpoint: "https://www.codebuddy.ai/v2/chat/completions",
	}
	if provider.RegionOf(domestic) != "domestic" {
		t.Fatalf("region=%s", provider.RegionOf(domestic))
	}
	if got := provider.ResolveProtocolDirectEndpoint(domestic); got != "https://copilot.tencent.com/v2/chat/completions" {
		t.Fatalf("domestic endpoint=%s", got)
	}
	if got := provider.ResolveProtocolDirectDomain(domestic); got != "copilot.tencent.com" {
		t.Fatalf("domestic domain=%s", got)
	}

	global := provider.ChatOptions{
		Site:    "global",
		BaseURL: "https://www.codebuddy.cn",
	}
	// Explicit global site wins even if BaseURL looks domestic.
	if provider.RegionOf(global) != "global" {
		t.Fatalf("region=%s", provider.RegionOf(global))
	}
	if got := provider.ResolveProtocolDirectEndpoint(global); got != "https://www.codebuddy.ai/v2/chat/completions" {
		t.Fatalf("global endpoint=%s", got)
	}
}

func TestNormalizeToolChoice(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want any
	}{
		{"string auto", "auto", "auto"},
		{"object auto", map[string]any{"type": "auto"}, "auto"},
		{"object none", map[string]any{"type": "none"}, "none"},
		{"object function", map[string]any{"type": "function", "function": map[string]any{"name": "bash"}}, "auto"},
		{"nil", nil, nil},
		{"empty string", "", nil},
	}
	for _, tc := range cases {
		got := provider.NormalizeToolChoice(tc.in)
		if fmt.Sprint(got) != fmt.Sprint(tc.want) {
			t.Fatalf("%s: got %#v want %#v", tc.name, got, tc.want)
		}
	}
}

func TestNormalizeModelsFromV3Shape(t *testing.T) {
	rows := provider.NormalizeModels(map[string]any{
		"data": map[string]any{
			"models": map[string]any{
				"availableModels": []any{"auto", "glm-5.2"},
				"models": []any{
					map[string]any{"id": "auto", "name": "Auto", "supportsToolCall": true},
					map[string]any{"id": "glm-5.2", "name": "GLM", "supportsToolCall": true},
					map[string]any{"id": "hidden", "name": "Hidden"},
				},
			},
		},
	})
	if len(rows) != 2 {
		t.Fatalf("expected 2 models, got %d (%v)", len(rows), rows)
	}
}

func TestMapSSEOpenAIDelta(t *testing.T) {
	events := provider.MapSSEEvent(map[string]any{
		"choices": []any{
			map[string]any{
				"delta": map[string]any{"content": "hello"},
			},
		},
	})
	if len(events) != 1 || events[0].Type != "text_delta" || events[0].Text != "hello" {
		t.Fatalf("unexpected events: %+v", events)
	}
}

func TestEnsureUpstreamMessagesDropsEmpty(t *testing.T) {
	out := provider.EnsureUpstreamMessages([]map[string]any{
		{"role": "user", "content": ""},
		{"role": "user", "content": "hi"},
	})
	if len(out) != 1 {
		t.Fatalf("expected 1 message, got %d", len(out))
	}
}

func TestEnsureUpstreamMessagesMapsDeveloper(t *testing.T) {
	out := provider.EnsureUpstreamMessages([]map[string]any{
		{"role": "developer", "content": "You are a coding agent."},
		{"role": "user", "content": "hi"},
	})
	if len(out) != 2 {
		t.Fatalf("expected 2 messages, got %d (%v)", len(out), out)
	}
	if out[0]["role"] != "system" {
		t.Fatalf("developer role=%v want system", out[0]["role"])
	}
}
