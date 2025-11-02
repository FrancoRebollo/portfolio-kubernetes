// internal/adapters/in/http/router.go
package http

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/FrancoRebollo/auth-security-svc/internal/adapters/in/http/middlewares"
	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
	"github.com/FrancoRebollo/auth-security-svc/internal/platform/config"
	configconstants "github.com/FrancoRebollo/auth-security-svc/internal/platform/config/constants"
)

// Interfaces mínimas que deben cumplir tus handlers
type VersionHandler interface {
	GetVersion(c *gin.Context)
}

type Router struct {
	*gin.Engine
}

func NewRouter(
	cfg *config.HTTP,
	versionHandler VersionHandler,
	healthcheckHandler HealthcheckHandler,
	securityHandler SecurityHandler,
) (*Router, error) {

	// Modo
	if cfg.Environment == configconstants.PRODUCCION {
		gin.SetMode(gin.ReleaseMode)
	}

	// CORS
	ginConfig := cors.DefaultConfig()
	originsList := strings.Split(cfg.AllowedOrigins, ",")
	ginConfig.AllowOrigins = originsList

	// Server
	r := gin.New()

	// Middlewares globales
	r.Use(gin.Recovery(), cors.New(ginConfig))
	r.Use(middlewares.CancelCheckMiddleware())
	r.Use(middlewares.LoggerMiddleware())

	// Swagger
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rutas
	api := r.Group("/api")
	{
		// Version
		api.Group("/version").
			GET("", versionHandler.GetVersion)

		// Healthcheck
		api.Group("/healthcheck").
			GET("", middlewares.ValidateGetHealthcheck, healthcheckHandler.GetHealthcheck)
	}

	sec := r.Group("/sec")
	{
		sec.Group("/validate-jwt").GET("", middlewares.NewRateLimiterMiddleware(), securityHandler.ValidateJWT)
		sec.Group("/log-in").POST("", middlewares.NewRateLimiterMiddleware(), securityHandler.Login)
		sec.Group("/get-jwt").POST("", middlewares.NewRateLimiterMiddleware(), securityHandler.GetJWT)
		sec.Group("/recovery-password").POST("", middlewares.NewRateLimiterMiddleware(), securityHandler.RecoveryPassword)
	}

	adm := r.Group("/adm")
	{
		adm.Group("/create-user").POST("", middlewares.NewRateLimiterMiddleware(), securityHandler.CreateUser)
		adm.Group("/create-method-auth").POST("", middlewares.NewRateLimiterMiddleware(), securityHandler.CreateCanalDigital)
		adm.Group("/unaccess-person").POST("", middlewares.NewRateLimiterMiddleware(), securityHandler.AccessPerson)
		adm.Group("/unaccess-digital-channel").POST("", middlewares.NewRateLimiterMiddleware(), securityHandler.AccessCanalDigital)
		adm.Group("/unaccess-digital-channel-person").POST("", middlewares.NewRateLimiterMiddleware(), securityHandler.CreateCanalDigital)
		adm.Group("/unaccess-api-key").POST("", middlewares.NewRateLimiterMiddleware(), securityHandler.AcessApiKey)
	}

	// 404
	r.NoRoute(func(c *gin.Context) {
		err := domain.HealthcheckError{
			Code:    domain.ErrCodeRouteNotFound,
			Message: "La ruta solicitada no existe en el servidor",
		}
		c.JSON(http.StatusNotFound, err)
	})

	return &Router{r}, nil
}

func (r *Router) Listen(addr string) error {
	return r.Run(addr)
}
