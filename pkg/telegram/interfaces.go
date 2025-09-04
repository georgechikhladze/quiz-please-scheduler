package telegram

type Notifier interface {
	SendMessage(string) error
}
