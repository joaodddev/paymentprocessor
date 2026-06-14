package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/joaodddev/paymentprocessor/internal/domain"
	kafkapkg "github.com/joaodddev/paymentprocessor/internal/kafka"
	"github.com/joaodddev/paymentprocessor/internal/repository"
)

type WebhookNotifier struct {
	consumer          *kafkapkg.Consumer
	webhookRepository repository.WebhookRepository
	httpClient        *http.Client
}

func NewWebhookNotifier(
	consumer *kafkapkg.Consumer,
	webhookRepository repository.WebhookRepository,
) *WebhookNotifier {
	return &WebhookNotifier{
		consumer:          consumer,
		webhookRepository: webhookRepository,
		httpClient:        &http.Client{Timeout: 10 * time.Second},
	}
}

func (n *WebhookNotifier) Run(ctx context.Context) {
	log.Println("🔔 Webhook Notifier iniciado, aguardando eventos payment.processed...")

	for {
		select {
		case <-ctx.Done():
			log.Println("webhook notifier encerrado")
			return
		default:
			msg, err := n.consumer.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("erro ao consumir mensagem: %v", err)
				continue
			}

			if err := n.notify(ctx, msg.Value); err != nil {
				log.Printf("erro ao enviar webhook: %v", err)
			} else {
				n.consumer.Commit(ctx, msg)
			}
		}
	}
}

func (n *WebhookNotifier) notify(ctx context.Context, data []byte) error {
	var event kafkapkg.PaymentProcessedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	if event.WebhookURL == "" {
		log.Printf("pagamento %s sem webhook_url, ignorando", event.PaymentID)
		return nil
	}

	paymentID, _ := uuid.Parse(event.PaymentID)
	payload, _ := json.Marshal(event)

	now := time.Now()
	webhook := &domain.Webhook{
		ID:          uuid.New(),
		PaymentID:   paymentID,
		URL:         event.WebhookURL,
		Status:      domain.WebhookPending,
		Attempt:     0,
		MaxAttempts: 5,
		Payload:     payload,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := n.webhookRepository.Create(ctx, webhook); err != nil {
		return fmt.Errorf("erro ao salvar webhook: %w", err)
	}

	// Retry com backoff exponencial
	for webhook.Attempt < webhook.MaxAttempts {
		webhook.Attempt++

		if webhook.Attempt > 1 {
			backoff := time.Duration(math.Pow(2, float64(webhook.Attempt-1))) * time.Second
			log.Printf("⏳ aguardando %v antes da tentativa %d...", backoff, webhook.Attempt)
			time.Sleep(backoff)
		}

		statusCode, err := n.dispatch(event.WebhookURL, payload)

		if err != nil {
			webhook.LastError = err.Error()
			log.Printf("⚠️  tentativa %d falhou: %v", webhook.Attempt, err)
		} else {
			webhook.ResponseCode = &statusCode

			if statusCode >= 200 && statusCode < 300 {
				delivered := time.Now()
				webhook.Status = domain.WebhookDelivered
				webhook.DeliveredAt = &delivered
				n.webhookRepository.UpdateStatus(ctx, webhook)
				log.Printf("✅ Webhook entregue para %s (tentativa %d)", event.WebhookURL, webhook.Attempt)
				return nil
			}

			webhook.LastError = fmt.Sprintf("status HTTP inesperado: %d", statusCode)
			log.Printf("⚠️  tentativa %d retornou status %d", webhook.Attempt, statusCode)
		}

		n.webhookRepository.UpdateStatus(ctx, webhook)
	}

	webhook.Status = domain.WebhookFailed
	n.webhookRepository.UpdateStatus(ctx, webhook)
	log.Printf("❌ Webhook falhou após %d tentativas para %s", webhook.MaxAttempts, event.WebhookURL)

	return nil
}

func (n *WebhookNotifier) dispatch(url string, payload []byte) (int, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
