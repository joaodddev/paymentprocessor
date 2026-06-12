package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joaodddev/paymentprocessor/internal/domain"
)

type PostgresPaymentRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPaymentRepository(db *pgxpool.Pool) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{
		db: db,
	}
}

func (r *PostgresPaymentRepository) Create(
	ctx context.Context,
	payment *domain.Payment,
) error {

	query := `
	INSERT INTO payments (
		id,
		amount,
		currency,
		customer_id,
		status,
		idempotency_key,
		created_at,
		updated_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		payment.ID,
		payment.Amount,
		payment.Currency,
		payment.CustomerID,
		payment.Status,
		payment.IdempotencyKey,
		payment.CreatedAt,
		payment.UpdatedAt,
	)

	return err
}

func (r *PostgresPaymentRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Payment, error) {

	query := `
	SELECT
		id,
		amount,
		currency,
		customer_id,
		status,
		idempotency_key,
		created_at,
		updated_at
	FROM payments
	WHERE id = $1
	`

	var payment domain.Payment

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&payment.ID,
		&payment.Amount,
		&payment.Currency,
		&payment.CustomerID,
		&payment.Status,
		&payment.IdempotencyKey,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *PostgresPaymentRepository) List(
	ctx context.Context,
) ([]domain.Payment, error) {

	query := `
	SELECT
		id,
		amount,
		currency,
		customer_id,
		status,
		idempotency_key,
		created_at,
		updated_at
	FROM payments
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var payments []domain.Payment

	for rows.Next() {

		var payment domain.Payment

		err := rows.Scan(
			&payment.ID,
			&payment.Amount,
			&payment.Currency,
			&payment.CustomerID,
			&payment.Status,
			&payment.IdempotencyKey,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PostgresPaymentRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status domain.PaymentStatus,
) error {

	query := `
	UPDATE payments
	SET
		status = $1,
		updated_at = NOW()
	WHERE id = $2
	`

	_, err := r.db.Exec(
		ctx,
		query,
		status,
		id,
	)

	return err
}
