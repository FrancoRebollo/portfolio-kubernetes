// internal/adapters/in/http/router.go
package http

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/FrancoRebollo/api-integration-svc/internal/adapters/in/http/middlewares"
	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
	"github.com/FrancoRebollo/api-integration-svc/internal/platform/config"
	configconstants "github.com/FrancoRebollo/api-integration-svc/internal/platform/config/constants"
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
