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

// TestLLMConnection validates LLM connectivity by listing available models.
// We intentionally avoid sending a chat completion here because many
// OpenAI-compatible gateways or reasoning models reject simple test prompts
// with opaque 400 errors (e.g. "Param Incorrect"). Listing models is a
// reliable, read-only way to verify the URL and API key.
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

	result := TestLLMResult{Success: true, Models: make([]LLMModelInfo, 0, len(models.Models))}
	for _, m := range models.Models {
		cap := llm.DetectCapability(m.ID)
		// Allow user overrides stored in config to take precedence.
		if override, ok := s.cfg.LLM.ModelCapabilities[m.ID]; ok {
			if override.Vision {
				cap.Vision = true
			}
			if override.Tools {
				cap.Tools = true
			}
			if override.Reasoning {
				cap.Reasoning = true
			}
			if override.MaxContext > 0 {
				cap.MaxContext = override.MaxContext
			}
		}
		result.Models = append(result.Models, LLMModelInfo{
			ID:        cap.ID,
			Vision:    cap.Vision,
			Tools:     cap.Tools,
			Reasoning: cap.Reasoning,
		})
	}
	return result
}
