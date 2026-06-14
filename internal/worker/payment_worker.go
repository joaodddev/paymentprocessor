package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/joaodddev/paymentprocessor/internal/domain"
	kafkapkg "github.com/joaodddev/paymentprocessor/internal/kafka"
	"github.com/joaodddev/paymentprocessor/internal/repository"
)

type PaymentWorker struct {
	consumer   *kafkapkg.Consumer
	repository repository.PaymentRepository
	producer   *kafkapkg.Producer
}

func NewPaymentWorker(
	consumer *kafkapkg.Consumer,
	repository repository.PaymentRepository,
	producer *kafkapkg.Producer,
) *PaymentWorker {
	return &PaymentWorker{
		consumer:   consumer,
		repository: repository,
		producer:   producer,
	}
}

func (w *PaymentWorker) Run(ctx context.Context) {
	log.Println("🔄 Worker iniciado, aguardando eventos payment.created...")

	for {
		select {
		case <-ctx.Done():
			log.Println("worker encerrado")
			return
		default:
			msg, err := w.consumer.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("erro ao consumir mensagem: %v", err)
				continue
			}

			if err := w.process(ctx, msg.Value); err != nil {
				log.Printf("erro ao processar pagamento: %v", err)
			} else {
				w.consumer.Commit(ctx, msg)
			}
		}
	}
}

func (w *PaymentWorker) process(ctx context.Context, data []byte) error {
	var event kafkapkg.PaymentCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	log.Printf("📦 Processando pagamento %s | método: %s | valor: %d %s",
		event.PaymentID, event.Method, event.Amount, event.Currency)

	time.Sleep(500 * time.Millisecond)

	status := resolveStatus(event.Method)

	id, err := uuid.Parse(event.PaymentID)
	if err != nil {
		return err
	}

	if err := w.repository.UpdateStatus(ctx, id, status); err != nil {
		return err
	}

	log.Printf("✅ Pagamento %s atualizado para %s", event.PaymentID, status)

	// Publica evento payment.processed para o webhook notifier
	processed := kafkapkg.PaymentProcessedEvent{
		PaymentID:   event.PaymentID,
		Status:      string(status),
		Method:      event.Method,
		Amount:      event.Amount,
		Currency:    event.Currency,
		CustomerID:  event.CustomerID,
		WebhookURL:  event.WebhookURL,
		ProcessedAt: time.Now(),
	}

	if err := w.producer.PublishPaymentProcessed(ctx, processed); err != nil {
		log.Printf("⚠️  erro ao publicar payment.processed: %v", err)
	}

	return nil
}

func resolveStatus(method string) domain.PaymentStatus {
	switch method {
	case "PIX":
		return domain.PaymentApproved
	case "CREDIT_CARD":
		return domain.PaymentApproved
	case "BOLETO":
		return domain.PaymentPending
	default:
		return domain.PaymentRejected
	}
}
