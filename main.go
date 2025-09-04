package main

import (
	"fmt"
	"log"

	"quiz-please-scheduler/config"
	"quiz-please-scheduler/gameprovider"
	"quiz-please-scheduler/service"
)

func main() {
	fmt.Println("Quiz Please Scheduler is starting...")

	var provider gameprovider.GameProvider = gameprovider.New()

	games := provider.GetGamesList()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error config loading: %v", err)
		return
	}

	telegramService, err := service.New(cfg.Telegram.BotToken, cfg.Telegram.ChatID)
	if err != nil {
		log.Fatalf("Error initialization Telegram: %v", err)
		return
	}

	telegramService.SendGames(games)
}
