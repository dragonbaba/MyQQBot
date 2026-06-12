package service

import (
	"context"
	"fmt"

	"github.com/dragonbaba/MyQQBot/internal/config"
	"github.com/sashabaranov/go-openai"
)

const configPath = "config.yaml"

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

// TestLLMConnection sends a minimal request to validate LLM connectivity.
func (s *ConfigService) TestLLMConnection(baseURL, apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("API key is empty")
	}
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = baseURL
	client := openai.NewClientWithConfig(cfg)

	_, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: s.cfg.LLM.Model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: "Hi"},
		},
	})
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	return nil
}
