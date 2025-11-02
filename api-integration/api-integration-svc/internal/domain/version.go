package domain

type Version struct {
	NombreApi    string `json:"nombre_api"`
	Cliente      string `json:"cliente"`
	Version      string `json:"version"`
	FechaStartUp string `json:"fecha_start_up"`
}
