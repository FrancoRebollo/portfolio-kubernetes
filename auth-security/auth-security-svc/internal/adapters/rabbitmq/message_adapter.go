package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
	"github.com/streadway/amqp"
)

type RabbitMQAdapter struct {
	conn  *amqp.Connection
	pubCh *amqp.Channel
	conCh *amqp.Channel
	// topology
	exchange string // e.g. "app_events"
}

func NewRabbitMQAdapter(amqpURL, exchange string) (*RabbitMQAdapter, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	pubCh, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open publish channel: %w", err)
	}
	conCh, err := conn.Channel()
	if err != nil {
		pubCh.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to open consume channel: %w", err)
	}

	return &RabbitMQAdapter{
		conn:     conn,
		pubCh:    pubCh,
		conCh:    conCh,
		exchange: "app_events",
	}, nil
}

// Publish with routing key (e.g., "config.updated", "user.created")
func (r *RabbitMQAdapter) Publish(ctx context.Context, event domain.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return r.pubCh.Publish(
		r.exchange, event.RoutingKey,
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
}

// Consume from a named queue that has already been bound to routing keys
func (r *RabbitMQAdapter) Consume(ctx context.Context, queue string, handler func(domain.Event)) error {
	// (Optional) QoS to avoid overwhelming consumers
	if err := r.conCh.Qos(50, 0, false); err != nil {
		return fmt.Errorf("qos: %w", err)
	}
	msgs, err := r.conCh.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}

				var event domain.Event
				if err := json.Unmarshal(d.Body, &event); err != nil {
					fmt.Printf("⚠️ Error unmarshalling message: %v\n", err)
					_ = d.Nack(false, false) // opcional: no reencolar
					continue
				}

				// 👉 Ahora sí, pasamos el evento de dominio al handler
				handler(event)

				// (Opcional) confirmar procesamiento exitoso
				_ = d.Ack(false)

			}
		}
	}()
	return nil
}

func (r *RabbitMQAdapter) Ack(d amqp.Delivery)                { _ = d.Ack(false) }
func (r *RabbitMQAdapter) Nack(d amqp.Delivery, requeue bool) { _ = d.Nack(false, requeue) }

func (r *RabbitMQAdapter) Close() {
	if r.pubCh != nil {
		_ = r.pubCh.Close()
	}
	if r.conCh != nil {
		_ = r.conCh.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
}
