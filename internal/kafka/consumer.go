package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joaodddev/paymentprocessor/internal/config"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg *config.Config, topic string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{cfg.KafkaBrokers},
		GroupID:        cfg.KafkaGroupID,
		Topic:          topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		MaxWait:        1 * time.Second,
		CommitInterval: time.Second,
	})

	return &Consumer{reader: reader}
}

func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	msg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return kafka.Message{}, fmt.Errorf("erro ao ler mensagem: %w", err)
	}
	return msg, nil
}

func (c *Consumer) Commit(ctx context.Context, msg kafka.Message) error {
	return c.reader.CommitMessages(ctx, msg)
}

func (c *Consumer) Close() error {
	log.Println("fechando consumer Kafka...")
	return c.reader.Close()
}
