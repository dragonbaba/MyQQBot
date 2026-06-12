package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// ModelCapability describes what a model can do.
type ModelCapability struct {
	ID         string `json:"id"`
	Vision     bool   `json:"vision"`
	Tools      bool   `json:"tools"`
	Reasoning  bool   `json:"reasoning"`
	MaxContext int    `json:"max_context"`
}

// KnownModelCapabilities contains heuristics for popular models.
// Unknown models will be detected by ID patterns.
var knownModelCapabilities = map[string]ModelCapability{
	// OpenAI
	"gpt-4o":              {ID: "gpt-4o", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"gpt-4o-mini":         {ID: "gpt-4o-mini", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"gpt-4o-latest":       {ID: "gpt-4o-latest", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"gpt-4-turbo":         {ID: "gpt-4-turbo", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"gpt-4-turbo-preview": {ID: "gpt-4-turbo-preview", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"gpt-4":               {ID: "gpt-4", Vision: false, Tools: true, Reasoning: false, MaxContext: 8192},
	"gpt-4-32k":           {ID: "gpt-4-32k", Vision: false, Tools: true, Reasoning: false, MaxContext: 32768},
	"gpt-3.5-turbo":       {ID: "gpt-3.5-turbo", Vision: false, Tools: true, Reasoning: false, MaxContext: 16385},
	"o1":                  {ID: "o1", Vision: true, Tools: true, Reasoning: true, MaxContext: 200000},
	"o1-mini":             {ID: "o1-mini", Vision: true, Tools: true, Reasoning: true, MaxContext: 128000},
	"o1-preview":          {ID: "o1-preview", Vision: true, Tools: true, Reasoning: true, MaxContext: 128000},
	"o3":                  {ID: "o3", Vision: true, Tools: true, Reasoning: true, MaxContext: 200000},
	"o3-mini":             {ID: "o3-mini", Vision: true, Tools: true, Reasoning: true, MaxContext: 200000},

	// Anthropic Claude
	"claude-3-opus":              {ID: "claude-3-opus", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"claude-3-opus-20240229":     {ID: "claude-3-opus-20240229", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"claude-3-opus-latest":       {ID: "claude-3-opus-latest", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"claude-3-sonnet":            {ID: "claude-3-sonnet", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"claude-3-sonnet-20240229":   {ID: "claude-3-sonnet-20240229", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"claude-3-5-sonnet":          {ID: "claude-3-5-sonnet", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"claude-3-5-sonnet-20240620": {ID: "claude-3-5-sonnet-20240620", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"claude-3-5-sonnet-20241022": {ID: "claude-3-5-sonnet-20241022", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"claude-3-5-sonnet-latest":   {ID: "claude-3-5-sonnet-latest", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"claude-3-haiku":             {ID: "claude-3-haiku", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},
	"claude-3-haiku-20240307":    {ID: "claude-3-haiku-20240307", Vision: true, Tools: true, Reasoning: false, MaxContext: 200000},

	// Google Gemini
	"gemini-1.5-pro":      {ID: "gemini-1.5-pro", Vision: true, Tools: true, Reasoning: false, MaxContext: 2000000},
	"gemini-1.5-flash":    {ID: "gemini-1.5-flash", Vision: true, Tools: true, Reasoning: false, MaxContext: 1000000},
	"gemini-1.5-pro-latest": {ID: "gemini-1.5-pro-latest", Vision: true, Tools: true, Reasoning: false, MaxContext: 2000000},
	"gemini-2.0-flash":    {ID: "gemini-2.0-flash", Vision: true, Tools: true, Reasoning: false, MaxContext: 1000000},
	"gemini-2.0-pro":      {ID: "gemini-2.0-pro", Vision: true, Tools: true, Reasoning: false, MaxContext: 2000000},
	"gemini-pro":          {ID: "gemini-pro", Vision: false, Tools: true, Reasoning: false, MaxContext: 1000000},
	"gemini-pro-vision":   {ID: "gemini-pro-vision", Vision: true, Tools: true, Reasoning: false, MaxContext: 1000000},

	// xAI Grok
	"grok-1":              {ID: "grok-1", Vision: false, Tools: true, Reasoning: false, MaxContext: 8192},
	"grok-2":              {ID: "grok-2", Vision: false, Tools: true, Reasoning: false, MaxContext: 128000},
	"grok-2-vision":       {ID: "grok-2-vision", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"grok-3":              {ID: "grok-3", Vision: true, Tools: true, Reasoning: true, MaxContext: 200000},
	"grok-3-mini":         {ID: "grok-3-mini", Vision: true, Tools: true, Reasoning: true, MaxContext: 200000},

	// 阿里通义千问 Qwen
	"qwen-turbo":          {ID: "qwen-turbo", Vision: false, Tools: true, Reasoning: false, MaxContext: 128000},
	"qwen-plus":           {ID: "qwen-plus", Vision: false, Tools: true, Reasoning: false, MaxContext: 128000},
	"qwen-max":            {ID: "qwen-max", Vision: false, Tools: true, Reasoning: false, MaxContext: 128000},
	"qwen-vl-plus":        {ID: "qwen-vl-plus", Vision: true, Tools: true, Reasoning: false, MaxContext: 32000},
	"qwen-vl-max":         {ID: "qwen-vl-max", Vision: true, Tools: true, Reasoning: false, MaxContext: 32000},
	"qwen2.5":             {ID: "qwen2.5", Vision: false, Tools: true, Reasoning: false, MaxContext: 128000},
	"qwen2.5-vl":          {ID: "qwen2.5-vl", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"qwen3":               {ID: "qwen3", Vision: false, Tools: true, Reasoning: true, MaxContext: 128000},

	// 百度文心 ERNIE
	"ernie-bot":           {ID: "ernie-bot", Vision: false, Tools: true, Reasoning: false, MaxContext: 8000},
	"ernie-bot-4":         {ID: "ernie-bot-4", Vision: false, Tools: true, Reasoning: false, MaxContext: 8000},
	"ernie-bot-8k":        {ID: "ernie-bot-8k", Vision: false, Tools: true, Reasoning: false, MaxContext: 8192},
	"ernie-speed":         {ID: "ernie-speed", Vision: false, Tools: true, Reasoning: false, MaxContext: 128000},
	"ernie-lite":          {ID: "ernie-lite", Vision: false, Tools: true, Reasoning: false, MaxContext: 128000},
	"ernie-tiny":          {ID: "ernie-tiny", Vision: false, Tools: true, Reasoning: false, MaxContext: 128000},

	// 智谱 GLM
	"chatglm3":            {ID: "chatglm3", Vision: false, Tools: true, Reasoning: false, MaxContext: 32000},
	"glm-4":               {ID: "glm-4", Vision: false, Tools: true, Reasoning: false, MaxContext: 128000},
	"glm-4v":              {ID: "glm-4v", Vision: true, Tools: true, Reasoning: false, MaxContext: 8000},
	"glm-4v-plus":         {ID: "glm-4v-plus", Vision: true, Tools: true, Reasoning: false, MaxContext: 8000},
	"glm-4-flash":         {ID: "glm-4-flash", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"glm-4-air":           {ID: "glm-4-air", Vision: false, Tools: true, Reasoning: false, MaxContext: 128000},
	"glm-4-airx":          {ID: "glm-4-airx", Vision: false, Tools: true, Reasoning: false, MaxContext: 8000},

	// Moonshot Kimi
	"moonshot-v1-8k":      {ID: "moonshot-v1-8k", Vision: false, Tools: true, Reasoning: false, MaxContext: 8192},
	"moonshot-v1-32k":     {ID: "moonshot-v1-32k", Vision: false, Tools: true, Reasoning: false, MaxContext: 32768},
	"moonshot-v1-128k":    {ID: "moonshot-v1-128k", Vision: false, Tools: true, Reasoning: false, MaxContext: 128000},
	"kimi-k1.5":           {ID: "kimi-k1.5", Vision: false, Tools: true, Reasoning: true, MaxContext: 256000},
	"kimi-k2":             {ID: "kimi-k2", Vision: false, Tools: true, Reasoning: true, MaxContext: 256000},

	// DeepSeek
	"deepseek-chat":       {ID: "deepseek-chat", Vision: false, Tools: true, Reasoning: false, MaxContext: 64000},
	"deepseek-coder":      {ID: "deepseek-coder", Vision: false, Tools: true, Reasoning: false, MaxContext: 64000},
	"deepseek-reasoner":   {ID: "deepseek-reasoner", Vision: false, Tools: true, Reasoning: true, MaxContext: 64000},
	"deepseek-v3":         {ID: "deepseek-v3", Vision: false, Tools: true, Reasoning: false, MaxContext: 64000},
	"deepseek-r1":         {ID: "deepseek-r1", Vision: false, Tools: true, Reasoning: true, MaxContext: 64000},

	// Meta Llama
	"llama-3.2-11b-vision": {ID: "llama-3.2-11b-vision", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"llama-3.2-90b-vision": {ID: "llama-3.2-90b-vision", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},

	// Mistral
	"mistral-large":       {ID: "mistral-large", Vision: false, Tools: true, Reasoning: false, MaxContext: 128000},
	"pixtral":             {ID: "pixtral", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
	"pixtral-large":       {ID: "pixtral-large", Vision: true, Tools: true, Reasoning: false, MaxContext: 128000},
}

// detectFamilyCapabilities returns default capabilities inferred from model ID patterns.
func detectFamilyCapabilities(id string) ModelCapability {
	cap := ModelCapability{ID: id, MaxContext: 128000}

	// Vision patterns
	if strings.Contains(id, "vision") ||
		strings.Contains(id, "vl") ||
		strings.Contains(id, "4o") ||
		strings.Contains(id, "claude-3") ||
		strings.Contains(id, "grok-2-vision") ||
		strings.Contains(id, "grok-3") ||
		strings.Contains(id, "glm-4v") ||
		strings.Contains(id, "glm-4-flash") ||
		strings.Contains(id, "gemini") ||
		strings.Contains(id, "llama-3.2") && strings.Contains(id, "vision") ||
		strings.Contains(id, "pixtral") {
		cap.Vision = true
	}

	// Tools patterns: most modern chat models support tools unless they are specialised.
	if !strings.Contains(id, "embedding") &&
		!strings.Contains(id, "tts") &&
		!strings.Contains(id, "whisper") &&
		!strings.Contains(id, "dall-e") &&
		!strings.Contains(id, "moderation") &&
		!strings.Contains(id, "rerank") {
		cap.Tools = true
	}

	// Reasoning patterns
	if strings.HasPrefix(id, "o1") ||
		strings.HasPrefix(id, "o3") ||
		strings.Contains(id, "deepseek-r") ||
		strings.Contains(id, "deepseek-reasoner") ||
		strings.Contains(id, "kimi-k") ||
		strings.Contains(id, "qwen3") ||
		strings.Contains(id, "grok-3") ||
		strings.Contains(id, "reasoner") {
		cap.Reasoning = true
	}

	// Context length patterns
	switch {
	case strings.Contains(id, "32k"):
		cap.MaxContext = 32768
	case strings.Contains(id, "8k"):
		cap.MaxContext = 8192
	case strings.Contains(id, "gpt-4") && !strings.Contains(id, "turbo") && !strings.Contains(id, "4o"):
		cap.MaxContext = 8192
	case strings.Contains(id, "claude"):
		cap.MaxContext = 200000
	case strings.Contains(id, "gemini-1.5-pro"):
		cap.MaxContext = 2000000
	case strings.Contains(id, "gemini"):
		cap.MaxContext = 1000000
	case strings.Contains(id, "kimi-k"):
		cap.MaxContext = 256000
	case strings.Contains(id, "deepseek"):
		cap.MaxContext = 64000
	}

	return cap
}

// DetectCapability guesses a model's capability from its ID.
func DetectCapability(modelID string) ModelCapability {
	id := strings.ToLower(modelID)
	if cap, ok := knownModelCapabilities[id]; ok {
		return cap
	}
	return detectFamilyCapabilities(id)
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
