package domain

import (
	"time"

	"github.com/google/uuid"
)

type WebhookStatus string

const (
	WebhookPending   WebhookStatus = "PENDING"
	WebhookDelivered WebhookStatus = "DELIVERED"
	WebhookFailed    WebhookStatus = "FAILED"
)

type Webhook struct {
	ID           uuid.UUID
	PaymentID    uuid.UUID
	URL          string
	Status       WebhookStatus
	Attempt      int
	MaxAttempts  int
	Payload      []byte
	ResponseCode *int
	LastError    string
	DeliveredAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
