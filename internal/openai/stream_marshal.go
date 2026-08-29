package openai

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// MarshalStreamChunk serializes a streaming chunk with a fixed JSON skeleton.
// Dynamic string fields (content, tool_calls, usage) still go through json.Marshal
// for correct escaping and OpenAI wire semantics.
func MarshalStreamChunk(c StreamChunk) ([]byte, error) {
	switch {
	case len(c.Choices) == 0 && c.Usage != nil:
		return marshalUsageStreamChunk(c)
	case len(c.Choices) == 1 && c.Usage == nil:
		return marshalChoiceStreamChunk(c, c.Choices[0])
	default:
		return json.Marshal(c)
	}
}

func marshalUsageStreamChunk(c StreamChunk) ([]byte, error) {
	usageRaw, err := json.Marshal(c.Usage)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	b.Grow(len(usageRaw) + len(c.ID) + len(c.Model) + 96)
	writeStreamChunkPrefix(&b, c)
	b.WriteString(`,"choices":[]`)
	b.WriteString(`,"usage":`)
	b.Write(usageRaw)
	b.WriteByte('}')
	return b.Bytes(), nil
}

func marshalChoiceStreamChunk(c StreamChunk, ch Choice) ([]byte, error) {
	if ch.Index != 0 {
		return json.Marshal(c)
	}
	if ch.Delta != nil && len(ch.Delta.ToolCalls) > 0 {
		return json.Marshal(c)
	}
	var b bytes.Buffer
	b.Grow(len(c.ID) + len(c.Model) + 128)
	if ch.Delta != nil && ch.Delta.Content != "" {
		contentRaw, err := json.Marshal(ch.Delta.Content)
		if err != nil {
			return nil, err
		}
		b.Grow(len(contentRaw))
		writeStreamChunkPrefix(&b, c)
		b.WriteString(`,"choices":[{"index":0,"delta":{"content":`)
		b.Write(contentRaw)
		b.WriteString(`}`)
		writeFinishReason(&b, ch.FinishReason)
		b.WriteString(`}]}`)
		return b.Bytes(), nil
	}
	if ch.Delta != nil && ch.Delta.Role != "" && ch.Delta.Content == "" && len(ch.Delta.ToolCalls) == 0 {
		writeStreamChunkPrefix(&b, c)
		b.WriteString(`,"choices":[{"index":0,"delta":{"role":`)
		writeJSONString(&b, ch.Delta.Role)
		b.WriteString(`}`)
		writeFinishReason(&b, ch.FinishReason)
		b.WriteString(`}]}`)
		return b.Bytes(), nil
	}
	if ch.Delta != nil && ch.Delta.Role == "" && ch.Delta.Content == "" && len(ch.Delta.ToolCalls) == 0 {
		writeStreamChunkPrefix(&b, c)
		b.WriteString(`,"choices":[{"index":0,"delta":{}`)
		writeFinishReason(&b, ch.FinishReason)
		b.WriteString(`}]}`)
		return b.Bytes(), nil
	}
	if ch.Delta == nil {
		writeStreamChunkPrefix(&b, c)
		b.WriteString(`,"choices":[{"index":0,"delta":null`)
		writeFinishReason(&b, ch.FinishReason)
		b.WriteString(`}]}`)
		return b.Bytes(), nil
	}
	return json.Marshal(c)
}

func writeStreamChunkPrefix(b *bytes.Buffer, c StreamChunk) {
	b.WriteString(`{"id":`)
	writeJSONString(b, c.ID)
	b.WriteString(`,"object":"chat.completion.chunk","created":`)
	b.WriteString(strconv.FormatInt(c.Created, 10))
	b.WriteString(`,"model":`)
	writeJSONString(b, c.Model)
}

func writeFinishReason(b *bytes.Buffer, finishReason *string) {
	b.WriteString(`,"finish_reason":`)
	if finishReason == nil {
		b.WriteString("null")
		return
	}
	writeJSONString(b, *finishReason)
}

func writeJSONString(b *bytes.Buffer, s string) {
	raw, err := json.Marshal(s)
	if err != nil {
		b.WriteString(`""`)
		return
	}
	b.Write(raw)
}

// MarshalSSE implements httputil.SSEMarshaler for hot-path stream chunks.
func (c StreamChunk) MarshalSSE() ([]byte, error) {
	return MarshalStreamChunk(c)
}
