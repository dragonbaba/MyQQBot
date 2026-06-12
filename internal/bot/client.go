package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dragonbaba/MyQQBot/internal/config"
	"github.com/dragonbaba/MyQQBot/internal/llm"
	"github.com/gorilla/websocket"
)

// Client is a OneBot v11 forward WebSocket client.
type Client struct {
	wsURL       string
	httpURL     string
	accessToken string
	conn        *websocket.Conn
	llmClient   *llm.Client
	executor    llm.ToolExecutor
	cfg         config.Config
	handler     *Handler
	running     atomic.Bool
	mu          sync.Mutex
	cancel      context.CancelFunc
}
func NewClient(cfg config.Config, llmClient *llm.Client, executor llm.ToolExecutor) *Client {
	client := &Client{
		wsURL:       cfg.Bot.OneBotWSURL,
		httpURL:     cfg.Bot.OneBotHTTPURL,
		accessToken: cfg.Bot.AccessToken,
		llmClient:   llmClient,
		executor:    executor,
		cfg:         cfg,
		handler:     NewHandler(cfg, llmClient, executor),
	}
	ClientConfig.HTTPURL = cfg.Bot.OneBotHTTPURL
	ClientConfig.AccessToken = cfg.Bot.AccessToken
	return client
}

// Connect dials the OneBot WebSocket endpoint and starts reading events.
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running.Load() {
		return fmt.Errorf("bot already running")
	}

	dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}
	headers := http.Header{}
	if c.accessToken != "" {
		headers.Set("Authorization", "Bearer "+c.accessToken)
	}
	conn, _, err := dialer.DialContext(ctx, c.wsURL, headers)
	if err != nil {
		return fmt.Errorf("connect to onebot websocket failed: %w", err)
	}

	c.conn = conn
	c.running.Store(true)
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	go c.readLoop(ctx)
	return nil
}

// Disconnect gracefully closes the WebSocket connection.
func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running.Load() {
		return nil
	}
	c.running.Store(false)
	if c.cancel != nil {
		c.cancel()
	}
	if c.conn != nil {
		_ = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		_ = c.conn.Close()
	}
	return nil
}

// IsRunning reports whether the bot is connected.
func (c *Client) IsRunning() bool {
	return c.running.Load()
}

func (c *Client) readLoop(ctx context.Context) {
	defer func() {
		c.mu.Lock()
		c.running.Store(false)
		if c.conn != nil {
			_ = c.conn.Close()
		}
		c.mu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if c.running.Load() {
				fmt.Printf("[bot] websocket read error: %v\n", err)
			}
			return
		}

		var event Event
		if err := json.Unmarshal(data, &event); err != nil {
			fmt.Printf("[bot] failed to decode event: %v\n", err)
			continue
		}

		if event.PostType != EventTypeMessage {
			continue
		}

		// Process messages concurrently but don't block the read loop.
		go c.handler.HandleMessage(ctx, event)
	}
}
