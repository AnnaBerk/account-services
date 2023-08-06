package app

import (
	"account-management/internal/config"
	v1 "account-management/internal/controller/http/v1"
	"account-management/internal/lib/hasher"
	sl "account-management/internal/lib/slog"
	"account-management/internal/repo"
	"account-management/internal/service"
	"account-management/pkg/httpserver"
	"account-management/pkg/psql"
	"account-management/pkg/validator"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run() {
	// Configuration
	cfg := config.MustLoad()

	// Logger
	log := setupLogger(cfg.Env)

	log.Info(
		"starting account-manager",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	// Create PostgreSQL connection string

	storage, err := psql.New(cfg.ConnectionString, psql.MaxPoolSize(cfg.MaxPoolSize))
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	defer storage.Close()

	log.Info("Initializing repositories...")
	repositories := repo.NewRepositories(storage)

	// Services dependencies
	log.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos:    repositories,
		Hasher:   hasher.NewSHA1Hasher(cfg.Hasher.Salt),
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}
	services := service.NewServices(deps, log)

	// Echo handler
	log.Info("Initializing handlers and routes...")
	handler := echo.New()

	// setup handler validator as lib validator
	handler.Validator = validator.NewCustomValidator()
	v1.NewRouter(handler, services, log)

	// HTTP server
	log.Info("Starting http server...")
	log.Debug("Server port", slog.String("port", fmt.Sprint(cfg.DBPort)))
	httpServer := httpserver.New(handler, httpserver.Port(cfg.Port))

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error("app - Run - httpServer.Notify: %w", sl.Err(err))
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error("app - Run - httpServer.Shutdown: %w", sl.Err(err))
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
