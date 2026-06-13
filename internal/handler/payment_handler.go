package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/joaodddev/paymentprocessor/internal/domain"
	"github.com/joaodddev/paymentprocessor/internal/dto"
	"github.com/joaodddev/paymentprocessor/internal/service"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id inválido"})
		return
	}

	payment, err := h.service.CreatePayment(
		c.Request.Context(),
		req.Amount,
		req.Currency,
		customerID,
		domain.PaymentMethod(req.Method),
		req.IdempotencyKey,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toResponse(payment))
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	payment, err := h.service.GetPaymentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pagamento não encontrado"})
		return
	}

	c.JSON(http.StatusOK, toResponse(payment))
}

func (h *PaymentHandler) ListPayments(c *gin.Context) {
	payments, err := h.service.ListPayments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []dto.PaymentResponse
	for _, p := range payments {
		responses = append(responses, toResponse(&p))
	}

	c.JSON(http.StatusOK, responses)
}

func toResponse(p *domain.Payment) dto.PaymentResponse {
	return dto.PaymentResponse{
		ID:             p.ID.String(),
		Amount:         p.Amount,
		Currency:       p.Currency,
		CustomerID:     p.CustomerID.String(),
		Method:         string(p.Method),
		Status:         string(p.Status),
		IdempotencyKey: p.IdempotencyKey,
		CreatedAt:      p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
