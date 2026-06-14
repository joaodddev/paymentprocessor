package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaodddev/paymentprocessor/internal/domain"
)

type PostgresWebhookRepository struct {
	db *pgxpool.Pool
}

func NewPostgresWebhookRepository(db *pgxpool.Pool) *PostgresWebhookRepository {
	return &PostgresWebhookRepository{db: db}
}

func (r *PostgresWebhookRepository) Create(ctx context.Context, webhook *domain.Webhook) error {
	query := `
	INSERT INTO webhooks (id, payment_id, url, status, attempt, max_attempts, payload, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		webhook.ID,
		webhook.PaymentID,
		webhook.URL,
		webhook.Status,
		webhook.Attempt,
		webhook.MaxAttempts,
		webhook.Payload,
		webhook.CreatedAt,
		webhook.UpdatedAt,
	)
	return err
}

func (r *PostgresWebhookRepository) UpdateStatus(ctx context.Context, webhook *domain.Webhook) error {
	query := `
	UPDATE webhooks
	SET status = $1, attempt = $2, response_code = $3, last_error = $4, delivered_at = $5, updated_at = $6
	WHERE id = $7
	`
	_, err := r.db.Exec(ctx, query,
		webhook.Status,
		webhook.Attempt,
		webhook.ResponseCode,
		webhook.LastError,
		webhook.DeliveredAt,
		time.Now(),
		webhook.ID,
	)
	return err
}
