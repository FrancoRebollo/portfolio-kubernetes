package in

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FrancoRebollo/auth-security-svc/internal/adapters/rabbitmq"
	"github.com/FrancoRebollo/auth-security-svc/internal/application"
	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
)

type UserEventConsumer struct {
	service *application.SecurityService
	rabbit  *rabbitmq.RabbitMQAdapter
}

func NewUserEventConsumer(service *application.SecurityService, rabbit *rabbitmq.RabbitMQAdapter) *UserEventConsumer {
	return &UserEventConsumer{
		service: service,
		rabbit:  rabbit,
	}
}

// Start escucha la cola y enruta los eventos según el RoutingKey
func (c *UserEventConsumer) Start(ctx context.Context, queue string) {
	handler := func(evt domain.Event) {
		fmt.Printf("📩 Received event: %s | RoutingKey: %s\n", evt.ID, evt.RoutingKey)

		switch evt.RoutingKey {
		case "user.created":
			c.handleUserCreated(ctx, evt)
		//case "user.deleted":
		//	c.handleUserDeleted(ctx, evt)
		default:
			fmt.Printf("⚠️ Unknown routing key: %s (ignored)\n", evt.RoutingKey)
		}
	}
	fmt.Println("STARTING CONSUMER")
	if err := c.rabbit.Consume(ctx, queue, handler); err != nil {
		fmt.Printf("❌ Error starting consumer: %v\n", err)
	}
}

// 🧩 Handler para user.created
func (c *UserEventConsumer) handleUserCreated(ctx context.Context, evt domain.Event) {
	var payload domain.UserCreated
	data, _ := json.Marshal(evt.Payload)
	if err := json.Unmarshal(data, &payload); err != nil {
		fmt.Printf("⚠️ Invalid payload for user.created: %v\n", err)
		return
	}
	fmt.Println("calling from handleUserCreated")
	fmt.Printf("DEBUG Payload type: %T\n", evt.Payload)
	fmt.Printf("DEBUG Payload value: %+v\n", evt.Payload)

	fmt.Printf("DEBUG payload value: %+v\n", payload)
	if _, err := c.service.CreateUserAPI(ctx, payload); err != nil {
		fmt.Printf("❌ Error creating user: %v\n", err)
		return
	}

	fmt.Printf("✅ User created successfully: %s\n", payload.LoginName)
}
