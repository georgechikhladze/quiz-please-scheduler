package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"quiz-please-scheduler/config"
	"quiz-please-scheduler/gameprovider"
	"quiz-please-scheduler/service"
	"quiz-please-scheduler/telegram"
)

func main() {
	fmt.Println("Quiz Please Scheduler is starting...")

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error config loading: %v", err)
		return
	}

	provider := gameprovider.NewInstance()
	sender, err := telegram.NewInstance(cfg.Telegram.BotToken, cfg.Telegram.ChatID)

	if err != nil {
		log.Fatalf("Error telegram initialization: %v", err)
		return
	}

	scheduler := service.NewScheduler(provider, sender, cfg.Schedule)
	scheduler.Start()

	log.Printf("Background Service has started. Schedule: %s", cfg.Schedule)
	log.Printf("Notificatinos will be send to the chat: %s", cfg.Telegram.ChatID)

	waitForShutdown(scheduler)

	log.Println("Quiz Please Scheduler is stopping...")
}

func waitForShutdown(scheduler *service.Scheduler) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Get stop signal: %v", sig)

	scheduler.Stop()
}
