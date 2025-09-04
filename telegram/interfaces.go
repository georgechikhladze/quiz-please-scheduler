package telegram

type TelegramSender interface {
	SendMessage(string) error
}
