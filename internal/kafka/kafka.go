package kafka

import (
	"L0-arch/internal/config"
	"github.com/IBM/sarama"
	"time"
)

type Kafka struct {
	broker string
}

func New(config config.KafkaConfig) *Kafka {
	return &Kafka{
		broker: config.Broker,
	}
}

func (k *Kafka) Produce(topic string, value []byte) (partition int32, offset int64, err error) {
	producer, err := sarama.NewSyncProducer([]string{k.broker}, nil)
	if err != nil {
		return 0, 0, err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(value),
		Timestamp: time.Now(),
	}
	return producer.SendMessage(msg)
}
