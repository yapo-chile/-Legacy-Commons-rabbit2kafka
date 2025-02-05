package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Yapo/logger"
	"github.mpi-internal.com/Yapo/rabbit2kafka/pkg/infrastructure"
	"github.mpi-internal.com/Yapo/rabbit2kafka/pkg/interfaces"
	"github.mpi-internal.com/Yapo/rabbit2kafka/pkg/usecases"
)

var kafkaProducer infrastructure.KafkaProducer

func main() {
	fmt.Println("Loading config")
	var conf infrastructure.Config
	infrastructure.LoadFromEnv(&conf)
	if jconf, err := json.MarshalIndent(conf, "", "    "); err == nil {
		fmt.Printf("Config: \n%s\n", jconf)
	} else {
		fmt.Printf("Config: \n%+v\n", conf)
	}
	fmt.Println("Setting up logger")
	loggerConf := logger.LogConfig{
		Syslog: logger.SyslogConfig{
			Enabled:  conf.LoggerConf.SyslogEnabled,
			Identity: conf.LoggerConf.SyslogIdentity,
		},
		Stdlog: logger.StdlogConfig{
			Enabled: conf.LoggerConf.StdlogEnabled,
		},
	}
	if err := logger.Init(loggerConf); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger.SetLogLevel(conf.LoggerConf.LogLevel)
	fmt.Printf("LogLevel: %d\n", conf.LoggerConf.LogLevel)

	k, err := infrastructure.NewKafkaProducer([]string{conf.KafkaConf.GetBroker()})
	if err != nil {
		logger.Error("Error starting kafka producer: %+v", err)
		os.Exit(1)
	}
	c := infrastructure.NewConsumer(
		conf.RabbitMQConf.Host,
		conf.RabbitMQConf.Port,
		conf.RabbitMQConf.Queue,
		conf.RabbitMQConf.User,
		conf.RabbitMQConf.Password,
		conf.RabbitMQConf.ConsumerTag,
		conf.RabbitMQConf.VHost,
		conf.RabbitMQConf.Exchange,
		conf.RabbitMQConf.ExchangeType,
	)
	storageRepo := interfaces.NewStorageRepo(c)
	messageRepo := interfaces.NewMessageRepo(k)
	messageRepo.Topic = conf.KafkaConf.Topic
	messageTransfer := usecases.NewMessageTransfer(storageRepo, messageRepo)
	messageTransfer.StartReader(false)
}
