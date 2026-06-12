package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

const maxToolRounds = 5

var ErrAPIKeyNotConfigured = errors.New("API key not configured")

// Chat sends a chat completion request and handles tool call loops.
// It returns the final assistant message (text-only).
func (c *Client) Chat(ctx context.Context, messages []openai.ChatCompletionMessage, tools []openai.Tool, executor ToolExecutor) (*openai.ChatCompletionMessage, error) {
	if c.cfg.APIKey == "" {
		return nil, ErrAPIKeyNotConfigured
	}

	currentMessages := make([]openai.ChatCompletionMessage, len(messages))
	copy(currentMessages, messages)

	for round := 0; round <= maxToolRounds; round++ {
		// Clamp messages to model context limit (approximate by tokens ~= chars/4).
		maxCtx := c.MaxContext()
		currentMessages = trimMessagesToContext(currentMessages, maxCtx)

		req := openai.ChatCompletionRequest{
			Model:    c.cfg.Model,
			Messages: currentMessages,
		}
		if len(tools) > 0 && executor != nil {
			req.Tools = tools
		}
		if c.cfg.ReasoningEffort != "" {
			req.ReasoningEffort = c.cfg.ReasoningEffort
		}

		resp, err := c.client.CreateChatCompletion(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("llm request failed: %w", err)
		}

		if len(resp.Choices) == 0 {
			return nil, errors.New("llm returned no choices")
		}

		choice := resp.Choices[0]
		msg := choice.Message

		if msg.Role == "" {
			msg.Role = openai.ChatMessageRoleAssistant
		}

		if len(msg.ToolCalls) == 0 {
			return &msg, nil
		}

		if round == maxToolRounds {
			return nil, errors.New("搜索轮次过多，请简化问题")
		}

		currentMessages = append(currentMessages, msg)

		for _, tc := range msg.ToolCalls {
			if tc.Function.Name == "" {
				continue
			}
			result, err := executor.Execute(tc.Function.Name, json.RawMessage(tc.Function.Arguments))
			if err != nil {
				result = fmt.Sprintf("执行工具 %s 出错: %v", tc.Function.Name, err)
			}
			currentMessages = append(currentMessages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    result,
				Name:       tc.Function.Name,
				ToolCallID: tc.ID,
			})
		}
	}

	return nil, errors.New("tool call recursion exceeded limit")
}

// trimMessagesToContext drops oldest messages until the total is under the context limit.
// Approximation: 1 token ~= 4 UTF-8 characters for CJK/English mixed text.
func trimMessagesToContext(messages []openai.ChatCompletionMessage, maxCtx int) []openai.ChatCompletionMessage {
	if maxCtx <= 0 {
		return messages
	}
	maxChars := maxCtx * 4
	total := 0
	for _, m := range messages {
		total += len(m.Content)
	}
	start := 0
	for total > maxChars && start < len(messages)-1 {
		total -= len(messages[start].Content)
		start++
	}
	if start == 0 {
		return messages
	}
	return messages[start:]
}
