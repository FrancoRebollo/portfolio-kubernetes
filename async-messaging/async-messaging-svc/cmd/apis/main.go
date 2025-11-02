package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpin "github.com/FrancoRebollo/async-messaging-svc/internal/adapters/in/http"
	pg "github.com/FrancoRebollo/async-messaging-svc/internal/adapters/out/postgres"
	"github.com/FrancoRebollo/async-messaging-svc/internal/adapters/rabbitmq"
	"github.com/FrancoRebollo/async-messaging-svc/internal/application"
	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/config"
	"github.com/FrancoRebollo/async-messaging-svc/internal/ports"
)

func main() {
	fmt.Println("➡️  Iniciando main()")

	if len(os.Args) > 1 && os.Args[1] == "--healthcheck" {
		resp, err := http.Get("http://localhost:3003/api/healthcheck")
		if err != nil {
			fmt.Println("Healthcheck error:", err)
			os.Exit(1)
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Healthcheck failed with status:", resp.StatusCode)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// 1️⃣ Configuración global
	fmt.Println("📦 Cargando configuración global...")
	cfg, err := config.GetGlobalConfiguration()
	if err != nil {
		fmt.Println("❌ Error en configuración:", err)
		os.Exit(1)
	}
	fmt.Println("✅ Configuración cargada correctamente")

	// 2️⃣ Conexión a Postgres
	fmt.Println("🐘 Conectando a Postgres...")
	var dbPostgres *pg.PostgresDB
	for _, conf := range cfg.DB {
		fmt.Println("🔍 Probando conexión:", conf.Connection)
		if conf.Connection == "POSTGRES" {
			dbPostgres, err = pg.GetInstance(conf)
			if err != nil {
				fmt.Println("❌ Error conectando a Postgres:", err)
				os.Exit(1)
			}
		}
	}
	defer dbPostgres.Close()
	fmt.Println("✅ Conexión a Postgres exitosa")

	// 3️⃣ Inicialización de repositorios
	fmt.Println("🧩 Inicializando repositorios...")
	versionRepository := pg.NewVersionRepository(*dbPostgres)
	healthcheckRepository := pg.NewHealthcheckRepository(dbPostgres)
	messageRepository := pg.NewMessageRepository(dbPostgres)
	fmt.Println("✅ Repositorios inicializados")

	// 4️⃣ RabbitMQ adapter (outbound port)
	fmt.Println("🐇 Iniciando conexión a RabbitMQ...")
	amqpURL := os.Getenv("RABBITMQ_URL")
	fmt.Println("🔗 URL RabbitMQ:", amqpURL)
	rabbitMQAdapter, err := rabbitmq.NewRabbitMQAdapter(amqpURL)
	if err != nil {
		fmt.Println("❌ Error iniciando RabbitMQ:", err)
		os.Exit(1)
	}
	defer rabbitMQAdapter.Close()
	fmt.Println("✅ RabbitMQ inicializado correctamente")

	if err = rabbitMQAdapter.InitializeTopology(); err != nil {
		fmt.Println("❌ Error iniciando Topologia de colas:", err)
		os.Exit(1)
	}

	var messageQueue ports.MessageQueue = rabbitMQAdapter

	// 5️⃣ Contexto con cancelación
	fmt.Println("⚙️  Creando contexto de cancelación...")
	ctx, cancel := context.WithCancel(context.Background())

	// Goroutine para señales del sistema
	go func() {
		fmt.Println("🕹️  Escuchando señales del sistema...")
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("🛑 Señal recibida, cerrando servicios...")
		cancel()
	}()

	// 6️⃣ Consumer RabbitMQ
	/*
		fmt.Println("📨 Lanzando consumer de RabbitMQ en goroutine...")
		go func() {
			fmt.Println("➡️  Iniciando consumo de mensajes...")
			_, err := messageQueue.PullEventFromQueue(ctx)
			if err != nil {
				fmt.Println("❌ Error consumiendo mensajes:", err)
			}
		}()
	*/
	// 7️⃣ Servicios de aplicación
	fmt.Println("🧠 Creando servicios de aplicación...")
	versionService := application.NewVersionService(versionRepository, *cfg.App)
	healthcheckService := application.NewHealthcheckService(healthcheckRepository, *cfg.App)
	messageService := application.NewMessageService(messageRepository, messageQueue, *cfg.App)
	fmt.Println("✅ Servicios creados")

	// 8️⃣ Handlers HTTP
	fmt.Println("🌐 Inicializando handlers HTTP...")
	versionHandler := httpin.NewVersionHandler(versionService)
	healthcheckHandler := httpin.NewHealthcheckHandler(healthcheckService)
	messageHandler := httpin.NewMessageHandler(messageService)
	fmt.Println("✅ Handlers listos")

	// 9️⃣ Router
	fmt.Println("🛣️  Creando router HTTP...")
	rt, err := httpin.NewRouter(cfg.HTTP, versionHandler, *healthcheckHandler, *messageHandler)
	if err != nil {
		fmt.Println("❌ Error creando router:", err)
		os.Exit(1)
	}
	fmt.Println("✅ Router creado correctamente")

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
