package main

import (
	"avitotest/internal/config"
	"log/slog"
	"os"
)

const (
	envDev = "dev"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("Что")
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
