package service

import (
	"fmt"
	"strings"

	"quiz-please-scheduler/gameprovider"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	bot    *tgbotapi.BotAPI
	chatID string
}

func New(botToken, chatID string) (*TelegramService, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)

	if err != nil {
		return nil, err
	}

	return &TelegramService{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func (t *TelegramService) SendGames(games map[int][]gameprovider.Game) error {
	if len(games[1]) == 0 && len(games[2]) == 0 {
		return t.sendMessage("На данный момент нет доступных игр 😔")
	}

	var sb strings.Builder
	sb.WriteString("*Есть места:\n\n*")

	for _, game := range games[1] {
		addGame(&sb, game)
	}

	sb.WriteString("*Резерв:\n\n*")

	for _, game := range games[2] {
		addGame(&sb, game)
	}

	return t.sendMessage(sb.String())
}

func addGame(sb *strings.Builder, game gameprovider.Game) {
	sb.WriteString(game.Date)
	sb.WriteString("\n")

	formatted := fmt.Sprintf("[%s](%s), %s %s", game.Number, game.Link, game.Place, game.Time)
	sb.WriteString(formatted)
	sb.WriteString("\n\n")
}

func (t *TelegramService) sendMessage(text string) error {
	text = strings.ReplaceAll(text, "#", "\\#")

	msg := tgbotapi.NewMessageToChannel(t.chatID, text)
	msg.ParseMode = "MarkdownV2"
	msg.DisableWebPagePreview = false

	_, err := t.bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}
