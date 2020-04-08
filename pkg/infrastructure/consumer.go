package infrastructure

import (
	"fmt"
	"time"

	"github.com/Yapo/logger"
	"github.com/streadway/amqp"
	"github.mpi-internal.com/Yapo/rabbit2kafka/pkg/domain"
	"github.mpi-internal.com/Yapo/rabbit2kafka/pkg/interfaces"
)

//RabbitMQDelivery wrapper type for the received messages
type RabbitMQDelivery amqp.Delivery

//GetMessage returns the content of the body as a string
func (r RabbitMQDelivery) GetMessage() string {
	return string(r.Body[:])
}

//Remove remove the received message from the queue
func (r RabbitMQDelivery) Remove(remove bool) {
	delivery := amqp.Delivery(r)
	delivery.Ack(remove)
}

//Consumer represents a connection that receives the messages
type Consumer struct {
	QueueData    interfaces.QueueData
	ConsumerTag  string
	VHost        string
	Exchange     string
	ExchangeType string
	reader       domain.Reader
	Connection   *amqp.Connection
	Queue        amqp.Queue
	Channel      *amqp.Channel
	Done         chan error
	closeError   chan *amqp.Error
	connected    bool
}

//RecoverIntervalTime interval between reconnection tries
const RecoverIntervalTime = 15 * time.Second

//NewConsumer constructor
func NewConsumer(
	host string,
	port string,
	name string,
	username string,
	password string,
	consumerTag string,
	vhost string,
	exchange string,
	exchangeType string,
) interfaces.StorageHandler {
	consumer := new(Consumer)
	queue := new(interfaces.QueueData)
	queue.Host = host
	queue.Port = port
	queue.Name = name
	queue.Username = username
	queue.Password = password
	consumer.QueueData = *queue
	consumer.ConsumerTag = consumerTag
	consumer.VHost = vhost
	consumer.Exchange = exchange
	consumer.ExchangeType = exchangeType
	consumer.Done = make(chan error)
	consumer.closeError = make(chan *amqp.Error)

	storageHandler := interfaces.StorageHandler(consumer)
	return storageHandler
}

//SetReader sets the reader function to manage the messages
func (c *Consumer) SetReader(reader domain.Reader) {
	c.reader = reader
}

//Connect connects the consumer to RabbitMQ
func (c *Consumer) Connect() {
	connectionURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/"+c.VHost,
		c.QueueData.Username,
		c.QueueData.Password,
		c.QueueData.Host,
		c.QueueData.Port,
	)
	if c.QueueData.Username == "" && c.QueueData.Password == "" {
		connectionURL = fmt.Sprintf(
			"amqp://%s:%s/",
			c.QueueData.Host,
			c.QueueData.Port,
		)
	}

	conn, err := amqp.Dial(connectionURL)
	for err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
		time.Sleep(RecoverIntervalTime)
		conn, err = amqp.Dial(connectionURL)
	}
	c.Connection = conn
	c.closeError = make(chan *amqp.Error)
	c.Connection.NotifyClose(c.closeError)

	ch, err := conn.Channel()
	if failOnError(err, "Failed to open a channel") {
		return
	}

	ch.Qos(10000, 0, false)

	c.Channel = ch
	err = ch.ExchangeDeclare(
		c.Exchange,     // name
		c.ExchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if failOnError(err, "Failed to declare an exchange") {
		return
	}
	q, err := ch.QueueDeclare(
		c.QueueData.Name, // name
		true,             // durable
		false,            // delete when usused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if failOnError(err, "Failed to declare a queue") {
		return
	}
	c.Queue = q

	err = ch.QueueBind(
		q.Name,     // queue name
		"",         // routing key
		c.Exchange, // exchange
		false,
		nil)
	if failOnError(err, "Failed to bind a queue") {
		return
	}
	c.connected = true
	c.RunOnce()
	return
}

//Reconnector function to reconnect the consumer
func (c *Consumer) Reconnector() {
	var rabbitErr *amqp.Error
	for {
		rabbitErr = <-c.closeError
		logger.Debug("received close error: %+v", rabbitErr)
		if rabbitErr != nil {
			logger.Info("try Connecting")
			c.Connect()
		}
	}
}

//RunOnce starts the consumer and listening through a go routine
func (c *Consumer) RunOnce() {
	if !c.connected {
		return
	}
	msgs, err := c.Channel.Consume(
		c.Queue.Name,  // queue
		c.ConsumerTag, // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if failOnError(err, "Failed to register a consumer") {
		return
	}
	go c.processMessages(msgs)
	logger.Info("queue: %s", c.Queue.Name)
	logger.Info(" [*] Waiting for messages")
}

func (c *Consumer) processMessages(deliveries <-chan amqp.Delivery) {
	logger.Info("worker: started")
	for d := range deliveries {
		delivery := RabbitMQDelivery(d)
		c.reader(delivery)
	}
	logger.Info("worker: deliveries channel closed")
	c.Done <- nil
	logger.Info("worker: ended")
}

//Start starts the consumer
func (c *Consumer) Start(async bool) {
	if async {
		go c.Reconnector()
		c.closeError <- amqp.ErrClosed
	} else {
		c.Connect()
		c.Reconnector()
	}
}

//Stop stops the consumer
func (c *Consumer) Stop() error {
	logger.Info("try closing the channel")
	if err := c.Channel.Cancel(c.ConsumerTag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}
	logger.Info("try clossing the connection")
	if err := c.Connection.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}
	logger.Info("capture the return")
	returnError := <-c.Done
	logger.Info("ended stop. Error: %+v", returnError)
	return returnError
}

func failOnError(err error, msg string) bool {
	if err != nil {
		logger.Error("ERROR: %s: %s", msg, err)
		return true
	}
	return false
}
