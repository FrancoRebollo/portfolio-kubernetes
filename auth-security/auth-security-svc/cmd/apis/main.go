// cmd/apis/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpin "github.com/FrancoRebollo/auth-security-svc/internal/adapters/in/http"
	eventin "github.com/FrancoRebollo/auth-security-svc/internal/adapters/in/rabbitmq" // 🧠 nuevo
	pg "github.com/FrancoRebollo/auth-security-svc/internal/adapters/out/postgres"
	"github.com/FrancoRebollo/auth-security-svc/internal/adapters/rabbitmq"
	"github.com/FrancoRebollo/auth-security-svc/internal/application"
	"github.com/FrancoRebollo/auth-security-svc/internal/platform/config"
	"github.com/FrancoRebollo/auth-security-svc/internal/platform/logger"
	"github.com/FrancoRebollo/auth-security-svc/internal/ports"
)

func startOutboxWorker(ctx context.Context, svc ports.SecurityService) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := svc.ProcessOutboxEvents(ctx); err != nil {
					fmt.Println("❌ Outbox worker error:", err)
				}
			}
		}
	}()
}

func main() {
	// 1️⃣ Configuración global
	cfg, err := config.GetGlobalConfiguration()
	if err != nil {
		logger.LoggerError().Error(err)
		os.Exit(1)
	}

	// 2️⃣ Conexión a bases de datos
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

	// 3️⃣ Conexión a RabbitMQ
	rmq, err := rabbitmq.NewRabbitMQAdapter(os.Getenv("RABBITMQ_URL"), "")
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()
	var messageQueue ports.MessageQueue = rmq

	// 4️⃣ Repositorios (outbound adapters)
	versionRepository := pg.NewVersionRepository(*dbPostgres)
	healthcheckRepository := pg.NewHealthcheckRepository(dbPostgres)
	securityRepository := pg.NewSecurityRepository(dbPostgres)

	// 5️⃣ Servicios (application layer)
	versionService := application.NewVersionService(versionRepository, *cfg.App)
	healthcheckService := application.NewHealthcheckService(healthcheckRepository, *cfg.App)
	securityService := application.NewSecurityService(securityRepository, *cfg.App, messageQueue)

	// 6️⃣ Handlers HTTP (inbound adapters)
	versionHandler := httpin.NewVersionHandler(versionService)
	healthcheckHandler := httpin.NewHealthcheckHandler(healthcheckService)
	securityHandler := httpin.NewSecurityHandler(securityService)

	// 7️⃣ Iniciar consumer RabbitMQ 🧠 NUEVO BLOQUE
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userConsumer := eventin.NewUserEventConsumer(securityService, rmq)
	queueName := os.Getenv("USER_CREATED_QUEUE")
	if queueName == "" {
		queueName = "user_created_q"
	}

	go userConsumer.Start(ctx, queueName)
	logger.LoggerInfo().Infof("🎧 Listening RabbitMQ queue: %s", queueName)

	// 8️⃣ Señales para cerrar graceful
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		cancel()
		rmq.Close()
		logger.LoggerInfo().Info("🛑 Graceful shutdown consumer")
		os.Exit(0)
	}()

	// 9️⃣ Inicializar Router HTTP
	rt, err := httpin.NewRouter(cfg.HTTP, versionHandler, *healthcheckHandler, *securityHandler)
	if err != nil {
		fmt.Println(err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTP.Url, cfg.HTTP.Port),
		Handler: rt, // el router de Gin
	}

	go func() {
		fmt.Println("🚀 Iniciando servidor HTTP en", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("❌ Error al iniciar servidor:", err)
			os.Exit(1)
		}
	}()

	var svc ports.SecurityService = securityService

	startOutboxWorker(ctx, svc)
	/*
		// 🔟 Servidor HTTP
		address := fmt.Sprintf("%s:%s", cfg.HTTP.Url, cfg.HTTP.Port)
		fmt.Println("🚀 Iniciando servidor HTTP en", address)
		if err := rt.Listen(address); err != nil {
			fmt.Println("❌ Error al iniciar servidor:", err)
			os.Exit(1)
		}
	*/
	fmt.Println("⌛ Esperando ctx.Done()...")
	<-ctx.Done()
	fmt.Println("✅ Microservicio finalizado correctamente")
}
