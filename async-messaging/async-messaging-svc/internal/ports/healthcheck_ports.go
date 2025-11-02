package ports

import (
	"context"

	"github.com/FrancoRebollo/async-messaging-svc/internal/domain"
)

type HealthcheckService interface {
	GetHealthcheck(ctx context.Context) (*domain.Healthcheck, error)
}

type HealthcheckRepository interface {
	GetDatabasesPing(ctx context.Context) ([]domain.Database, error)
}
