package telegram

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID string
}

func NewInstance(botToken, chatID string) (TelegramSender, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)

	if err != nil {
		return nil, err
	}

	log.Printf("Notificatinos will be send to the chat: %s", chatID)

	return &TelegramNotifier{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func (t TelegramNotifier) SendMessage(text string) error {
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
