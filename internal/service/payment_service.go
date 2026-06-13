package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/joaodddev/paymentprocessor/internal/domain"
	"github.com/joaodddev/paymentprocessor/internal/repository"
)

type PaymentService struct {
	repository repository.PaymentRepository
}

func NewPaymentService(
	repository repository.PaymentRepository,
) *PaymentService {
	return &PaymentService{
		repository: repository,
	}
}

func (s *PaymentService) CreatePayment(
	ctx context.Context,
	amount int64,
	currency string,
	customerID uuid.UUID,
	idempotencyKey string,
) (*domain.Payment, error) {

	now := time.Now()

	payment := &domain.Payment{
		ID:             uuid.New(),
		Amount:         amount,
		Currency:       currency,
		CustomerID:     customerID,
		Status:         domain.PaymentPending,
		IdempotencyKey: idempotencyKey,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repository.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) GetPaymentByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Payment, error) {

	return s.repository.FindByID(ctx, id)
}

func (s *PaymentService) ListPayments(
	ctx context.Context,
) ([]domain.Payment, error) {

	return s.repository.List(ctx)
}

func (s *PaymentService) UpdatePaymentStatus(
	ctx context.Context,
	id uuid.UUID,
	status domain.PaymentStatus,
) error {

	return s.repository.UpdateStatus(ctx, id, status)
}
