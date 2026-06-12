package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "PENDING"
	PaymentProcessing PaymentStatus = "PROCESSING"
	PaymentApproved   PaymentStatus = "APPROVED"
	PaymentRejected   PaymentStatus = "REJECTED"
)

type Payment struct {
	ID             uuid.UUID
	Amount         int64
	Currency       string
	CustomerID     uuid.UUID
	Status         PaymentStatus
	IdempotencyKey string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
