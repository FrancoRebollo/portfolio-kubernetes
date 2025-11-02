package http

import (
	"net/http"

	"github.com/FrancoRebollo/api-integration-svc/internal/adapters/in/http/dto"
	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
	"github.com/FrancoRebollo/api-integration-svc/internal/platform/logger"
	"github.com/FrancoRebollo/api-integration-svc/internal/ports"

	"github.com/gin-gonic/gin"
)

type ApiIntegrationHandler struct {
	serv ports.ApiIntegrationService
}

func NewApiIntegrationHandler(serv ports.ApiIntegrationService) *ApiIntegrationHandler {
	return &ApiIntegrationHandler{
		serv,
	}
}

func (h *ApiIntegrationHandler) MakeRequest(c *gin.Context) {
	var req dto.ExternalAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	domainReq := domain.ExternalAPIRequest{
		Method: req.Method,
		URL:    req.URL,
		Params: req.Params,
		Body:   req.Body,
	}

	resp, err := h.serv.ForwardRequest(domainReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to call external API"})
		return
	}

	c.JSON(http.StatusOK, dto.ExternalAPIResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Data:       resp.Data,
	})
}

func (hh *ApiIntegrationHandler) PushEventToQueue(c *gin.Context) {
	ctx := c.Request.Context()

	var reqPushEvent dto.RequestPushEvent
	if err := c.BindJSON(&reqPushEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainEvent := domain.Event{
		ID:         reqPushEvent.ID,
		Type:       reqPushEvent.Type,
		RoutingKey: reqPushEvent.RoutingKey,
		Origin:     reqPushEvent.Origin,
		Timestamp:  reqPushEvent.Timestamp,
		Payload:    reqPushEvent.Payload, // por defecto
	}

	err := hh.serv.PushEventToQueueAPI(ctx, domainEvent)
	if err != nil {
		logger.LoggerError().Error(err)
		errorResponse(c, err)
		return
	}

	responseDefault := dto.ResponseDefault{
		Message: "Mensajo encolado exitosamente",
	}

	c.JSON(http.StatusOK, responseDefault)
}
