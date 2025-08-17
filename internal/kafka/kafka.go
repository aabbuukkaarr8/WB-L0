package kafka

import (
	"L0-arch/internal/config"
)

type Kafka struct {
	broker string
}

func New(config config.KafkaConfig) *Kafka {
	return &Kafka{
		broker: config.Broker,
	}
}
