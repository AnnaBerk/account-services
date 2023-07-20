package app

import (
	"account-management/internal/config"
	"account-management/internal/lib/hasher"
	sl "account-management/internal/lib/slog"
	"account-management/internal/repo"
	"account-management/internal/service"
	"account-management/pkg/psql"
	"golang.org/x/exp/slog"
	"os"
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
	_ = deps

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
