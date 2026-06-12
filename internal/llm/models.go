package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// ModelCapability describes what a model can do.
type ModelCapability struct {
	ID           string `json:"id"`
	Vision       bool   `json:"vision"`
	Tools        bool   `json:"tools"`
	Reasoning    bool   `json:"reasoning"`
	MaxContext   int    `json:"max_context"`
	DefaultModel bool   `json:"default_model"`
}

// KnownModelCapabilities contains heuristics for popular models.
// Unknown models will be detected by ID patterns.
var knownModelCapabilities = map[string]ModelCapability{
	"gpt-4o":          {ID: "gpt-4o", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"gpt-4o-mini":     {ID: "gpt-4o-mini", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"gpt-4-turbo":     {ID: "gpt-4-turbo", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"gpt-4":           {ID: "gpt-4", Vision: false, Tools: true, Reasoning: false, MaxContext: 8192},
	"gpt-3.5-turbo":   {ID: "gpt-3.5-turbo", Vision: false, Tools: true, Reasoning: false, MaxContext: 16385},
	"o1":              {ID: "o1", Vision: false, Tools: true, Reasoning: true, MaxContext: 200000},
	"o1-mini":         {ID: "o1-mini", Vision: false, Tools: true, Reasoning: true, MaxContext: 128000},
	"o3":              {ID: "o3", Vision: false, Tools: true, Reasoning: true, MaxContext: 200000},
	"o3-mini":         {ID: "o3-mini", Vision: false, Tools: true, Reasoning: true, MaxContext: 200000},
	"claude-3-opus":   {ID: "claude-3-opus", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"claude-3-sonnet": {ID: "claude-3-sonnet", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"claude-3-haiku":  {ID: "claude-3-haiku", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"gemini-pro":      {ID: "gemini-pro", Vision: false, Tools: true, Reasoning: false, MaxContext: 1000000},
	"gemini-pro-vision": {ID: "gemini-pro-vision", Vision: true, Tools: true, Reasoning: false, MaxContext: 1000000},
}

// DetectCapability guesses a model's capability from its ID.
func DetectCapability(modelID string) ModelCapability {
	id := strings.ToLower(modelID)
	if cap, ok := knownModelCapabilities[id]; ok {
		return cap
	}

	cap := ModelCapability{ID: modelID, MaxContext: 128000}
	if strings.Contains(id, "vision") || strings.Contains(id, "4o") || strings.Contains(id, "claude-3") || strings.Contains(id, "gemini-pro-vision") {
		cap.Vision = true
	}
	if !strings.Contains(id, "embedding") && !strings.Contains(id, "tts") && !strings.Contains(id, "whisper") && !strings.Contains(id, "dall-e") {
		cap.Tools = true
	}
	if strings.HasPrefix(id, "o1") || strings.HasPrefix(id, "o3") {
		cap.Reasoning = true
	}
	if strings.Contains(id, "32k") {
		cap.MaxContext = 32768
	} else if strings.Contains(id, "gpt-4") && !strings.Contains(id, "turbo") && !strings.Contains(id, "4o") {
		cap.MaxContext = 8192
	} else if strings.Contains(id, "claude") {
		cap.MaxContext = 200000
	}
	return cap
}

// ListModels fetches available models from the OpenAI-compatible API.
func (c *Client) ListModels(ctx context.Context) ([]ModelCapability, error) {
	if c.cfg.APIKey == "" {
		return nil, ErrAPIKeyNotConfigured
	}

	models, err := c.client.ListModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("list models failed: %w", err)
	}

	result := make([]ModelCapability, 0, len(models.Models))
	for _, m := range models.Models {
		cap := DetectCapability(m.ID)
		result = append(result, cap)
	}
	return result, nil
}

// HasVision reports whether the configured model supports vision.
func (c *Client) HasVision() bool {
	return DetectCapability(c.cfg.Model).Vision
}

// HasTools reports whether the configured model supports function calling.
func (c *Client) HasTools() bool {
	return DetectCapability(c.cfg.Model).Tools
}

// MaxContext returns the configured model's context limit.
func (c *Client) MaxContext() int {
	return DetectCapability(c.cfg.Model).MaxContext
}

// SetModel temporarily changes the active model for a single request context.
// This is useful for vision fallback.
func (c *Client) SetModel(model string) *Client {
	cfg := c.cfg
	cfg.Model = model
	ocfg := openai.DefaultConfig(cfg.APIKey)
	ocfg.BaseURL = cfg.BaseURL
	return &Client{cfg: cfg, client: openai.NewClientWithConfig(ocfg)}
}
