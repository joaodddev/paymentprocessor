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

	consumer := kafkapkg.NewConsumer(cfg, cfg.KafkaTopicPaymentCreated)
	defer consumer.Close()

	paymentRepository := repository.NewPostgresPaymentRepository(db)

	paymentWorker := worker.NewPaymentWorker(consumer, paymentRepository)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go paymentWorker.Run(ctx)

	<-quit
	log.Println("⏳ Encerrando worker...")
	cancel()
}
