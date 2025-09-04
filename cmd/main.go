package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"quiz-please-scheduler/internal/config"
	"quiz-please-scheduler/internal/services"
	"quiz-please-scheduler/pkg/gameprovider"
	"quiz-please-scheduler/pkg/telegram"
)

func main() {
	fmt.Println("Quiz Please Scheduler is starting...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error configuration loading: %v", err)
		return
	}

	provider := gameprovider.NewInstance()
	sender, err := telegram.NewInstance(cfg.Telegram.BotToken, cfg.Telegram.ChatID)

	if err != nil {
		log.Fatalf("Error telegram initialization: %v", err)
		return
	}

	scheduler := services.NewInstance(provider, sender, cfg.Schedule)
	scheduler.Start()
	waitForShutdown(scheduler)

	log.Println("Quiz Please Scheduler is stopping...")
}

func waitForShutdown(scheduler *services.Scheduler) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received stop signal: %v", sig)

	scheduler.Stop()
}
