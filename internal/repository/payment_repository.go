package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/joaodddev/paymentprocessor/internal/domain"
)

type PaymentRepository interface {
	Create(
		ctx context.Context,
		payment *domain.Payment,
	) error

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Payment, error)

	List(
		ctx context.Context,
	) ([]domain.Payment, error)

	UpdateStatus(
		ctx context.Context,
		id uuid.UUID,
		status domain.PaymentStatus,
	) error
}
