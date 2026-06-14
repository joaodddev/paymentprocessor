package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/joaodddev/paymentprocessor/internal/domain"
	kafkapkg "github.com/joaodddev/paymentprocessor/internal/kafka"
	"github.com/joaodddev/paymentprocessor/internal/repository"
)

const idempotencyTTL = 24 * time.Hour

type PaymentService struct {
	repository repository.PaymentRepository
	redis      *redis.Client
	producer   *kafkapkg.Producer
}

func NewPaymentService(
	repository repository.PaymentRepository,
	redis *redis.Client,
	producer *kafkapkg.Producer,
) *PaymentService {
	return &PaymentService{
		repository: repository,
		redis:      redis,
		producer:   producer,
	}
}

func (s *PaymentService) CreatePayment(
	ctx context.Context,
	amount int64,
	currency string,
	customerID uuid.UUID,
	method domain.PaymentMethod,
	idempotencyKey string,
) (*domain.Payment, error) {

	// Idempotência via Redis
	if idempotencyKey != "" {
		redisKey := fmt.Sprintf("idempotency:%s", idempotencyKey)

		cached, err := s.redis.Get(ctx, redisKey).Result()
		if err == nil {
			var payment domain.Payment
			if err := json.Unmarshal([]byte(cached), &payment); err == nil {
				return &payment, nil
			}
		}
	}

	switch method {
	case domain.PaymentMethodPIX, domain.PaymentMethodCreditCard, domain.PaymentMethodBoleto:
	default:
		return nil, fmt.Errorf("método de pagamento inválido: %s", method)
	}

	now := time.Now()
	payment := &domain.Payment{
		ID:             uuid.New(),
		Amount:         amount,
		Currency:       currency,
		CustomerID:     customerID,
		Method:         method,
		Status:         domain.PaymentPending,
		IdempotencyKey: idempotencyKey,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repository.Create(ctx, payment); err != nil {
		return nil, err
	}

	// Salva no Redis para idempotência
	if idempotencyKey != "" {
		redisKey := fmt.Sprintf("idempotency:%s", idempotencyKey)
		data, _ := json.Marshal(payment)
		s.redis.Set(ctx, redisKey, data, idempotencyTTL)
	}

	// Publica evento no Kafka
	event := kafkapkg.PaymentCreatedEvent{
		PaymentID:      payment.ID.String(),
		Amount:         payment.Amount,
		Currency:       payment.Currency,
		CustomerID:     payment.CustomerID.String(),
		Method:         string(payment.Method),
		IdempotencyKey: payment.IdempotencyKey,
		CreatedAt:      payment.CreatedAt,
	}

	if err := s.producer.PublishPaymentCreated(ctx, event); err != nil {
		fmt.Printf("⚠️  erro ao publicar evento Kafka (pagamento criado): %v\n", err)
	}

	return payment, nil
}

func (s *PaymentService) GetPaymentByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *PaymentService) ListPayments(ctx context.Context) ([]domain.Payment, error) {
	return s.repository.List(ctx)
}

func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus) error {
	return s.repository.UpdateStatus(ctx, id, status)
}
