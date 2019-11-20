package interfaces

//MessageHandler represents a handler that interacts with the message infrastructure
type MessageHandler interface {
	SendMessage(topic string, message string) error
}

//MessageRepo struct that contains a message handler
type MessageRepo struct {
	messageHandler MessageHandler
	Topic          string
}

//NewMessageRepo constructor for a MessageRepo
func NewMessageRepo(messageHandler MessageHandler) *MessageRepo {
	messageRepo := new(MessageRepo)
	messageRepo.messageHandler = messageHandler
	return messageRepo
}

//SendMessage Stores the message with the given topic
func (m *MessageRepo) SendMessage(message string) error {
	return m.messageHandler.SendMessage(m.Topic, message)
}
