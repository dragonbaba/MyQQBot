package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dragonbaba/MyQQBot/internal/bot"
	"github.com/dragonbaba/MyQQBot/internal/config"
	"github.com/dragonbaba/MyQQBot/internal/llm"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// BotStatus is the runtime status of the bot.
type BotStatus struct {
	Running     bool   `json:"running"`
	ConnectedAt string `json:"connectedAt"`
	OneBotURL   string `json:"oneBotURL"`
}

// BotEvent is emitted to the frontend.
type BotEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// BotService exposes bot control to the frontend.
type BotService struct {
	botClient *bot.Client
	llmClient *llm.Client
	cfg       config.Config
	app       *application.App
}

// NewBotService creates a bot service.
func NewBotService(cfg config.Config, botClient *bot.Client, llmClient *llm.Client) *BotService {
	return &BotService{
		botClient: botClient,
		llmClient: llmClient,
		cfg:       cfg,
	}
}

// SetApplication stores the Wails app reference.
func (s *BotService) SetApplication(app *application.App) {
	s.app = app
}

// StartBot connects to the OneBot WebSocket endpoint.
func (s *BotService) StartBot() error {
	if s.botClient.IsRunning() {
		return fmt.Errorf("bot is already running")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return s.botClient.Connect(ctx)
}

// StopBot disconnects from the OneBot endpoint.
func (s *BotService) StopBot() error {
	return s.botClient.Disconnect()
}

// GetBotStatus returns the current bot status.
func (s *BotService) GetBotStatus() BotStatus {
	return BotStatus{
		Running:   s.botClient.IsRunning(),
		OneBotURL: s.cfg.Bot.OneBotWSURL,
	}
}

// OnBotEvent emits a bot event to the frontend.
func (s *BotService) OnBotEvent(event BotEvent) {
	if s.app != nil {
		s.app.Event.Emit("bot:event", event)
	}
}
