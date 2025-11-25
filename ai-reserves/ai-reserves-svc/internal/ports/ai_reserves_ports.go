package ports

import (
	"context"
	"database/sql"

	"github.com/FrancoRebollo/ai-reserves-svc/internal/domain"
)

type AiReservesService interface {
	CreatePersonaAPI(ctx context.Context, req domain.PersonCreatedPayload) error
	UpdAtributoPersonaAPI(ctx context.Context, req domain.PersonaParcial) error
	UpdPersonaAPI(ctx context.Context, req domain.Persona) error
	UpsertConfigPersonaAPI(ctx context.Context, req domain.ConfigPersona) error

	CreateUnidadReservaAPI(ctx context.Context, req *domain.UnidadReserva) error
	CreateTipoUnidadReservaAPI(ctx context.Context, req domain.TipoUnidadReserva) error
	CreateSubTipoUnidadReservaAPI(ctx context.Context, req domain.SubTipoUnidadReserva) error

	ModifUnidadReservaAPI(ctx context.Context, req domain.UpdUnidadReserva) error
	UpdAtributeUnidadReservaAPI(ctx context.Context, req domain.UpdAtributeUnidadReserva) error

	ModifTipoUnidadReservaAPI(ctx context.Context, req domain.UpdTipoUnidadReserva) error
	ModifSubTipoUnidadReservaAPI(ctx context.Context, req domain.UpdSubTipoUnidadReserva) error
	UpdAtributeTipoUnidadReservaAPI(ctx context.Context, req domain.UpdAtributeTipoUnidadReserva) error
	UpdAtributeSubTipoUnidadReservaAPI(ctx context.Context, req domain.UpdAtributeSubTipoUnidadReserva) error

	CreateReserveAPI(ctx context.Context, req domain.Reserva) error
	CancelReserveAPI(ctx context.Context, idReserva int) error
	SearchReserveAPI(ctx context.Context, req domain.SearchReserve) error
	InitAgendaAPI(ctx context.Context, req domain.Agenda) error

	GetInfoPersonaAPI(ctx context.Context, idPersona int) error
	GetReservasPersonaAPI(ctx context.Context, req domain.GetReservaPersona) error
	GetReservasUnidadReservaAPI(ctx context.Context, req domain.GetReservaUnidadReserva) ([]domain.Reserva, error)
}

type AiReservesRepository interface {
	CreatePersona(ctx context.Context, req domain.PersonCreatedPayload) error
	UpdAtributoPersona(ctx context.Context, req domain.PersonaParcial) error
	UpdPersona(ctx context.Context, req domain.Persona) error
	UpsertConfigPersona(ctx context.Context, req domain.ConfigPersona) error

	CreateUnidadReserva(ctx context.Context, req *domain.UnidadReserva) (int, error)
	CreateTipoUnidadReserva(ctx context.Context, req domain.TipoUnidadReserva) (int, error)
	CreateSubTipoUnidadReserva(ctx context.Context, req domain.SubTipoUnidadReserva) (int, error)

	ModifUnidadReserva(ctx context.Context, req domain.UpdUnidadReserva) error
	UpdAtributeUnidadReserva(ctx context.Context, req domain.UpdAtributeUnidadReserva) error

	ModifTipoUnidadReserva(ctx context.Context, req domain.UpdTipoUnidadReserva) error
	UpdAtributeTipoUnidadReserva(ctx context.Context, req domain.UpdAtributeTipoUnidadReserva) error
	ModifSubTipoUnidadReserva(ctx context.Context, req domain.UpdSubTipoUnidadReserva) error
	UpdAtributeSubTipoUnidadReserva(ctx context.Context, req domain.UpdAtributeSubTipoUnidadReserva) error

	CreateReserve(ctx context.Context, req domain.Reserva) error
	CancelReserve(ctx context.Context, idReserva int) error
	SearchReserve(ctx context.Context, req domain.SearchReserve) error
	InitAgenda(ctx context.Context, req domain.Agenda) error

	GetInfoPersona(ctx context.Context, idPersona int) error
	GetReservasPersona(ctx context.Context, req domain.GetReservaPersona) error
	GetReservasUnidadReserva(ctx context.Context, req domain.GetReservaUnidadReserva) ([]domain.Reserva, error)

	PushEventToQueue(ctx context.Context, tx *sql.Tx, event domain.Event) error
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error
}

type MessageQueue interface {
	Publish(ctx context.Context, event domain.Event) error
}
