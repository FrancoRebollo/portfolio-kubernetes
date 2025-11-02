package application

import (
	"context"

	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/config"

	"github.com/FrancoRebollo/async-messaging-svc/internal/domain"
	"github.com/FrancoRebollo/async-messaging-svc/internal/ports"
)

type VersionService struct {
	ports.VersionRepository
	config.App
}

func NewVersionService(repo ports.VersionRepository, app config.App) *VersionService {
	return &VersionService{
		VersionRepository: repo,
		App:               app,
	}
}

func (s *VersionService) GetVersionAPI(ctx context.Context) (*domain.Version, error) {

	newVersion := domain.Version{
		NombreApi:    s.App.Name,
		Cliente:      s.App.Client,
		Version:      s.App.Version,
		FechaStartUp: s.App.FechaStartUp,
	}

	return &newVersion, nil
}
