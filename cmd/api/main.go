package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/joaodddev/paymentprocessor/internal/config"
	"github.com/joaodddev/paymentprocessor/internal/repository"
)

func main() {

	log.Println("🚀 Iniciando Payment Processor...")

	// Carrega configurações
	cfg := config.Load()

	log.Println("✅ Configurações carregadas")

	// Define porta padrão caso não exista
	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	// Conecta ao PostgreSQL
	db, err := repository.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("✅ PostgreSQL conectado")

	// Inicializa Repository
	paymentRepository := repository.NewPostgresPaymentRepository(db)

	// Apenas para evitar erro de variável não utilizada
	_ = paymentRepository

	// Inicializa Gin
	router := gin.Default()

	// Health Check
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
