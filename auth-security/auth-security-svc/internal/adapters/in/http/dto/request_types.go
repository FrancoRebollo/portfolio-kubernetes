package dto

type RequestAltaUser struct {
	IdPersona    int    `json:"id_persona"`
	CanalDigital string `json:"canal_digital"`
	LoginName    string `json:"login_name"`
	Password     string `json:"password"`
	MailPersona  string `json:"mail_persona"`
	TePersona    string `json:"tel_persona"`
}

type ReqCrearCanalDigital struct {
	CanalDigital string `json:"canal_digital"`
}

type ReqAccessPerson struct {
	IdPersona int    `json:"id_persona"`
	Revoke    string `json:"revoke"`
}

type ReqAccessDigitalChannel struct {
	CanalDigital string `json:"canal_digital"`
	Revoke       string `json:"revoke"`
}

type ReqAccessApiKey struct {
	ApiKey string `json:"api_key"`
	Revoke string `json:"revoke"`
}

type ReqAccessPerMethodAuth struct {
	IdPersona  int    `json:"id_persona"`
	MethodAuth string `json:"method_auth"`
	Revoke     string `json:"revoke"`
}

type ReqLogin struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ApiKey       string
	CanalDigital string `json:"canal_digital"`
}

type ReqValidateJWT struct {
	Jwt string `json:"jwt"`
}

type ReqGetJWT struct {
	RefreshToken string `json:"refresh_token"`
}

type ReqRecoveryPassword struct {
	LoginName string `json:"login_name"`
	ApiKey    string `json:"api_key"`
}

type ReqRecoveryPasswordDos struct {
	LoginName string `json:"login_name"`
}
