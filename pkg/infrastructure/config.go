package infrastructure

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ServiceConf holds configuration for this Service
type ServiceConf struct {
	Host      string `env:"HOST" envDefault:":8080"`
	Profiling bool   `env:"PROFILING" envDefault:"true"`
}

// LoggerConf holds configuration for logging
// LogLevel definition:
//   0 - Debug
//   1 - Info
//   2 - Warning
//   3 - Error
//   4 - Critic
type LoggerConf struct {
	SyslogIdentity string `env:"SYSLOG_IDENTITY"`
	SyslogEnabled  bool   `env:"SYSLOG_ENABLED" envDefault:"false"`
	StdlogEnabled  bool   `env:"STDLOG_ENABLED" envDefault:"true"`
	LogLevel       int    `env:"LOG_LEVEL" envDefault:"0"`
}

//RabbitMQConf conf for a RabbitMQ Consumer
type RabbitMQConf struct {
	Host         string `env:"HOST" envDefault:"rabbit"`
	Port         string `env:"PORT" envDefault:"5672"`
	Queue        string `env:"QUEUE" envDefault:"backend_event"`
	VHost        string `env:"VHOST" envDefault:"backend_event"`
	Exchange     string `env:"EXCHANGE" envDefault:"backend_event"`
	ExchangeType string `env:"EXCHANGE_TYPE" envDefault:"topic"`
	User         string `env:"USER" envDefault:"yapo"`
	Password     string `env:"PASSWORD" envDefault:"yapo2014"`
	ConsumerTag  string `env:"CONSUMER_TAG"`
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

// LoadFromEnv loads the config data from the environment variables
func LoadFromEnv(data interface{}) {
	load(reflect.ValueOf(data), "", "")
}

// valueFromEnv lookup the best value for a variable on the environment
func valueFromEnv(envTag, envDefault string) string {
	// Maybe it's a secret and <envTag>_FILE points to a file with the value
	// https://rancher.com/docs/rancher/v1.6/en/cattle/secrets/#docker-hub-images
	if fileName, ok := os.LookupEnv(fmt.Sprintf("%s_FILE", envTag)); ok {
		// filepath.Clean() will clean the input path and remove some unnecessary things
		// like multiple separators doble "." and others
		// if for some reason you are having troubles reaching your file, check the
		// output of the Clean function and test if its what you expect
		// you can find more info here: https://golang.org/pkg/path/filepath/#Clean
		b, err := ioutil.ReadFile(filepath.Clean(fileName))
		if err == nil {
			return string(b)
		}
		fmt.Print(err)
	}
	// The value might be set directly on the environment
	if value, ok := os.LookupEnv(envTag); ok {
		return value
	}
	// Nothing to do, return the default
	return envDefault
}

// load the variable defined in the envTag into Value
func load(conf reflect.Value, envTag, envDefault string) {
	if conf.Kind() == reflect.Ptr {
		reflectedConf := reflect.Indirect(conf)
		// Only attempt to set writeable variables
		if reflectedConf.IsValid() && reflectedConf.CanSet() {
			value := valueFromEnv(envTag, envDefault)
			// Print message if config is missing
			if envTag != "" && value == "" && !strings.HasSuffix(envTag, "_") {
				fmt.Printf("Config for %s missing\n", envTag)
			}
			switch reflectedConf.Interface().(type) {
			case int:
				if value, err := strconv.ParseInt(value, 10, 32); err == nil {
					reflectedConf.Set(reflect.ValueOf(int(value)))
				}
			case int64:
				if value, err := strconv.ParseInt(value, 10, 64); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case uint32:
				if value, err := strconv.ParseUint(value, 10, 32); err == nil {
					reflectedConf.Set(reflect.ValueOf(uint32(value)))
				}
			case float64:
				if value, err := strconv.ParseFloat(value, 64); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case string:
				reflectedConf.Set(reflect.ValueOf(value))
			case bool:
				if value, err := strconv.ParseBool(value); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case time.Time:
				if value, err := time.Parse(time.RFC3339, value); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case time.Duration:
				if t, err := time.ParseDuration(value); err == nil {
					reflectedConf.Set(reflect.ValueOf(t))
				}
			}
			if reflectedConf.Kind() == reflect.Struct {
				// Recursively load inner struct fields
				for i := 0; i < reflectedConf.NumField(); i++ {
					if tag, ok := reflectedConf.Type().Field(i).Tag.Lookup("env"); ok {
						def, _ := reflectedConf.Type().Field(i).Tag.Lookup("envDefault")
						load(reflectedConf.Field(i).Addr(), envTag+tag, def)
					}
				}
			}
		}
	}
}
