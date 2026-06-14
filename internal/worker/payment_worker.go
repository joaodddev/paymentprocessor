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
}

func NewPaymentWorker(
	consumer *kafkapkg.Consumer,
	repository repository.PaymentRepository,
) *PaymentWorker {
	return &PaymentWorker{
		consumer:   consumer,
		repository: repository,
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

	// Simula tempo de processamento
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
	return nil
}

func resolveStatus(method string) domain.PaymentStatus {
	switch method {
	case "PIX":
		// PIX: aprovação imediata
		return domain.PaymentApproved
	case "CREDIT_CARD":
		// Cartão: simula aprovação (em produção chamaria gateway)
		return domain.PaymentApproved
	case "BOLETO":
		// Boleto: fica pendente até pagamento confirmado
		return domain.PaymentPending
	default:
		return domain.PaymentRejected
	}
}
