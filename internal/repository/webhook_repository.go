package repository

import (
	"context"

	"github.com/joaodddev/paymentprocessor/internal/domain"
)

type WebhookRepository interface {
	Create(ctx context.Context, webhook *domain.Webhook) error
	UpdateStatus(ctx context.Context, webhook *domain.Webhook) error
}
