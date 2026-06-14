package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/joaodddev/paymentprocessor/internal/cache"
	"github.com/joaodddev/paymentprocessor/internal/config"
	"github.com/joaodddev/paymentprocessor/internal/handler"
	kafkapkg "github.com/joaodddev/paymentprocessor/internal/kafka"
	"github.com/joaodddev/paymentprocessor/internal/logger"
	"github.com/joaodddev/paymentprocessor/internal/repository"
	"github.com/joaodddev/paymentprocessor/internal/service"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.AppEnv)

	log.Info().Msg("🚀 Iniciando Payment Processor API...")

	db, err := repository.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("erro ao conectar no postgres")
	}
	defer db.Close()

	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("erro ao conectar no redis")
	}
	defer redisClient.Close()

	producer := kafkapkg.NewProducer(cfg)
	defer producer.Close()

	paymentRepository := repository.NewPostgresPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepository, redisClient, producer)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})

	router.POST("/payments", paymentHandler.CreatePayment)
	router.GET("/payments", paymentHandler.ListPayments)
	router.GET("/payments/:id", paymentHandler.GetPayment)

	log.Info().Str("port", cfg.AppPort).Msg("servidor iniciado")

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal().Err(err).Msg("erro ao iniciar servidor")
	}
}
