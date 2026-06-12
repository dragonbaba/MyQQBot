package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dragonbaba/MyQQBot/internal/config"
	"github.com/dragonbaba/MyQQBot/internal/llm"
	"github.com/sashabaranov/go-openai"
)

// ChatMessage represents a message in the GUI chat history.
type ChatMessage struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// ChatService exposes direct LLM chat to the frontend.
type ChatService struct {
	llmClient *llm.Client
	executor  llm.ToolExecutor
	cfg       config.LLMConfig
	history   []ChatMessage
	mu        sync.RWMutex
	tools     []openai.Tool
}

// NewChatService creates a chat service.
func NewChatService(llmClient *llm.Client, executor llm.ToolExecutor) *ChatService {
	var tools []openai.Tool
	if executor != nil {
		tools = executor.Definitions()
	}
	return &ChatService{
		llmClient: llmClient,
		executor:  executor,
		cfg:       llmClient.Config(),
		tools:     tools,
		history:   make([]ChatMessage, 0),
	}
}

// UpdateConfig refreshes the service's config copy.
func (s *ChatService) UpdateConfig(cfg config.LLMConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cfg = cfg
	s.llmClient = llm.NewClient(cfg)
}

// SendMessage sends a text message to the LLM and returns the assistant reply.
func (s *ChatService) SendMessage(content string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.history = append(s.history, ChatMessage{
		Role:      openai.ChatMessageRoleUser,
		Content:   content,
		Timestamp: time.Now().UnixMilli(),
	})

	messages := s.buildLLMMessages()
	reply, err := s.llmClient.Chat(context.Background(), messages, s.tools, s.executor)
	if err != nil {
		return "", err
	}

	s.history = append(s.history, ChatMessage{
		Role:      openai.ChatMessageRoleAssistant,
		Content:   reply.Content,
		Timestamp: time.Now().UnixMilli(),
	})

	return reply.Content, nil
}

// SendVisionMessage sends a message with image URLs, using the vision model fallback if needed.
func (s *ChatService) SendVisionMessage(text string, imageURLs []string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	client := s.llmClient
	useFallback := len(imageURLs) > 0 && !client.HasVision() && s.cfg.VisionModel != ""
	if useFallback {
		client = client.SetModel(s.cfg.VisionModel)
	}

	contentParts := make([]openai.ChatMessagePart, 0, len(imageURLs)+1)
	contentParts = append(contentParts, openai.ChatMessagePart{
		Type: openai.ChatMessagePartTypeText,
		Text: text,
	})
	for _, url := range imageURLs {
		contentParts = append(contentParts, openai.ChatMessagePart{
			Type: openai.ChatMessagePartTypeImageURL,
			ImageURL: &openai.ChatMessageImageURL{
				URL: url,
			},
		})
	}

	s.history = append(s.history, ChatMessage{
		Role:      openai.ChatMessageRoleUser,
		Content:   fmt.Sprintf("[图片消息] %s", text),
		Timestamp: time.Now().UnixMilli(),
	})

	messages := s.buildLLMMessages()
	messages = append(messages, openai.ChatCompletionMessage{
		Role:         openai.ChatMessageRoleUser,
		MultiContent: contentParts,
	})

	tools := s.tools
	if useFallback {
		// Vision fallback typically does not need tools.
		tools = nil
	}

	reply, err := client.Chat(context.Background(), messages, tools, nil)
	if err != nil {
		return "", err
	}

	s.history = append(s.history, ChatMessage{
		Role:      openai.ChatMessageRoleAssistant,
		Content:   reply.Content,
		Timestamp: time.Now().UnixMilli(),
	})

	return reply.Content, nil
}

// GetHistory returns the GUI chat history.
func (s *ChatService) GetHistory() []ChatMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]ChatMessage, len(s.history))
	copy(out, s.history)
	return out
}

// ClearHistory clears the GUI chat history.
func (s *ChatService) ClearHistory() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = make([]ChatMessage, 0)
}

func (s *ChatService) buildLLMMessages() []openai.ChatCompletionMessage {
	messages := make([]openai.ChatCompletionMessage, 0, len(s.history))
	for _, m := range s.history {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return messages
}
