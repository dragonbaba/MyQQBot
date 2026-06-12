package main

import (
	"embed"
	"log"

	"github.com/dragonbaba/MyQQBot/internal/config"
	"github.com/wailsapp/wails/v3/pkg/application"
)
//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app, err := App(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "MyQQBot",
		Width:            1200,
		Height:           800,
		MinWidth:         900,
		MinHeight:        600,
		BackgroundColour: application.NewRGB(2, 6, 23),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
