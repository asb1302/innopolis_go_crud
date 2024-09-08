package app

import (
	"context"
	"crud/internal/config"
	"crud/internal/handler"
	"crud/internal/pkg/authclient"
	"crud/internal/pkg/server"
	"crud/internal/repository/cache"
	"crud/internal/service"
	"crud/pkg/tracer"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
)

var serviceName = "Crud Service"

func Run() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup

	// initialize dbs
	DB, err := cache.RecipeCacheInit(ctx, &wg)
	if err != nil {
		log.Fatalf("ERROR failed to initialize user database: %v", err)
	}

	config.InitConfig()
	cfg := config.GetConfig()

	traceProv, err := tracer.InitTracer("http://jaeger:14268/api/traces", serviceName)
	if err != nil {
		log.Fatal("init tracer", err)
	}

	authclient.Init(cfg.AuthServiceHost, cfg.AuthServiceTLS)

	// initialize service
	service.Init(DB)

	go func() {
		err := server.Run(":8000", handler.ServerHandler)

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("ERROR server run ", err)
		}
	}()

	log.Println("INFO CRUD service is running")

	<-ctx.Done()

	if err = server.Stop(); err != nil {
		log.Fatal("ERROR server was not gracefully shutdown", err)
	}

	if err := traceProv.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	}

	wg.Wait()

	log.Println("INFO CRUD service was gracefully shutdown")
}
