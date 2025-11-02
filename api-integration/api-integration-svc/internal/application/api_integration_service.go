package application

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/FrancoRebollo/api-integration-svc/internal/platform/config"

	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
	"github.com/FrancoRebollo/api-integration-svc/internal/ports"
)

type ApiIntegrationService struct {
	hr         ports.ApiIntegrationRepository
	conf       config.App
	rmq        ports.MessageQueue
	httpClient *http.Client
}

func NewApiIntegrationService(hr ports.ApiIntegrationRepository, conf config.App, rmq ports.MessageQueue, httpClient *http.Client) *ApiIntegrationService {
	return &ApiIntegrationService{
		hr,
		conf,
		rmq,
		httpClient,
	}
}

func (hs *ApiIntegrationService) ForwardRequest(req domain.ExternalAPIRequest) (domain.ExternalAPIResponse, error) {
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

func (hs *ApiIntegrationService) PushEventToQueueAPI(ctx context.Context, event domain.Event) error {
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
