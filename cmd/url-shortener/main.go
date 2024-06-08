package main

import (
	"URLShortener/auth"
	"URLShortener/internal/config"
	"URLShortener/internal/server/handlers/authhandle"
	"URLShortener/internal/server/handlers/deleteurl"
	"URLShortener/internal/server/handlers/redirect"
	"URLShortener/internal/server/handlers/save"
	mwLogger "URLShortener/internal/server/middleware/logger"
	"URLShortener/internal/storage"
	database "URLShortener/internal/storage/db"
	logUtils "URLShortener/internal/utils/logger"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	SIGNIN = "login"
	SIGNUP = "register"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(os.Getenv("ENV"))
	log.Info("starting url shortener", slog.String("env", os.Getenv("ENV")))
	log.Debug("debug messages are enabled")

	db, err := database.ConnectDB(cfg).InitNewPostgreSQLStorage()
	st := storage.NewStorage(db)

	if err != nil {
		log.Error("failed to init DB ", logUtils.Err(err))
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Error closing db: ", err.Error())
		}
	}()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", auth.WithAuth(save.New(log, st)))
	router.Get("/{alias}", redirect.New(log, st))
	router.Delete("/url/{alias}", auth.WithAuth(deleteurl.New(log, st)))

	router.Route("/users", func(r chi.Router) {
		r.Post("/signup", authhandle.New(log, st, SIGNUP))
		r.Post("/signin", authhandle.New(log, st, SIGNIN))
	})

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	srv := http.Server{
		Addr:         cfg.HTTPServer.Address,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		Handler:      router,
	}

	srv.ListenAndServe()

	log.Error("server stopped")
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
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
