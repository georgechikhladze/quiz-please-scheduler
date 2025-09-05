package services

import (
	"log"

	"quiz-please-scheduler/pkg/gameprovider"
	"quiz-please-scheduler/pkg/telegram"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	provider gameprovider.Provider
	telegram telegram.Notifier
	cron     *cron.Cron
	schedule string
}

func NewInstance(
	provider gameprovider.Provider,
	telegram telegram.Notifier,
	schedule string) *Scheduler {
	return &Scheduler{
		provider: provider,
		telegram: telegram,
		schedule: schedule,
		cron:     cron.New(),
	}
}

func (s *Scheduler) Start() {
	_, err := s.cron.AddFunc(s.schedule, s.sendGamesMessage)
	if err != nil {
		log.Fatalf("Error adding task to cron: %v", err)
		return
	}

	s.cron.Start()
	log.Printf("Scheduler has started. Cron: %s", s.schedule)
	//go s.sendGamesMessage()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler has stopped")
}

func (s *Scheduler) sendGamesMessage() {
	games := s.provider.GetGamesList()

	log.Printf("Received %d opened and %d reserved games", len(games[1]), len(games[2]))

	if s.telegram != nil {
		message := GetGamesMessage(games)

		err := s.telegram.SendMessage(message)
		if err != nil {
			log.Printf("Error sending message to Telegram: %v", err)
		} else {
			log.Println("Schedule has successfully sent to Telegram")
		}
	}
}
