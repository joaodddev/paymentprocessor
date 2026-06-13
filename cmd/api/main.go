package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/joaodddev/paymentprocessor/internal/config"
	"github.com/joaodddev/paymentprocessor/internal/handler"
	"github.com/joaodddev/paymentprocessor/internal/repository"
	"github.com/joaodddev/paymentprocessor/internal/service"
)

func main() {

	log.Println("🚀 Iniciando Payment Processor...")

	cfg := config.Load()

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	db, err := repository.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	paymentRepository := repository.NewPostgresPaymentRepository(db)

	paymentService := service.NewPaymentService(paymentRepository)

	paymentHandler := handler.NewPaymentHandler(paymentService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	router.POST("/payments", paymentHandler.CreatePayment)
	router.GET("/payments", paymentHandler.ListPayments)
	router.GET("/payments/:id", paymentHandler.GetPayment)

	log.Printf("Servidor iniciado na porta %s", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
