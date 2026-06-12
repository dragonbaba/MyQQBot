package service

import (
	"context"
	"fmt"

	"github.com/dragonbaba/MyQQBot/internal/config"
	"github.com/dragonbaba/MyQQBot/internal/llm"
	"github.com/sashabaranov/go-openai"
)

const configPath = "config.yaml"

// LLMModelInfo is returned to the frontend.
type LLMModelInfo struct {
	ID        string `json:"id"`
	Vision    bool   `json:"vision"`
	Tools     bool   `json:"tools"`
	Reasoning bool   `json:"reasoning"`
}

// TestLLMResult combines connection test with model list.
type TestLLMResult struct {
	Success bool           `json:"success"`
	Models  []LLMModelInfo `json:"models"`
	Error   string         `json:"error"`
}

// ConfigService exposes configuration management to the frontend.
type ConfigService struct {
	cfg *config.Config
}

// NewConfigService creates a config service.
func NewConfigService(cfg *config.Config) *ConfigService {
	return &ConfigService{cfg: cfg}
}

// GetConfig returns the current configuration.
func (s *ConfigService) GetConfig() config.Config {
	return *s.cfg
}

// SaveConfig persists the configuration and updates the in-memory copy.
func (s *ConfigService) SaveConfig(cfg config.Config) error {
	if err := config.SaveConfig(configPath, &cfg); err != nil {
		return fmt.Errorf("save config failed: %w", err)
	}
	*s.cfg = cfg
	return nil
}

// TestLLMConnection validates LLM connectivity and returns available models.
func (s *ConfigService) TestLLMConnection(baseURL, apiKey string) TestLLMResult {
	if apiKey == "" {
		return TestLLMResult{Success: false, Error: "API key is empty"}
	}
	ocfg := openai.DefaultConfig(apiKey)
	ocfg.BaseURL = baseURL
	client := openai.NewClientWithConfig(ocfg)

	ctx := context.Background()
	models, err := client.ListModels(ctx)
	if err != nil {
		return TestLLMResult{Success: false, Error: fmt.Sprintf("list models failed: %v", err)}
	}

	// Also test a minimal chat completion with the current model.
	_, err = client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.cfg.LLM.Model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: "Hi"},
		},
	})
	if err != nil {
		return TestLLMResult{Success: false, Error: fmt.Sprintf("chat test failed: %v", err)}
	}

	result := TestLLMResult{Success: true, Models: make([]LLMModelInfo, 0, len(models.Models))}
	for _, m := range models.Models {
		cap := llm.DetectCapability(m.ID)
		result.Models = append(result.Models, LLMModelInfo{
			ID:        cap.ID,
			Vision:    cap.Vision,
			Tools:     cap.Tools,
			Reasoning: cap.Reasoning,
		})
	}
	return result
}
