package interfaces

//QueueData represents the data to connect to rabbitMQ
type QueueData struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}
