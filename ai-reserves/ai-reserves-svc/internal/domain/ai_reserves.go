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
	Nombre          string
	Descripcion     string
}

type TipoUnidadReserva struct {
	IDUnidadReserva     int
	IDTipoUnidadReserva int
	UnidadReserva       string
	Nombre              string
	Descripcion         string
}

type SubTipoUnidadReserva struct {
	IDUnidadReserva        int
	UnidadReserva          string
	IDTipoUnidadReserva    int
	TipoUnidadReserva      string
	IDSubTipoUnidadReserva int
	Nombre                 string
	Descripcion            string
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
	IDSubTipoUnidadReserva int
	FechaDesde             time.Time
	FechaHasta             time.Time
}

type PersonCreatedPayload struct {
	ID        int
	Email     string
	TePersona string
}

type UpdAtributeUnidadReserva struct {
	IDUnidadReserva int
	Atribute        string
	Value           string
}

type UpdUnidadReserva struct {
	IDUnidadReserva int
	Nombre          string
	Descripcion     string
}

type UpdAtributeTipoUnidadReserva struct {
	IDUnidadReserva     int
	IDTipoUnidadReserva int
	Atribute            string
	Value               string
}

type UpdAtributeSubTipoUnidadReserva struct {
	IDUnidadReserva        int
	IDTipoUnidadReserva    int
	IDSubTipoUnidadReserva int
	Atribute               string
	Value                  string
}

type UpdTipoUnidadReserva struct {
	IDUnidadReserva     int
	IDTipoUnidadReserva int
	Nombre              string
	Descripcion         string
}

type UpdSubTipoUnidadReserva struct {
	IDUnidadReserva        int
	IDTipoUnidadReserva    int
	IDSubTipoUnidadReserva int
	Nombre                 string
	Descripcion            string
}

type Reserva struct {
	ID                     int
	IDAgenda               int
	Fecha                  time.Time
	HoraInicio             string
	HoraFin                string
	IDPaciente             *int
	Estado                 string
	Observaciones          *string
	IDSubTipoUnidadReserva int
}
