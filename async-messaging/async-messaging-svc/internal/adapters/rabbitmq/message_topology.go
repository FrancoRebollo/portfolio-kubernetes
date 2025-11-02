package rabbitmq

import "fmt"

func (r *RabbitMQAdapter) InitializeTopology() error {
	r.exchange = "app_events" // ✅ setea el nombre del exchange

	if err := r.pubCh.ExchangeDeclare(
		r.exchange, "topic", true, false, false, false, nil,
	); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	// Declarar colas
	r.pubCh.QueueDeclare("config_updates_q", true, false, false, false, nil)
	r.pubCh.QueueDeclare("user_created_q", true, false, false, false, nil)

	// Vincular colas
	r.pubCh.QueueBind("config_updates_q", "config.updated", r.exchange, false, nil)
	r.pubCh.QueueBind("user_created_q", "user.created", r.exchange, false, nil)

	fmt.Println("✅ Topología creada: 1 exchange, 2 colas, 2 bindings")
	return nil
}
