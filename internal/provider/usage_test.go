package provider

import "testing"

func TestEstimateUsageUsesPromptLength(t *testing.T) {
	prompt := string(make([]byte, 40))
	usage := estimateUsage(prompt, "xxxx")
	if usage.PromptTokens != 10 { // (40+3)/4 = 10
		t.Fatalf("prompt tokens=%d want 10", usage.PromptTokens)
	}
	if usage.CompletionTokens < 1 {
		t.Fatal("completion tokens should be >= 1")
	}
}

func TestSnapshotUsesPromptInsteadOfEmpty(t *testing.T) {
	acc := newAccumulator()
	acc.push(Event{Type: "text_delta", Text: "hi"})
	turn := acc.snapshot(string(make([]byte, 40)))
	if turn.Usage.PromptTokens != 10 {
		t.Fatalf("prompt_tokens=%d want 10 (not empty-prompt floor)", turn.Usage.PromptTokens)
	}
}

func TestSnapshotPrefersUpstreamUsage(t *testing.T) {
	acc := newAccumulator()
	acc.push(Event{Type: "text_delta", Text: "hi"})
	acc.push(Event{Type: "usage", Usage: &Usage{PromptTokens: 123, CompletionTokens: 4, TotalTokens: 127}})
	turn := acc.snapshot("ignored")
	if turn.Usage.PromptTokens != 123 || turn.Usage.CompletionTokens != 4 {
		t.Fatalf("unexpected usage %+v", turn.Usage)
	}
}

func TestParseUsagePreservesCacheFields(t *testing.T) {
	usage := ParseUsage(map[string]any{
		"prompt_tokens":     float64(1000),
		"completion_tokens": float64(50),
		"total_tokens":      float64(1050),
		"prompt_tokens_details": map[string]any{
			"cached_tokens": float64(800),
		},
		"cache_read_input_tokens":     float64(800),
		"cache_creation_input_tokens": float64(20),
		"prompt_cache_hit_tokens":     float64(800),
	})
	if usage.PromptTokens != 1000 || usage.CompletionTokens != 50 || usage.TotalTokens != 1050 {
		t.Fatalf("base fields %+v", usage)
	}
	if usage.CachedTokens() != 800 {
		t.Fatalf("cached=%d", usage.CachedTokens())
	}
	if usage.PromptTokensDetails == nil || usage.PromptTokensDetails.CachedTokens != 800 {
		t.Fatalf("details=%+v", usage.PromptTokensDetails)
	}
	if usage.CacheCreationInputTokens != 20 {
		t.Fatalf("cache creation=%d", usage.CacheCreationInputTokens)
	}
}

func TestParseUsageDeepSeekAliases(t *testing.T) {
	usage := ParseUsage(map[string]any{
		"prompt_tokens":            float64(200),
		"completion_tokens":        float64(10),
		"prompt_cache_hit_tokens":  float64(150),
		"prompt_cache_miss_tokens": float64(50),
	})
	if usage.CachedTokens() != 150 {
		t.Fatalf("cached=%d", usage.CachedTokens())
	}
	if usage.PromptTokensDetails == nil || usage.PromptTokensDetails.CachedTokens != 150 {
		t.Fatalf("normalized details=%+v", usage.PromptTokensDetails)
	}
}

func TestUsageEventFromPayload(t *testing.T) {
	ev := usageEventFromPayload(map[string]any{
		"usage": map[string]any{
			"prompt_tokens":         float64(10),
			"completion_tokens":     float64(2),
			"prompt_tokens_details": map[string]any{"cached_tokens": float64(7)},
		},
	})
	if ev == nil || ev.Usage == nil || ev.Usage.CachedTokens() != 7 {
		t.Fatalf("event=%+v", ev)
	}
}
