package application

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/FrancoRebollo/auth-security-svc/internal/adapters/in/http/dto"
	"github.com/FrancoRebollo/auth-security-svc/internal/platform/config"
	"github.com/FrancoRebollo/auth-security-svc/internal/platform/utils"
	"github.com/FrancoRebollo/auth-security-svc/internal/ports"

	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
)

type SecurityService struct {
	hr   ports.SecurityRepository
	conf config.App
}

func NewSecurityService(hr ports.SecurityRepository, conf config.App) *SecurityService {
	return &SecurityService{
		hr,
		conf,
	}
}

func (hs *SecurityService) CreateUserAPI(ctx context.Context, reqAltaUser domain.UserCreated) (*domain.UserCreated, error) {
	//var serviceErr error
	fmt.Println(reqAltaUser.CanalDigital)
	fmt.Println(reqAltaUser.LoginName)
	userCreated, err := hs.hr.CreateUser(ctx, reqAltaUser)

	if err != nil {
		/*
			serviceErr = mapServiceError(err)
			logger.LoggerError().WithError(err).Error(serviceErr)
			return &domain.UserCreated{}, serviceErr
		*/
		return &domain.UserCreated{}, err
	}

	return userCreated, nil
}

func (s *SecurityService) CrearCanalDigitalAPI(ctx context.Context, crearCanalDigital domain.CanalDigital, apiKey string) error {

	if err := s.hr.CrearCanalDigital(ctx, crearCanalDigital, apiKey); err != nil {
		return err
	}

	return nil
}

func (s *SecurityService) AccessPersonAPI(ctx context.Context, accesPerson domain.AccessPerson, apiKey string) error {

	if err := s.hr.AccessPerson(ctx, accesPerson, apiKey); err != nil {
		return err
	}

	return nil
}

func (s *SecurityService) AccessCanalDigitalAPI(ctx context.Context, accessCanaldigital domain.AccessCanalDigital, apiKey string) error {

	if err := s.hr.AccessCanalDigital(ctx, accessCanaldigital, apiKey); err != nil {
		return err
	}

	return nil
}

func (s *SecurityService) AccessApiKeyAPI(ctx context.Context, accessApiKey domain.AccessApiKey, apiKey string) error {

	if err := s.hr.AccessApiKey(ctx, accessApiKey, apiKey); err != nil {
		return err
	}

	return nil
}

func (s *SecurityService) AccessPersonMethodAuthAPI(ctx context.Context, accessAccessPersonMethodAuth domain.AccessPersonMethodAuth, apiKey string) error {

	if err := s.hr.AccessPersonMethodAuth(ctx, accessAccessPersonMethodAuth, apiKey); err != nil {
		return err
	}

	return nil
}

func (s *SecurityService) LoginAPI(ctx context.Context, reqLogin domain.Login) (domain.UserStatus, error) {

	resp := &domain.UserStatus{
		Username:     reqLogin.Username,
		Status:       "error",
		RefreshToken: "",
		AccessToken:  "",
		Hash2FA:      "",
	}

	idPersona, seed2FA, err := s.hr.LoginValidations(ctx, reqLogin)

	if err != nil {
		return *resp, err
	}

	if seed2FA != nil {

		encrypted2FA, err := utils.EncryptTwo(reqLogin.Username+":"+reqLogin.Password, *seed2FA)
		if err != nil {
			return *resp, err
		}

		resp := &domain.UserStatus{
			Username:     reqLogin.Username,
			Status:       "Ingrese el codigo de seguridad de su aplicacion",
			RefreshToken: "",
			AccessToken:  "",
			Hash2FA:      encrypted2FA,
		}

		return *resp, nil

	}

	credentials := domain.Credentials{
		IdPersona:    idPersona,
		CanalDigital: reqLogin.CanalDigital,
		ApiKey:       reqLogin.ApiKey,
	}

	ctdMins, err := s.hr.GetAccessTokenDuration(ctx, credentials.ApiKey)

	if err != nil {
		return *resp, err
	}
	//
	accessToken, err := utils.JWTCreate(ctdMins, credentials, "ACCESS")

	if err != nil {
		accessToken = "error en creacion"
	}

	refreshDuration, err := strconv.Atoi(os.Getenv("REF_TOKEN_DURATION"))

	if err != nil {
		return *resp, err
	}

	refreshToken, err := utils.JWTCreate(refreshDuration, credentials, "REFRESH")

	if err != nil {
		refreshToken = "error en creacion"
	}

	upsertAccessToken := &domain.UpsertAccessToken{
		IdPersona:    credentials.IdPersona,
		CanalDigital: credentials.CanalDigital,
		ApiKey:       credentials.ApiKey,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if err := s.hr.UpsertAccessToken(ctx, upsertAccessToken); err != nil {
		return *resp, err
	}

	resp = &domain.UserStatus{
		Username:     reqLogin.Username,
		Status:       "Logged",
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		Hash2FA:      "",
	}

	return *resp, nil
}

func (s *SecurityService) ValidateJWTAPI(ctx context.Context, token string) (*domain.CheckJWT, error) {

	checkJWTResponse, err := utils.CheckJWTAccessToken(token)

	if err != nil {
		return nil, err
	}

	return checkJWTResponse, nil
}

func (s *SecurityService) GetJWTAPI(ctx context.Context, refreshToken string, accessTokenParam string) (string, error) {

	expirationTime, err := utils.GetTokenExpiration(refreshToken, "REFRESH")

	if err != nil {
		return "", err
	}

	if expirationTime != nil && expirationTime.Before(time.Now()) {
		return "", fmt.Errorf("inicie sesion nuevamente")
	}

	claims, err := utils.GetClaimsFromToken(refreshToken, "REFRESH")

	if err != nil {
		return "", err
	}

	credentials := domain.Credentials{
		IdPersona:    int(claims["id_persona"].(float64)),
		ApiKey:       claims["api_key"].(string),
		CanalDigital: claims["canal_digital"].(string),
	}

	if err := s.hr.CheckTokenCreation(ctx, credentials); err != nil {
		return "", err
	}

	if err := s.hr.CheckLastRefreshToken(ctx, refreshToken, credentials); err != nil {
		return "", err
	}

	if err := s.hr.CheckLastAccessToken(ctx, accessTokenParam, credentials); err != nil {
		return "", err
	}

	ctdMins, err := s.hr.GetAccessTokenDuration(ctx, credentials.ApiKey)

	if err != nil {
		return "", err
	}

	accessToken, err := utils.JWTCreate(ctdMins, credentials, "ACCESS")

	if err != nil {
		accessToken = "error en creacion"
	}

	credentialsToken := domain.CredentialsToken{
		IdPersona:    credentials.IdPersona,
		ApiKey:       credentials.ApiKey,
		CanalDigital: credentials.CanalDigital,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if err := s.hr.PersistToken(ctx, credentialsToken); err != nil {
		return "", err
	}

	return accessToken, nil

}

func (s *SecurityService) CheckApiKeyExpiradaAPI(ctx context.Context, apiKey string) (bool, error) {

	bool, err := s.hr.CheckApiKeyExpirada(ctx, apiKey)

	if err != nil {
		return false, err
	}

	if !bool {
		return false, nil
	}

	return true, nil
}

func (s *SecurityService) RecuperacionPasswordAPI(ctx context.Context, recuperacionPassword dto.ReqRecoveryPasswordDos) error {

	newPassword, err := utils.GenerateRandomPassword(16)

	if err != nil {
		return fmt.Errorf("no fue posible generar una nueva contraseña")
	}

	if err := s.hr.CambioPasswordByLogin(ctx, recuperacionPassword.LoginName, newPassword); err != nil {
		return err
	}

	return nil
}
