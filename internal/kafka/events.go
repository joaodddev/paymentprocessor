package kafka

import "time"

type PaymentCreatedEvent struct {
	PaymentID      string    `json:"payment_id"`
	Amount         int64     `json:"amount"`
	Currency       string    `json:"currency"`
	CustomerID     string    `json:"customer_id"`
	Method         string    `json:"method"`
	IdempotencyKey string    `json:"idempotency_key"`
	WebhookURL     string    `json:"webhook_url"`
	CreatedAt      time.Time `json:"created_at"`
}

type PaymentProcessedEvent struct {
	PaymentID   string    `json:"payment_id"`
	Status      string    `json:"status"`
	Method      string    `json:"method"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	CustomerID  string    `json:"customer_id"`
	WebhookURL  string    `json:"webhook_url"`
	ProcessedAt time.Time `json:"processed_at"`
}
