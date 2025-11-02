package dto

type AltaUserResponse struct {
	IdPersona    int    `json:"id_persona"`
	CanalDigital string `json:"canal_digital"`
	Message      string `json:"message"`
}

type DefaultResponse struct {
	Message string `json:"Message"`
}

type LoginResponse struct {
	Username     string `json:"username"`
	Status       string `json:"status"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Hash2FA      string `json:"hash_2fa"`
}

type GetJWTResponse struct {
	AccessToken string `json:"access_token"`
}
