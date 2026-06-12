package bot

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dragonbaba/MyQQBot/internal/config"
	"github.com/dragonbaba/MyQQBot/internal/llm"
	"github.com/sashabaranov/go-openai"
)

const maxHistoryPerSession = 20
const maxSessions = 50

// Handler encapsulates message processing and conversation history.
type Handler struct {
	llmClient     *llm.Client
	tools         []openai.Tool
	executor      llm.ToolExecutor
	cfg           config.Config
	history       map[string][]openai.ChatCompletionMessage
	lastAccessed  map[string]time.Time
	mu            sync.RWMutex
}

// NewHandler creates a new message handler.
func NewHandler(cfg config.Config, llmClient *llm.Client, executor llm.ToolExecutor) *Handler {
	var tools []openai.Tool
	if executor != nil {
		tools = executor.Definitions()
	}
	return &Handler{
		llmClient:    llmClient,
		tools:        tools,
		executor:     executor,
		cfg:          cfg,
		history:      make(map[string][]openai.ChatCompletionMessage),
		lastAccessed: make(map[string]time.Time),
	}
}

// HandleMessage processes an incoming OneBot message event.
func (h *Handler) HandleMessage(ctx context.Context, event Event) {
	msgText := strings.TrimSpace(event.RawMessage)
	if msgText == "" {
		return
	}

	sessionKey := makeSessionKey(event)
	history := h.getHistory(sessionKey)

	userName := event.Sender.Nickname
	if userName == "" {
		userName = fmt.Sprintf("%d", event.UserID)
	}

	history = append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: fmt.Sprintf("[%s] %s", userName, msgText),
	})

	messages := h.buildMessages(history)
	reply, err := h.llmClient.Chat(ctx, messages, h.tools, h.executor)

	var replyText string
	if err != nil {
		replyText = fmt.Sprintf("抱歉，处理消息时出错了：%v", err)
	} else {
		replyText = reply.Content
		history = append(history, *reply)
	}

	h.setHistory(sessionKey, history)

	var msgType string
	var targetID int64
	if event.MessageType == MessageTypePrivate {
		msgType = MessageTypePrivate
		targetID = event.UserID
	} else {
		msgType = MessageTypeGroup
		targetID = event.GroupID
	}

	if err := sendMessage(msgType, targetID, event.GroupID, replyText); err != nil {
		// Log and continue; don't crash the bot.
		fmt.Printf("[bot] failed to send message: %v\n", err)
	}
}

func (h *Handler) buildMessages(history []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	messages := make([]openai.ChatCompletionMessage, 0, len(history)+1)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: h.cfg.LLM.SystemPrompt,
	})
	messages = append(messages, history...)
	return messages
}

func (h *Handler) getHistory(key string) []openai.ChatCompletionMessage {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastAccessed[key] = time.Now()
	return h.history[key]
}

func (h *Handler) setHistory(key string, history []openai.ChatCompletionMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(history) > maxHistoryPerSession {
		history = history[len(history)-maxHistoryPerSession:]
	}
	h.history[key] = history
	h.lastAccessed[key] = time.Now()
	h.evictIfNeeded()
}

func (h *Handler) evictIfNeeded() {
	if len(h.history) <= maxSessions {
		return
	}
	var oldestKey string
	var oldestTime time.Time
	first := true
	for k, t := range h.lastAccessed {
		if first || t.Before(oldestTime) {
			oldestKey = k
			oldestTime = t
			first = false
		}
	}
	if oldestKey != "" {
		delete(h.history, oldestKey)
		delete(h.lastAccessed, oldestKey)
	}
}

func makeSessionKey(event Event) string {
	if event.MessageType == MessageTypeGroup {
		return fmt.Sprintf("group_%d", event.GroupID)
	}
	return fmt.Sprintf("private_%d", event.UserID)
}
