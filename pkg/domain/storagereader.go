package domain

//StorageReader reader for reading and processing messages
type StorageReader interface {
	Start(async bool)
	SetReader(reader Reader)
}

//Reader function used to manipulate the received messages
type Reader func(message StorageMessage)

//StorageMessage interface that represents a message received from the storage reader
type StorageMessage interface {
	GetMessage() string
	Remove(bool)
}
