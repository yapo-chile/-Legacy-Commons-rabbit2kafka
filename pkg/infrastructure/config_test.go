package infrastructure

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestConfigLoad(t *testing.T) {
	configVariables := []string{
		"RABBITMQ_HOST",
		"RABBITMQ_PORT",
		"RABBITMQ_QUEUE",
		"RABBITMQ_EXCHANGE",
		"LOGGER_SYSLOG_ENABLED",
		"LOGGER_SYSLOG_IDENTITY",
		"LOGGER_STDLOG_ENABLED",
		"LOGGER_LOG_LEVEL",
		"KAFKA_HOST",
		"KAFKA_PORT",
		"KAFKA_TOPIC",
	}
	storedValues := make([]string, len(configVariables))
	testConfigVariables := []string{
		"localhost",
		"3785",
		"test_queue",
		"test_exchange",
		"true",
		"test_log",
		"true",
		"1",
		"test-kafka",
		"9092",
		"test-topic",
	}
	//store the environ variables and set the desired values
	for index, variable := range configVariables {
		storedValues[index] = os.Getenv(variable)
		os.Setenv(variable, testConfigVariables[index])
	}
	//Load the test enviroment
	c := LoadConfig()
	for index, variable := range configVariables {
		os.Setenv(variable, storedValues[index])
	}
	//Validate the load
	assert.Equal(t, fmt.Sprintf("%+v", c), `&{LoggerConf:{SyslogEnabled:true SyslogIdentity:test_log StdlogEnabled:true LogLevel:1} RabbitMQConf:{Host:localhost Port:3785 Queue:test_queue Exchange:test_exchange} KafkaConf:{Host:test-kafka Port:9092 Topic:test-topic}}`)
	assert.Equal(t, c.KafkaConf.GetBroker(), `test-kafka:9092`)
}
