package infrastructure

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

//LoggerConf configuration for the logger struct used for logging
// LogLevel definition:
// 0 - Debug
// 1 - Info
// 2 - Warning
// 3 - Error
// 4 - Critic
type LoggerConf struct {
	SyslogEnabled  bool   `env:"SYSLOG_ENABLED" envDefault:"false"`
	SyslogIdentity string `env:"SYSLOG_IDENTITY"`
	StdlogEnabled  bool   `env:"STDLOG_ENABLED" envDefault:"true"`
	LogLevel       int    `env:"LOG_LEVEL" envDefault:"0"`
}

//RabbitMQConf conf for a RabbitMQ Consumer
type RabbitMQConf struct {
	Host        string `env:"HOST" envDefault:"rabbit"`
	Port        string `env:"PORT" envDefault:"5672"`
	Queue       string `env:"QUEUE" envDefault:"backend_event`
	Exchange    string `env:"EXCHANGE" envDefault:"/`
	User        string `env:"USER" envDefault:"guest"`
	Password    string `env:"PASSWORD" envDefault:"guest"`
	ConsumerTag string `env:"CONSUMER_TAG"`
}

//KafkaConf conf for a Kafka SyncProducer
type KafkaConf struct {
	Host  string `env:"HOST" envDefault:"kafka"`
	Port  string `env:"PORT" envDefault:"9092"`
	Topic string `env:"TOPIC" envDefault:"events-queue"`
}

//GetBroker returns the broker to be used for Kafka
func (k *KafkaConf) GetBroker() string {
	return fmt.Sprintf("%s:%s", k.Host, k.Port)
}

//Config struct that represents the config of the microservice
type Config struct {
	LoggerConf   LoggerConf   `env:"LOGGER_"`
	RabbitMQConf RabbitMQConf `env:"RABBITMQ_"`
	KafkaConf    KafkaConf    `env:"KAFKA_"`
}

//This function was copied from here: https://github.com/JorgePoblete/goenv
//load loads the configuracion on the reflect.Value with the envTag anf the envDefault value
func load(conf reflect.Value, envTag, envDefault string) {
	// here conf could be either a struct or just a variable
	// if it's a variable we just set its value to the value of the
	// environment variable referenced by its tag, or its default, otherwise we recursively
	// set the struct value to the value returned by load(...) of each of its
	// individual fields

	if conf.Kind() == reflect.Ptr {
		reflectedConf := reflect.Indirect(conf)
		// we should only keep going if we can set values
		if reflectedConf.IsValid() && reflectedConf.CanSet() {
			value, ok := os.LookupEnv(envTag)
			// if the env variable is not set we just use the envDefault
			if !ok {
				value = envDefault
			}
			switch reflectedConf.Kind() {
			case reflect.Struct:
				for i := 0; i < reflectedConf.NumField(); i++ {
					if tag, ok := reflectedConf.Type().Field(i).Tag.Lookup("env"); ok {
						def, _ := reflectedConf.Type().Field(i).Tag.Lookup("envDefault")
						load(reflectedConf.Field(i).Addr(), envTag+tag, def)
					}
				}
				break
			// Here for each type we should make a cast of the env variable and then set the value
			case reflect.String:
				reflectedConf.SetString(value)
				break
			case reflect.Int:
				value, _ := strconv.Atoi(value)
				reflectedConf.Set(reflect.ValueOf(value))
				break
			case reflect.Bool:
				value, _ := strconv.ParseBool(value)
				reflectedConf.Set(reflect.ValueOf(value))
			}
		}
	}
}

//LoadConfig load the microservice conf using the enviroment variables
func LoadConfig() *Config {
	conf := Config{}
	load(reflect.ValueOf(&conf), "", "")
	return &conf
}
