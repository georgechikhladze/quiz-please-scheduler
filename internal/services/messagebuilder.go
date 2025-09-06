package services

import (
	"fmt"
	"strings"

	"quiz-please-scheduler/pkg/gameprovider"
)

func GetGamesMessage(games map[int][]gameprovider.Game) string {
	if len(games[1]) == 0 && len(games[2]) == 0 {
		return "На данный момент нет доступных игр 😔"
	}

	var sb strings.Builder

	if len(games[1]) > 0 {
		sb.WriteString("*Есть места:\n\n*")

		for _, game := range games[1] {
			addGame(&sb, game)
		}
	}

	if len(games[2]) > 0 {
		sb.WriteString("*Резерв:\n\n*")

		for _, game := range games[2] {
			addGame(&sb, game)
		}
	}

	return sb.String()
}

func addGame(sb *strings.Builder, game gameprovider.Game) {
	sb.WriteString(game.Date)
	sb.WriteString("\n")

	formatted := fmt.Sprintf("[%s](%s), %s %s", game.Number, game.Link, game.Place, game.Time)
	sb.WriteString(formatted)
	sb.WriteString("\n\n")
}
