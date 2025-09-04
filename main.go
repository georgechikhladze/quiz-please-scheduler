package main

import (
	"fmt"

	"quiz-please-scheduler/gameprovider"
)

func main() {
	fmt.Println("Quiz Please Scheduler is starting...")

	var provider gameprovider.GameProvider = gameprovider.New()

	games := provider.GetGamesList()
	fmt.Printf("Get games %v", len(games))
}
