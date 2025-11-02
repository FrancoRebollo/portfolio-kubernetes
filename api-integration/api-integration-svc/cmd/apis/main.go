// cmd/apis/main.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	httpin "github.com/FrancoRebollo/api-integration-svc/internal/adapters/in/http"
	pg "github.com/FrancoRebollo/api-integration-svc/internal/adapters/out/postgres"
	"github.com/FrancoRebollo/api-integration-svc/internal/adapters/rabbitmq"
	"github.com/FrancoRebollo/api-integration-svc/internal/application"
	"github.com/FrancoRebollo/api-integration-svc/internal/platform/config"
	"github.com/FrancoRebollo/api-integration-svc/internal/platform/logger"
	"github.com/FrancoRebollo/api-integration-svc/internal/ports"
)

func main() {
	// 1) Configuración
	cfg, err := config.GetGlobalConfiguration()
	if err != nil {
		logger.LoggerError().Error(err)
		os.Exit(1)
	}

	// 2) Conexiones a bases de datos (según cfg.DB[*].Connection)
	var dbPostgres *pg.PostgresDB

	for _, conf := range cfg.DB {
		switch conf.Connection {
		case "POSTGRES":
			dbPostgres, err = pg.GetInstance(conf)
			if err != nil {
				logger.LoggerError().Errorf("Error conectando a Postgres: %s", err)
				os.Exit(1)
			}
		}
	}

	if dbPostgres != nil {
		logger.LoggerInfo().Info("Conexión a Postgres exitosa")
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 3) Repositorios (adapters out)
	//    Si tus constructores reales difieren, cambiá estas 2 líneas únicamente:
	versionRepository := pg.NewVersionRepository(*dbPostgres)        // <- AJUSTAR si tu firma real difiere
	healthcheckRepository := pg.NewHealthcheckRepository(dbPostgres) // <- AJUSTAR si tu firma real difiere
	apiIntegrationRepository := pg.NewApiIntegrationRepository(dbPostgres)

	fmt.Println("🐇 Iniciando conexión a RabbitMQ...")
	amqpURL := os.Getenv("RABBITMQ_URL")
	fmt.Println("🔗 URL RabbitMQ:", amqpURL)
	rabbitMQAdapter, err := rabbitmq.NewRabbitMQAdapter(amqpURL, os.Getenv("RABBITMQ_QUEUE_EXCHANGE"))
	if err != nil {
		fmt.Println("❌ Error iniciando RabbitMQ:", err)
		os.Exit(1)
	}
	defer rabbitMQAdapter.Close()
	fmt.Println("✅ RabbitMQ inicializado correctamente")

	var messageQueue ports.MessageQueue = rabbitMQAdapter

	// 4) Servicios (application)
	versionService := application.NewVersionService(versionRepository, *cfg.App)             // <- AJUSTAR a tu firma real
	healthcheckService := application.NewHealthcheckService(healthcheckRepository, *cfg.App) // <- AJUSTAR a tu firma real
	apiIntegrationService := application.NewApiIntegrationService(apiIntegrationRepository, *cfg.App, messageQueue, httpClient)

	// 5) Handlers (adapters in/http)
	versionHandler := httpin.NewVersionHandler(versionService) // debe cumplir la interface del router
	healthcheckHandler := httpin.NewHealthcheckHandler(healthcheckService)
	apiIntegrationHandler := httpin.NewApiIntegrationHandler(apiIntegrationService)

	// 6) Router
	rt, err := httpin.NewRouter(cfg.HTTP, versionHandler, *healthcheckHandler, *apiIntegrationHandler)
	if err != nil {
		fmt.Println(err)
	}

	// 7) Server
	address := fmt.Sprintf("%s:%s", cfg.HTTP.Url, cfg.HTTP.Port)
	if err := rt.Listen(address); err != nil {
		logger.LoggerError().Errorf("No se pudo iniciar el servidor: %s", err.Error())
		os.Exit(1)
	}
}
