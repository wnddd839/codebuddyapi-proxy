package openai

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestMarshalStreamChunkMatchesJSON(t *testing.T) {
	finishStop := "stop"
	finishTools := "tool_calls"
	cases := []StreamChunk{
		StreamChunkOf("id1", "auto", Delta{Role: "assistant"}, nil),
		StreamChunkOf("id1", "auto", Delta{Content: "hello"}, nil),
		StreamChunkOf("id1", "auto", Delta{Content: "引号\"换行\n制表\t"}, nil),
		StreamChunkOf("id1", "auto", Delta{}, &finishStop),
		StreamUsageChunk("id1", "auto", Usage{
			PromptTokens:     10,
			CompletionTokens: 4,
			TotalTokens:      14,
			PromptTokensDetails: &PromptTokensDetails{
				CachedTokens: 3,
			},
		}),
		StreamChunkOf("id1", "auto", Delta{ToolCalls: []ToolCall{{
			Index: 0,
			ID:    "call_x",
			Type:  "function",
			Function: ToolFunction{
				Name:      "lookup",
				Arguments: `{"q":"test"}`,
			},
		}}}, nil),
		StreamChunkOf("id2", "hy4-preview", Delta{Content: "字"}, nil),
		StreamChunkOf("id2", "hy4-preview", Delta{}, &finishTools),
	}
	for i, chunk := range cases {
		want, err := json.Marshal(chunk)
		if err != nil {
			t.Fatalf("case %d marshal want: %v", i, err)
		}
		got, err := MarshalStreamChunk(chunk)
		if err != nil {
			t.Fatalf("case %d marshal got: %v", i, err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("case %d mismatch\nwant: %s\ngot:  %s", i, want, got)
		}
	}
}

func TestMarshalStreamChunkFinishReasonNull(t *testing.T) {
	chunk := StreamChunkOf("id", "m", Delta{Content: "x"}, nil)
	raw, err := MarshalStreamChunk(chunk)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(raw, []byte(`"finish_reason":null`)) {
		t.Fatalf("missing finish_reason:null: %s", raw)
	}
}

func BenchmarkMarshalStreamChunkContent(b *testing.B) {
	chunk := StreamChunkOf("chatcmpl_x", "auto", Delta{Content: "这段是中文流式内容"}, nil)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MarshalStreamChunk(chunk)
	}
}

func BenchmarkJSONMarshalStreamChunkContent(b *testing.B) {
	chunk := StreamChunkOf("chatcmpl_x", "auto", Delta{Content: "这段是中文流式内容"}, nil)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = json.Marshal(chunk)
	}
}
