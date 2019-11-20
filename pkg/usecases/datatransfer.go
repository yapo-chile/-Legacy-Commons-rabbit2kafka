package usecases

import "github.mpi-internal.com/Yapo/rabbit2kafka/pkg/domain"

//MessageTransfer struct that represents the transfer from the storage reader to the message sender
type MessageTransfer struct {
	reader domain.StorageReader
	sender domain.MessageSender
}

//NewMessageTransfer constructor for the MessageTransfer struct
func NewMessageTransfer(reader domain.StorageReader, sender domain.MessageSender) *MessageTransfer {
	messageTransfer := new(MessageTransfer)
	messageTransfer.reader = reader
	messageTransfer.sender = sender
	return messageTransfer
}

//StartReader starts the storage reader
func (mt *MessageTransfer) StartReader(async bool) {
	mt.reader.SetReader(mt.ReaderFunction)
	mt.reader.Start(async)
}

//ReaderFunction function to be used to receive and transfer the messages
func (mt *MessageTransfer) ReaderFunction(message domain.StorageMessage) {
	err := mt.sender.SendMessage(message.GetMessage())
	message.Remove(err == nil)
}
