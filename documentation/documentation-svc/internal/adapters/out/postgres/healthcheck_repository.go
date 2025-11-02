package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
	"github.com/FrancoRebollo/api-integration-svc/internal/platform/logger"
)

type HealthcheckRepository struct {
	dbPost *PostgresDB
}

func NewHealthcheckRepository(dbPost *PostgresDB) *HealthcheckRepository {
	return &HealthcheckRepository{
		dbPost: dbPost,
	}
}

func (hr *HealthcheckRepository) GetDatabasesPing(ctx context.Context) ([]domain.Database, error) {
	databases := []domain.Database{}
	var fechaUltimaActividad string
	var mappedErr error
	var repoErr error

	query := `SELECT NOW()`

	rows, err := hr.dbPost.GetDB().QueryContext(ctx, query)
	if err != nil {
		mappedErr = hr.dbPost.MapPostgresError(err)
		repoErr = getRepoErr(mappedErr)
		logger.LoggerError().WithError(err).Error(repoErr)
		return databases, repoErr
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&fechaUltimaActividad)
		if err != nil {
			mappedErr = hr.dbPost.MapPostgresError(err)
			repoErr = getRepoErr(mappedErr)
			logger.LoggerError().WithError(err).Error(repoErr)
			return databases, repoErr
		}
	}

	if err = rows.Err(); err != nil {
		mappedErr = hr.dbPost.MapPostgresError(err)
		repoErr = getRepoErr(mappedErr)
		logger.LoggerError().WithError(err).Error(repoErr)
		return databases, repoErr
	}

	databases = append(databases, domain.Database{
		Base:                     "POSTGRES",
		FechaHoraUltimaActividad: fechaUltimaActividad,
	})

	return databases, nil
}

func getRepoErr(err error) error {
	if errors.Is(err, domain.ErrDeadlockDetected) {
		return fmt.Errorf("%w: error de healthcheck repository: %s", domain.ErrDeadlockDetected, err.Error())
	} else if errors.Is(err, domain.ErrTableOrViewDoesNotExist) {
		return fmt.Errorf("%w: error de healthcheck repository: %s", domain.ErrTableOrViewDoesNotExist, err.Error())
	} else if errors.Is(err, domain.ErrEndOfCommunication) {
		return fmt.Errorf("%w: error de healthcheck repository: %s", domain.ErrEndOfCommunication, err.Error())
	} else if errors.Is(err, domain.ErrConnectionTimeout) {
		return fmt.Errorf("%w: error de healthcheck repository: %s", domain.ErrConnectionTimeout, err.Error())
	}
	return fmt.Errorf("%w: error de healthcheck repository: %s", domain.ErrInternalServer, err.Error())
}
