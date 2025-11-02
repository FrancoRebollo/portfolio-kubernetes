package http

import (
	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
	"github.com/FrancoRebollo/auth-security-svc/internal/platform/logger"

	"github.com/FrancoRebollo/auth-security-svc/internal/ports"

	"github.com/gin-gonic/gin"
)

var _ domain.Version

type Versionhandler struct {
	ports.VersionService
}

func NewVersionHandler(service ports.VersionService) *Versionhandler {
	return &Versionhandler{service}
}

// GetVersion retorna la versión actual de la API
// @Summary Obtiene la versión de la API
// @Description Devuelve un JSON con la versión actual de la API.
// @Tags version
// @Accept json
// @Produce json
// @Success 200 {object} domains.Version "Versión de la API"
// @Security BearerAuth
// @Router /api/version/ [get]
func (h *Versionhandler) GetVersion(c *gin.Context) {
	version_api, err := h.GetVersionAPI(c)
	if err != nil {
		logger.LoggerError().Error(err)
		c.JSON(500, gin.H{"message": "Error que viene del servicio"})
		return
	}

	c.JSON(200, version_api)
}
