package middlewares

import (
	"bytes"
	"io"
	"net/http"

	"github.com/FrancoRebollo/api-integration-svc/internal/adapters/in/http/validators"

	"github.com/FrancoRebollo/api-integration-svc/internal/platform/logger"

	"github.com/gin-gonic/gin"
)

func ValidateGetHealthcheck(c *gin.Context) {
	query := c.Request.URL.Query()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.LoggerError().Error(err)
	}

	err = validators.ValidateEmptyQuery(query)
	if err != nil {
		logger.LoggerError().Error(err)
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
	}

	err = validators.ValidateEmptyBody(body)
	if err != nil {
		logger.LoggerError().Error(err)
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
	}

	// Reestablecer formato de body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	c.Next()
}
