package ports

import (
	"context"

	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
)

type VersionService interface {
	GetVersionAPI(ctx context.Context) (*domain.Version, error)
}

type VersionRepository interface {
}
