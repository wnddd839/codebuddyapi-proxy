package openai

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wnddd839/codebuddy-proxy/internal/provider"
	"github.com/wnddd839/codebuddy-proxy/internal/strutil"
)

type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    any    `json:"code,omitempty"`
}

func NewError(message, typ string) ErrorBody {
	return ErrorBody{Error: ErrorDetail{Message: message, Type: typ}}
}

type ChatCompletion struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int      `json:"index"`
	Message      *Message `json:"message,omitempty"`
	Delta        *Delta   `json:"delta,omitempty"`
	FinishReason *string  `json:"finish_reason"`
}

type Message struct {
	Role      string     `json:"role"`
	Content   any        `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type Delta struct {
	Role      string     `json:"role,omitempty"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	Index    int          `json:"index,omitempty"`
	ID       string       `json:"id,omitempty"`
	Type     string       `json:"type,omitempty"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type StreamChunk struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

func FromTurn(turn provider.Turn, id, model string) ChatCompletion {
	if id == "" {
		id = fmt.Sprintf("chatcmpl_%d", time.Now().UnixMilli())
	}
	toolCalls := make([]ToolCall, 0, len(turn.ToolUses))
	for _, tool := range turn.ToolUses {
		toolCalls = append(toolCalls, ToolCall{
			ID:   strutil.First(tool.ID, fmt.Sprintf("call_%d", time.Now().UnixNano())),
			Type: "function",
			Function: ToolFunction{
				Name:      strutil.First(tool.Name, "tool"),
				Arguments: mustJSON(tool.Input),
			},
		})
	}
	finish := "stop"
	if len(toolCalls) > 0 {
		finish = "tool_calls"
	}
	var content any = turn.Text
	if turn.Text == "" {
		content = nil
	}
	msg := &Message{Role: "assistant", Content: content}
	if len(toolCalls) > 0 {
		msg.ToolCalls = toolCalls
	}
	return ChatCompletion{
		ID:      id,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []Choice{{
			Index:        0,
			Message:      msg,
			FinishReason: &finish,
		}},
		Usage: Usage{
			PromptTokens:     turn.Usage.PromptTokens,
			CompletionTokens: turn.Usage.CompletionTokens,
			TotalTokens:      turn.Usage.TotalTokens,
		},
	}
}

func StreamChunkOf(id, model string, delta Delta, finishReason *string) StreamChunk {
	return StreamChunk{
		ID:      id,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []Choice{{
			Index:        0,
			Delta:        &delta,
			FinishReason: finishReason,
		}},
	}
}

func mustJSON(value map[string]any) string {
	if value == nil {
		return "{}"
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(raw)
}
