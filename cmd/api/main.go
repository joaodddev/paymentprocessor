package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/joaodddev/paymentprocessor/internal/config"
	"github.com/joaodddev/paymentprocessor/internal/repository"
)

func main() {

	log.Println("🚀 Iniciando Payment Processor...")

	cfg := config.Load()

	log.Println("✅ Configurações carregadas")

	ctx := context.Background()

	db, err := repository.NewPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("✅ PostgreSQL conectado")

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	log.Printf("🌐 Servidor iniciado na porta %s\n", cfg.AppPort)

	router.Run(":" + cfg.AppPort)
}
