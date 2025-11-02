package ports

import (
	"context"

	"github.com/FrancoRebollo/async-messaging-svc/internal/domain"
)

type MessageService interface {
}

type MessageRepository interface {
	GetDatabasesPing(ctx context.Context) ([]domain.Database, error)
}

type MessageQueue interface {
	Publish(ctx context.Context, event domain.Event) error
	Consume(ctx context.Context, queue string, handler func(domain.Event)) error
}
