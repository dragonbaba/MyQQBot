package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// LLMConfig holds the OpenAI-compatible LLM settings.
type LLMConfig struct {
	BaseURL           string                 `mapstructure:"base_url" json:"base_url"`
	APIKey            string                 `mapstructure:"api_key" json:"api_key"`
	Model             string                 `mapstructure:"model" json:"model"`
	VisionModel       string                 `mapstructure:"vision_model" json:"vision_model"`
	ReasoningEffort   string                 `mapstructure:"reasoning_effort" json:"reasoning_effort"`
	SystemPrompt      string                 `mapstructure:"system_prompt" json:"system_prompt"`
	ModelCapabilities map[string]ModelCapability `mapstructure:"model_capabilities" json:"model_capabilities"`
}

// ModelCapability stores per-model overrides detected by the user.
type ModelCapability struct {
	Vision     bool `mapstructure:"vision" json:"vision"`
	Tools      bool `mapstructure:"tools" json:"tools"`
	Reasoning  bool `mapstructure:"reasoning" json:"reasoning"`
	MaxContext int  `mapstructure:"max_context" json:"max_context"`
}

// BotConfig holds the OneBot adapter connection settings.
type BotConfig struct {
	OneBotWSURL string `mapstructure:"onebot_ws_url" json:"onebot_ws_url"`
	AdminQQ     string `mapstructure:"admin_qq" json:"admin_qq"`
}

// SearchConfig holds the web search provider settings.
type SearchConfig struct {
	Provider     string `mapstructure:"provider" json:"provider"`
	TavilyAPIKey string `mapstructure:"tavily_api_key" json:"tavily_api_key"`
}

// Config is the top-level application configuration.
type Config struct {
	LLM    LLMConfig    `mapstructure:"llm" json:"llm"`
	Bot    BotConfig    `mapstructure:"bot" json:"bot"`
	Search SearchConfig `mapstructure:"search" json:"search"`
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() *Config {
	return &Config{
		LLM: LLMConfig{
			BaseURL:           "https://api.openai.com/v1",
			APIKey:            "",
			Model:             "gpt-4o",
			VisionModel:       "",
			ReasoningEffort:   "",
			SystemPrompt:      "你是一个有用的 QQ 机器人助手。当用户询问时事、天气、新闻或需要实时信息时，请使用 search_web 工具获取最新信息后再回答。",
			ModelCapabilities: make(map[string]ModelCapability),
		},
		Bot: BotConfig{
			OneBotWSURL: "ws://127.0.0.1:3001",
			AdminQQ:     "",
		},
		Search: SearchConfig{
			Provider:     "duckduckgo",
			TavilyAPIKey: "",
		},
	}
}

// LoadConfig reads configuration from the given YAML file path.
// If the file does not exist, it returns the default configuration without error.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix("MYQQBOT")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// If the file is simply missing, return defaults.
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return cfg, nil
		}
		// For explicit file paths, viper returns a different error type when the file is not found.
		// Treat any "not found" error as defaults.
		if path != "" && err.Error() == fmt.Sprintf("open %s: The system cannot find the file specified.", path) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("failed to read config: %w", err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// SaveConfig writes the configuration to the given YAML file path.
func SaveConfig(path string, cfg *Config) error {
	v := viper.New()
	v.Set("llm", cfg.LLM)
	v.Set("bot", cfg.Bot)
	v.Set("search", cfg.Search)
	return v.WriteConfigAs(path)
}
