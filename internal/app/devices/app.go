package devicesapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/PritOriginal/nethub-back/internal/config"
	"github.com/PritOriginal/nethub-back/internal/handler"
	deviceshandler "github.com/PritOriginal/nethub-back/internal/handler/devices"
	"github.com/PritOriginal/nethub-back/internal/service/devices"
	"github.com/PritOriginal/nethub-back/internal/storage/postgres"
	slogger "github.com/PritOriginal/nethub-back/pkg/logger"
	"github.com/gin-gonic/gin"
)

type App struct {
	server *http.Server
	log    *slog.Logger
	db     *postgres.Postgres
	router *gin.Engine
	port   int
}

func New(log *slog.Logger, cfg *config.Config) *App {
	postgresDB, err := postgres.New(cfg.DB)
	if err != nil {
		log.Error("failed connection to database", slogger.Err(err))
		panic(err)
	}
	log.Info("PostgreSQL connected!")

	r := handler.GetRouter(log, cfg.Env)

	devicesStorage := postgres.NewDeviceStorage(postgresDB.DB)
	devicesService := devices.New(log, devicesStorage)
	deviceshandler.Register(r, log, devicesService)

	server := &http.Server{
		Addr:         cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.Timeout.Read,
		WriteTimeout: cfg.Server.Timeout.Write,
		IdleTimeout:  cfg.Server.Timeout.Idle,
	}

	return &App{
		server: server,
		log:    log,
		db:     postgresDB,
		router: r,
		port:   cfg.Server.Port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "rest.Run"

	a.log.Info("server started", slog.String("address", ":"+strconv.Itoa(a.port)))
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.log.Error("failed to start server")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "rest.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping REST server", slog.Int("port", a.port))

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.log.Error("an error occurred while stopping the server", slogger.Err(err))
	}

	if err := a.db.DB.Close(); err != nil {
		a.log.Error("an error occurred while closing the connection to the database", slogger.Err(err))
	}
}
