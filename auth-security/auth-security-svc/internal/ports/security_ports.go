package ports

import (
	"context"

	"github.com/FrancoRebollo/auth-security-svc/internal/adapters/in/http/dto"
	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
)

type SecurityService interface {
	CreateUserAPI(ctx context.Context, reqAltaUser domain.UserCreated) (*domain.UserCreated, error)
	CrearCanalDigitalAPI(ctx context.Context, crearCanalDigital domain.CanalDigital, apiKey string) error
	AccessPersonAPI(ctx context.Context, accessPerson domain.AccessPerson, apikey string) error
	AccessCanalDigitalAPI(ctx context.Context, accessCanaldigital domain.AccessCanalDigital, apikey string) error
	AccessApiKeyAPI(ctx context.Context, accessApiKey domain.AccessApiKey, apikey string) error
	AccessPersonMethodAuthAPI(ctx context.Context, accesPerMethodAuth domain.AccessPersonMethodAuth, apikey string) error
	LoginAPI(ctx context.Context, reqLogin domain.Login) (domain.UserStatus, error)
	ValidateJWTAPI(ctx context.Context, token string) (*domain.CheckJWT, error)
	GetJWTAPI(ctx context.Context, refreshToken string, accessTokenParam string) (string, error)
	CheckApiKeyExpiradaAPI(ctx context.Context, apiKey string) (bool, error)
	RecuperacionPasswordAPI(ctx context.Context, recuperacionPassword dto.ReqRecoveryPasswordDos) error
}

type SecurityRepository interface {
	CreateUser(ctx context.Context, reqAltaUser domain.UserCreated) (*domain.UserCreated, error)
	CrearCanalDigital(ctx context.Context, crearCanalDigital domain.CanalDigital, apiKey string) error
	AccessPerson(ctx context.Context, accessPerson domain.AccessPerson, apikey string) error
	AccessCanalDigital(ctx context.Context, accessCanaldigital domain.AccessCanalDigital, apikey string) error
	AccessApiKey(ctx context.Context, accessApiKey domain.AccessApiKey, apikey string) error
	AccessPersonMethodAuth(ctx context.Context, accesPersonMethodAuth domain.AccessPersonMethodAuth, apikey string) error
	LoginValidations(ctx context.Context, reqLogin domain.Login) (int, *string, error)
	GetAccessTokenDuration(ctx context.Context, ApiKey string) (int, error)
	UpsertAccessToken(ctx context.Context, requestUpsert *domain.UpsertAccessToken) error
	CheckLastAccessToken(ctx context.Context, token string, credentials domain.Credentials) error
	CheckLastRefreshToken(ctx context.Context, token string, credentials domain.Credentials) error
	CheckTokenCreation(ctx context.Context, credentials domain.Credentials) error
	PersistToken(ctx context.Context, credentials domain.CredentialsToken) error
	CheckApiKeyExpirada(ctx context.Context, apiKey string) (bool, error)
	CambioPasswordByLogin(ctx context.Context, loginName string, newPassword string) error
}
