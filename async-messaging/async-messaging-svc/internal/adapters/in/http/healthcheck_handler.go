package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/FrancoRebollo/async-messaging-svc/internal/domain"
	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/logger"
	"github.com/FrancoRebollo/async-messaging-svc/internal/ports"

	"github.com/gin-gonic/gin"
)

type HealthcheckHandler struct {
	serv ports.HealthcheckService
}

func NewHealthcheckHandler(serv ports.HealthcheckService) *HealthcheckHandler {
	return &HealthcheckHandler{
		serv,
	}
}

// GetHealthcheck verifica el estado del servicio
// @Summary Verifica el estado del servicio
// @Description Devuelve un JSON indicando si el servicio está activo.
// @Tags healthcheck
// @Accept json
// @Produce json
// @Success 200 {object} domain.Healthcheck "Estado del servicio"
// @Failure 400 {object} domain.HealthcheckError "Bad Request"
// @Failure 401 {object} domain.HealthcheckError "Unauthorized"
// @Failure 404 {object} domain.HealthcheckError "Not found"
// @Failure 409 {object} domain.HealthcheckError "Conflict"
// @Failure 500 {object} domain.HealthcheckError "Internal Server Error"
// @Failure 503 {object} domain.HealthcheckError "Service Unavailable"
// @Failure 504 {object} domain.HealthcheckError "Timeout"
// @Security BearerAuth
// @Router /api/healthcheck/ [get]
func (hh *HealthcheckHandler) GetHealthcheck(c *gin.Context) {
	ctx := c.Request.Context()

	healthcheck, err := hh.serv.GetHealthcheck(ctx)
	if err != nil {
		logger.LoggerError().Error(err)
		errorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, healthcheck)
}

func errorResponse(c *gin.Context, err error) {
	if errors.Is(err, context.Canceled) {
		c.JSON(http.StatusRequestTimeout, domain.HealthcheckError{
			Code:    domain.ErrCodeRequestTimeout,
			Message: "Solicitud cancelada"})
		return
	}

	var handlerErr *domain.HealthcheckError
	if errors.As(err, &handlerErr) {
		switch handlerErr.Code {
		case domain.ErrCodeConnectionTimeout:
			c.JSON(http.StatusGatewayTimeout, handlerErr)
			return
		case domain.ErrCodeDeadlockDetected:
			c.JSON(http.StatusServiceUnavailable, handlerErr)
			return
		case domain.ErrCodeDuplicateKey:
			c.JSON(http.StatusConflict, handlerErr)
			return
		case domain.ErrCodeEndOfCommunication:
			c.JSON(http.StatusServiceUnavailable, handlerErr)
			return
		case domain.ErrCodeForeignKeyViolation:
			c.JSON(http.StatusConflict, handlerErr)
			return
		case domain.ErrCodeInvalidInput:
			c.JSON(http.StatusBadRequest, handlerErr)
			return
		case domain.ErrCodeNotNullViolation:
			c.JSON(http.StatusBadRequest, handlerErr)
			return
		case domain.ErrCodeTableOrViewDoesNotExist:
			c.JSON(http.StatusNotFound, handlerErr)
			return
		case domain.ErrCodeInternalServer:
			c.JSON(http.StatusInternalServerError, handlerErr)
			return
		}
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err})
}
