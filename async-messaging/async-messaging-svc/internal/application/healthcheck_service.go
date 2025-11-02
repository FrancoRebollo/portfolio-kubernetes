package application

import (
	"context"
	"errors"

	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/config"
	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/logger"

	"github.com/FrancoRebollo/async-messaging-svc/internal/domain"
	"github.com/FrancoRebollo/async-messaging-svc/internal/ports"
)

type HealthcheckService struct {
	hr   ports.HealthcheckRepository
	conf config.App
}

func NewHealthcheckService(hr ports.HealthcheckRepository, conf config.App) *HealthcheckService {
	return &HealthcheckService{
		hr,
		conf,
	}
}

func (hs *HealthcheckService) GetHealthcheck(ctx context.Context) (*domain.Healthcheck, error) {
	var serviceErr error
	listDBPing, err := hs.hr.GetDatabasesPing(ctx)
	if err != nil {
		serviceErr = mapServiceError(err)
		logger.LoggerError().WithError(err).Error(serviceErr)
		return &domain.Healthcheck{}, serviceErr
	}

	return &domain.Healthcheck{
		NombreApi:     hs.conf.Name,
		Cliente:       hs.conf.Client,
		Version:       hs.conf.Version,
		VersionModelo: "",
		FechaStartUp:  hs.conf.FechaStartUp,
		BasesDeDatos:  listDBPing,
	}, nil
}

func mapServiceError(err error) error {
	if errors.Is(err, domain.ErrDeadlockDetected) {
		return &domain.HealthcheckError{Code: domain.ErrCodeDeadlockDetected, Message: "Error generado por un lockeo"}
	} else if errors.Is(err, domain.ErrTableOrViewDoesNotExist) {
		return &domain.HealthcheckError{Code: domain.ErrCodeTableOrViewDoesNotExist, Message: "Error generado por tabla o vista inexistente"}
	} else if errors.Is(err, domain.ErrEndOfCommunication) {
		return &domain.HealthcheckError{Code: domain.ErrCodeEndOfCommunication, Message: "Error generado por corte en la comunicación con la DB"}
	} else if errors.Is(err, domain.ErrConnectionTimeout) {
		return &domain.HealthcheckError{Code: domain.ErrCodeConnectionTimeout, Message: "Error generado por un timeout de conexión"}
	}
	return &domain.HealthcheckError{Code: domain.ErrCodeInternalServer, Message: "Error genérico interno"}
}
