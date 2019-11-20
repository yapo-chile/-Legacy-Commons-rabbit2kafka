package infrastructure

import (
	"log"
	"os"
	"time"

	"github.com/Yapo/logger"
	"github.mpi-internal.com/Yapo/rabbit2kafka/pkg/interfaces"
	sarama "gopkg.in/Shopify/sarama.v1"
)

//KafkaProducer struct representing a message producer for kafka
type KafkaProducer struct {
	producer sarama.SyncProducer
}

//NewKafkaProducer creates a new KafkaProducer with the given brokers
func NewKafkaProducer(brokerList []string) (interfaces.MessageHandler, error) {

	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.ClientID = "rabbitmq-to-kafka"
	config.Metadata.Retry.Max = 5
	config.Metadata.Retry.Backoff = 350 * time.Millisecond
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	sarama.Logger = log.New(os.Stdout, "[Kafka-client] ", log.LstdFlags)

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		logger.Info("Failed to start Sarama producer:", err)
		return nil, err
	}
	kafkaProducer := KafkaProducer{producer: producer}
	messageSender := interfaces.MessageHandler(kafkaProducer)
	return messageSender, nil
}

//Close close the KafkaProducer
func (k KafkaProducer) Close() error {
	err := k.producer.Close()
	if err != nil {
		logger.Info("Failed to shut down kafka producer cleanly: %+v", err)
	}
	return err
}

//SendMessage sends a message with the specified topic
func (k KafkaProducer) SendMessage(topic string, message string) error {
	// We are not setting a message key, which means that all messages will
	// be distributed randomly over the different partitions.
	_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("1"),
		Value: sarama.StringEncoder(message),
	})

	if err != nil {
		logger.Info("Failed to send the message %s: %s", message, err)
	} else {
		// The tuple (topic, partition, offset) can be used as a unique identifier
		// for a message in a Kafka cluster.
		// logger.Info("message sent with unique identifier %s/%d/%d", topic, partition, offset)
	}
	return err
}
