// internal/adapters/in/http/router.go
package http

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/FrancoRebollo/ai-reserves-svc/internal/adapters/in/http/middlewares"
	"github.com/FrancoRebollo/ai-reserves-svc/internal/domain"
	"github.com/FrancoRebollo/ai-reserves-svc/internal/platform/config"
	configconstants "github.com/FrancoRebollo/ai-reserves-svc/internal/platform/config/constants"
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
	AiReservesHandler AiReservesHandler,
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

	ai_res := r.Group("/reserves")
	{
		//ai_res.Group("/create-person").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.CreatePersona)
		ai_res.Group("/upd-atribute-person").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.UpdAtributoPersona)
		ai_res.Group("/upd-person").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.UpdPersona)

		ai_res.Group("/upsert-config-person").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.UpsertConfigPersona)
		ai_res.Group("/create-unidad-reserva").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.CreateUnidadReserva)
		ai_res.Group("/create-tipo-unidad-reserva").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.CreateTipoUnidadReserva)

		ai_res.Group("/create-sub-tipo-unidad-reserva").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.CreateSubTipoUnidadReserva)
		ai_res.Group("/upd-unidad-reserva").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.ModifUnidadReserva)
		ai_res.Group("/upd-tipo-unidad-reserva").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.ModifTipoUnidadReserva)

		ai_res.Group("/upd-sub-tipo-unidad-reserva").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.ModifSubTipoUnidadReserva)
		ai_res.Group("/upd-atribute-unidad-reserva").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.ModifUnidadReservaParcial)
		ai_res.Group("/upd-atribute-tipo-unidad-reserva").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.ModifTipoUnidadReservaParcial)
		ai_res.Group("/upd-atribute-sub-tipo-unidad-reserva").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.ModifSubTipoUnidadReservaParcial)

		//
		ai_res.Group("/create-reserve").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.CreateReserve)
		ai_res.Group("/cancel-reserve").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.CancelReserve)

		ai_res.Group("/search-reserve").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.SearchReserve)
		ai_res.Group("/init-agenda").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.InitAgenda)
		ai_res.Group("/get-info-person").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.GetInfoPersona)

		ai_res.Group("/get-reserves-person").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.GetReservasPersona)
		ai_res.Group("/get-reserves-unidad-reserva").POST("", middlewares.SecurityMiddleware(), middlewares.NewRateLimiterMiddleware(), AiReservesHandler.GetReservasUnidadReserva)

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
