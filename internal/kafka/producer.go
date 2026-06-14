package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/joaodddev/paymentprocessor/internal/config"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writers map[string]*kafka.Writer
	cfg     *config.Config
}

func NewProducer(cfg *config.Config) *Producer {
	makeWriter := func(topic string) *kafka.Writer {
		return &kafka.Writer{
			Addr:         kafka.TCP(cfg.KafkaBrokers),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
		}
	}

	return &Producer{
		cfg: cfg,
		writers: map[string]*kafka.Writer{
			cfg.KafkaTopicPaymentCreated:   makeWriter(cfg.KafkaTopicPaymentCreated),
			cfg.KafkaTopicPaymentProcessed: makeWriter(cfg.KafkaTopicPaymentProcessed),
		},
	}
}

func (p *Producer) publish(ctx context.Context, topic string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar evento: %w", err)
	}

	writer, ok := p.writers[topic]
	if !ok {
		return fmt.Errorf("writer não encontrado para tópico: %s", topic)
	}

	return writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", time.Now().UnixNano())),
		Value: data,
	})
}

func (p *Producer) PublishPaymentCreated(ctx context.Context, payload any) error {
	return p.publish(ctx, p.cfg.KafkaTopicPaymentCreated, payload)
}

func (p *Producer) PublishPaymentProcessed(ctx context.Context, payload any) error {
	return p.publish(ctx, p.cfg.KafkaTopicPaymentProcessed, payload)
}

func (p *Producer) Close() error {
	for _, w := range p.writers {
		w.Close()
	}
	return nil
}
