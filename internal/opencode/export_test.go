package opencode_test

import (
	"testing"

	"github.com/wnddd839/codebuddy-proxy/internal/models"
	"github.com/wnddd839/codebuddy-proxy/internal/opencode"
)

func TestVariantsForGLMFlash(t *testing.T) {
	variants := opencode.VariantsForModel(models.Model{
		ID:                "glm-5.3-flash",
		SupportsReasoning: true,
		Reasoning: map[string]any{
			"supportedEfforts":   []any{"low", "high", "max"},
			"canDisableThinking": true,
		},
	})
	if variants["low"]["reasoningEffort"] != "low" {
		t.Fatalf("low=%v", variants["low"])
	}
	if variants["max"]["reasoningEffort"] != "max" {
		t.Fatalf("max=%v", variants["max"])
	}
	if variants["none"]["thinking"] == nil {
		t.Fatalf("expected none thinking disable, got %v", variants["none"])
	}
}

func TestModelListFields(t *testing.T) {
	fields := opencode.ModelListFields(models.Model{
		ID:                "glm-5.3-flash",
		SupportsReasoning: true,
		Reasoning: map[string]any{
			"supportedEfforts": []any{"low", "high"},
			"defaultEffort":    "high",
		},
	})
	if fields["reasoning"] != true {
		t.Fatalf("reasoning=%v want true", fields["reasoning"])
	}
	cfg, ok := fields["reasoning_config"].(map[string]any)
	if !ok || cfg["defaultEffort"] != "high" {
		t.Fatalf("reasoning_config=%v", fields["reasoning_config"])
	}
	if fields["interleaved"] == nil {
		t.Fatalf("missing interleaved")
	}
	variants, ok := fields["variants"].(map[string]map[string]any)
	if !ok || variants["high"]["reasoningEffort"] != "high" {
		t.Fatalf("variants=%v", fields["variants"])
	}
}

func TestLiteLLMModelInfoEntry(t *testing.T) {
	entry := opencode.LiteLLMModelInfoEntry(models.Model{
		ID:                "glm-5.3-flash",
		SupportsReasoning: true,
		Reasoning: map[string]any{
			"supportedEfforts":   []any{"low", "high", "max"},
			"canDisableThinking": true,
		},
	})
	info := entry["model_info"].(map[string]any)
	if info["supports_reasoning"] != true {
		t.Fatalf("supports_reasoning=%v", info["supports_reasoning"])
	}
	if info["supports_max_reasoning_effort"] != true {
		t.Fatalf("supports_max_reasoning_effort=%v", info["supports_max_reasoning_effort"])
	}
	if info["supports_none_reasoning_effort"] != true {
		t.Fatalf("supports_none_reasoning_effort=%v", info["supports_none_reasoning_effort"])
	}
}
