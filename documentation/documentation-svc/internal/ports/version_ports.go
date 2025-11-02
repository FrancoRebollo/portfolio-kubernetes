package ports

import (
	"context"

	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
)

type VersionService interface {
	GetVersionAPI(ctx context.Context) (*domain.Version, error)
}

type VersionRepository interface {
}
