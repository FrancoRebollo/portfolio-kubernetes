package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
	"github.com/FrancoRebollo/api-integration-svc/internal/platform/logger"
)

type ApiIntegrationRepository struct {
	dbPost *PostgresDB
}

func NewApiIntegrationRepository(dbPost *PostgresDB) *ApiIntegrationRepository {
	return &ApiIntegrationRepository{
		dbPost: dbPost,
	}
}

func (hr *ApiIntegrationRepository) GetDatabasesPing(ctx context.Context) ([]domain.Database, error) {
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

func (hr *ApiIntegrationRepository) PushEventToQueue(ctx context.Context, tx *sql.Tx, event domain.Event) error {
	query := `
		INSERT INTO api_int.message_event (
			id_event,
			source_system,
			destiny_system,
			payload,
			status,
			fecha_recepcion,
			actualizado_por
		)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6)
		ON CONFLICT (id_event, source_system)
		DO NOTHING;
	`

	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}

	res, err := tx.ExecContext(ctx, query,
		event.ID,
		event.Origin,     // → source_system
		event.RoutingKey, // → queue_name
		payloadJSON,
		"RECEIVED",
		"SYSTEM",
	)
	if err != nil {
		return fmt.Errorf("error inserting event: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrDuplicateEvent
	}

	return nil
}

func (hr *ApiIntegrationRepository) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := hr.dbPost.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
