package provider

import (
	"encoding/json"
	"testing"
)

func BenchmarkEstimateTokenCountASCII(b *testing.B) {
	s := string(make([]byte, 4000))
	b.ReportAllocs()
	for b.Loop() {
		_ = estimateTokenCount(s)
	}
}

func BenchmarkEstimateTokenCountCJK(b *testing.B) {
	s := ""
	for range 500 {
		s += "中文测试内容"
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = estimateTokenCount(s)
	}
}

func BenchmarkUnmarshalOpenAIChunkTyped(b *testing.B) {
	raw := []byte(`{"id":"x","choices":[{"index":0,"delta":{"content":"hi"},"finish_reason":null}]}`)
	b.ReportAllocs()
	for b.Loop() {
		var chunk openAISSEChunk
		if err := json.Unmarshal(raw, &chunk); err != nil {
			b.Fatal(err)
		}
		_ = eventsFromOpenAIChunk(chunk)
	}
}

func BenchmarkUnmarshalOpenAIChunkAny(b *testing.B) {
	raw := []byte(`{"id":"x","choices":[{"index":0,"delta":{"content":"hi"},"finish_reason":null}]}`)
	b.ReportAllocs()
	for b.Loop() {
		var payload any
		if err := json.Unmarshal(raw, &payload); err != nil {
			b.Fatal(err)
		}
		_ = MapSSEEvent(payload)
	}
}

func BenchmarkEstimatePromptText(b *testing.B) {
	msgs := make([]map[string]any, 0, 50)
	for i := range 50 {
		msgs = append(msgs, map[string]any{
			"role":    "user",
			"content": "这段是模拟长上下文消息内容，编号=" + string(rune('A'+i%26)),
		})
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = estimatePromptText(msgs)
	}
}
