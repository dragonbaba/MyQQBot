package llm

import (
	"github.com/dragonbaba/MyQQBot/internal/config"
	"github.com/sashabaranov/go-openai"
)

// Client wraps the OpenAI-compatible HTTP client.
type Client struct {
	cfg    config.LLMConfig
	client *openai.Client
}

// NewClient creates a new LLM client from configuration.
// If the API key is empty, the client still initialises but calls will fail with a clear error.
func NewClient(cfg config.LLMConfig) *Client {
	ocfg := openai.DefaultConfig(cfg.APIKey)
	ocfg.BaseURL = cfg.BaseURL
	return &Client{
		cfg:    cfg,
		client: openai.NewClientWithConfig(ocfg),
	}
}

// Config returns the client's LLM configuration.
func (c *Client) Config() config.LLMConfig {
	return c.cfg
}
