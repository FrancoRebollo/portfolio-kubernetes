package dto

import "time"

type ReqCaptureEvent struct {
	IdEvent      int    `json:"id_event"`
	EventType    string `json:"event_type"`
	EventContent string `json:"event_content"`
}

type ExternalAPIRequest struct {
	Method string            `json:"method"` // "GET" or "POST"
	URL    string            `json:"url"`
	Params map[string]string `json:"params,omitempty"` // For GET
	Body   map[string]any    `json:"body,omitempty"`   // For POST
}

type RequestPushEvent struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`       // semántica de dominio
	RoutingKey string      `json:"routingKey"` // para RabbitMQ
	Origin     string      `json:"origin"`
	Timestamp  time.Time   `json:"timestamp"`
	Payload    interface{} `json:"payload"`
}

type Persona struct {
	ID                  int    `json:"id"`
	Nombre              string `json:"nombre"`
	ApellidoRazonSocial string `json:"apellido_razon_social"`
	PersonaJuridia      string `json:"persona_juridica"`
	TipoDocPersona      string `json:"tipo_doc_persona"`
	NroDocPersona       string `json:"nro_doc_persona"`
	Email               string `json:"email"`
	TelPersona          string `json:"tel_persona"`
}

type PersonaParcial struct {
	ID       int    `json:"id"`
	Atribute string `json:"atribute"`
	Value    string `json:"value"`
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
	IDPersona int    `json:"id_persona"`
	Atribute  string `json:"atribute"`
	Value     string `json:"value"`
}

type UnidadReserva struct {
	IDUnidadReserva int
	Nombre          string `json:"nombre"`
	Descripcion     string `json:"descripcion"`
}

type UpdAtributeUnidadReserva struct {
	IDUnidadReserva int
	Atribute        string `json:"atribute"`
	Value           string `json:"value"`
}

type UpdUnidadReserva struct {
	IDUnidadReserva int    `json:"id_unidad_reserva"`
	Nombre          string `json:"nombre"`
	Descripcion     string `json:"descripcion"`
}

type TipoUnidadReserva struct {
	IDUnidadReserva     int `json:"id_unidad_reserva"`
	IDTipoUnidadReserva int
	UnidadReserva       string `json:"unidad_reserva"`
	Nombre              string `json:"nombre"`
	Descripcion         string `json:"descripcion"`
}

type SubTipoUnidadReserva struct {
	IDUnidadReserva        int    `json:"id_unidad_reserva"`
	UnidadReserva          string `json:"unidad_reserva"`
	IDTipoUnidadReserva    int    `json:"id_tipo_unidad_reserva"`
	TipoUnidadReserva      string `json:"tipo_unidad_reserva"`
	IDSubTipoUnidadReserva int
	Nombre                 string `json:"nombre"`
	Descripcion            string `json:"descripcion"`
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
