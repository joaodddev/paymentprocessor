# 💳 Payment Processor

A payment processing system built in **Go**, focused on modern architecture, asynchronous processing, scalability, and backend best practices.

> Created for study and portfolio purposes, simulating a payment system used by fintechs and companies that need to process transactions reliably and resiliently.

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22 |
| API | Gin |
| Database | PostgreSQL 16 |
| Cache / Idempotency | Redis 7 |
| Messaging | Apache Kafka |
| Migrations | golang-migrate |
| Logging | zerolog |
| Infrastructure | Docker Compose |

## 💳 Payment Methods

| Method | Behavior |
|---|---|
| PIX | Approved immediately |
| CREDIT_CARD | Approved after gateway simulation |
| BOLETO | Stays pending until confirmed |