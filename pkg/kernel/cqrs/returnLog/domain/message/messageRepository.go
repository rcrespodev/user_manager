package message

type MessageRepository interface {
	GetMessageText(id MessageId, messagePkg string) string
}
