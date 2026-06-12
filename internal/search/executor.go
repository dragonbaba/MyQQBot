package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dragonbaba/MyQQBot/internal/config"
	"github.com/dragonbaba/MyQQBot/internal/llm"
	"github.com/sashabaranov/go-openai"
)

// Executor implements llm.ToolExecutor for web search.
type Executor struct {
	provider Provider
}

// NewExecutor creates a search executor from configuration.
func NewExecutor(cfg config.SearchConfig) *Executor {
	var p Provider
	switch strings.ToLower(cfg.Provider) {
	case "tavily":
		p = NewTavilyProvider(cfg.TavilyAPIKey)
	default:
		p = NewDuckDuckGoProvider()
	}
	return &Executor{provider: p}
}

// NewExecutorWithProvider creates an executor with an explicit provider.
func NewExecutorWithProvider(p Provider) *Executor {
	return &Executor{provider: p}
}

// Execute runs the named tool.
func (e *Executor) Execute(name string, args json.RawMessage) (string, error) {
	switch name {
	case "search_web":
		var params struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("search_web: invalid arguments: %w", err)
		}
		if params.Query == "" {
			return "", fmt.Errorf("search_web: query is empty")
		}
		results, err := e.provider.Search(context.Background(), params.Query)
		if err != nil {
			return "", err
		}
		return formatResults(results), nil
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

// Definitions returns the OpenAI tool definitions.
func (e *Executor) Definitions() []openai.Tool {
	return llm.SearchToolDefinitions()
}

func formatResults(results []Result) string {
	if len(results) == 0 {
		return "未找到相关搜索结果。"
	}
	var b strings.Builder
	b.WriteString("搜索结果：\n")
	for i, r := range results {
		b.WriteString(fmt.Sprintf("%d. [%s](%s)\n   %s\n", i+1, r.Title, r.URL, r.Snippet))
	}
	return b.String()
}
