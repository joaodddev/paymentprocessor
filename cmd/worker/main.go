package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/joaodddev/paymentprocessor/internal/config"
	kafkapkg "github.com/joaodddev/paymentprocessor/internal/kafka"
	"github.com/joaodddev/paymentprocessor/internal/logger"
	"github.com/joaodddev/paymentprocessor/internal/repository"
	"github.com/joaodddev/paymentprocessor/internal/worker"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.AppEnv)

	log.Info().Msg("🚀 Iniciando Payment Processor Worker...")

	db, err := repository.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("erro ao conectar no postgres")
	}
	defer db.Close()

	paymentConsumer := kafkapkg.NewConsumer(cfg, cfg.KafkaTopicPaymentCreated)
	defer paymentConsumer.Close()

	webhookConsumer := kafkapkg.NewConsumer(cfg, cfg.KafkaTopicPaymentProcessed)
	defer webhookConsumer.Close()

	producer := kafkapkg.NewProducer(cfg)
	defer producer.Close()

	paymentRepository := repository.NewPostgresPaymentRepository(db)
	webhookRepository := repository.NewPostgresWebhookRepository(db)

	paymentWorker := worker.NewPaymentWorker(paymentConsumer, paymentRepository, producer)
	webhookNotifier := worker.NewWebhookNotifier(webhookConsumer, webhookRepository)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go paymentWorker.Run(ctx)
	go webhookNotifier.Run(ctx)

	log.Info().Msg("workers rodando, aguardando eventos...")

	<-quit
	log.Info().Msg("⏳ Encerrando workers...")
	cancel()
}
