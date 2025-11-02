package middlewares

import (
	"bytes"
	"net/http"
	"os"

	"github.com/FrancoRebollo/async-messaging-svc/internal/domain"
	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/logger"
	"github.com/ulule/limiter/v3"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memoryStore "github.com/ulule/limiter/v3/drivers/store/memory"

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
