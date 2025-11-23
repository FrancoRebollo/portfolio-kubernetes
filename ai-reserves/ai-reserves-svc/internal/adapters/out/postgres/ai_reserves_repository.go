package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/FrancoRebollo/ai-reserves-svc/internal/domain"
	"github.com/FrancoRebollo/ai-reserves-svc/internal/platform/logger"
)

type AiReservesRepository struct {
	dbPost *PostgresDB
}

func NewAiReservesRepository(dbPost *PostgresDB) *AiReservesRepository {
	return &AiReservesRepository{
		dbPost: dbPost,
	}
}

func (hr *AiReservesRepository) GetDatabasesPing(ctx context.Context) ([]domain.Database, error) {
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

func (hr *AiReservesRepository) CreatePersona(ctx context.Context, req domain.PersonCreatedPayload) error {

	// Primero verificamos si la persona existe
	var exists bool

	err := hr.dbPost.GetDB().QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM ai_res.personas WHERE id = $1)`,
		req.ID,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("checking persona existence: %w", err)
	}

	if !exists {
		// INSERT si la persona NO existe
		_, err = hr.dbPost.GetDB().ExecContext(ctx,
			`INSERT INTO ai_res.personas(id, email, telefono,created_at,created_by,updated_at)
			 VALUES ($1, $2, $3,CURRENT_TIMESTAMP,'auth_security',null)`,
			req.ID,
			req.Email,
			req.TePersona,
		)
		if err != nil {
			return fmt.Errorf("insert persona: %w", err)
		}

		fmt.Printf("👤 Persona creada ID=%d\n", req.ID)
		return nil
	}

	// UPDATE si ya existe
	_, err = hr.dbPost.GetDB().ExecContext(ctx,
		`UPDATE ai_res.personas
		    SET email = $2,
		        telefono = $3,
		        updated_at = CURRENT_TIMESTAMP,
				updated_by = 'auth_security'
		  WHERE id = $1`,
		req.ID,
		req.Email,
		req.TePersona,
	)
	if err != nil {
		return fmt.Errorf("update persona: %w", err)
	}

	fmt.Printf("🔄 Persona actualizada ID=%d\n", req.ID)
	return nil
}

func (hr *AiReservesRepository) UpdAtributoPersona(ctx context.Context, req domain.PersonaParcial) error {

	// 1) Validar que la persona exista
	var exists bool
	err := hr.dbPost.GetDB().QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM ai_res.personas WHERE id = $1)`,
		req.ID,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("checking persona existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("persona with ID=%d does not exist", req.ID)
	}

	// 2) Lista blanca de columnas permitidas (evita SQL injection)
	validColumns := map[string]bool{
		"email":                 true,
		"telefono":              true,
		"nombre":                true,
		"apellido_razon_social": true,
		"tipo_doc_persona":      true,
		"nro_doc_persona":       true,
		"persona_juridica":      true,
	}

	if !validColumns[req.Atribute] {
		return fmt.Errorf("invalid attribute '%s' for update", req.Atribute)
	}

	// 3) Armar el SQL dinámico de forma segura
	query := fmt.Sprintf(`
        UPDATE ai_res.personas
           SET %s = $1,
               updated_at = CURRENT_TIMESTAMP,
               updated_by = 'auth_security'
         WHERE id = $2`, req.Atribute)

	// 4) Ejecutar el update
	_, err = hr.dbPost.GetDB().ExecContext(ctx, query, req.Value, req.ID)
	if err != nil {
		return fmt.Errorf("updating persona %s: %w", req.Atribute, err)
	}

	fmt.Printf("📝 Persona %d actualizada: %s = %v\n", req.ID, req.Atribute, req.Value)
	return nil
}

func (hr *AiReservesRepository) UpdPersona(ctx context.Context, req domain.Persona) error {

	// 1) Validar existencia de la persona
	var exists bool
	err := hr.dbPost.GetDB().QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM ai_res.personas WHERE id = $1)`,
		req.ID,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("checking persona existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("persona with ID=%d does not exist", req.ID)
	}

	// 2) Mapa de columnas válidas (protección anti SQL injection)
	validFields := map[string]interface{}{
		"nombre":                req.Nombre,
		"apellido_razon_social": req.ApellidoRazonSocial,
		"persona_juridica":      req.PersonaJuridia,
		"tipo_doc_persona":      req.TipoDocPersona,
		"nro_doc_persona":       req.NroDocPersona,
		"email":                 req.Email,
		"telefono":              req.TelPersona,
	}

	// 3) Construcción dinámica del UPDATE
	setClauses := []string{}
	values := []interface{}{}
	paramIndex := 1

	for col, val := range validFields {
		// actualizar solo si el campo viene con valor
		if v, ok := val.(string); ok && v != "" {
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, paramIndex))
			values = append(values, v)
			paramIndex++
		}
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields provided for update")
	}

	// agregamos metadata
	setClauses = append(setClauses,
		fmt.Sprintf("updated_at = CURRENT_TIMESTAMP"),
	)
	// updated_by también
	setClauses = append(setClauses,
		fmt.Sprintf("updated_by = 'ai_reserves'"),
	)

	// 4) Armamos el query final
	query := fmt.Sprintf(`
        UPDATE ai_res.personas
           SET %s
         WHERE id = $%d`,
		strings.Join(setClauses, ", "),
		paramIndex,
	)

	values = append(values, req.ID)

	// 5) Ejecutar UPDATE
	_, err = hr.dbPost.GetDB().ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("updating persona: %w", err)
	}

	fmt.Printf("📝 Persona %d actualizada con éxito\n", req.ID)
	return nil
}

func (hr *AiReservesRepository) UpsertConfigPersona(ctx context.Context, req domain.ConfigPersona) error {

	// --- 1) Validar que exista la persona ---
	var personaExists bool
	err := hr.dbPost.GetDB().QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM ai_res.personas WHERE id = $1)`,
		req.IDPersona,
	).Scan(&personaExists)
	if err != nil {
		return fmt.Errorf("checking persona existence: %w", err)
	}
	if !personaExists {
		return fmt.Errorf("persona id=%d does not exist", req.IDPersona)
	}

	// --- 2) Lista blanca completa de columnas actualizables ---
	validCols := map[string]string{
		"notificar_por_mail":     "boolean",
		"notificar_por_sms":      "boolean",
		"dias_visibles_adelante": "int",
		"id_agenda":              "int",
	}

	// Determinar si la columna es válida
	colType, ok := validCols[req.Atribute]
	if !ok {
		return fmt.Errorf("attribute '%s' cannot be updated", req.Atribute)
	}

	// --- 3) Convertir valor según tipo ---
	var castValue interface{}

	switch colType {
	case "boolean":
		if req.Value == "true" || req.Value == "1" {
			castValue = true
		} else if req.Value == "false" || req.Value == "0" {
			castValue = false
		} else {
			return fmt.Errorf("invalid boolean value '%s'", req.Value)
		}

	case "int":
		n, errConv := strconv.Atoi(req.Value)
		if errConv != nil {
			return fmt.Errorf("invalid integer value '%s'", req.Value)
		}
		castValue = n

	default:
		return fmt.Errorf("unsupported type '%s' for column '%s'", colType, req.Atribute)
	}

	// --- 4) Ver si existe config para persona ---
	var exists bool
	err = hr.dbPost.GetDB().QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM ai_res.conf_agenda_persona WHERE id_persona = $1)`,
		req.IDPersona,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("checking config existence: %w", err)
	}

	// --- 5) INSERT si no existe ---
	if !exists {

		queryInsert := fmt.Sprintf(`
            INSERT INTO ai_res.conf_agenda_persona
                (id_persona, %s, created_by)
            VALUES ($1, $2, 'auth_security')
        `, req.Atribute)

		_, err = hr.dbPost.GetDB().ExecContext(ctx, queryInsert, req.IDPersona, castValue)
		if err != nil {
			return fmt.Errorf("insert config: %w", err)
		}

		fmt.Printf("🆕 Config creada para persona %d -> %s = %v\n",
			req.IDPersona, req.Atribute, castValue)
		return nil
	}

	// --- 6) UPDATE dinámico (solo con columnas seguras) ---
	queryUpdate := fmt.Sprintf(`
        UPDATE ai_res.conf_agenda_persona
           SET %s = $1,
               updated_at = CURRENT_TIMESTAMP,
               updated_by = 'auth_security'
         WHERE id_persona = $2
    `, req.Atribute)

	_, err = hr.dbPost.GetDB().ExecContext(ctx, queryUpdate, castValue, req.IDPersona)
	if err != nil {
		return fmt.Errorf("update config: %w", err)
	}

	fmt.Printf("🔄 Config actualizada para persona %d -> %s = %v\n",
		req.IDPersona, req.Atribute, castValue)

	return nil
}

func (hr *AiReservesRepository) CreateUnidadReserva(ctx context.Context, req *domain.UnidadReserva) (int, error) {

	// 1) Verificar si existe por nombre
	var idExistente int
	err := hr.dbPost.GetDB().QueryRowContext(ctx,
		`SELECT id FROM ai_res.unidad_reserva WHERE nombre = $1`,
		req.Nombre,
	).Scan(&idExistente)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("error verificando existencia de unidad_reserva: %w", err)
	}

	// --------------------------------------------------------------------
	// 2) SI YA EXISTE → actualizar descripcion
	// --------------------------------------------------------------------
	if idExistente > 0 {
		_, err := hr.dbPost.GetDB().ExecContext(ctx,
			`UPDATE ai_res.unidad_reserva
			   SET descripcion = $2,
			       updated_at = CURRENT_TIMESTAMP,
			       updated_by = 'ai_reserves'
			 WHERE id = $1`,
			idExistente,
			req.Descripcion,
		)
		if err != nil {
			return 0, fmt.Errorf("error actualizando unidad_reserva: %w", err)
		}

		fmt.Printf("🔄 UnidadReserva actualizada (ID=%d)\n", idExistente)
		return idExistente, nil
	}

	// --------------------------------------------------------------------
	// 3) SI NO EXISTE → Insertar y devolver ID generado
	// --------------------------------------------------------------------
	var newID int
	err = hr.dbPost.GetDB().QueryRowContext(ctx,
		`INSERT INTO ai_res.unidad_reserva
			(nombre, descripcion, created_by)
		 VALUES ($1, $2, 'ai_reserves')
		 RETURNING id`,
		req.Nombre,
		req.Descripcion,
	).Scan(&newID)

	if err != nil {
		return 0, fmt.Errorf("error insertando unidad_reserva: %w", err)
	}

	fmt.Printf("🆕 UnidadReserva creada (ID=%d)\n", newID)
	return newID, nil
}

func (hr *AiReservesRepository) CreateTipoUnidadReserva(ctx context.Context, req domain.TipoUnidadReserva) (int, error) {

	// ---------------------------------------------------------------------
	// 1) Verificar que la UnidadReserva exista
	// ---------------------------------------------------------------------
	var existsUnidad bool

	err := hr.dbPost.GetDB().QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1 
			  FROM ai_res.unidad_reserva 
			 WHERE id = $1
		)`,
		req.IDUnidadReserva,
	).Scan(&existsUnidad)

	if err != nil {
		return 0, fmt.Errorf("error verificando unidad_reserva: %w", err)
	}

	if !existsUnidad {
		return 0, fmt.Errorf("unidad_reserva con ID=%d no existe", req.IDUnidadReserva)
	}

	// ---------------------------------------------------------------------
	// 2) Verificar si el TipoUnidadReserva ya existe dentro de esa unidad
	// ---------------------------------------------------------------------
	var idTipoExistente int

	err = hr.dbPost.GetDB().QueryRowContext(ctx,
		`SELECT id
		   FROM ai_res.tipo_unidad_reserva
		  WHERE id_unidad_reserva = $1
		    AND nombre = $2`,
		req.IDUnidadReserva,
		req.Nombre,
	).Scan(&idTipoExistente)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("error verificando tipo_unidad_reserva: %w", err)
	}

	// ---------------------------------------------------------------------
	// 3) SI YA EXISTE → actualizar
	// ---------------------------------------------------------------------
	if idTipoExistente > 0 {
		_, err := hr.dbPost.GetDB().ExecContext(ctx,
			`UPDATE ai_res.tipo_unidad_reserva
			   SET descripcion = $2,
			       updated_at = CURRENT_TIMESTAMP,
			       updated_by = 'ai_reserves'
			 WHERE id = $1`,
			idTipoExistente,
			req.Descripcion,
		)
		if err != nil {
			return 0, fmt.Errorf("error actualizando tipo_unidad_reserva: %w", err)
		}

		fmt.Printf("🔄 TipoUnidadReserva actualizado (ID=%d)\n", idTipoExistente)
		return idTipoExistente, nil
	}

	// ---------------------------------------------------------------------
	// 4) SI NO EXISTE → insertar nuevo registro
	// ---------------------------------------------------------------------
	var newID int

	err = hr.dbPost.GetDB().QueryRowContext(ctx,
		`INSERT INTO ai_res.tipo_unidad_reserva
		    (id_unidad_reserva, nombre, descripcion, created_by)
		 VALUES ($1, $2, $3, 'ai_reserves')
		 RETURNING id`,
		req.IDUnidadReserva,
		req.Nombre,
		req.Descripcion,
	).Scan(&newID)

	if err != nil {
		return 0, fmt.Errorf("error insertando tipo_unidad_reserva: %w", err)
	}

	fmt.Printf("🆕 TipoUnidadReserva creado (ID=%d)\n", newID)
	return newID, nil
}

func (hr *AiReservesRepository) CreateSubTipoUnidadReserva(
	ctx context.Context,
	req domain.SubTipoUnidadReserva,
) (int, error) {

	// ---------------------------------------------------------------------
	// 1) Verificar que la UnidadReserva exista
	// ---------------------------------------------------------------------
	var existsUnidad bool

	err := hr.dbPost.GetDB().QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1 
			  FROM ai_res.unidad_reserva 
			 WHERE id = $1
		)`,
		req.IDUnidadReserva,
	).Scan(&existsUnidad)

	if err != nil {
		return 0, fmt.Errorf("error verificando unidad_reserva: %w", err)
	}

	if !existsUnidad {
		return 0, fmt.Errorf("unidad_reserva con ID=%d no existe", req.IDUnidadReserva)
	}

	// ---------------------------------------------------------------------
	// 2) Verificar que el TipoUnidad exista y pertenezca a esa Unidad
	// ---------------------------------------------------------------------
	var existsTipo bool

	err = hr.dbPost.GetDB().QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1
			  FROM ai_res.tipo_unidad_reserva
			 WHERE id = $1
			   AND id_unidad_reserva = $2
		)`,
		req.IDTipoUnidadReserva,
		req.IDUnidadReserva,
	).Scan(&existsTipo)

	if err != nil {
		return 0, fmt.Errorf("error verificando tipo_unidad_reserva: %w", err)
	}

	if !existsTipo {
		return 0, fmt.Errorf("tipo_unidad_reserva ID=%d no pertenece a unidad_reserva ID=%d",
			req.IDTipoUnidadReserva, req.IDUnidadReserva)
	}

	// ---------------------------------------------------------------------
	// 3) Verificar si el SubTipoUnidad ya existe
	// ---------------------------------------------------------------------
	var idExistente int

	err = hr.dbPost.GetDB().QueryRowContext(ctx,
		`SELECT id
		   FROM ai_res.sub_tipo_unidad_reserva
		  WHERE id_tipo_unidad_reserva = $1
		    AND nombre = $2`,
		req.IDTipoUnidadReserva,
		req.Nombre,
	).Scan(&idExistente)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("error verificando sub_tipo_unidad_reserva: %w", err)
	}

	// ---------------------------------------------------------------------
	// 4) SI EXISTE → actualizar
	// ---------------------------------------------------------------------
	if idExistente > 0 {

		_, err := hr.dbPost.GetDB().ExecContext(ctx,
			`UPDATE ai_res.sub_tipo_unidad_reserva
			   SET descripcion = $2,
			       updated_at = CURRENT_TIMESTAMP,
			       updated_by = 'ai_reserves'
			 WHERE id = $1`,
			idExistente,
			req.Descripcion,
		)

		if err != nil {
			return 0, fmt.Errorf("error actualizando sub_tipo_unidad_reserva: %w", err)
		}

		fmt.Printf("🔄 SubTipoUnidadReserva actualizado (ID=%d)\n", idExistente)
		return idExistente, nil
	}

	// ---------------------------------------------------------------------
	// 5) SI NO EXISTE → insertar
	// ---------------------------------------------------------------------
	var newID int

	err = hr.dbPost.GetDB().QueryRowContext(ctx,
		`INSERT INTO ai_res.sub_tipo_unidad_reserva
		    (id_tipo_unidad_reserva, nombre, descripcion, duracion_reserva_minutos, created_by)
		 VALUES ($1, $2, $3, 0, 'ai_reserves')
		 RETURNING id`,
		req.IDTipoUnidadReserva,
		req.Nombre,
		req.Descripcion,
	).Scan(&newID)

	if err != nil {
		return 0, fmt.Errorf("error insertando sub_tipo_unidad_reserva: %w", err)
	}

	fmt.Printf("🆕 SubTipoUnidadReserva creado (ID=%d)\n", newID)
	return newID, nil
}

func (hr *AiReservesRepository) ModifUnidadReserva(ctx context.Context, req domain.UnidadReserva) error {
	return nil
}

func (hr *AiReservesRepository) ModifTipoUnidadReserva(ctx context.Context, req domain.TipoUnidadReserva) error {
	return nil
}

func (hr *AiReservesRepository) ModifSubTipoUnidadReserva(ctx context.Context, req domain.SubTipoUnidadReserva) error {
	return nil
}

func (hr *AiReservesRepository) CreateReserve(ctx context.Context, req domain.Reserva) error {
	return nil
}

func (hr *AiReservesRepository) CancelReserve(ctx context.Context, idReserva int) error {
	return nil
}

func (hr *AiReservesRepository) SearchReserve(ctx context.Context, req domain.SearchReserve) error {
	return nil
}

func (hr *AiReservesRepository) InitAgenda(ctx context.Context, req domain.Agenda) error {
	return nil
}

func (hr *AiReservesRepository) GetInfoPersona(ctx context.Context, idPersona int) error {
	return nil
}

func (hr *AiReservesRepository) GetReservasPersona(ctx context.Context, req domain.GetReservaPersona) error {
	return nil
}

func (hr *AiReservesRepository) GetReservasUnidadReserva(ctx context.Context, req domain.GetReservaUnidadReserva) error {
	return nil
}

func (hr *AiReservesRepository) PushEventToQueue(ctx context.Context, tx *sql.Tx, event domain.Event) error {
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

func (hr *AiReservesRepository) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
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
