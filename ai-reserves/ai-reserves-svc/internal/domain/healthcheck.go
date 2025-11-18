package domain

type Healthcheck struct {
	NombreApi     string     `json:"nombre_api"`
	Cliente       string     `json:"cliente"`
	Version       string     `json:"version"`
	VersionModelo string     `json:"version_modelo"`
	FechaStartUp  string     `json:"fecha_start_up"`
	BasesDeDatos  []Database `json:"bases_de_datos"`
}

type Database struct {
	Base                     string `json:"base"`
	FechaHoraUltimaActividad string `json:"fecha_hora_ultima_actividad"`
}
