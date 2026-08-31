package models

import (
	"testing"

	"github.com/wnddd839/codebuddy-proxy/internal/provider"
)

func TestMergeModelsByID(t *testing.T) {
	merged := mergeModelsByID(
		[]map[string]any{{"id": "gpt-5.4"}, {"id": "gemini-3.5-flash"}},
		[]map[string]any{{"id": "hy4-preview"}, {"id": "gpt-5.4"}},
	)
	if len(merged) != 3 {
		t.Fatalf("len=%d want 3: %+v", len(merged), merged)
	}
	ids := make([]string, len(merged))
	for i, m := range merged {
		ids[i] = modelRowID(m)
	}
	if ids[0] != "gpt-5.4" || ids[1] != "gemini-3.5-flash" || ids[2] != "hy4-preview" {
		t.Fatalf("order/ids=%v", ids)
	}
}

func TestV3ConfigCandidateBasesGlobalMergesCopilot(t *testing.T) {
	bases := v3ConfigCandidateBases(provider.ChatOptions{
		Site:    "global",
		BaseURL: "https://www.codebuddy.ai",
	})
	if len(bases) != 2 || bases[0] != "https://www.codebuddy.ai" || bases[1] != "https://copilot.tencent.com" {
		t.Fatalf("global bases=%v", bases)
	}
	domestic := v3ConfigCandidateBases(provider.ChatOptions{Site: "domestic"})
	if len(domestic) != 1 || domestic[0] != "https://copilot.tencent.com" {
		t.Fatalf("domestic bases=%v", domestic)
	}
}

func TestPublicModelIDStripsCodeBuddyPrefix(t *testing.T) {
	cases := map[string]string{
		"":                          "auto",
		"default":                   "auto",
		"auto":                      "auto",
		"codebuddy/auto":            "auto",
		"codebuddy:auto":            "auto",
		"codebuddy/deepseek-v4-pro": "deepseek-v4-pro",
		"CODEBUDDY/glm-5.3":         "glm-5.3",
		"minimax-m2.5":              "minimax-m2.5",
	}
	for in, want := range cases {
		if got := PublicModelID(in); got != want {
			t.Fatalf("PublicModelID(%q)=%q want %q", in, got, want)
		}
	}
}
