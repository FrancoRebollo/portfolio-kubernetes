package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/FrancoRebollo/auth-security-svc/internal/adapters/in/http/dto"
	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
	"github.com/FrancoRebollo/auth-security-svc/internal/platform/utils"

	"github.com/FrancoRebollo/auth-security-svc/internal/ports"
	"github.com/gin-gonic/gin"
)

type SecurityHandler struct {
	serv ports.SecurityService
}

func NewSecurityHandler(serv ports.SecurityService) *SecurityHandler {
	return &SecurityHandler{
		serv,
	}
}

func (hh *SecurityHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var altaUser dto.RequestAltaUser

	if err := c.BindJSON(&altaUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainUser := domain.UserCreated{
		IdPersona:    altaUser.IdPersona,
		CanalDigital: altaUser.CanalDigital,
		LoginName:    altaUser.LoginName,
		Password:     altaUser.Password,
		MailPersona:  altaUser.MailPersona,
		TePersona:    altaUser.TePersona,
	}

	userCreated, err := hh.serv.CreateUserAPI(ctx, domainUser)
	if err != nil {
		/*
			logger.LoggerError().Error(err)
			errorResponse(c, err)
		*/
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userCreated)
}

func (hh *SecurityHandler) CreateCanalDigital(c *gin.Context) {
	var crearCanalDigital dto.ReqCrearCanalDigital

	apiKey := c.GetHeader("Api-Key")

	if err := c.BindJSON(&crearCanalDigital); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainCanal := domain.CanalDigital{
		CanalDigital: crearCanalDigital.CanalDigital,
	}

	if err := hh.serv.CrearCanalDigitalAPI(c, domainCanal, apiKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := dto.DefaultResponse{
		Message: "Canal digital creado",
	}

	c.JSON(200, resp)
}

func (hh *SecurityHandler) AccessPerson(c *gin.Context) {
	var accesPerson dto.ReqAccessPerson

	apiKey := c.GetHeader("Api-Key")

	if err := c.BindJSON(&accesPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domAccessPerson := domain.AccessPerson{
		IdPersona: accesPerson.IdPersona,
		Revoke:    accesPerson.Revoke,
	}

	if err := hh.serv.AccessPersonAPI(c, domAccessPerson, apiKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resp dto.DefaultResponse

	if domAccessPerson.Revoke == "S" {

		resp = dto.DefaultResponse{
			Message: "Revoke access to person",
		}
	} else {
		resp = dto.DefaultResponse{
			Message: "Revoke unaccess to person",
		}

	}

	c.JSON(200, resp)
}

func (hh *SecurityHandler) AccessCanalDigital(c *gin.Context) {
	var accesCanalDigital dto.ReqAccessDigitalChannel

	apiKey := c.GetHeader("Api-Key")

	if err := c.BindJSON(&accesCanalDigital); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainAccessCanalDigital := domain.AccessCanalDigital{
		CanalDigital: accesCanalDigital.CanalDigital,
		Revoke:       accesCanalDigital.Revoke,
	}

	if err := hh.serv.AccessCanalDigitalAPI(c, domainAccessCanalDigital, apiKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resp dto.DefaultResponse

	if domainAccessCanalDigital.Revoke == "S" {

		resp = dto.DefaultResponse{
			Message: "Revoke access to digital channel",
		}
	} else {
		resp = dto.DefaultResponse{
			Message: "Revoke unaccess to digital channel",
		}

	}
	c.JSON(200, resp)
}

func (hh *SecurityHandler) AcessApiKey(c *gin.Context) {
	var accessApiKey dto.ReqAccessApiKey

	apiKey := c.GetHeader("Api-Key")

	if err := c.BindJSON(&accessApiKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainAccesApiKey := domain.AccessApiKey{
		ApiKey: accessApiKey.ApiKey,
		Revoke: accessApiKey.Revoke,
	}

	if err := hh.serv.AccessApiKeyAPI(c, domainAccesApiKey, apiKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resp dto.DefaultResponse

	if domainAccesApiKey.Revoke == "S" {

		resp = dto.DefaultResponse{
			Message: "Revoke access to api key",
		}
	} else {
		resp = dto.DefaultResponse{
			Message: "Revoke unaccess to api key",
		}

	}

	c.JSON(200, resp)
}

func (hh *SecurityHandler) AccessPerMethodAuth(c *gin.Context) {
	var accessPerMethodAuth dto.ReqAccessPerMethodAuth

	apiKey := c.GetHeader("Api-Key")

	if err := c.BindJSON(&accessPerMethodAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainPerMethodAuth := domain.AccessPersonMethodAuth{
		IdPersona:  accessPerMethodAuth.IdPersona,
		Revoke:     accessPerMethodAuth.Revoke,
		MethodAuth: accessPerMethodAuth.MethodAuth,
	}

	if err := hh.serv.AccessPersonMethodAuthAPI(c, domainPerMethodAuth, apiKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resp dto.DefaultResponse

	if domainPerMethodAuth.Revoke == "S" {

		resp = dto.DefaultResponse{
			Message: "Revoke access to person by digital channel",
		}
	} else {
		resp = dto.DefaultResponse{
			Message: "Revoke unaccess to person by digital channel",
		}

	}

	c.JSON(200, resp)
}

func (hh *SecurityHandler) Login(c *gin.Context) {

	var reqLogin dto.ReqLogin

	if err := c.BindJSON(&reqLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainLogin := &domain.Login{
		Username:     reqLogin.Username,
		Password:     reqLogin.Password,
		ApiKey:       c.GetHeader("Api-Key"),
		CanalDigital: reqLogin.CanalDigital,
	}

	domainUserStatus, err := hh.serv.LoginAPI(c, *domainLogin)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, domainUserStatus)
}

func (h *SecurityHandler) ValidateJWT(c *gin.Context) {
	fmt.Println("Entra validate JWT handler")
	var idPersona int

	accessToken := c.GetHeader("Authorization")
	accessBear := strings.TrimPrefix(accessToken, "Bearer ")

	claims, err := utils.GetClaimsFromToken(accessBear, "ACCESS")

	switch v := claims["id_persona"].(type) {
	case float64:
		idPersona = int(v)
	case int:
		idPersona = v
	case string:
		idPersona, _ = strconv.Atoi(v)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read claims"})
		return
	}

	checkJWTResponse, err := h.serv.ValidateJWTAPI(c, accessBear)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	checkJWTResponse.IdPersona = idPersona

	c.JSON(200, checkJWTResponse)
}

func (h *SecurityHandler) GetJWT(c *gin.Context) {

	var reqGetJWT dto.ReqGetJWT

	if err := c.BindJSON(&reqGetJWT); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken := c.GetHeader("Authorization")
	accessBear := strings.TrimPrefix(accessToken, "Bearer ")
	/*
		credentialsExt := domain.CredentialsExtended{
			IdPersona:    0,
			ApiKey:       "n/a",
			CanalDigital: "",
			IpAddress:    c.ClientIP(),
			Endpoint:     c.FullPath(),
		}
	*/
	jwt, err := h.serv.GetJWTAPI(c, reqGetJWT.RefreshToken, accessBear)

	if err != nil {
		//h.LogProcedure(c, credentialsExt, err.Error(), reqGetJWT.RefreshToken, 0)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	resp := &dto.GetJWTResponse{
		AccessToken: jwt,
	}

	c.JSON(200, resp)
}

func (h *SecurityHandler) RecoveryPassword(c *gin.Context) {

	var recuperacionPasswordDos dto.ReqRecoveryPassword
	var recuperacionPassword dto.ReqRecoveryPasswordDos

	if err := c.BindJSON(&recuperacionPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recuperacionPasswordDos.LoginName = recuperacionPassword.LoginName
	recuperacionPasswordDos.ApiKey = c.GetHeader("Api-Key")

	if recuperacionPasswordDos.ApiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You must to provide one valid api-key"})
		return
	}

	apiKeyExpirada, err := h.serv.CheckApiKeyExpiradaAPI(c, recuperacionPasswordDos.ApiKey)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if apiKeyExpirada {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "api key expirada o revocada"})
		return
	}

	if err := h.serv.RecuperacionPasswordAPI(c, recuperacionPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := dto.DefaultResponse{
		Message: "Una nueva contraseña fue enviada, verifique su correo electronico",
	}

	c.JSON(200, resp)
}
