package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string
type PaymentMethod string

const (
	PaymentPending    PaymentStatus = "PENDING"
	PaymentProcessing PaymentStatus = "PROCESSING"
	PaymentApproved   PaymentStatus = "APPROVED"
	PaymentRejected   PaymentStatus = "REJECTED"
)

const (
	PaymentMethodPIX        PaymentMethod = "PIX"
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	PaymentMethodBoleto     PaymentMethod = "BOLETO"
)

type Payment struct {
	ID             uuid.UUID
	Amount         int64
	Currency       string
	CustomerID     uuid.UUID
	Method         PaymentMethod
	Status         PaymentStatus
	IdempotencyKey string
	WebhookURL     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
