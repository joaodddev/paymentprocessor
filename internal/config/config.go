package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	AppPort          string
}

func Load() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("arquivo .env não encontrado, utilizando variáveis do ambiente")
	}

	return &Config{
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		AppPort:          os.Getenv("APP_PORT"),
	}
}
