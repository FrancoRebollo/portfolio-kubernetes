package domain

import "time"

type ExternalAPIRequest struct {
	Method string
	URL    string
	Params map[string]string
	Body   map[string]any
}

type ExternalAPIResponse struct {
	Status     string
	StatusCode int
	Data       any
}

type Persona struct {
	ID                  int
	Nombre              string
	ApellidoRazonSocial string
	PersonaJuridia      string
	TipoDocPersona      string
	NroDocPersona       string
	Email               string
	TelPersona          string
}

type PersonaParcial struct {
	ID       int
	Atribute string
	Value    string
}

type Event struct {
	ID         string
	Type       string
	RoutingKey string
	Origin     string
	Timestamp  time.Time
	Payload    interface{}
}

type ConfigPersona struct {
	IDPersona int
	Atribute  string
	Value     string
}

type UnidadReserva struct {
	IDUnidadReserva int
	Descripcion     string
}

type TipoUnidadReserva struct {
	IDUnidadReserva int
	Descripcion     string
}

type SubTipoUnidadReserva struct {
	IDUnidadReserva int
	Descripcion     string
}

type UnidadReservaModif struct {
	IDUnidadReserva int
	Atribute        string
	Value           string
}

type TipoUnidadReservaModif struct {
	IDUnidadReserva int
	Atribute        string
	Value           string
}

type SubTipoUnidadReservaModif struct {
	IDUnidadReserva int
	Atribute        string
	Value           string
}

type Reserva struct {
	IDUReserva      int
	IDPersona       int
	IDUnidadReserva int
	FechaInicio     time.Time
	FechaFin        time.Time
	Estado          string
}

type ReservaCancel struct {
	IDReserva int
	Status    string
}

type SearchReserve struct {
	IDProfesional   int
	IDUnidadReserva int
	FechaDesde      time.Time
	FechaHasta      time.Time
}

type Agenda struct {
	IDProfesional   int
	IDUnidadReserva int
	FechaDesde      time.Time
	FechaHasta      time.Time
}

type GetReservaPersona struct {
	IDPersona  int
	FechaDesde time.Time
	FechaHasta time.Time
}

type GetReservaUnidadReserva struct {
	IDUnidadReserva int
	FechaDesde      time.Time
	FechaHasta      time.Time
}

type PersonCreatedPayload struct {
	ID        int
	Email     string
	TePersona string
}
