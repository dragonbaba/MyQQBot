package llm

import (
	"encoding/json"

	"github.com/sashabaranov/go-openai"
)

// ToolExecutor executes LLM function tools and provides their definitions.
type ToolExecutor interface {
	Execute(name string, args json.RawMessage) (string, error)
	Definitions() []openai.Tool
}

// SearchToolExecutor is the minimal executor interface for search tools.
// It matches the broader ToolExecutor contract.
type SearchToolExecutor interface {
	ToolExecutor
}

// SearchToolDefinitions returns the OpenAI function definitions for web search.
func SearchToolDefinitions() []openai.Tool {
	return []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "search_web",
				Description: "Search the web for real-time information when the user asks about current events, weather, news, or facts that may change over time.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "The search query in the user's language",
						},
					},
					"required": []string{"query"},
				},
			},
		},
	}
}
