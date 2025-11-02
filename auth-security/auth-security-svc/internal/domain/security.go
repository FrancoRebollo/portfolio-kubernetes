package domain

import "time"

type UserCreated struct {
	IdPersona    int    `json:"id_persona"`
	CanalDigital string `json:"canal_digital"`
	LoginName    string `json:"login_name"`
	Password     string `json:"password"`
	MailPersona  string `json:"mail_persona"`
	TePersona    string `json:"tel_persona"`
}

type CanalDigital struct {
	CanalDigital string
}

type AccessPerson struct {
	IdPersona int
	Revoke    string
}

type AccessCanalDigital struct {
	CanalDigital string
	Revoke       string
}

type AccessApiKey struct {
	ApiKey        string
	FechaVigencia string
	Revoke        string
}

type AccessPersonMethodAuth struct {
	IdPersona  int
	MethodAuth string
	Revoke     string
}

type Login struct {
	Username     string
	Password     string
	ApiKey       string
	CanalDigital string
}

type Credentials struct {
	IdPersona    int
	ApiKey       string
	CanalDigital string
}

type CheckJWT struct {
	IdPersona   int
	TokenStatus string
}

type UserStatus struct {
	Username     string
	Status       string
	RefreshToken string
	AccessToken  string
	Hash2FA      string
}

type UpsertAccessToken struct {
	IdPersona    int
	CanalDigital string
	ApiKey       string
	AccessToken  string
	RefreshToken string
}

type JWT struct {
	JWT string
}

type CredentialsExtended struct {
	IdPersona    int
	ApiKey       string
	CanalDigital string
	IpAddress    string
	Endpoint     string
}

type CredentialsToken struct {
	IdPersona    int
	ApiKey       string
	CanalDigital string
	AccessToken  string
	RefreshToken string
}

type Event struct {
	ID         string
	Type       string
	RoutingKey string
	Origin     string
	Timestamp  time.Time
	Payload    interface{}
}
