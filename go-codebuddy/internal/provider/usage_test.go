package provider

import "testing"

func TestEstimateUsageUsesPromptLength(t *testing.T) {
	prompt := "abcd" // 4 bytes -> 1 token floor via (4+3)/4 = 1, but longer:
	prompt = string(make([]byte, 40))
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
