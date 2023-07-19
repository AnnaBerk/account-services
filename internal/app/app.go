package app

import (
	"account-management/internal/config"
	sl "account-management/internal/lib/slog"
	"account-management/internal/repo"
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
	repositories := repo.NewRepositories(storage)
	_ = repositories

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
