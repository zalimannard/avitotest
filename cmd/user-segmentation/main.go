package main

import (
	"avitotest/internal/config"
	"avitotest/internal/http-server/handlers/history"
	"avitotest/internal/http-server/handlers/segment"
	users_segment "avitotest/internal/http-server/handlers/users-segment"
	mwLogger "avitotest/internal/http-server/middleware/logger"
	"avitotest/internal/http-server/validators"
	"avitotest/internal/lib/logger/sl"
	"avitotest/internal/storage/postgres/schema"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
)

const (
	envDev = "dev"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("Logger initialized")

	storage, err := schema.New(cfg.DbUrl)
	if err != nil {
		log.Error("Failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}
	log.Info("Storage initialized")

	v := validator.New()
	validators.RegisterCustomValidators(v)

	router := setupRouter(log, storage, cfg)
	log.Info("Router initialized")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", sl.Err(err))
		}
	}()

	log.Info("Server started")

	<-done
	log.Info("Stopping server")
}

func setupRouter(log *slog.Logger, storage *schema.Storage, cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/api", func(r chi.Router) {
		r.Route("/segments", func(r chi.Router) {
			r.Post("/", segment.Insert(*log, storage))
			r.Delete("/", segment.Delete(*log, storage))
		})
		r.Route("/users", func(r chi.Router) {
			r.Route("/segments/percent", func(r chi.Router) {
				r.Post("/", users_segment.AssignPercent(*log, storage))
			})
			r.Route("/{userId}", func(r chi.Router) {
				r.Route("/segments", func(r chi.Router) {
					r.Post("/", users_segment.AssignSlugs(*log, storage))
					r.Delete("/", users_segment.DischargeSlugs(*log, storage))
					r.Get("/", users_segment.ReadSlugs(*log, storage))
				})
			})
		})
		r.Route("/reports", func(r chi.Router) {
			r.Get("/", history.ReadHistoryRecordsForMonth(*log, storage, cfg.Address, cfg.ReportDir))
		})
	})
	router.Route("/reports", func(r chi.Router) {
		r.Get("/{fileName}", func(w http.ResponseWriter, r *http.Request) {
			fileName := chi.URLParam(r, "fileName") + ".csv"
			http.ServeFile(w, r, path.Join(cfg.ReportDir, fileName))
		})
	})

	return router
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
