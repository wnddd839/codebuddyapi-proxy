package openai

import (
	"encoding/json"
	"testing"

	"github.com/wnddd839/codebuddy-proxy/internal/provider"
)

func TestUsageFromProviderPreservesCache(t *testing.T) {
	u := UsageFromProvider(provider.Usage{
		PromptTokens:         100,
		CompletionTokens:     20,
		TotalTokens:          120,
		PromptTokensDetails:  &provider.PromptTokensDetails{CachedTokens: 70},
		CacheReadInputTokens: 70,
	})
	if u.PromptTokensDetails == nil || u.PromptTokensDetails.CachedTokens != 70 {
		t.Fatalf("%+v", u)
	}
	raw, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	if !contains(s, `"cached_tokens":70`) || !contains(s, `"prompt_tokens":100`) {
		t.Fatalf("json=%s", s)
	}
}

func TestStreamUsageChunk(t *testing.T) {
	chunk := StreamUsageChunk("id1", "auto", Usage{PromptTokens: 1, CompletionTokens: 2, TotalTokens: 3})
	if chunk.Usage == nil || chunk.Usage.TotalTokens != 3 {
		t.Fatalf("%+v", chunk)
	}
	if len(chunk.Choices) != 0 {
		t.Fatalf("choices should be empty for include_usage chunk")
	}
	raw, _ := json.Marshal(chunk)
	if !contains(string(raw), `"usage"`) {
		t.Fatalf("%s", raw)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || (len(s) > 0 && (indexOf(s, sub) >= 0)))
}
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
