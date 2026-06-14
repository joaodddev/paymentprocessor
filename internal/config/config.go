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
	AppEnv           string

	RedisAddr     string
	RedisPassword string
	RedisDB       string

	KafkaBrokers               string
	KafkaGroupID               string
	KafkaTopicPaymentCreated   string
	KafkaTopicPaymentProcessed string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("arquivo .env não encontrado, usando variáveis do ambiente")
	}

	return &Config{
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		AppPort:          os.Getenv("APP_PORT"),
		AppEnv:           os.Getenv("APP_ENV"),

		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       os.Getenv("REDIS_DB"),

		KafkaBrokers:               os.Getenv("KAFKA_BROKERS"),
		KafkaGroupID:               os.Getenv("KAFKA_GROUP_ID"),
		KafkaTopicPaymentCreated:   os.Getenv("KAFKA_TOPIC_PAYMENT_CREATED"),
		KafkaTopicPaymentProcessed: os.Getenv("KAFKA_TOPIC_PAYMENT_PROCESSED"),
	}
}
