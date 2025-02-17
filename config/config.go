package cnf

import (
	"log"

	"github.com/caarlos0/env/v10"
)

type Mongodb struct {
	Mongourl string `env:"MONGO_URL" envDefault:"mongodb://localhost:27017/"`
}

type Config struct {
	Mongodb Mongodb
}

func LoadConfig() (Config, error) { 
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Printf("Failed to load config: %v", err)
		return Config{}, err 
	}
	return cfg, nil
}
type ConsumerConfig struct {
	Url                string `env:"URL" validate:"required" envDefault:"amqp://guest:guest@localhost:5672/"`
	Exchange           string `env:"EXCHANGE_NAME"  envDefault:"coures_add_exchange"`
	ExchangeType       string `env:"EXCHANGE_TYPE"  envDefault:"direct"`
	PrefetchCount      int    `env:"PREFETCH_COUNT"  envDefault:"100"`
	ConnectionPoolSize int    `env:"CONNECTIONPOOL_SIZE"  envDefault:"10"`
	QueueName          string `env:"QUEUE_NAME" envDefault:"coures_add"`
	BindingKeyName     string `env:"BINDING_KEY_NAME" envDefault:"coures_add_bindkey"`
	DelayedQueueName   string `env:"DELAYED_QUEUE_NAME" envDefault:"coures_add_delay_queue"`
	QueueTaskName      string `env:"COURES_QUEUE_TASK"  envDefault:"coures"`
}

var Consumerconfig ConsumerConfig

func LoadConsumer() {
	if err := env.Parse(&Consumerconfig); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
}
