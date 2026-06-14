package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joaodddev/paymentprocessor/internal/config"
	kafkapkg "github.com/joaodddev/paymentprocessor/internal/kafka"
	"github.com/joaodddev/paymentprocessor/internal/repository"
	"github.com/joaodddev/paymentprocessor/internal/worker"
)

func main() {
	log.Println("🚀 Iniciando Payment Processor Worker...")

	cfg := config.Load()

	db, err := repository.NewPostgresPool(cfg)
	if err != nil {
		log.Fatalf("erro ao conectar no postgres: %v", err)
	}
	defer db.Close()

	// Consumers
	paymentConsumer := kafkapkg.NewConsumer(cfg, cfg.KafkaTopicPaymentCreated)
	defer paymentConsumer.Close()

	webhookConsumer := kafkapkg.NewConsumer(cfg, cfg.KafkaTopicPaymentProcessed)
	defer webhookConsumer.Close()

	// Producer (para publicar payment.processed)
	producer := kafkapkg.NewProducer(cfg)
	defer producer.Close()

	// Repositories
	paymentRepository := repository.NewPostgresPaymentRepository(db)
	webhookRepository := repository.NewPostgresWebhookRepository(db)

	// Workers
	paymentWorker := worker.NewPaymentWorker(paymentConsumer, paymentRepository, producer)
	webhookNotifier := worker.NewWebhookNotifier(webhookConsumer, webhookRepository)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go paymentWorker.Run(ctx)
	go webhookNotifier.Run(ctx)

	<-quit
	log.Println("⏳ Encerrando workers...")
	cancel()
}
