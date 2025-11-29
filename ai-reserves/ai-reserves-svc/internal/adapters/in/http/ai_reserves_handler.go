package http

import (
	"net/http"
	"strconv"

	"github.com/FrancoRebollo/ai-reserves-svc/internal/adapters/in/http/dto"
	"github.com/FrancoRebollo/ai-reserves-svc/internal/domain"
	"github.com/FrancoRebollo/ai-reserves-svc/internal/ports"
	"github.com/gin-gonic/gin"
)

type AiReservesHandler struct {
	serv ports.AiReservesService
}

func NewAiReservesHandler(serv ports.AiReservesService) *AiReservesHandler {
	return &AiReservesHandler{
		serv,
	}
}

func newErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, dto.DefaultResponse{
		Code:    400,
		Message: err.Error(),
	})
}

func newSuccessResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, dto.DefaultResponse{
		Code:    200,
		Message: msg,
	})
}

/*
	func (h *AiReservesHandler) CreatePersona(c *gin.Context) {
		var req dto.Persona
		if err := c.BindJSON(&req); err != nil {
			newErrorResponse(c, err)
			return
		}
		domainReq := domain.Persona(req)

		if err := h.serv.CreatePersonaAPI(c, domainReq); err != nil {
			newErrorResponse(c, err)
			return
		}
		newSuccessResponse(c, "Persona creada correctamente")
	}
*/
func (h *AiReservesHandler) UpdAtributoPersona(c *gin.Context) {
	var req dto.PersonaParcial
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.PersonaParcial(req)

	if err := h.serv.UpdAtributoPersonaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Atributo de persona actualizado")
}

func (h *AiReservesHandler) InsertFullConfigPersona(c *gin.Context) {
	var req dto.ConfigPersonaFull
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.ConfigPersonaFull(req)

	if err := h.serv.InsertFullConfigPersonaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Configuration created succesfully")
}

func (h *AiReservesHandler) InsertConfigPersonalSubTipo(c *gin.Context) {
	var req dto.ConfigPersonalSubTipo
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.ConfigPersonalSubTipo(req)

	if err := h.serv.InsertConfigPersonalSubTipoAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Especialities saved")
}

func (h *AiReservesHandler) InsertOrUpdateConfEstablecimiento(c *gin.Context) {
	var req dto.ConfEstablecimiento
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.ConfEstablecimiento(req)

	if err := h.serv.InsertOrUpdateConfEstablecimientoAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Configurations saved")
}

func (h *AiReservesHandler) UpdateConfEstablecimientoField(c *gin.Context) {
	var req dto.ConfigEstablecimiento
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.ConfigEstablecimiento(req)

	if err := h.serv.UpdateConfEstablecimientoFieldAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Configurations saved")
}

func (h *AiReservesHandler) UpdPersona(c *gin.Context) {
	var req dto.Persona
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.Persona(req)

	if err := h.serv.UpdPersonaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Persona actualizada")
}

func (h *AiReservesHandler) UpsertConfigPersona(c *gin.Context) {
	var req dto.ConfigPersona
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}

	domainReq := domain.ConfigPersona(req)

	if err := h.serv.UpsertConfigPersonaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Configuración de persona actualizada")
}

func (h *AiReservesHandler) CreateUnidadReserva(c *gin.Context) {
	var req dto.UnidadReserva
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.UnidadReserva(req)

	if err := h.serv.CreateUnidadReservaAPI(c, &domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Unidad de reserva creada")
}

func (h *AiReservesHandler) CreateTipoUnidadReserva(c *gin.Context) {
	var req dto.TipoUnidadReserva
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.TipoUnidadReserva(req)

	if err := h.serv.CreateTipoUnidadReservaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Tipo de unidad creado")
}

func (h *AiReservesHandler) CreateSubTipoUnidadReserva(c *gin.Context) {
	var req dto.SubTipoUnidadReserva
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.SubTipoUnidadReserva(req)

	if err := h.serv.CreateSubTipoUnidadReservaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Subtipo de unidad creado")
}

func (h *AiReservesHandler) ModifUnidadReserva(c *gin.Context) {
	var req dto.UpdUnidadReserva
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.UpdUnidadReserva(req)

	if err := h.serv.ModifUnidadReservaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Unidad reserva modificada")
}

func (h *AiReservesHandler) ModifUnidadReservaParcial(c *gin.Context) {
	var req dto.UpdAtributeUnidadReserva
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.UpdAtributeUnidadReserva(req)

	if err := h.serv.UpdAtributeUnidadReservaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Unidad reserva modificada")
}

func (h *AiReservesHandler) ModifTipoUnidadReserva(c *gin.Context) {
	var req dto.UpdTipoUnidadReserva
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.UpdTipoUnidadReserva(req)

	if err := h.serv.ModifTipoUnidadReservaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Tipo unidad modificado")
}

func (h *AiReservesHandler) ModifTipoUnidadReservaParcial(c *gin.Context) {
	var req dto.UpdAtributeTipoUnidadReserva
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.UpdAtributeTipoUnidadReserva(req)

	if err := h.serv.UpdAtributeTipoUnidadReservaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Unidad reserva modificada")
}

func (h *AiReservesHandler) ModifSubTipoUnidadReserva(c *gin.Context) {
	var req dto.UpdSubTipoUnidadReserva
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.UpdSubTipoUnidadReserva(req)

	if err := h.serv.ModifSubTipoUnidadReservaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Subtipo unidad modificado")
}

func (h *AiReservesHandler) ModifSubTipoUnidadReservaParcial(c *gin.Context) {
	var req dto.UpdAtributeSubTipoUnidadReserva
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.UpdAtributeSubTipoUnidadReserva(req)

	if err := h.serv.UpdAtributeSubTipoUnidadReservaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Unidad reserva modificada")
}

func (h *AiReservesHandler) CreateReserve(c *gin.Context) {
	var req dto.Reserva
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}
	domainReq := domain.Reserva(req)

	if err := h.serv.CreateReserveAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Reserva creada")
}

func (h *AiReservesHandler) CancelReserve(c *gin.Context) {
	idStr := c.Query("idReserva")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	if err := h.serv.CancelReserveAPI(c, id); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Reserva cancelada")
}

func (h *AiReservesHandler) SearchReserve(c *gin.Context) {
	var req dto.SearchReserve
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}

	domainReq := domain.SearchReserve(req)

	if err := h.serv.SearchReserveAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Búsqueda realizada")
}

func (h *AiReservesHandler) InitAgenda(c *gin.Context) {
	var req dto.Agenda
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}

	domainReq := domain.Agenda(req)

	if err := h.serv.InitAgendaAPI(c, domainReq); err != nil {
		newErrorResponse(c, err)
		return
	}
	newSuccessResponse(c, "Agenda inicializada")
}

func (h *AiReservesHandler) GetInfoPersona(c *gin.Context) {
	idStr := c.Query("idPersona")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	persona, err := h.serv.GetInfoPersonaAPI(c, id)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    persona,
	})
}

func (h *AiReservesHandler) GetReservasPersona(c *gin.Context) {
	var req dto.GetReservaPersona
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, err)
		return
	}

	domainReq := domain.GetReservaPersona(req)

	reservas, err := h.serv.GetReservasPersonaAPI(c, domainReq)

	if err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(200, reservas)
}

func (h *AiReservesHandler) GetReservasUnidadReserva(c *gin.Context) {
	var req dto.GetReservaUnidadReserva
	if err := c.BindQuery(&req); err != nil {
		newErrorResponse(c, err)
		return
	}

	domainReq := domain.GetReservaUnidadReserva(req)

	reservas, err := h.serv.GetReservasUnidadReservaAPI(c, domainReq)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(200, reservas)
}

/*
func (h *AiReservesHandler) MakeRequest(c *gin.Context) {
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

func (hh *AiReservesHandler) PushEventToQueue(c *gin.Context) {
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
*/
