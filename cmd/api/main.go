package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/joaodddev/paymentprocessor/internal/config"
	"github.com/joaodddev/paymentprocessor/internal/repository"
	"github.com/joaodddev/paymentprocessor/internal/service"
)

func main() {

	log.Println("🚀 Iniciando Payment Processor...")

	cfg := config.Load()

	log.Println("✅ Configurações carregadas")

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	db, err := repository.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("✅ PostgreSQL conectado")

	paymentRepository := repository.NewPostgresPaymentRepository(db)

	paymentService := service.NewPaymentService(paymentRepository)

	_ = paymentService

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	log.Printf("🌐 Servidor iniciado na porta %s\n", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
