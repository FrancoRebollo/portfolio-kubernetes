package application

import (
	"context"
	"net/http"

	"github.com/FrancoRebollo/ai-reserves-svc/internal/platform/config"

	"github.com/FrancoRebollo/ai-reserves-svc/internal/domain"
	"github.com/FrancoRebollo/ai-reserves-svc/internal/ports"
)

type AiReservesService struct {
	hr         ports.AiReservesRepository
	conf       config.App
	rmq        ports.MessageQueue
	httpClient *http.Client
}

func NewAiReservesService(hr ports.AiReservesRepository, conf config.App, rmq ports.MessageQueue, httpClient *http.Client) *AiReservesService {
	return &AiReservesService{
		hr,
		conf,
		rmq,
		httpClient,
	}
}

func (hs *AiReservesService) CreatePersonaAPI(ctx context.Context, req domain.PersonCreatedPayload) error {

	if err := hs.hr.CreatePersona(ctx, req); err != nil {
		return err
	}

	return nil
}

func (hs *AiReservesService) UpdAtributoPersonaAPI(ctx context.Context, req domain.PersonaParcial) error {

	if err := hs.hr.UpdAtributoPersona(ctx, req); err != nil {
		return err
	}

	return nil
}
func (hs *AiReservesService) UpdPersonaAPI(ctx context.Context, req domain.Persona) error {

	if err := hs.hr.UpdPersona(ctx, req); err != nil {
		return err
	}

	return nil
}
func (hs *AiReservesService) UpsertConfigPersonaAPI(ctx context.Context, req domain.ConfigPersona) error {

	if err := hs.hr.UpsertConfigPersona(ctx, req); err != nil {
		return err
	}

	return nil
}

func (hs *AiReservesService) CreateUnidadReservaAPI(ctx context.Context, req domain.UnidadReserva) error {
	return nil
}
func (hs *AiReservesService) CreateTipoUnidadReservaAPI(ctx context.Context, req domain.TipoUnidadReserva) error {
	return nil
}
func (hs *AiReservesService) CreateSubTipoUnidadReservaAPI(ctx context.Context, req domain.SubTipoUnidadReserva) error {
	return nil
}

func (hs *AiReservesService) ModifUnidadReservaAPI(ctx context.Context, req domain.UnidadReserva) error {
	return nil
}
func (hs *AiReservesService) ModifTipoUnidadReservaAPI(ctx context.Context, req domain.TipoUnidadReserva) error {
	return nil
}
func (hs *AiReservesService) ModifSubTipoUnidadReservaAPI(ctx context.Context, req domain.SubTipoUnidadReserva) error {
	return nil
}

func (hs *AiReservesService) CreateReserveAPI(ctx context.Context, req domain.Reserva) error {
	return nil
}
func (hs *AiReservesService) CancelReserveAPI(ctx context.Context, idReserva int) error {
	return nil
}
func (hs *AiReservesService) SearchReserveAPI(ctx context.Context, req domain.SearchReserve) error {
	return nil
}
func (hs *AiReservesService) InitAgendaAPI(ctx context.Context, req domain.Agenda) error {
	return nil
}

func (hs *AiReservesService) GetInfoPersonaAPI(ctx context.Context, idPersona int) error {
	return nil
}
func (hs *AiReservesService) GetReservasPersonaAPI(ctx context.Context, req domain.GetReservaPersona) error {
	return nil
}
func (hs *AiReservesService) GetReservasUnidadReservaAPI(ctx context.Context, req domain.GetReservaUnidadReserva) error {
	return nil
}

/*
	func (hs *AiReservesService) ForwardRequest(req domain.ExternalAPIRequest) (domain.ExternalAPIResponse, error) {
		fullURL := req.URL

		// Armar query params para GET
		if req.Method == "GET" && len(req.Params) > 0 {
			query := url.Values{}
			for k, v := range req.Params {
				query.Add(k, v)
			}
			fullURL = fmt.Sprintf("%s?%s", fullURL, query.Encode())
		}

		// Armar body si es POST
		var body io.Reader
		if req.Method == "POST" && req.Body != nil {
			jsonBody, _ := json.Marshal(req.Body)
			body = bytes.NewBuffer(jsonBody)
		}

		httpReq, err := http.NewRequest(req.Method, fullURL, body)
		if err != nil {
			return domain.ExternalAPIResponse{}, err
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := hs.httpClient.Do(httpReq)
		if err != nil {
			return domain.ExternalAPIResponse{}, err
		}
		defer resp.Body.Close()

		var result any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return domain.ExternalAPIResponse{}, err
		}

		return domain.ExternalAPIResponse{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Data:       result,
		}, nil
	}

func (hs *AiReservesService) PushEventToQueueAPI(ctx context.Context, event domain.Event) error {
	fmt.Println("🧩 Iniciando transacción controlada para PushEventToQueueAPI...")

	err := hs.hr.WithTransaction(ctx, func(tx *sql.Tx) error {
		if err := hs.hr.PushEventToQueue(ctx, tx, event); err != nil {
			if errors.Is(err, domain.ErrDuplicateEvent) {
				fmt.Println("⚠️ Evento duplicado detectado, no se publicará en la cola")
				// 👉 devolvemos nil para que NO haya rollback
				return nil
			}
			return err
		}

		fmt.Println("✅ Evento persistido correctamente en DB")

		if err := hs.rmq.Publish(ctx, event); err != nil {
			fmt.Printf("❌ Error al publicar evento %s en cola: %v\n", event.ID, err)
			return err // rollback automático
		}

		fmt.Println("📨 Evento publicado en RabbitMQ correctamente")
		return nil
	})

	if err != nil {
		fmt.Println("🔻 Transacción revertida por error")
		return err
	}

	fmt.Println("✅ Transacción completada con éxito")
	return nil
}
*/
