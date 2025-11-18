package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
	"github.com/FrancoRebollo/api-integration-svc/internal/platform/logger"
	"github.com/ulule/limiter/v3"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memoryStore "github.com/ulule/limiter/v3/drivers/store/memory"

	"github.com/gin-gonic/gin"
)

type SecurityResponse struct {
	IdPersona   int    `json:"id_persona"`
	TokenStatus string `json:"token_status"`
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capturar la respuesta
		responseRecorder := &logger.ResponseRecorder{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseRecorder

		c.Next()

		// Loguear la solicitud y la respuesta juntas después de que se haya procesado
		logger.LoggerHttp(c, responseRecorder.Body.String())
	}
}

func CancelCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		select {
		case <-ctx.Done():
			c.JSON(http.StatusRequestTimeout, domain.HealthcheckError{
				Code:    domain.ErrCodeRequestTimeout,
				Message: "Solicitud cancelada"})
			c.Abort()
			return
		default:
			// Continua al siguiente middleware o handler si no ha sido cancelado
		}

		c.Next()
	}
}

func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		fmt.Println("En security middleware")
		fmt.Println(accessToken)
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, domain.HealthcheckError{
				Code:    domain.ErrCodeUnauthorized,
				Message: "Token de autenticación inválido o ausente"})
			c.Abort()
			return
		}

		client := &http.Client{}

		urlSeg := os.Getenv("HTTP_SECURITY_URL")

		req, err := http.NewRequest("GET", urlSeg, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.HealthcheckError{
				Code:    domain.ErrCodeInternalServer,
				Message: "Ocurrió un error inesperado " + err.Error()})
			c.Abort()
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", accessToken)

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.HealthcheckError{
				Code:    domain.ErrCodeInternalServer,
				Message: "Ocurrió un error inesperado" + err.Error()})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			c.JSON(http.StatusUnauthorized, domain.HealthcheckError{
				Code:    domain.ErrCodeUnauthorized,
				Message: "Token de autenticación inválido o ausente"})
			c.Abort()
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.HealthcheckError{
				Code:    domain.ErrCodeInternalServer,
				Message: "Ocurrió un error inesperado: " + err.Error()})
			c.Abort()
			return
		}

		var securityResp SecurityResponse
		err = json.Unmarshal(body, &securityResp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.HealthcheckError{
				Code:    domain.ErrCodeInternalServer,
				Message: "Ocurrió un error inesperado" + err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func NewRateLimiterMiddleware() gin.HandlerFunc {

	limit := os.Getenv("RATE_LIMITATING")

	rate, err := limiter.NewRateFromFormatted(limit)
	if err != nil {
		panic(err)
	}

	store := memoryStore.NewStore()
	instance := limiter.New(store, rate)

	return ginLimiter.NewMiddleware(instance)
}
