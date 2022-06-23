package message

type MessageRepository interface {
	GetMessageData(id MessageId, messagePkg string) (text string, clientErrorType ClientErrorType)
}
