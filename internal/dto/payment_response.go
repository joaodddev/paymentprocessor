package dto

type PaymentResponse struct {
	ID             string `json:"id"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	CustomerID     string `json:"customer_id"`
	Method         string `json:"method"`
	Status         string `json:"status"`
	IdempotencyKey string `json:"idempotency_key"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
