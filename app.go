package main

import (
	"github.com/dragonbaba/MyQQBot/internal/bot"
	"github.com/dragonbaba/MyQQBot/internal/config"
	"github.com/dragonbaba/MyQQBot/internal/llm"
	"github.com/dragonbaba/MyQQBot/internal/search"
	"github.com/dragonbaba/MyQQBot/internal/service"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// App creates the Wails application and wires all services.
func App(cfg *config.Config) (*application.App, error) {
	llmClient := llm.NewClient(cfg.LLM)
	searchExecutor := search.NewExecutor(cfg.Search)
	botClient := bot.NewClient(*cfg, llmClient, searchExecutor)

	botService := service.NewBotService(*cfg, botClient, llmClient)
	chatService := service.NewChatService(llmClient, searchExecutor)
	configService := service.NewConfigService(cfg)

	app := application.New(application.Options{
		Name:        "MyQQBot",
		Description: "Go + Wails v3 QQ Bot with OpenAI function calling",
		Services: []application.Service{
			application.NewService(botService),
			application.NewService(chatService),
			application.NewService(configService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	botService.SetApplication(app)

	return app, nil
}
