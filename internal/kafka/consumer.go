package kafka

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	srvorders "L0-arch/internal/service/orders"
	"L0-arch/pkg/metrics"
	vld "L0-arch/pkg/validator"

	"github.com/IBM/sarama"
)

type OrderSaver interface {
	SaveOrder(ctx context.Context, o srvorders.Model) error
}

func (k *Kafka) StartConsumerGroup(ctx context.Context, topic, groupID string, saver OrderSaver) error {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	cg, err := sarama.NewConsumerGroup([]string{k.broker}, groupID, cfg)
	if err != nil {
		return err
	}
	defer cg.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		cancel()
	}()

	handler := &groupHandler{saver: saver}

	baseDelay := 200 * time.Millisecond
	maxDelay := 3 * time.Second
	delay := baseDelay

	for {
		if err := cg.Consume(ctx, []string{topic}, handler); err != nil {
			log.Printf("kafka consume error: %v", err)
			metrics.KafkaConsumeErrors.Inc()
			jitter := time.Duration(rand.Int63n(int64(delay)))
			if jitter > maxDelay {
				jitter = maxDelay
			}
			select {
			case <-time.After(jitter):
			case <-ctx.Done():
				return nil
			}
			if delay < maxDelay/2 {
				delay *= 2
			}
		} else {
			delay = baseDelay
		}
		if ctx.Err() != nil {
			return nil
		}
	}
}

type groupHandler struct {
	saver OrderSaver
}

func (h *groupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *groupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *groupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var o srvorders.Model
		if err := json.Unmarshal(msg.Value, &o); err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}
		if err := vld.Validate(&o); err != nil {
			log.Printf("validation error: %v", err)
			continue
		}
		if err := h.saver.SaveOrder(context.Background(), o); err != nil {
			log.Printf("save error: %v", err)
			continue
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
