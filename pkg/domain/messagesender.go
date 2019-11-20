package domain

//MessageSender interface for sending a string message
type MessageSender interface {
	SendMessage(message string) error
}
