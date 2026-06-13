package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("arquivo .env não encontrado, usando variáveis do ambiente")
	}

	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("falha ao inicializar migrate: %v", err)
	}
	defer m.Close()

	switch direction {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migration up falhou: %v", err)
		}
		fmt.Println("✅ Migrations aplicadas com sucesso")
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migration down falhou: %v", err)
		}
		fmt.Println("✅ Migrations revertidas com sucesso")
	default:
		log.Fatalf("direção inválida: %s (use 'up' ou 'down')", direction)
	}
}
