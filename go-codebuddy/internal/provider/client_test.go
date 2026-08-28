package provider_test

import (
	"testing"

	"github.com/wnddd839/codebuddy-proxy/internal/provider"
)

func TestResolveProtocolDirectDomestic(t *testing.T) {
	endpoint := provider.ResolveProtocolDirectEndpoint(provider.ChatOptions{
		Site:    "domestic",
		BaseURL: "https://www.codebuddy.cn",
	})
	if endpoint != "https://copilot.tencent.com/v2/chat/completions" {
		t.Fatalf("unexpected endpoint: %s", endpoint)
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
