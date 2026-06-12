package service

import (
	"context"
	"sync"
	"time"

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
		tools:     tools,
		history:   make([]ChatMessage, 0),
	}
}

// SendMessage sends a message to the LLM and returns the assistant reply.
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
