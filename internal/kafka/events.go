package kafka

import "time"

type PaymentCreatedEvent struct {
	PaymentID      string    `json:"payment_id"`
	Amount         int64     `json:"amount"`
	Currency       string    `json:"currency"`
	CustomerID     string    `json:"customer_id"`
	Method         string    `json:"method"`
	IdempotencyKey string    `json:"idempotency_key"`
	CreatedAt      time.Time `json:"created_at"`
}
