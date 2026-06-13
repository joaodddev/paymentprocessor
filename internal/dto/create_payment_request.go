package dto

type CreatePaymentRequest struct {
	Amount         int64  `json:"amount" binding:"required,min=1"`
	Currency       string `json:"currency" binding:"required,len=3"`
	CustomerID     string `json:"customer_id" binding:"required,uuid"`
	Method         string `json:"method" binding:"required,oneof=PIX CREDIT_CARD BOLETO"`
	IdempotencyKey string `json:"idempotency_key"`
}
