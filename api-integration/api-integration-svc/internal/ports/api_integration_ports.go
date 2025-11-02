package ports

import (
	"context"
	"database/sql"

	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
)

type ApiIntegrationService interface {
	ForwardRequest(req domain.ExternalAPIRequest) (domain.ExternalAPIResponse, error)
	PushEventToQueueAPI(ctx context.Context, reqEvent domain.Event) error
}

type ApiIntegrationRepository interface {
	PushEventToQueue(ctx context.Context, tx *sql.Tx, event domain.Event) error
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error
}

type MessageQueue interface {
	Publish(ctx context.Context, event domain.Event) error
}
