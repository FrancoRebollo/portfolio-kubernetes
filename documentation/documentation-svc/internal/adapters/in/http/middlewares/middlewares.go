package middlewares

import (
	"bytes"
	"net/http"

	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
	"github.com/FrancoRebollo/api-integration-svc/internal/platform/logger"

	"github.com/gin-gonic/gin"
)

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
