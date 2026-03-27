package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	devicesapp "github.com/PritOriginal/nethub-back/internal/app/devices"
	"github.com/PritOriginal/nethub-back/internal/config"
	slogger "github.com/PritOriginal/nethub-back/pkg/logger"
)

//	@title			Нетхаб REST API
//	@version		1.0
//	@description	Это документация REST сервиса для тестового задания компании Нетхаб.

//	@BasePath	/api/

//	@tag.name			devices
//	@tag.description	Operations with devices

func main() {
	cfg := config.MustLoad()

	logger, err := slogger.SetupLogger(cfg.Env)
	if err != nil {
		log.Fatalf("error init logger: %v", err)
	}

	app := devicesapp.New(logger, cfg)

	go func() {
		app.MustRun()
	}()

	// Graceful shutdown

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done

	app.Stop()

	logger.Info("server stopped")
}
