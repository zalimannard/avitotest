package main

import (
	"avitotest/internal/config"
	"avitotest/internal/lib/logger/sl"
	"avitotest/internal/storage/postgres"
	"log/slog"
	"os"
)

const (
	envDev = "dev"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("Логгер инициализирован")

	storage, err := postgres.New(cfg.DbUrl)
	if err != nil {
		log.Error("Ошибка инициализации хранилища", sl.Err(err))
		os.Exit(1)
	}
	_ = storage
	log.Info("Хранилище инициализировано")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}
